# Secure Boot & the signed-distro model

**Short version.** We do not sign our own bootloader chain. `manx-iso` boot layer is
inherited, Microsoft-signature-verified, and shipped by the upstream distro. What
*this repo* signs is **release integrity** (checksums + a cosign signature over
`SHA256SUMS`) so that a downloaded ISO is provably the same bytes CI produced,
*and* CI proves secrets never walked into an artifact. That is the whole of the
"Secure Boot for manx" story, explained properly below.

## 1. Two different "signatures" — do not conflate them

| Layer | What it answers | Who signs it | Mechanism |
|---|---|---|---|
| **UEFI Secure Boot chain** | "Will the *machine* boot this image at all?" | Microsoft (via shim) + the upstream distro (shim→grub→kernel) | `shim-signed` + `grub-efi-amd64-signed` + distro-signed `vmlinuz` from SystemRescue's base (Debian) |
| **Release integrity** | "Are these bytes what CI produced?" | us (`MTG-Thomas/manx` CI) | `SHA256SUMS` + cosign/gpg signature, reproducible Docker build |

Secure Boot only ever cares about #1. Release integrity is a supply-chain matter;
nobody's firmware checks it.

## 2. Why the base-distro chain is the design (not a compromise)

- The stock Debian kernel that SystemRescue boots is **already** SHA/signed to the
  Microsoft 3rd-party CA via `shim-signed` and `grub-efi-amd64-signed`. Every
  OEM `db` in the wild already carries Microsoft's UEFI CA.
- Everything that makes MANX useful at rescue-time is **userspace** — `chntpw`,
  `hivexsh`, `libguestfs`, `dislocker`, `ddrescue`, our overlay scripts, and our
  Go binaries run in RAM after the kernel is up. None of them are Secure-Boot
  constrained. Kernel modules loaded by the stock kernel are signed by the
  distro; we don't need to add our own.
- Consequence: **a target with factory Secure Boot boots `manx-iso` with zero
  MOK enrollment**, because it is already configured to trust the Microsoft CA
  and the distro's key for the stock kernel.

## 3. The only future case that raises a new requirement: out-of-tree kernel modules

If a later release introduces a custom kernel module (there is *no* current
requirement), Secure Boot must verify it. The standing plan, in order:

1. **Attempt to avoid it forever.** Userspace accounts for 99% of rescue needs
   (hive, imaging, guestfish, ddrescue). This is not a "we'll get to it" note —
   it is the actual policy: a module-use feature must demonstrate in the
   contributing doc that it *cannot* be done from userspace before a module is
   accepted.
2. If/when required: ship a `manx-mok.pub` key + a **one-time MokManager**
   enrollment screen (boot once, enroll once, per box). CI produces the
   `*.ko.signed` via `sbsign` and the chain-of-custody extension to this
   document. Key ceremony lives in GH Actions secrets (see the policy in
   `SECURITY.md`).
3. For fleet-scale deployment where we control the fleet's BIOS DB (via iDRAC/
   vendor profiles), pre-provision the custom KEK in the target's `db` and skip
   MokManager — a documented option, not the default.

## 4. What the release builds (this is the "signed distro" deliverable)

Per tag (`v0.x.y`) the release workflow does the following; every step in a
reproducible Docker builder with `SOURCE_DATE_EPOCH` fixed from the tag:

1. Build `bin/manx`, `bin/manx-tui`; assemble `manx-iso` (SystemRescue remaster
   with the `toolkit-overlay` + compiled binaries).
2. Generate `SHA256SUMS` over every artifact (`ISO`, `bin/*`, overlay zip).
3. Sign `SHA256SUMS` with **cosign** (keyless via GitHub OIDC, or a long-term
   key when public verification matters more than sigstore availability: see
   `SECURITY.md` §2 for details). Publish `SHA256SUMS.sig` next to it.
4. Publish the **boot chain files** (`.efi` binaries) as release assets, so a
   user can `sbverify --list` them *without booting*: the point of this section
   is that someone can prove, offline, that the boot chain they'd run is
   Microsoft-linked. This is the "chain-of-custody" the org can audit.
5. CI runs the secrets scan (see `SECURITY.md`), then the release job publishes
   `SHA256SUMS`, `SHA256SUMS.sig`, and the per-artifact hashes.

## 5. Verification UX (how a consumer or an on-box incident operator proves it)

```sh
# 1. verify the ISO bytes are what CI signed
sha256sum -c SHA256SUMS
cosign verify --key manx-release.pub manx-iso-<ver>.iso   # or: gpg --verify SHA256SUMS.sig

# 2. prove the boot chain inside the ISO is Microsoft-linked (before deploying it)
sbverify --list /mnt/iso/EFI/BOOT/BOOTX64.EFI          # shim: Microsoft signer
sbverify --list /mnt/iso/EFI/BOOT/grubx64.efi          # grub: distro signer
sbverify --list /mnt/iso/boot/vmlinuz                  # kernel: distro signer

# 3. on a live pre-installed target, confirm the platform state
mokutil --sb-state       # Secure Boot enabled? MOK size?
```

## 6. Verification what changes when the *boot* chain changes

- Update this file in the same PR as the change (no silent chain mutation).
- If a new kernel is adopted from upstream, the distro's new signer appears —
  that is *normal*, because the trust front is Microsoft CA + vendor key. What
  is *not* normal, and requires a spec PR + a CI chain-of-custody update, is:
  - introducing our own signed-from-scratch kernel
  - enrolling a private MOK into release artifacts
  - shipping an unsigned shim
- Anything in those three categories must be paired with the §12 maturity
  checklist updated and an advisory released alongside the fix.

## 7. Why this model was chosen (history)

The 2026-10-02 incident's tooling failures became the spec's requirements. One of
them - "build a fully self-signed boot chain" - would have required:
- a private-key ceremony that this team does not staff,
- per-box MOK enrollment on every Dell (iDRAC8/SecureBoot-profile dependent),
- and no infrastructural immunity from bad actors who *can* compromise the CI
  runner; our signing key would have been the same attack surface as the
  cosign public key but without publicly-auditable keyless transparency logs.

The base-distro chain gives us a Microsoft-grade trust front with no new
security surface, no new key ceremony, and no per-box enrollment work, in
exchange for depending on an upstream we already trust (SystemRescue/Debian).
That trade is the right one for a rescue medium.
