// Package actions is the single action surface for MANX.
// Every capability is implemented exactly once here, and is surfaced by the CLI
// (cmd/manx), the whiptail menu (manx-iso overlay), and the TUI (cmd/manx-tui).
// The verbs and their order are the spec's parity contract: additions require a
// spec PR first, then an implementation PR.
//
// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

import "fmt"

// Action is one spec-contract verb. Destructive verbs must require an explicit
// `--i-know` in every view, and must write an audit.Row with Danger set.
type Action struct {
	Verb     string   // e.g. "status"
	Args     []string // display of accepted args/flags
	Danger   bool     // destructive? (requires --i-know)
	Runs     func(args []string) error
	Summary  string   // one-line help shown in menu/TUI/CLI help
}

// Registry is the ordered set of spec-contract actions. The spec is authoritative
// for ordering and names (see docs/rescue/MANX-SPEC.md, function parity contract).
type Registry struct {
	actions []Action
}

func (r *Registry) Add(a Action) error {
	for _, x := range r.actions {
		if x.Verb == a.Verb {
			return fmt.Errorf("duplicate verb %q", a.Verb)
		}
	}
	r.actions = append(r.actions, a)
	return nil
}

func (r *Registry) Get(verb string) (Action, bool) {
	for _, a := range r.actions {
		if a.Verb == verb {
			return a, true
		}
	}
	return Action{}, false
}

// List returns the ordered verbs (spec-parity contract surface).
func (r *Registry) List() []Action { return r.actions }

// Default is the spec-defined set. Wrappers (menu, TUI, harness) must call into
// these, not re-implement.
func Default() *Registry {
	r := &Registry{}

	// Non-destructive / inspection verbs (spec contract, non-dangerous by default).
	must(r, Action{Verb: "status", Summary: "one-line runtime summary (net/mesh/drivers/scrash)",
		Runs: func(args []string) error { fmt.Println("status: ok (skeleton)") ; return nil }})
	must(r, Action{Verb: "detect-hw", Summary: "inventory unclaimed/missing-driver hardware",
		Runs: func(args []string) error { fmt.Println("detect-hw: skeleton")     ; return nil }})
	must(r, Action{Verb: "collect", Summary: "harvest logs/evtx/minidumps/hive listings (read-only)",
		Runs: func(args []string) error { fmt.Println("collect: skeleton")       ; return nil }})
	must(r, Action{Verb: "img-in", Summary: "read-only disk/volume inventory (lsblk/blkid/SMART)",
		Runs: func(args []string) error { fmt.Println("img-in: skeleton")        ; return nil }})

	// Destructive verbs: require --i-know in every view + audit row with Danger.
	must(r, Action{Verb: "hive-edit", Danger: true, Summary: "SAFE hive edit template (snapshot first; edit copy; stage back)",
		Runs: func(args []string) error { return notImpl("hive-edit") }})
	must(r, Action{Verb: "img-out", Danger: true, Summary: "image a volume/partition to a file (ntfsclone/partclone/ddrescue)",
		Runs: func(args []string) error { return notImpl("img-out") }})
	must(r, Action{Verb: "bootstrap-drivers", Danger: true, Summary: "runtime driver fetch ladder (L1 local/L2 distro/L3 vendor/L4 report)",
		Runs: func(args []string) error { return notImpl("bootstrap-drivers") }})
	must(r, Action{Verb: "bringup-windows-vm", Danger: true, Summary: "qemu/OVMF bringup of a Windows guest against an attached disk",
		Runs: func(args []string) error { return notImpl("bringup-windows-vm") }})

	return r
}

func notImpl(verb string) error { return fmt.Errorf("%s: not implemented (spec §10.8 parity contract)", verb) }

func must(r *Registry, a Action) {
	if err := r.Add(a); err != nil {
		panic(err) // compile/registration safety only; never reached in normal runs
	}
}
