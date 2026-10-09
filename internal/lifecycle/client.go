// SPDX-License-Identifier: AGPL-3.0-or-later
// Package lifecycle implements restricted device enrollment, never rescue execution.
package lifecycle

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/MTG-Thomas/manx/internal/lifecycle/contract"
)

const Version = "bifrost.device-lifecycle/v1"

type Checkpoint struct {
	JournalVersion int               `json:"journal_version"`
	InstallationID string            `json:"installation_id"`
	Endpoint       string            `json:"endpoint"`
	Registration   json.RawMessage   `json:"registration"`
	Session        *contract.Session `json:"session"`
	Pending        []json.RawMessage `json:"pending"`
	Acknowledged   int64             `json:"acknowledged"`
	BootID         string            `json:"boot_id"`
}

type Client struct {
	Directory string
	State     Checkpoint
	key       *ecdsa.PrivateKey
	HTTP      *http.Client
}

func UUID() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return ""
	}
	value[6] = (value[6] & 15) | 64
	value[8] = (value[8] & 63) | 128
	return fmt.Sprintf("%x-%x-%x-%x-%x", value[:4], value[4:6], value[6:8], value[8:10], value[10:])
}

func Open(directory, endpoint string) (*Client, error) {
	if endpoint == "" {
		data, err := os.ReadFile(filepath.Join(directory, "checkpoint.json"))
		var state Checkpoint
		if err != nil || json.Unmarshal(data, &state) != nil {
			return nil, errors.New("enrollment endpoint required")
		}
		endpoint = state.Endpoint
	}
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || (parsed.Path != "" && parsed.Path != "/") {
		return nil, errors.New("HTTPS enrollment origin required")
	}
	if err := os.MkdirAll(directory, 0700); err != nil {
		return nil, errors.New("session directory unavailable")
	}
	info, err := os.Lstat(directory)
	if err != nil || !info.IsDir() || info.Mode().Perm()&0077 != 0 {
		return nil, errors.New("session directory must be private (0700)")
	}
	c := &Client{Directory: directory, HTTP: &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}}
	statePath := filepath.Join(directory, "checkpoint.json")
	data, err := os.ReadFile(statePath)
	if err == nil {
		if err := json.Unmarshal(data, &c.State); err != nil {
			return nil, errors.New("invalid session checkpoint")
		}
		if c.State.JournalVersion != 1 {
			return nil, errors.New("unsupported journal version")
		}
		if c.State.Endpoint != strings.TrimRight(endpoint, "/") {
			return nil, errors.New("endpoint change requires review")
		}
	} else if os.IsNotExist(err) {
		c.State = Checkpoint{JournalVersion: 1, InstallationID: UUID(), Endpoint: strings.TrimRight(endpoint, "/"), BootID: UUID()}
		if err := c.save(); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("session checkpoint unavailable")
	}
	keyPath := filepath.Join(directory, "key.der")
	keyData, err := os.ReadFile(keyPath)
	if os.IsNotExist(err) {
		if len(c.State.Registration) != 0 || c.State.Session != nil {
			return nil, errors.New("session key missing; manual intervention required")
		}
		c.key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, errors.New("key generation failed")
		}
		keyData, err = x509.MarshalPKCS8PrivateKey(c.key)
		if err != nil {
			return nil, errors.New("key encoding failed")
		}
		if err := atomicWrite(keyPath, keyData); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, errors.New("session key unavailable")
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(keyData)
	if err != nil {
		return nil, errors.New("session key invalid")
	}
	var ok bool
	c.key, ok = parsedKey.(*ecdsa.PrivateKey)
	if !ok || c.key.Curve != elliptic.P256() {
		return nil, errors.New("session key invalid")
	}
	if len(c.State.Registration) == 0 {
		public, err := x509.MarshalPKIXPublicKey(&c.key.PublicKey)
		if err != nil {
			return nil, errors.New("public key unavailable")
		}
		architecture := map[string]string{"amd64": "x64", "arm64": "arm64", "386": "x86"}[runtime.GOARCH]
		if architecture == "" {
			architecture = "unknown"
		}
		observed := func(path string) string { value, _ := os.ReadFile(path); return strings.TrimSpace(string(value)) }
		registration := contract.Registration{ContractVersion: Version, InstallationId: c.State.InstallationID, ClientKind: "manx_recovery", PublicKey: base64.StdEncoding.EncodeToString(public),
			SerialNumber: observed("/sys/class/dmi/id/product_serial"), SmbiosUuid: observed("/sys/class/dmi/id/product_uuid"), Manufacturer: observed("/sys/class/dmi/id/sys_vendor"), Model: observed("/sys/class/dmi/id/product_name"),
			Architecture: architecture, BootEnvironment: "linux_rescue", ClientVersion: "0.1.0", ObservedAt: time.Now().UTC().Format(time.RFC3339Nano)}
		if err := registration.Validate(); err != nil {
			return nil, err
		}
		c.State.Registration, err = json.Marshal(registration)
		if err != nil {
			return nil, err
		}
		if err := c.save(); err != nil {
			return nil, err
		}
	}
	if boot, err := os.ReadFile("/proc/sys/kernel/random/boot_id"); err == nil {
		current := strings.TrimSpace(string(boot))
		if c.State.BootID != current {
			c.State.BootID = current
			if err := c.save(); err != nil {
				return nil, err
			}
		}
	}
	return c, nil
}

