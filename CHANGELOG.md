# Changelog

All notable changes to MANX are documented here (Keep a Changelog, semver from v0).

## [Unreleased]

### Changed
- `manx-tui` polish pass: ASCII-cat header with a **live status strip**
  (UTC clock, host, kernel, uptime, toolkit state, net (non-loopback
  preferred), audit-row count) refreshed every 2 s; a **rotating truth line**
  (the boot motd's incident-response + creature-care truths, themed via the
  `t` key); a compact **audit-tail pane** (last rows, view-tagged, refusals
  shown); output viewport trimmed to 5 lines; alt-screen rendering so the
  ticker doesn't spam scrollback.
- `manx status` / TUI strip now prefers a non-loopback IPv4 interface when
  reporting the net line (loopback is shown honestly if it is all there is).
- `actions.StatusSnapshot()` extracts the same gathered values the CLI verb
  renders — one implementation, three views (CLI text unchanged apart from
  the net preference).
- `docs/VENDOR-RESTORE-KIT.md`: Cove/N-able restore-kit manifest — public links
  where they exist, provenance where they don't (update-only installer is
  console-derived), versions + SHA256s; Rule Zero stated; licensed-test-device
  gap tracked internally, not as a public issue.

## [0.0.3] - 2026-10-03

### Added
- `manx setup` verb: stages /toolkit from the boot media (idempotent; writes
  `/toolkit/.setup-complete`), with motd + PATH wiring — the in-Go twin of the
  ISO autorun script. Spec amended (PR #1157) to add `setup` + `detect-hw` to
  the function parity contract.
- First-boot auto-setup on the ISO: `sysrescue.d/10-manx-autorun.yaml` ships
  YAML-config `autorun.exec` entries (SystemRescue >= 9.05) for the setup
  script — restoring the kit at boot with zero keystrokes. Root cause of the
  old manual kick: name-based autorun only executes scripts whose names begin
  with `autorun`; `00-manx-setup.sh` was silently skipped.
- `sign-release` CI workflow: keyless (OIDC) cosign signing of `SHA256SUMS`
  on every published release, bundle uploaded as `SHA256SUMS.sigstore.json`
  (self-verified in-lane before upload).
- Bench recipe YAMLs canonicalized in-repo (`manx-iso/recipe/{base,bench}/`)
  so the CT build cycle no longer drifts from reviewable source.

### Changed
- `manx-tui`: Enter runs the selected verb through `GateRunner.RunView("tui", …)`
  with verb stdout captured into the viewport; destructive verbs show the
  §10.8 refusal plus the exact CLI command that can perform them. Fixed a verb
  mismatch (`collect-drivers` → `bootstrap-drivers`) and the garbled help line.

### Removed
- ~49 MB of committed raw `.ppm` screendumps (bench evidence kept as small
  `.png` only; `*.ppm` now gitignored).

## [0.0.2] - 2026-10-03

### Added
- Real verb implementations: `status` (live host/kernel/net/toolkit), `detect-hw`
  (lspci UNCLAIMED report to `/toolkit/out/drivers-report.txt`), `img-in`
  (lsblk/blkid inventory to `/toolkit/out/img-in-report.txt`).
- The motd: ASCII cat + incident-response truths + human truths at boot and on
  every new shell (`/etc/manx-motd`).
- GateRunner view contract across CLI/menu/TUI (audit rows tagged per view).
- `docs/BENCH-REPORT-v0.0.2.md`; hardening: shellcheck-clean overlay, CodeQL,
  SHA-pinned actions, `lint` lane (shellcheck + actionlint).

### Fixed
- Driver-leak bug in `detect-hw` parsing; POSIX shellcheck fixes across overlay.

## [0.0.1] - 2026-10-03

### Added
- bootstrap: spec, LICENSE (AGPL-3.0-or-later), SECURITY.md, CONTRIBUTING.md, AGENTS.md,
  CODEOWNERS, CI + govulncheck + dependabot, Makefile with codex-swarm-parity targets.
- CLI skeleton: `manx --list-actions` + `manx status` (one action surface, verb registry).
