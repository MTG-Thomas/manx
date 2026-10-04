---
title: Windows server bootloop triage (0xC000021A class)
triggers: 0xC000021A, 0xC0000001, bootloop, servicing, components hive
first-verbs: status, collect, img-in
spec: §10.2, §10.2b
---

# Windows server bootloop triage (0xC000021A class)

Origin lesson: the 2026-10 incident where a *servicing repair* put a healthy
DC into a deeper hole. The ladder below encodes what that incident cost so
the next box is cheaper.

## Rules before steps

1. **Collect BEFORE you change anything.** `manx collect` + disk image if the
   target disk allows. If you cannot rename a file without imaging it, image it.
2. **One variable at a time.** Say what you are about to do out loud (audit
   row/comment) BEFORE doing it.
3. **Repairs you cannot undo deeper than their upside are net-negative.**
   An in-place servicing repair can break more boot components than it fixes.

## The ladder (stop at the first rung that holds)

1. **Know the box** — `manx status`, `manx img-in`. Identify the offline OS
   volume; mount it **read-only** first (`mount-ro` lane).
2. **Know the failure** — collect the boot evidence into the audit trail:
   - `%SystemRoot%\Logs\CBS\CBS.log` tail (latest servicing session),
   - `SetupAPI.dev.log`, panther logs (`C:\$WINDOWS.~BT\Sources\Panther\`),
   - event logs (`winevt/Logs`, esp. `System`, `Microsoft-Windows-* servicing`),
   - last uptime context: was the box healthy BEFORE a patch/reboot pair?
3. **Version-skew theorem (look here FIRST when a repair preceded death).**
   If store/manifests reference component versions newer than the installed
   binaries (e.g. `10.0.1.x` bundle vs installed `…9503`-class), the repair
   stack itself is inconsistent — an in-place repair will chase its own tail.
   Check: `dism /image:<mnt> /get-packages | findstr <pending/feature>` +
   matching `Component *}Version` in the COMPONENTS hive.
4. **COMPONENTS hive cap theorem.** A ≥1.9–2.0 GiB `COMPONENTS` hive is a
   corruption-adjacent condition: servicing writes failing/wedged long before
   the label error shows. Record hive size during collect; treat it as a
   separate finding, not background noise.
5. **Conservative rungs** (in order; re-verify boot after each):
   - `bcdedit` sanity on the offline hive (bootmgr path/device, safeboot flags off),
   - Remove-not-rollback of the last servicing package ONLY when its
     `installed` state is consistent (see #3; never blind `dism /rollback-reply`),
   - Boot recovery env (§10.2 Windows guest lane) for `sfc /scannow` +
     `dism /restorehealth` against a KNOWN-good source (mount the store from
     identical-version media; a version-skewed source poisons the victim).
6. **Stop condition.** Two consecutive rungs failing for the same root reason
   = stop medical attempts, pivot to restore (see COVE-RESTORE-LANE).
   Record the pivot decision in the audit trail with the reason.

## Never

- Blind in-place upgrade as the first repair on a domain controller.
- Registry edits on a live suspect hive (snapshot first; edit a COPY).
- Two repair mechanisms interleaved (GUI updater + DISM + manual hive edits
  at once — that is three suspects wearing one hat).
