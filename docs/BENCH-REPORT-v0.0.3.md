# MANX v0.0.3 — manual bench acceptance report (pve-t340 bench VM 120, UEFI/OVMF)

Date: 2026-10-03. Operator-driven per Rule Zero (MTG-owned bench hardware only).
PR under test: MTG-Thomas/manx#4 (head `958a58d`).

## Artifacts under test

| ISO | sha256 |
|---|---|
| `manx-iso-v0.0.3.iso` (base release) | `f3e351b09e39f40495ebdb0fdf6ccab49bbd850fb977a18ba0cf85983e160da6` |
| `manx-bench-iso-v0.0.3.iso` (bench variant) | `7cea8bb8268445c8db218bcdb9680b4156568e79594c8ace0d445e41552e6291` |

Binaries rebuilt from PR head `958a58d` (linux/amd64, CGO off, trimpath). Guest:
DHCP `172.16.15.171`, kernel `6.18.34-1-lts`.

## Verdict summary

| # | Check | Verdict | Evidence |
|---|-------|---------|----------|
| 1 | UEFI boot (OVMF/q35), zero keystrokes | PASS | reached rescue shell untouched |
| 2 | **First-boot auto-setup (no manual kick)** | PASS | motd + toolkit on screen at fresh boot with ZERO keystrokes (`v003-boot.png`, `v003-final-boot.png`); `/toolkit/.setup-complete` = `setup ok at 2026-10-03T18:15:50Z (bash autorun)` |
| 3 | **Benchmark sshd auto-up** (bench variant) | PASS | YAML `autorun.exec` entry starts sshd at boot; laptop SSH-in worked with no console driving (previous releases needed a manual script run) |
| 4 | motd content | PASS | ASCII cat + IR truths + human truths on console; persisted to `/etc/manx-motd` |
| 5 | `manx status` real | PASS | host/kernel/toolkit/net/uptime/audit-rows live |
| 6 | `manx collect` real (new verb) | PASS | `collected 4 targets into /toolkit/out/collect-2026-10-03T18:17:43Z.tar.gz` (232,914 bytes) |
| 7 | `manx setup` idempotent re-run | PASS | from `/toolkit/bin/manx`: "MANX toolkit ready at /toolkit"; marker rewritten by the Go verb (self-overwrite skip fix verified) |
| 8 | Verb discovery | PASS | 9 verbs, destructive marked `[requires --i-know]` |
| 9 | **TUI Enter-runs a real verb** | PASS | `manx-tui` → Enter on *status* → verb stdout captured in viewport (`v003-tui-status.png`) |
| 10 | **TUI destructive refusal + CLI hint** | PASS | img-out → "refused: destructive action requires --i-know in any view (spec §10.8) / run \`manx img-out --i-know\` from the rescue shell" (`v003-tui-refusal.png`) |
| 11 | **Audit parity across views (spec §10.8 #1/#3)** | PASS | audit.log rows identical shape for `view=cli` and `view=tui`, refusals recorded `danger:"denied"`, `result:"error:refused"` (`v003-tui-audit.png`) |
| 12 | Gate (CLI) refusal | PASS | `img-out` without flag → §10.8 text, rc=1 |

## Bench-caught bug (fixed in this PR)

- `manx setup` run from `/toolkit/bin/manx` aborted with `text file busy`
  (staging over the running executable). Fix: `copyAll` skips the running
  binary; re-run verified PASS (check #7). The failed attempt was honestly
  recorded in the audit log before the fix (`view=cli`, `result=error:...busy`).

## Root-cause note (first-boot)

v0.0.1/v0.0.2 needed a manual `bash 00-manx-setup.sh` kick: SystemRescue's
*name-based* autorun only executes scripts whose filenames begin with
`autorun` — `00-manx-setup.sh` was silently ignored. v0.0.3 ships YAML-config
`autorun.exec` entries (`sysrescue.d/10-manx-autorun.yaml`, documented
mechanism, SystemRescue >= 9.05) on both variants. Canonical recipe YAMLs now
live in-repo under `manx-iso/recipe/`.

## Console evidence

- `v003-boot.png` / `v003-final-boot.png` — zero-keystroke fresh boot, motd up
- `v003-tui-status.png` — TUI Enter-run with captured verb output
- `v003-tui-refusal.png` — TUI destructive refusal with CLI hint
- `v003-tui-audit.png` — cross-view audit parity tail
- `bench-v003-ssh.txt` — raw SSH acceptance transcript

## prerelease → candidate status

All three v0.0.2 Done-ness gaps are closed and bench-verified: first-boot
auto-setup (zero keystrokes), keyless-signed SHA256SUMS (`sign-release` lane,
verified on the v0.0.3 release), TUI as a real verb surface. v0.0.3 remains a
prerelease by version policy (v1.0.0 gated on a formal bench acceptance
review per spec §8b.6/§10.7/§10.8 acceptance sets).
