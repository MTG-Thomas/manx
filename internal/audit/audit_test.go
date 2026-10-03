// Package audit - test the row contract.
// SPDX-License-Identifier: AGPL-3.0-or-later
package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteRowJSON(t *testing.T) {
	dir := t.TempDir()
	w := Writer{Path: filepath.Join(dir, "audit.log")}

	rows := []Row{
		{Actor: "operator", View: "cli", Verb: "status", Result: "ok"},
		{Actor: "agent:ops", View: "harness", Verb: "img-out", Danger: "denied", Result: "error:refused"},
		{Actor: "operator", View: "menu", Verb: "img-out", Danger: "i-know", Result: "ok"},
	}
	for _, r := range rows {
		if err := w.Write(r); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	b, err := os.ReadFile(w.Path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if got := w.Count(); got != len(rows) {
		t.Fatalf("Count = %d, want %d", got, len(rows))
	}

	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) != len(rows) {
		t.Fatalf("file lines %d, want %d", len(lines), len(rows))
	}
	var row Row
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &row); err != nil {
		t.Fatalf("last row JSON: %v", err)
	}
	if row.Verb != "img-out" || row.View != "menu" || row.Danger != "i-know" {
		t.Fatalf("last row mismatch: %+v", row)
	}
}

func TestWriterRefusesNestedPath(t *testing.T) {
	dir := t.TempDir()
	w := Writer{Path: filepath.Join(dir, "nested", "audit.log")}
	err := w.Write(Row{Actor: "x"})
	if err == nil {
		t.Fatal("expected error for nested path whose parent does not exist")
	}
	t.Logf("correct refusal: %v", err)
}
