---
title: Restricted recovery enrollment
triggers: enrollment, pairing, offline synchronization
first-verbs: session
spec: https://github.com/MTG-Thomas/bifrost-workspace/pull/1239
---

# Recovery enrollment

Cloudflare is a coordinator, never permission to execute recovery actions. The
existing GateRunner confirmations remain independent of session approval.

```sh
manx session register --endpoint https://<nonproduction-worker-origin>
manx session status
manx session record --diagnostic hardware_inventory
manx session sync
```

`MANX_LIFECYCLE_ENDPOINT` selects the public HTTPS origin. No token is needed.
The menu exposes the same subcommands. Select `session` in TUI: Enter/status,
`r`/register, `s`/status, `y`/sync, `d`/diagnostic. All invoke the one action body.

The restricted 0700 journal defaults to `/toolkit/out/lifecycle`; key/checkpoint
files use 0600. This scratch directory may be volatile on a live ISO. For reboot
persistence, explicitly select an operator-owned persistent scratch volume with
`--state-dir`; never discover, mount or modify a customer's disk for journaling.
Loss of scratch means a new recovery identity and fresh technician claim, not
reuse of Windows identity based on serial. Process restart preserves the journal.

Offline `record` works before registration. Sync fills the eventual server
locator, preserves logical event IDs/sequences, and retries pending events.
Phase 1 implements `hardware_inventory` diagnostics only: it reads sysfs and
stores fresh hardware observations in a restricted `evidence-<event-id>.json`
file before journaling the checkpoint. Registration publishes hardware identity
observations; synchronization publishes the typed collection result. Full
diagnostic artifact upload and other diagnostic categories remain later work.
Conflicting acknowledgment leaves the journal intact and reports intervention.
No secrets or server response bodies appear in status, audit or release artifacts.

Disposable VM acceptance: run `tests/acceptance/lifecycle.sh register <lab-origin>
<operator-scratch>`, claim/approve the displayed code in Workspace, then run its
`offline` and `sync` phases. Capture the ISO/release hash, VM identity and
Workspace event timeline. The script simulates lost HTTPS connectivity without
modifying network interfaces or disks. Client assertions alone do not establish
Workspace acceptance. Never publish the private `state/key.der` file.

Rollback: stop session sync and retain restricted journal for investigation.
The rescue toolkit and existing local operations remain available offline. Do
not erase an accepted event merely to obtain a passing synchronization result.