func atomicWrite(path string, value []byte) error {
	file, err := os.OpenFile(path+".pending", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return errors.New("journal write failed")
	}
	if _, err = file.Write(value); err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return errors.New("journal write failed")
	}
	if err := os.Rename(path+".pending", path); err != nil {
		return errors.New("journal commit failed")
	}
	directory, err := os.Open(filepath.Dir(path))
	if err != nil {
		return errors.New("journal directory unavailable")
	}
	defer directory.Close()
	if err := directory.Sync(); err != nil {
		return errors.New("journal directory flush failed")
	}
	return nil
}
func (c *Client) save() error {
	value, err := json.Marshal(c.State)
	if err != nil {
		return err
	}
	return atomicWrite(filepath.Join(c.Directory, "checkpoint.json"), value)
}

func (c *Client) request(ctx context.Context, method, path string, body []byte, registration bool) ([]byte, error) {
	nonce := UUID()
	if !registration {
		if c.State.Session == nil {
			return nil, errors.New("session not registered")
		}
		request, _ := http.NewRequestWithContext(ctx, "POST", c.State.Endpoint+"/v1/sessions/"+c.State.Session.SessionId+"/challenge", nil)
		response, err := c.HTTP.Do(request)
		if err != nil {
			return nil, errors.New("enrollment network unavailable")
		}
		value, readErr := io.ReadAll(io.LimitReader(response.Body, 16385))
		response.Body.Close()
		if readErr != nil || response.StatusCode != 200 || len(value) > 16384 {
			return nil, errors.New("challenge unavailable")
		}
		var challenge struct {
			Challenge string `json:"challenge"`
		}
		if json.Unmarshal(value, &challenge) != nil || challenge.Challenge == "" {
			return nil, errors.New("challenge invalid")
		}
		nonce = challenge.Challenge
	}
	stamp := fmt.Sprintf("%d", time.Now().Unix())
	hash := sha256.Sum256(body)
	frame := strings.Join([]string{Version, method, path, c.State.InstallationID, stamp, nonce, hex.EncodeToString(hash[:])}, "\n")
	hashed := sha256.Sum256([]byte(frame))
	r, s, err := ecdsa.Sign(rand.Reader, c.key, hashed[:])
	if err != nil {
		return nil, errors.New("proof signing failed")
	}
	signature := make([]byte, 64)
	r.FillBytes(signature[:32])
	s.FillBytes(signature[32:])
	request, err := http.NewRequestWithContext(ctx, method, c.State.Endpoint+path, bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("request invalid")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Lifecycle-Installation", c.State.InstallationID)
	request.Header.Set("X-Lifecycle-Timestamp", stamp)
	request.Header.Set("X-Lifecycle-Nonce", nonce)
	request.Header.Set("X-Lifecycle-Signature", base64.StdEncoding.EncodeToString(signature))
	response, err := c.HTTP.Do(request)
	if err != nil {
		return nil, errors.New("enrollment network unavailable")
	}
	defer response.Body.Close()
	value, err := io.ReadAll(io.LimitReader(response.Body, 4*1048576+1))
	if err != nil || len(value) > 4*1048576 {
		return nil, errors.New("response exceeded bound")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("session request rejected (HTTP %d)", response.StatusCode)
	}
	return value, nil
}

func (c *Client) Register(ctx context.Context) error {
	response, err := c.request(ctx, "POST", "/v1/sessions", c.State.Registration, true)
	if err != nil {
		return err
	}
	var session contract.Session
	decoder := json.NewDecoder(bytes.NewReader(response))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&session); err != nil {
		return errors.New("session response invalid")
	}
	if err := session.Validate(); err != nil {
		return err
	}
	var registration contract.Registration
	if json.Unmarshal(c.State.Registration, &registration) != nil {
		return errors.New("local registration invalid")
	}
	if !reflect.DeepEqual(session.Registration, registration) {
		return errors.New("session identity conflict")
	}
	if c.State.Session != nil && c.State.Session.SessionId != session.SessionId {
		return errors.New("session identity conflict")
	}
	c.State.Session = &session
	return c.save()
}

