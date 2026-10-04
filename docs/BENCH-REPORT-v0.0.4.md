# MANX v0.0.4 — bench acceptance report (pve-t340 VM 120, UEFI/OVMF)

Date: 2026-10-04. Operator-driven per Rule Zero (MTG-owned bench hardware only).
Under test: main `c6362d1` (TUI polish pass, PR #5 merged).

## Artifacts under test

| ISO | sha256 |
|---|---|
| `manx-iso-v0.0.4.iso` (base release) | `921aebb01a3b920c0b4050902e1151b3546439b8cad1cfcead5fbc0eaada4c9a` |
| `manx-bench-iso-v0.0.4.iso` (bench variant) | `6e28cd4ee8971807e420a55d475f164e8e704fae5679bdba14e25b5599565221` |

Binaries rebuilt from merged main `c6362d1` (linux/amd64, trimpath).

## What v0.0.4 adds (TUI polish per operator wishlist)

1. **ASCII-cat header** + `MANX · one tool, three views, zero panic` title.
2. **Live status strip** on a 2 s ticker: UTC clock · host · kernel · uptime ·
   toolkit state · **net (non-loopback preferred)** · audit-row count.
3. **Rotating truths**: the boot motd's incident-response + creature-care
   lines, one quoted line at a time (14 s rotation; `t` bumps manually).
4. **Audit-tail pane**: last 2 audit rows (`HH:MM · view · verb · result`),
   refusals visible (`error:refused (denied)`); cross-session rows flow
   through — prior-session `cli` rows appear next to fresh `tui` rows.
5. Output viewport trimmed (5 lines, `...` head when longer), tty-safe `>`
   cursor glyph, **alt-screen rendering** so the ticker never spams scrollback.
6. Shared gatherer refactor: `actions.StatusSnapshot()` extracted so the CLI
   verb, TUI strip, and future harness surfaces read one implementation
   (view-parity remains: strip/audit refresh are renders — no audit rows).
7. `status` net line prefers a non-loopback IPv4 interface (fixes the "lo"
   display seen in v0.0.3 bench).

## Verdict summary

| # | Check | Verdict | Evidence |
|---|-------|---------|----------|
| 1 | UEFI boot (OVMF/q35), zero keystrokes | PASS | rescue shell untouched |
| 2 | First-boot auto-setup | PASS | motd + marker `setup ok at 2026-10-04T02:36:54Z (bash autorun)` |
| 3 | Bench sshd auto-up (bench variant) | PASS | laptop SSH-in with no console driving |
| 4 | motd content | PASS | cat + truths on console; `/etc/manx-motd` persisted |
| 5 | TUI render (alt-screen) | PASS | `v004-tui-head.png`: cat, live strip, truth, verbs, audit tail |
| 6 | Live strip correctness | PASS | `ens18 172.16.15.171/24` shows the real NIC (was loopback in v0.0.3) |
| 7 | Truth rotation | PASS | three different truths across head/status/refusal shots |
| 8 | Enter-runs verb in TUI (status) | PASS | output captured in viewport; audit row `tui status ok` lands |
| 9 | Destructive refusal in TUI | PASS | `error:refused (denied)` row lands; refusal text + CLI hint in view |
| 10 | Audit cross-view parity | PASS | `cli` + `tui` rows identical shape (`v004-tui-status.png` pane) |
| 11 | CLI status (net preference) | PASS | real host/kernel/toolkit block (ens18 line) |
| 12 | go build/vet/test (CI) | PASS | all five lanes green on PR #5 |

## Console evidence

- `v004-boot.png` — fresh-boot zero-key motd
- `v004-tui-head.png` / `v004-tui-status.png` / `v004-tui-refusal.png` — TUI states on the live guest
- `drive_v004_tui.py` — the console-driving script (sendkey + screendump lane)

## prerelease status

Still prerelease by policy. v1.0.0 remains gated on a formal spec §8b.6 /
§10.7 / §10.8 acceptance review. The Cove restore lanes (spec §10.2b) are
designed + spec'd; their bench validation awaits a licensed sacrificial test
device (tracked internally, not as a public issue).
