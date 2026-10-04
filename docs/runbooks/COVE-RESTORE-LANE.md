---
title: Cove restore lane without the vendor BMR ISO
triggers: restore, cove, bare-metal, backup, volume restore
first-verbs: status, img-in, detect-hw
spec: §10.2b
---

# Cove restore lane without the vendor BMR ISO

The vendor's Windows BMR boot media is the wrong environment to think in.
MANX's rule: never parse the `.cfs` format; orchestrate the vendor engine in
an environment we own. Two lanes (spec §10.2b has the full contract).

## Lane 0 — preconditions (check before anything)

- [ ] Rule Zero: you are on MTG bench hardware or a sacrificial device.
- [ ] The vendor kit is in scope: see `docs/VENDOR-RESTORE-KIT.md` (hashes
      recorded). Run the update-only installer against the media first
      (freshness-first: media engine ≥ fleet agent version class).
- [ ] Licensed restore target exists for the restore point being pulled —
      never point the engine at production backups from a box you are not
      prepared to license/authorize.
- [ ] Credentials (device User/Password + EncryptionKey) arrive at runtime
      (operator paste/env/secret store). NEVER in the ISO, git, or logs.

## Lane 1 — data pulls (files, databases; no VM needed)

1. `manx status` + `manx net-up` (or mesh-up): the box needs reach to the
   Cove repo.
2. Fetch the Linux Backup Manager via the fetch ladder (operator-declared
   hash), install **restore-only mode** (`[General] ReadOnlyMode=1`; config
   from operator-supplied env — `User`, `Password`, `EncryptionKey`).
3. Drive the restore headlessly: `/opt/MXB/bin/ClientTool`:
   - `control.session.list -datasource FileSystem` (pick start time),
   - `control.restore.start -datasource FileSystem -restore-to
     /toolkit/out/backup-pull/ -selection <path> -time "<ts>"` (+ existing
     files policy per case).
4. Land objects under `/toolkit/out/backup-pull/`; audit rows per object.
   Network-share targets are in-place only — do not hand them a scratch file
   path and expect success.

## Lane 2 — volumes (bare metal or P2V)

1. Their WinPE BMR media boots as a **KVM guest** (it ships VirtIO drivers +
   `ClientTool` + CEF wizard): attach vendor ISO as cdrom + blank
   virtio scratch disk (≥ source disk size) + virtio NIC. UI over VNC.
2. The wizard restores the volume datasource **onto the scratch disk only**.
   The guest must never be able to see the real target disks.
3. After restore: MANX byte-level copy scratch → target (partclone/ddrescue
   — `img-restore` lane, `--i-know`, surface summary printed first), or boot
   the scratch image in a throwaway VM for validation first.
4. Post-copy boot repair is OUR lane (not the vendor's): `detect-hw` +
   `bootstrap-drivers` + §10.2 Windows-guest `bcdedit`/`dism`.

## Known limits

- Headless volume restore from the LINUX agent is unproven; volumes go
  through the VDR/guest media until validated on a sacrificial device.
- A deduplicated cloud repository cannot be read offline without their
  engine's metadata; item-level pulls need the engine too.
