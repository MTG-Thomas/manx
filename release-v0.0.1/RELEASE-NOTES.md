## MANX v0.0.1 — bootstrap release

First artifact release: the offline rescue ISO built from the MANX spec
(`docs/rescue/MANX-SPEC.md`, `docs/BUILD-PIPELINE.md`).

### What's on the ISO

- Toolkit at `manx/toolkit/`, set up to `/toolkit` at boot by the autorun banner:
  - `bin/manx` — CLI action surface (spec §10.8 verb registry, audit rows, `--i-know`
    gates)
  - `bin/manx-tui` — Tier-2 Bubble Tea view (same verbs; tty-capability gate with
    fallback to Tier 1)
  - `manx-iso-overlay/rescue-menu.sh` — Tier-1 whiptail menu
  - `manx-iso-overlay/detect-hw.sh` + `bootstrap-drivers.sh` — unclaimed/missing-driver
    detection + runtime fetch ladder (L1 local / L2 distro / L3 vendor / L4 report)
- Per-file SHA256 inventory baked in at `manx/MANIFEST.txt` (chain-of-custody).
- Boot chain: BIOS (`isolinux.bin`) + UEFI (`EFI/archiso/efiboot.img`) — verified by
  xorriso's report on the produced image (see release-verification block below).

### Verify

```sh
sha256sum -c SHA256SUMS
# then either boot (UEFI or legacy; Secure Boot supported via the base chain,
# see docs/SECURE-BOOT.md) or inspect the ISO with `xorriso -indev <file> -report_el_torito plain`
```

### Build provenance (roughly)

- Produced by the **official SystemRescue `sysrescue-customize --auto`** lane with a
  recipe directory (see `manx-iso/build/`) on a dedicated unprivileged Linux build
  container; base image = **SystemRescue 13.01**, hash-pinned
  (`56289b690bc87c85d2b9eb35790319b2d42cbdafbeae476b601dc0576b040b65`).
- Not yet reproducible (squashfs block-order varies); see
  `docs/BUILD-PIPELINE.md` "Reproducibility: attempted, honestly caveated".

### Known constraints

- Only `manx status` / `detect-hw` / `collect` / `img-in` are real verbs; the
  destructive verbs (`img-out`, `hive-edit`, `bootstrap-drivers`, `bringup-windows-vm`)
  are spec-registered skeletons (they refuse with `--i-know` and write audit rows).
- Bench acceptance (§8b.6 / §10.7 / §10.8) on MTG-owned hardware is **still pending**
  for this artifact (booting the ISO via the pve bench VM has not been run yet as a
  test case); treat the release as **pre-release / bench-unverified** until then.
- `docs/SECURE-BOOT.md` describes the trust chain; release-manifest signing is planned
  for the v0.0.2+ CI lane (`cosign` keyless via GitHub OIDC, per `docs/BUILD-PIPELINE.md`).
- Builds come from a dedicated LXC on MTG infrastructure (not customer hardware, per
  the spec's BENCH RULE `hard`).

Full changelog and repo: see the repo `CHANGELOG.md` + the
[`RESCUE-TOOLKIT-SPEC`](https://github.com/MTG-Thomas/bifrost-workspace/blob/main/docs/rescue/MANX-SPEC.md)
(the design authority).
