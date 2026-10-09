// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

// SPEC_VERBS: the implemented subset of the MANX-SPEC.md function parity
// contract (amended 2026-10-03 in the spec PR to add `setup` + `detect-hw`).
// This list locks the registry to exactly these verbs, in this order.
var SPEC_VERBS = []string{
	"status",
	"detect-hw",
	"collect",
	"setup",
	"img-in",
	"session",
	"hive-edit",
	"img-out",
	"bootstrap-drivers",
	"bringup-windows-vm",
}

// SPEC_DESTRUCTIVE marks the destructive verbs (gated behind --i-know).
var SPEC_DESTRUCTIVE = map[string]bool{
	"img-out":            true,
	"hive-edit":          true,
	"bootstrap-drivers":  true,
	"bringup-windows-vm": true,
}
