// MANX CLI: `manx <action> [args...]` is the action surface.
// The spec's parity contract requires every view (menu, CLI, TUI, harness) to
// express the same verbs; this file is the CLI view, driving verbs through
// internal/actions GateRunner so the audit rows and the --i-know gate are
// enforced in exactly one place.
//
// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MTG-Thomas/manx/internal/actions"
)

const (
	auditPath = "/toolkit/out/audit.log" // spec: append-only, shared by every view
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}

	verb, rest := args[0], args[1:]
	if verb == "help" {
		usage()
		os.Exit(2)
	}

	reg := actions.Default()
	if _, ok := reg.Get(verb); !ok {
		fmt.Fprintf(os.Stderr, "unknown action %q\n", verb)
		usage()
		os.Exit(2)
	}

	// Single audit path + single gate: the GateRunner refuses destructive verbs
	// without --i-know and writes exactly one audit row for every invocation.
	gr := reg.NewAuditor(auditPath)
	if err := gr.Run(verb, rest); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", verb, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: manx <action> [args...]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "actions (spec §10.8 parity contract):")
	for _, a := range actions.Default().List() {
		marker := ""
		if a.Danger {
			marker = " [requires --i-know]"
		}
		fmt.Fprintf(os.Stderr, "  %-20s %s%s\n", a.Verb, a.Summary, marker)
	}
}
