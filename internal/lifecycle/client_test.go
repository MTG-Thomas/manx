// SPDX-License-Identifier: AGPL-3.0-or-later
package lifecycle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/MTG-Thomas/manx/internal/lifecycle/contract"
)

type transport func(*http.Request) (*http.Response, error)

func (fn transport) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }

func TestOfflineJournalSurvivesRestart(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "private")
	client, err := Open(dir, "https://lab.invalid")
	if err != nil {
		t.Fatal(err)
	}
	client.HTTP.Transport = transport(func(*http.Request) (*http.Response, error) { return nil, errors.New("offline") })
	if err := client.Record("hardware_inventory"); err != nil {
		t.Fatal(err)
	}
	if err := client.Sync(context.Background()); err == nil {
		t.Fatal("offline sync succeeded")
	}
	restarted, err := Open(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	if restarted.State.InstallationID != client.State.InstallationID || len(restarted.State.Pending) != 1 {
		t.Fatal("offline evidence/identity lost")
	}
	if _, err := Open(dir, "https://other.invalid"); err == nil {
		t.Fatal("endpoint silently changed")
	}
	if err := restarted.Record("execute_command"); err == nil {
		t.Fatal("untyped operation accepted")
	}
	info, _ := os.Stat(filepath.Join(dir, "key.der"))
	if info.Mode().Perm() != 0600 {
		t.Fatal("key permissions")
	}
}

func TestUnknownAcknowledgmentNeverErasesEvent(t *testing.T) {
	client, err := Open(filepath.Join(t.TempDir(), "private"), "https://lab.invalid")
	if err != nil {
		t.Fatal(err)
	}
	var registration contract.Registration
	_ = json.Unmarshal(client.State.Registration, &registration)
	session := contract.Session{ContractVersion: Version, SessionId: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Registration: registration, IdentitySha256: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", PairingCode: "ABCDEFGH", State: "approved", Revision: 2, CreatedAt: registration.ObservedAt, ExpiresAt: "2026-10-20T12:00:00Z", LastSeenAt: registration.ObservedAt, Events: []contract.ClientEvent{}, OperatorEvents: []contract.OperatorEvent{}}
	client.State.Session = &session
	if err := client.Record("hardware_inventory"); err != nil {
		t.Fatal(err)
	}
	var event contract.ClientEvent
	_ = json.Unmarshal(client.State.Pending[0], &event)
	event.EventId = UUID()
	session.Events = []contract.ClientEvent{event}
	session.LastSequence = 1
	response, _ := json.Marshal(session)
	if err := client.reconcile(response); err == nil {
		t.Fatal("conflicting ack accepted")
	}
	if len(client.State.Pending) != 1 {
		t.Fatal("pending event erased")
	}
}

func TestCanonicalFixtures(t *testing.T) {
	data, err := os.ReadFile("contract/fixtures.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures []struct {
		Name  string
		Model string
		Valid bool
		Value json.RawMessage
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range fixtures {
		t.Run(fixture.Name, func(t *testing.T) {
			var value interface{ Validate() error }
			switch fixture.Model {
			case "Registration":
				value = &contract.Registration{}
			case "ClientEvent":
				value = &contract.ClientEvent{}
			case "Session":
				value = &contract.Session{}
			default:
				t.Fatal("fixture model")
			}
			decoder := json.NewDecoder(bytes.NewReader(fixture.Value))
			decoder.DisallowUnknownFields()
			err := decoder.Decode(value)
			if err == nil {
				err = value.Validate()
			}
			if (err == nil) != fixture.Valid {
				t.Fatalf("valid=%v error=%v", fixture.Valid, err)
			}
		})
	}
}

func TestRegistrationRetryPreservesBodyAndIdentity(t *testing.T) {
	client, err := Open(filepath.Join(t.TempDir(), "private"), "https://lab.invalid")
	if err != nil {
		t.Fatal(err)
	}
	original := append([]byte(nil), client.State.Registration...)
	client.HTTP.Transport = transport(func(request *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(request.Body)
		if !bytes.Equal(body, original) {
			t.Fatal("registration changed")
		}
		return nil, errors.New("response lost after acceptance")
	})
	if client.Register(context.Background()) == nil {
		t.Fatal("lost response accepted")
	}
	restarted, err := Open(client.Directory, "")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(restarted.State.Registration, original) {
		t.Fatal("retry created new installation identity")
	}
}
