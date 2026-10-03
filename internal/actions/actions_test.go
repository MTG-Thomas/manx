// Package actions - internal test for parity + audit contract.
// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

// The parity lists live in parity.go (single source; shared with non-test code).
// (SPEC_VERBS / SPEC_DESTRUCTIVE defined there.)

func TestParity(t *testing.T) {
	r := Default()
	got := []string{}
	for _, a := range r.List() {
		got = append(got, a.Verb)
	}
	if strings.Join(got, ",") != strings.Join(SPEC_VERBS, ",") {
		t.Fatalf("registry verbs %v differ from spec %v", got, SPEC_VERBS)
	}
	for _, verb := range SPEC_VERBS {
		a, ok := r.Get(verb)
		if !ok {
			t.Fatalf("verb %s missing from registry", verb)
		}
		if a.Danger != SPEC_DESTRUCTIVE[verb] {
			t.Fatalf("verb %s: registered Danger=%v, spec says %v", verb, a.Danger, SPEC_DESTRUCTIVE[verb])
		}
	}
}

// Gate test: destructive verbs must be refused WITHOUT an explicit --i-know,
// and every invocation (refusal or proceed) must write exactly one audit row.
func TestGate(t *testing.T) {
	dir := t.TempDir()
	r := Default()
	aw := r.NewAuditor(filepath.Join(dir, "audit.log"))

	rowCount := func() int {
		n := aw.Count()
		t.Logf("audit rows so far: %d", n)
		return n
	}

	// destructive verb refused without --i-know
	err := aw.Run("img-out", []string{})
	if !errors.Is(err, ErrGateRefused) {
		t.Fatalf("img-out without i-know: got error %v, want ErrGateRefused", err)
	}
	n0 := rowCount()
	if n0 < 1 {
		t.Fatalf("refusal must write an audit row")
	}

	// with --i-know it proceeds into the action body (not-implemented is the
	// expected skeleton behavior and still writes its own audit row)
	err2 := aw.Run("img-out", []string{"--i-know"})
	if err2 == nil {
		t.Fatalf("img-out body should not be implemented yet; got nil error")
	}
	if !strings.Contains(err2.Error(), "not implemented") {
		t.Fatalf("img-out with i-know: got %v", err2)
	}
	n1 := rowCount()
	if n1 < n0+1 {
		t.Fatalf("i-know path must write a second audit row; saw %d then %d", n0, n1)
	}
}

func TestRegistryNotEmpty(t *testing.T) {
	r := Default()
	if len(r.List()) == 0 {
		t.Fatal("registry empty: parity contract violation")
	}
}
