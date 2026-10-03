// Package audit implements the audit-row contract for MANX action verbs.
//
// SPDX-License-Identifier: AGPL-3.0-or-later
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Row is one line of the append-only audit log.
// Every destructive action (and hence every invocation of any view: whiptail
// menu, CLI, TUI, or agent harness) writes exactly one row, in the same shape,
// regardless of which view caused the action.
type Row struct {
	TS        string `json:"ts"`         // RFC3339Nano, UTC
	Actor     string `json:"actor"`      // "operator", "agent:<name>", "harness"
	View      string `json:"view"`       // "cli", "menu", "tui", "harness", "unknown"
	Verb      string `json:"verb"`       // canonical action name, spec contract
	Danger    string `json:"danger"`     // "", "i-know", "denied"
	Result    string `json:"result"`     // "", "ok", "error:<code>"
	DetailOID string `json:"detail_oid"` // hash-like reference; for compactness, no per-line bodies
}

// Writer appends Rows to a file. The file must not be required to exist before
// first write; the parent dir must already exist.
type Writer struct {
	Path string
}

// Count returns the number of complete audit rows currently in the file
// (used by the parity-contract tests).
func (w Writer) Count() int {
	b, err := os.ReadFile(w.Path)
	if err != nil {
		return 0
	}
	n := 0
	for _, ln := range b {
		if ln == '\n' {
			n++
		}
	}
	return n
}

func (w Writer) Write(r Row) error {
	if r.TS == "" {
		r.TS = time.Now().UTC().Format(time.RFC3339Nano)
	}
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(w.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", b)
	return err
}
