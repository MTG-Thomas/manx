// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"github.com/MTG-Thomas/manx/internal/actions"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestViewsUseRegistryAndSessionAudit(t *testing.T) {
	registry := actions.Default()
	views := initialVerbs()
	if len(views) != len(registry.List()) {
		t.Fatal("TUI action surface differs from registry")
	}
	for index, action := range registry.List() {
		if views[index].verb != action.Verb || views[index].danger != action.Danger {
			t.Fatal("TUI action metadata drift")
		}
	}
	directory := t.TempDir()
	auditPath := filepath.Join(directory, "audit.log")
	m := model{gr: registry.NewAuditor(auditPath)}
	output, err := m.runVerbArgs("session", []string{"status", "--endpoint", "https://synthetic.invalid", "--state-dir", filepath.Join(directory, "state")})
	if err != nil || !strings.Contains(output, "not_registered") {
		t.Fatalf("session status failed: %v", err)
	}
	data, err := os.ReadFile(auditPath)
	if err != nil || !strings.Contains(string(data), `"view":"tui"`) || !strings.Contains(string(data), `"verb":"session"`) {
		t.Fatalf("missing shared audit evidence: %s, %v", data, err)
	}
}
