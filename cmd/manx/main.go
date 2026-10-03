// MANX CLI: `manx <action> [args...]` is the action surface.
// The spec's parity contract requires every view (menu, CLI, TUI, harness) to
// express the same verbs; this file is the CLI view.
//
// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MTG-Thomas/manx/internal/actions"
	"github.com/MTG-Thomas/manx/internal/audit"
)

const (
	auditPath = "/toolkit/out/audit.log" // spec: append-only, path exposed to views
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
	reg := actions.Default()
	a, ok := reg.Get(verb)
	if !ok || verb == "help" {
		usage()
		os.Exit(2)
	}

	// Audit contract: one row per verb invocation, in every view.
	// (CLI view: Actor is "operator" by convention; the agent harness sets its
	// own Actor string in upstream calls.)
	view := "cli"
	defer func() {
		_ = audit.Writer{Path: auditPath}.Write(audit.Row{
			Actor: "operator", View: view, Verb: verb,
			Danger: dangerFor(a, rest), Result: "ok",
		})
	}()

	if a.Danger && !hasIKnow(rest) {
		fmt.Fprintf(os.Stderr, "%s: destructive action requires --i-know in any view (spec \u00a710.8)\n", verb)
		view = "cli"
		_ = audit.Writer{Path: auditPath}.Write(audit.Row{
			Actor: "operator", View: view, Verb: verb, Danger: "denied", Result: "error:refused",
		})
		os.Exit(1)
	}

	if err := a.Runs(rest); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", verb, err)
		os.Exit(1)
	}
}

func hasIKnow(args []string) bool {
	for _, a := range args {
		if a == "--i-know" {
			return true
		}
	}
	return false
}

func dangerFor(a actions.Action, args []string) string {
	if a.Danger && hasIKnow(args) {
		return "i-know"
	}
	return ""
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: manx <action> [args...]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "actions (spec \u00a710.8 parity contract):")
	for _, a := range actions.Default().List() {
		marker := ""
		if a.Danger {
			marker = " [requires --i-know]"
		}
		fmt.Fprintf(os.Stderr, "  %-20s %s%s\n", a.Verb, a.Summary, marker)
	}
}
