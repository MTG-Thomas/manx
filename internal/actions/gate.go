// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

import (
	"errors"
	"fmt"

	"github.com/MTG-Thomas/manx/internal/audit"
)

// ErrGateRefused is returned when a destructive verb was invoked without the
// explicit --i-know confirmation marker in ANY view (spec §10.8).
// All gated refusals still write exactly one audit row with Danger="denied".
var ErrGateRefused = errors.New("destructive action requires --i-know in any view (spec §10.8)")

// auditor writes contracts-aware audit rows for gate decisions and verb results.
type auditor struct {
	w audit.Writer
}

// NewAuditor creates an action-surface auditor for the given audit file path.
func (r *Registry) NewAuditor(path string) *GateRunner {
	return &GateRunner{reg: r, w: audit.Writer{Path: path}}
}

// GateRunner drives verbs through the --i-know gate, always auditing.
// Every view (CLI, whiptail menu, TUI, harness) must express itself through
// this runner, so the parity contract (one row per invocation, either fate)
// can't be bypassed accidentally.
type GateRunner struct {
	reg *Registry
	w   audit.Writer
}

// Count reports the number of audit rows written so far (test-side helper for
// the "exactly one row per invocation" contract).
func (g *GateRunner) Count() int { return g.w.Count() }

// Run executes a verb in the CLI view (spec default); other views use RunView.
func (g *GateRunner) Run(verb string, args []string) error {
	return g.RunView("cli", verb, args)
}

// RunView executes a verb for a named view (cli/menu/tui/harness), enforcing
// the --i-know gate for destructive verbs and writing exactly one audit row
// tagged with the view that caused the invocation.
func (g *GateRunner) RunView(view, verb string, args []string) error {
	row := audit.Row{Actor: "operator", View: view, Verb: verb}

	a, ok := g.reg.Get(verb)
	if !ok {
		row.Result = "error:unknown-verb"
		_ = g.w.Write(row)
		return fmt.Errorf("unknown verb %q (see manx --list-actions)", verb)
	}

	if a.Danger && !iknow(args) {
		row.Danger = "denied"
		row.Result = "error:refused"
		_ = g.w.Write(row)
		return ErrGateRefused
	}
	if a.Danger {
		row.Danger = "i-know"
	}

	err := a.body(args)
	if err != nil {
		row.Result = "error:" + err.Error()
	} else {
		row.Result = "ok"
	}
	_ = g.w.Write(row)
	return err
}

func iknow(args []string) bool {
	for _, a := range args {
		if a == "--i-know" {
			return true
		}
	}
	return false
}
