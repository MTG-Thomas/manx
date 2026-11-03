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

## Secret gate before publication

Before publishing the current PR branch, commit the complete candidate and run
`python3 -I scripts/secret_gate.py --base origin/main` (native Windows: `python -I`). Require exit 0;
missing tooling/authentication, scan errors and skipped scans block publication.
Then run the existing repository gates. Repeat after changing outgoing commits.
Push exactly the scanned HEAD. This PR scan covers the current branch against its
fetched intended base; other refs, multi-ref pushes and rewritten history require
`--pre-push` with Git's complete ref-update stream. Stop if any pushed ref is
uncovered. See [secret scanning](SECRET-SCANNING.md) for base selection, new branches,
pre-push input and prerequisites. Preserve GitGuardian App, CI and Gitleaks checks.

Local tooling: `make pre-push`.

For a non-default PR target, fetch its destination and pass that exact ref to
`--base`; the remote default is not a substitute for a release/maintenance base.

The secret gate requires Python 3.11+ and approved ggshield >=1.56.1,<2.
