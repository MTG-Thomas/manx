// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectLines(t *testing.T) {
	lspci := "" +
		`00:00.0 Host bridge: Intel Corp [8086:1234]` + "\n" +
		"	Subsystem: Dell [1028:5678]\n" +
		"	Kernel driver in use: skylake_pcie\n" +
		`03:00.0 RAID bus controller: LSI [1000:5678]` + "\n" +
		"	Subsystem: Dell [1028:5678]\n"
	got := unclaimedLines(lspci)
	if len(got) != 1 {
		t.Fatalf("want exactly 1 unclaimed line (the LSI without driver), got %v", got)
	}
	if !strings.Contains(got[0], "UNCLAIMED") || !strings.Contains(got[0], "RAID") {
		t.Fatalf("unclaimed line missing the device text: %q", got[0])
	}
	// a fully-driven device must NOT appear (skylake has a driver)
	if strings.Contains(strings.Join(got, "\n"), "Host bridge") {
		t.Fatal("driven devices must be excluded from the UNCLAIMED report")
	}
}

func TestComposeStatus(t *testing.T) {
	s := composeStatus("hostA", "14393.9503", "ens18 172.16.15.171/24", "1234 seconds", 4, true)
	for _, want := range []string{"hostA", "14393.9503", "ens18 172.16.15.171/24", "1234 seconds", "audit rows: 4", "present (/toolkit)"} {
		if !strings.Contains(s, want) {
			t.Fatalf("status missing %q:\n%s", want, s)
		}
	}
	s2 := composeStatus("hostA", "windows", "no ip", "-", 0, false)
	if strings.Contains(s2, "present (/toolkit)") {
		t.Fatal("toolkit reported present while flag says absent")
	}
}

func TestDefaultAuditPathEnv(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("MANX_AUDIT", filepath.Join(dir, "a.log"))
	if got := DefaultAuditPath(); got != filepath.Join(dir, "a.log") {
		t.Fatalf("env override not honored: %q", got)
	}
}

func TestRunViewWritesRowInView(t *testing.T) {
	dir := t.TempDir()
	auditPath := filepath.Join(dir, "audit.log")
	gr := Default().NewAuditor(auditPath)
	if err := gr.RunView("tui", "status", nil); err != nil {
		t.Fatalf("status via tui: %v", err)
	}
	b, err := os.ReadFile(auditPath)
	if err != nil {
		t.Fatal(err)
	}
	row := strings.TrimSpace(string(b))
	if !strings.Contains(row, `"view":"tui"`) {
		t.Fatalf("row not tagged with tui view: %s", row)
	}
}

// NOTE: the Writer row-contract is covered in internal/audit/audit_test.go
// (Count, JSON shape, nested-path refusal) to keep one authority per package.
