# Vendor restore kit manifest (Cove Data Protection / N-able)

Fetch-ladder documentation for the vendor restore binaries MANX uses per
[`MANX-SPEC.md` §10.2b](https://github.com/MTG-Thomas/bifrost-workspace/blob/main/docs/rescue/MANX-SPEC.md)
(backup-restore sidecar; never read the `.cfs` format ourselves).

**Rule Zero:** all retrieval and validation happens on MTG-owned bench/build
hardware (pve-t340 CT lanes) or an operator workstation. **Nothing is staged
or executed on practice/customer equipment**, and nothing here is committed to
git — only hashes and links.

## Source of truth

- Vendor downloads page (public, no login):
  <https://www.n-able.com/products/cove-data-protection/downloads>
- Public CDN: `https://cdn.cloudbackup.management/maxdownloads/`
- Some artifacts (the update-only installer, per-device agent installers, and
  the Linux Backup Manager installer) are served from the **partner Recovery
  Console / Management Console downloads pages** (login required). Where a
  link is login-derived we record provenance and version instead of a URL.

## Artifacts (recorded 2026-10-03)

| Artifact | Ver | Bytes | sha256 | Link / provenance |
|---|---|---|---:|---|
| `mxb-bmr-windows.exe` | 26.7.0.26243 | 637,721,608 | `173cddded2d0e6b0d27b7df16c8e602e86795171b7b4e9d189561b932a36d08f` | <https://cdn.cloudbackup.management/maxdownloads/mxb-bmr-windows.exe> |
| `mxb-bmr-windows.iso` | 26.7.0.26243 | 678,082,560 | `1c6cb84ba1a5587d08b541f63f750cdc4c77d59ef310d0bcd1325a022bceb17a` | <https://cdn.cloudbackup.management/maxdownloads/mxb-bmr-windows.iso> |
| `cove#update#binariesonly.exe` | 26.7.0.26244 | 165,999,744 | `4fbc135bc62f4ecd0ecddea1bcb88e518b1e76a98d100b9e0d872077ae33ac0d` | Partner Recovery Console downloads page ("Update (software only)") — login-derived; no stable public URL found (CDN HEAD probes 403, 2026-10-03) |
| `mxb-vd-windows-x64.exe` | 26.7.0.26187 | 1,258,536 | `844e837fd4992ca8484c4ac366798b5a03d2fc4af600a8dba3f16c2edfc3429d` | <https://cdn.cloudbackup.management/maxdownloads/mxb-vd-windows-x64.exe> |

Authenticode: builder EXE verified `CN=N-ABLE TECHNOLOGIES LTD, O=N-ABLE
TECHNOLOGIES LTD, L=Dundee, C=GB` (valid). Verify before any use.

## Roles (see spec §10.2b for the full design)

- **`mxb-bmr-windows.iso`** — the Windows BMR boot media (WinPE). Inspection
  (2026-10-03): contains the Backup Manager engine, `ClientTool.exe`
  (scriptable CLI) and baked VirtIO drivers (NetKVM, viostor, vioscsi,
  vioserial) → runs as a **KVM guest** with a virtio scratch disk for the
  volume-lane restore. Never booted on target metal; guests see only scratch
  disks.
- **`mxb-bmr-windows.exe`** — self-contained builder for that media (ADK-free)
  when a rebuild or driver injection is needed.
- **`cove#update#binariesonly.exe`** — binaries-only update installer
  (v26.7.0.26244 — one build ahead of the recorded media). **Run the update
  against the media before any real restore** so the engine matches current
  fleet agent versions (freshness-first rule).
- **`mxb-vd-windows-x64.exe`** — Virtual Drive module (mount backup sessions
  as a browsable drive on a Windows helper; Windows-only convenience).
- The **Linux Backup Manager installer** (data lane: restore-only mode +
  `/opt/MXB/bin/ClientTool control.restore.start`) is issued from the console
  downloads page (per-device installation token flow). Record its hash at
  first fetch; same Rule Zero rules.

## Bench prerequisites (internal tracking)

A **licensed sacrificial test device** does not exist yet (2026-10-03); until
one is allocated (never practice equipment), the guest-lane and data-lane
restore validations stay on the dry/structural level above. Tracked in MTG
operational memory — deliberately **not** a public repo issue.
