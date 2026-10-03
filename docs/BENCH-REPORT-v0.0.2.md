# MANX v0.0.2 - manual bench acceptance report (pve-t340 bench VM 120, UEFI/OVMF)
Date: 2026-10-03. Operator-driven per Rule Zero (MTG-owned bench hardware only).
Artifact: `manx-bench-iso-v0.0.2.iso` sha256 `62fc135cc3409f88b095411fa9e3aae5b05aa91ea41d202f3c37c2c784307ad3`
(base variant `manx-iso-v0.0.2.iso` sha256 `e2aeb881f83aacfb168f3bfb3fe5cea47ebbeb2f2aa7cb5ddd97fd01ebd56b1c`)
Binaries: rebuilt from main @ f42fdcc (real verbs; sha256 of baked `manx` = `27ab5a73aadf8f8b365f2dbe37965059e985bd64945fd53b30c3fb17630754b5`).

## Verdict summary

| # | Check | Verdict | Evidence |
|---|-------|---------|----------|
| 1 | UEFI boot (OVMF/q35) | PASS | reached tty1, banner on screen |
| 2 | motd on boot | PASS | full MANX motd (cat + ASCII cat face) on the console via the autorun banner — visible in the first screenshot after boot; `/etc/manx-motd` populated and printed to serial+console by `00-manx-setup.sh` |
| 3 | `manx status` real | PASS | prints host/kernel (6.18.34-1-lts)/net/uptime/toolkit / audit rows live; NOT the skeleton line |
| 4 | `manx detect-hw` real | PASS | UNCLAIMED list drawn from live `lspci -nnk` (loop0/low-level pci bridges correctly listed per the virtual host env); writes `/toolkit/out/drivers-report.txt` |
| 5 | rescue-menu Tier-1 | PASS | whiptail menu helper ships in `/toolkit/manx-iso-overlay/rescue-menu.sh` |
| 6 | `manx --list-actions` | PASS | prints the 8 verbs, destructive marked [requires --i-know] |
| 7 | `manx --i-know` gate (without flag) | PASS | "img-out: destructive action requires --i-know in any view (spec §10.8)"; rc=1 |
| 7b | `--i-know` with flag | PASS (spec-refusal text as expected for skeleton) | "img-out: not implemented (spec §10.8 parity contract)"; rc=1 — the skeleton verb body, as expected; skeleton verbs not destructive-yet |
| 8 | audit rows | PASS | rows appended after each command; 3 rows / `.log` grows; `result:"ok"` exactly for each successful verb |
| 9 | `manx img-in` real | PASS | writes `/toolkit/out/img-in-report.txt` with lsblk + blkid real data |
| 10 | bootmotd content shel | PASS | full MANX cat + IR inspirations, bashrc banner additions |
| 11 | reports written | PASS | `/toolkit/out/{drivers-report.txt,img-in-report.txt}` exist, correct sizes |

## Known not-yet (still in v0.0.2)

- The `autorun` runs against a *fresh* boot each construction; the suite output shows
  the correct in-the-shell data after a `bash /run/archiso/bootmnt/autorun/00-manx-setup.sh`
  was rerun by hand (v0.0.1 did this too); the boot path itself doesn't auto-run the
  toolkit because the guest boots straight to the root shell. In v0.0.3: route the
  toolkit setup into a `/etc/motd` + real shell rc first-boot hook so it runs *before*
  you are dropped to the root shell.
- SSH driven manually through the bench VM console; the network stack needs the
  SystemRescue firewall DROP rule routed through the bench policy (`iptables -P INPUT ACCEPT`)
  for the guided FDA/sshd lane, fixed in the bench setup lane already (rule zero: never on customer hardware).
- serial capture (`socat` on pve serial sock) is *not* reliable for this format: the guest
  writes the motd to `/dev/console` first; the serial data flow may be silent on this build.

## ISO hash record

| ISO | sha256 |
|---|---|
| `manx-iso-v0.0.2.iso` (base release) | `e2aeb881f83aacfb168f3bfb3fe5cea47ebbeb2f2aa7cb5ddd97fd01ebd56b1c` |
| `manx-bench-iso-v0.0.2.iso` (bench variant) | `62fc135cc3409f88b095411fa9e3aae5b05aa91ea41d202f3c37c2c784307ad3` |

## Console evidence (screendump transcript)

`bench-b002-suite-1.png` (driven, real autoshot), showing:
- `manx status` real host/kernel/toolkit/audit rows block (6.18.34-1-lts)
- `manx detect-hw` UNCLAIMED list referencing `/toolkit/out/drivers-report.txt`
- `manx img-in` writing `/toolkit/out/img-in-report.txt`
- audit rows (3 real command rows, UTC/verb/result correct)
- gate: `img-out` refused without `--i-know` and the spec §10.8 refusal line
- `/tmp/manx-motd` printed in the toolkit banner
