# MANX v0.0.1 ISO — build provenance (MTG infra, MTG-owned hardware)

How the `manx-iso-v0.0.1.iso` artifact was produced, reproducibility scope, and
what to do if the build needs to happen again. Does NOT describe any customer
hardware or customer-identifying details; **bench/build is MTG-infrastructure
only** (spec `MANX-SPEC.md` BENCH RULE hard).

## Build environment

Dedicated unprivileged LXC on the fleet Proxmox host (Debian 13 template):

```
pct create <id> /var/lib/vz/template/cache/debian-13-standard_13.1-2_amd64.tar.zst \
  --hostname manx-build --memory 6144 --swap 2048 --cores 4 \
  --rootfs local-lvm:20 --net0 name=eth0,bridge=vmbr0,ip=dhcp \
  --features nesting=1 --unprivileged 1 --tags manx --onboot 1
```

Dependencies installed in the CT:

```
apt-get install -y --no-install-recommends xorriso squashfs-tools rsync curl wget file patch
```

(`patch` is required by `sysrescue-customize`'s dependency check even when no
patch step runs; keep it in the pack list.)

## Remaster (official lane)

**`sysrescue-customize --auto`**, not a hand-rolled xorriso pipeline. Recipe
directory (this repo, `manx-iso/`):

```
recipe/iso_add/autorun/00-manx-setup.sh    # ISO-root autorun banner; sets up /toolkit
recipe/iso_add/manx/toolkit/bin/{manx,manx-tui}     # compiled Go binaries (linux/amd64, static)
recipe/iso_add/manx/toolkit/manx-iso-overlay/       # Tier-1 whiptail + helper scripts
recipe/iso_add/manx/toolkit/README.ISO.md
```

Invocation:

```
bash sysrescue-customize \
  --auto \
  --source=sysrescue-13.01-amd64.iso \
  --dest=manx-iso-v0.0.1.iso \
  --recipe-dir=<recipe-dir> \
  --work-dir=<empty dir> \
  --overwrite
```

`sysrescue-customize` itself:
`https://gitlab.com/systemrescue/systemrescue-sources/-/raw/main/airootfs/usr/share/sysrescue/bin/sysrescue-customize?inline=false`
(shipped on the ISO too; use the upstream URL so the pipeline is reproducible).

## Base ISO pin

`BASE-URL.sha256` (in `manx-iso/build/`) pins URL + sha256 of the SystemRescue
release; the build script refuses to proceed on hash mismatch. v0.0.1 pin:

```
https://fastly-cdn.system-rescue.org/releases/13.01/systemrescue-13.01-amd64.iso
56289b690bc87c85d2b9eb35790319b2d42cbdafbeae476b601dc0576b040b65
```

## Result (v0.0.1)

- `manx-iso-v0.0.1.iso` — 1,364,000,768 B, SHA256
  `03a5c3c1a4656472a3cc7483006712860c87622c42ccaae034507189d9d2f0c6`
- El Torito (re-verified by xorriso on the produced ISO):
  - image 1: BIOS, `isolinux/isolinux.bin` (boot-info-table, isohybrid-suitable)
  - image 2: UEFI, `EFI/archiso/efiboot.img`
- `manx/MANIFEST.txt` is baked into the ISO with per-file SHA256 of every staged
  file — chain-of-custody, not a reproducibility claim.

## Reproducibility scope

Not reproducible end-to-end at v0.0.1 (squashfs/xorriso ordering varies between
runs even with fixed input). The CI `iso-linux` lane, when it lands, will attempt
`SOURCE_DATE_EPOCH`-style reproducibility and record the honest boundary in
`docs/BUILD-PIPELINE.md`. Validation for a re-build = the **hash of the
delivered artifact** + the `MANIFEST.txt` inside it, not byte-for-byte re-build.

## Bench VM (pve, MTG-owned)

`pve-create-bench-vm.sh` in this directory provisions the acceptance bench VM:
Debian guest VM (UEFI/OVMF) with the ISO as the boot CD, 4 GB RAM, 10 GB scratch
disk, no customer data. The bench lane (`bench-e2e.yml`) drives it via the
`manx-bench` label only when an operator dispatches it.