func (c *Client) Record(diagnostic string) error {
	if diagnostic != "hardware_inventory" {
		return errors.New("unsupported diagnostic")
	}
	if len(c.State.Pending) >= 128 {
		return errors.New("local journal full")
	}
	// Offline observations can exist before registration; assign the server session
	// locator once registration succeeds, without changing logical event identity.
	sessionID := strings.Repeat("0", 64)
	if c.State.Session != nil {
		sessionID = c.State.Session.SessionId
	}
	boot := c.State.BootID
	if value, err := os.ReadFile("/proc/sys/kernel/random/boot_id"); err == nil {
		boot = strings.TrimSpace(string(value))
	}
	event := contract.ClientEvent{ContractVersion: Version, SessionId: sessionID, InstallationId: c.State.InstallationID, Sequence: c.State.Acknowledged + int64(len(c.State.Pending)) + 1, EventId: UUID(), EventType: "diagnostic", Timestamp: time.Now().UTC().Format(time.RFC3339Nano), Payload: contract.EventPayload{BootId: boot, Checkpoint: "diagnostic_collected", Result: "ok", Diagnostic: diagnostic}}
	observations := map[string]string{"observed_at": event.Timestamp, "boot_id": boot}
	for name, path := range map[string]string{"serial_number": "/sys/class/dmi/id/product_serial", "smbios_uuid": "/sys/class/dmi/id/product_uuid", "manufacturer": "/sys/class/dmi/id/sys_vendor", "model": "/sys/class/dmi/id/product_name"} {
		if data, err := os.ReadFile(path); err == nil {
			observations[name] = strings.TrimSpace(string(data))
		}
	}
	if len(observations) == 2 {
		event.Payload.Result = "unavailable"
	}
	if err := event.Validate(); err != nil {
		return err
	}
	evidence, err := json.Marshal(observations)
	if err != nil {
		return errors.New("hardware evidence encoding failed")
	}
	if err := atomicWrite(filepath.Join(c.Directory, "evidence-"+event.EventId+".json"), evidence); err != nil {
		return err
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		return err
	}
	c.State.Pending = append(c.State.Pending, encoded)
	return c.save()
}

func (c *Client) reconcile(response []byte) error {
	var session contract.Session
	decoder := json.NewDecoder(bytes.NewReader(response))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&session) != nil || session.Validate() != nil {
		return errors.New("session response invalid")
	}
	if c.State.Session == nil || session.SessionId != c.State.Session.SessionId || !reflect.DeepEqual(session.Registration, c.State.Session.Registration) {
		return errors.New("session identity conflict")
	}
	for len(c.State.Pending) > 0 {
		var pending contract.ClientEvent
		if json.Unmarshal(c.State.Pending[0], &pending) != nil {
			return errors.New("local event invalid")
		}
		if pending.Sequence > session.LastSequence {
			break
		}
		matched := false
		for _, event := range session.Events {
			if event.Sequence == pending.Sequence && reflect.DeepEqual(event, pending) {
				matched = true
				break
			}
		}
		if !matched {
			return errors.New("event reconciliation conflict")
		}
		c.State.Acknowledged = pending.Sequence
		c.State.Pending = c.State.Pending[1:]
	}
	if session.LastSequence != c.State.Acknowledged {
		return errors.New("event journal conflict")
	}
	c.State.Session = &session
	return c.save()
}

func (c *Client) Sync(ctx context.Context) error {
	if c.State.Session == nil {
		if err := c.Register(ctx); err != nil {
			return err
		}
	}
	for index, raw := range c.State.Pending {
		var event contract.ClientEvent
		if json.Unmarshal(raw, &event) != nil {
			return errors.New("local event invalid")
		}
		if event.SessionId == strings.Repeat("0", 64) {
			event.SessionId = c.State.Session.SessionId
			encoded, err := json.Marshal(event)
			if err != nil {
				return err
			}
			c.State.Pending[index] = encoded
		}
	}
	if err := c.save(); err != nil {
		return err
	}
	path := "/v1/sessions/" + c.State.Session.SessionId
	response, err := c.request(ctx, "GET", path, nil, false)
	if err != nil {
		return err
	}
	if err := c.reconcile(response); err != nil {
		return err
	}
	for len(c.State.Pending) > 0 {
		response, err := c.request(ctx, "POST", path+"/events", c.State.Pending[0], false)
		if err != nil {
			return err
		}
		if err := c.reconcile(response); err != nil {
			return err
		}
	}
	heartbeat, _ := json.Marshal(contract.Heartbeat{ContractVersion: Version, InstallationId: c.State.InstallationID, BootId: c.State.BootID})
	response, err = c.request(ctx, "POST", path+"/heartbeat", heartbeat, false)
	if err != nil {
		return err
	}
	return c.reconcile(response)
}

func (c *Client) Status() map[string]any {
	result := map[string]any{"installation_id": c.State.InstallationID, "state": "not_registered", "pending_events": len(c.State.Pending), "acknowledged_sequence": c.State.Acknowledged}
	if c.State.Session != nil {
		result["session_id"] = c.State.Session.SessionId
		result["pairing_code"] = c.State.Session.PairingCode
		result["state"] = c.State.Session.State
		result["last_seen_at"] = c.State.Session.LastSeenAt
	}
	return result
}
