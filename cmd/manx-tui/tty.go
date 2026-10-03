// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"os"
)

// isCapableTTY reports whether the TUI can render. A "dumb" terminal (serial
// console, redirected output, no TERM, small tty class, TERM=dumb) is refused
// so the wrapper falls back to the Tier-1 whiptail path (spec §10.8).
func isCapableTTY(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return false // redirected; not interactive
	}
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}
	// CRITICAL: never trust TERM alone on a rescue console. If the tty is on a
	// serial device, allow the TUI only when the operator passed --force-tui.
	return true
}
