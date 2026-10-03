# Repository guidance for AI agents working on MANX.

This repo is deliberately written to be both human- and agent-navigable.

## What this repo is

`github.com/MTG-Thomas/manx` implements the **MANX** rescue toolkit.
The authoritative design lives in `MTG-Thomas/bifrost-workspace`:
`docs/rescue/MANX-SPEC.md`. Read that spec first for any design question; do not
duplicate its content here.

## Non-negotiables for agent work

- AGPL-3.0-or-later — `SPDX-License-Identifier: AGPL-3.0-or-later` header on new
  source files.
- ONE action surface, THREE views. Every new capability added here must be:
  - a single verb under `internal/actions/` (one implementation),
  - surfaced via `cmd/manx` (CLI),
  - exposed in the `whiptail` menu (`manx-iso/toolkit-overlay/bin/rescue-menu.sh`),
  - exposed in `manx-tui` (Bubble Tea view of the same surface),
  - ordered by the audit contract: every destructive verb writes an audit row to
    `internal/audit`, and every row is identical no matter which view invoked it.
- No `panic` outside `main()`. Destructive verbs require an explicit `--i-know`
  confirmation in every view.
- No secrets in artifacts. Never bake meshes keys/credentials into the ISO,
  manx-tui strings, or the audit log. CI runs a secret-shape scan before release.
- Correctness over cleverness: `make check` (gofmt/vet/test) plus
  `-race` in CI (see codex-swarm's CI as the maturity reference).

## Conventions

- `Makefile` targets follow `MTG-Thomas/codex-swarm`: `all/check/build/test/vet/fmt/fmt-check/vulncheck/check/clean`.
- `CHANGELOG.md` is Keep-a-Changelog (Always update on PRs that change behavior).
- Prefer `testify` over hand-rolled asserts; match codex-swarm's CI shape.
- Cross-repo references resolve to `bifrost-workspace` docs via absolute URLs.

## Where the "what" lives

- SPEC: authority (`bifrost-workspace/docs/rescue/MANX-SPEC.md`)
- CLI: `cmd/manx/main.go`, operations in `internal/actions/*.go`
- TUI: `cmd/manx-tui/main.go` (Bubble Tea), view-only, same verbs
- ISO overlay scripts: `manx-iso/toolkit-overlay/bin/*` (whiptail + shell)
- Audit/state/self: `internal/{audit,state}/*.go`
