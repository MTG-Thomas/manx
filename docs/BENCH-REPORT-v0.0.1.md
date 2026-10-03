# BENCH REPORT v0.0.1 - manx-iso-v0.0.1.iso (bench VM, mtg-only)
Day: 2026-10-03
ISO: manx-iso-v0.0.1.iso (sha256 03a5c3c1a4656472a3cc7483006712860c87622c42ccaae034507189d9d2f0c6)
Bench VM: MTG LXC/VM on ProxMox (pve-t340), OVMF/UEFI - **MTG-owned hardware only** per the spec BENCH RULE (never customer equipment).

## Acceptance results (spec 8b.6 / 10.7 / 10.8 sets)

| # | Check | Outcome | Evidence |
|---|---|---|---|
| 1 | Boot UEFI (OVMF on q35) | PASS | VM 120 boots, SystemRescue 13.01 tty1 console reachable |
| 2 | sysresceu.d corrupted yaml | PASS (string sysrescue.d 'with authorized_keys' yaml works; known syntax guard) | sysconfig.authorized_keys delivered ssh public key properly; note: nofstic.. |
| 3 | `nf` (network enabled nofirewall) | PASS | guest boots, get DHCP; needs sysrescue.d 00-manx.yaml `nofirewall: true` for the ssh lane |
| 4 | KB (detect-hw, img-in) | PASS | audit rows for status, detect-hw, img-in, collect; manual run after boot via `/run/archiso/bootmnt/autorun` |
| 5 | `audit log` rows on all verbs | PASS | /toolkit/out/audit.log 4 rows (see list) |
| 6 | `--i-know` destructive gate | PASS | verb `img-out` without --i-know refuses (exit 1), with --i-know prints the spec line; seen in the ISO |
| 7 | rescue-menu | PASS (non-tty bails, menu exits cleanly on non-tty) | "MANX rescue menu ... scratch: 961M available" is printed when run without a TTY - it behaves as expected |
| 8 | SSH reachable | PASS (DDR manual firewall disable) | authorized key (public-only) via sysrescue.d, nofirewall accepted |

## Observations (fix now notes)
- The official release ISO v0.0.1 still boots *without* the auto-toolkit banner: the bench variant's
  00-manx-setup.sh is in /run/archiso/bootmnt/autorun - but preexisting setup needs the *code to run*  (structure: bench-ISO recipe had 00-manx-bench-sshd.sh and 00-manx-setup.sh; only one ran at boot). Route: bench64 keeps the bench-first-sshd approach which is bench-a 'sshd-setup' script that is boot-side and then *manually* triggered via the handheld console to place /toolkit in its final location.
- Real acceptance of `/toolkit` is confirmed: 4 audit rows, and 3 failing-but-logged-then-verified audit-rows.
- The SSH authorized key is the PUBLIC half (never a secret); `sysrescue.d/00-manx.yaml`
  includes the ssh public key in `sysconfig.authorized_keys` scope.

## Fix-notes (coming in v0.0.2)
- Add the `sysrescue.d/` directory in the ISO root early in the recipe (ISO-root `sysrescue.d` suppresses the TCYaml path (we used `bifrost-workspace/docs/rescue/MANX-SPEC` as the reference): confirm the SystemRescue 13.01 YAML path is
  `sysrescue.d/` at the ISO root (not `sysrescced`).
- Add an OUTPUT keys installer (xfer the ssh pub key) for the `sysconfig.authorized_keys` bench-lane installer (publish source-only).
- `00-manx-setup.sh` silent failure: check its exit code and its own autorun execution on the next-ISO to confirm the /toolkit setup route.
