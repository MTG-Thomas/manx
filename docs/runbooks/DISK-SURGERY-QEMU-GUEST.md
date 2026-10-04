---
title: Disk surgery via the qemu guest (never the metal)
triggers: disk surgery, qemu, hive edit, snapshot, bcdedit, dism
first-verbs: status, img-in, collect
spec: §10.2, §10.4
---

# Disk surgery via the qemu guest (never the metal)

Default: do NOT drive Windows on bare metal. Boot Linux, image the disk,
operate on a COPY in a throwaway Windows guest; put changed bytes back as a
deliberate, gated copy.

## The loop

1. **Snapshot first.** Image the disk (or dd the partition) BEFORE any
   surgery. `manx img-in` tells you what you are looking at.
2. **qcow2 overlay** onto the read-only base image for throwaway boots — the
   base never changes; overlays are disposable.
3. **Guest boots** (qemu/OVMF per §10.2): run `dism`, `sfc`, `bcdedit` INSIDE
   the guest against the attached disk, with exactly one goal per boot.
4. **Boot the overlay** to verify the surgery did what it claimed, BEFORE
   promoting anything.
5. **Promote deliberately**: copy overlay→base (or disk→metal) only as a
   separate, stated, gated step. `img-restore` = `--i-know` and a
   pre-printed surface summary. If you cannot print the summary, you are
   not ready to copy.

## Registry surgery (hive-edit lane)

- NEVER edit the suspect hive in place. Copy it, open the copy
  (hivex/chntpw class tools), and stage it back with the same
  copy-back gate as block data.
- A hive you cannot load offline is evidence, not cargo: collect it and
  stop changing things.

## Agent-managed overlays

The agent harness (§10.3) can drive this lane end-to-end: snapshot, overlay,
guest up, tool runs, boot check, promotion — each step an audit row. If any
step's row lands `error:*`, the promotion step is forbidden by construction.
