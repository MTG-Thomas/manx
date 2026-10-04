---
title: Agent-first operations on a live rescue box
triggers: agent, harness, first steps, triage
first-verbs: status, detect-hw, collect, img-in
spec: §10.3, §10.8
---

# Agent-first operations on a live rescue box

The operating card for any agent (LLM-driven or scripted) dropped onto a
rescue environment: what to run, in what order, and what you owe the audit
trail.

## The first five commands (always, in this order)

1. `manx status` — host, kernel, net, toolkit state, audit row count.
   Know the box before touching it.
2. `manx detect-hw` — unknown hardware inventory (UNCLAIMED report). This is
   the driver map for `bootstrap-drivers` later.
3. `manx collect` — gather the evidence bundle (logs/fstab/mounts/hive
   listings) into the output dir. Evidence BEFORE experiments.
4. `manx img-in` — disk inventory. The imaging map (sizes, filesystems,
   mount layout) that the snapshot decision hangs on.
5. Report what you found (status line + the UNCLAIMED list + row count) to
   the human before proposing any verb with a `--i-know` gate.

## Non-negotiables

- **Destructive verbs are CLI-only from your shell, never from a GUI**: run
  `manx <verb> --i-know` and only after you could PRINT the surface summary
  of what it will touch. If the summary is not printable, do not run it.
- **One audit row per action, every view identical.** If you run verbs
  through menus, the TUI, or a script, the rows must look the same. They do;
  verify with the audit tail you just watched.
- **Document before mutating.** The audit row records the verb AFTER it ran;
  your written intent documents what you BELIEVED it would do. Both go to
  the same human.
- **Mount read-only first.** `mount-ro` lane before anything R/W touches a
  volume you are not jotting down.
- **Secrets arrive at runtime** (env/paste/secret store — §10.3). Never:
  baked into an ISO, in git, or in audit logs.

- **Care clause (agent and human both).** The motd truths are not decoration:
  water, stretch, one change at a time. An agent that finds itself churning
  (two consecutive rungs failing for the same root reason) reports the churn
  and waits for a human instead of escalating on its own.
