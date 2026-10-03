#!/bin/bash
# MANX remaster builder: (base ISO + overlay + manx binaries) -> out.iso
# Contract (docs/rescue/MANX-SPEC.md §8b.2, §11.8):
#   - SystemRescue base (Debian, MS-signed shim/grub/kernel chain via its base)
#   - overlay unpacked to /toolkit on the live system, plus manx binaries at
#     /toolkit/bin
#   - no secrets, no baked drivers (self-bootstrapping via bootstrap-drivers.sh)
#   - boot menu: Linux primary entry; optional WinPE entry if a boot.wim is given
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail

usage() { echo "usage: $0 <base-iso> <overlay-dir> <bin-dir> <out-dir> <version>"; exit 2; }
[ $# -eq 5 ] || usage
BASE="$1"; OVL="$2"; BIND="$3"; OUTD="$4"; VER="$5"

need() { command -v "$1" >/dev/null 2>&1 || { echo "missing tool: $1" >&2; exit 1; }; }
need xorriso; need unsquashfs; need mksquashfs; need mount; need umount

WORK=$(mktemp -d -p "$OUTD" work.XXXX)
ISO_MNT="$WORK/iso"; LVMNT="$WORK/lv"; OVMNT="$WORK/ov"
mkdir -p "$ISO_MNT" "$LVMNT" "$OVMNT" "$OUTD"

trap 'umount "$ISO_MNT" 2>/dev/null || true; umount "$LVMNT" 2>/dev/null || true' EXIT

echo "[remaster] mount base iso"
mount -o loop,ro "$BASE" "$ISO_MNT"

# 1. copy the squashfs root
echo "[remaster] copy live root"
cp -a "$ISO_MNT/." "$LVMNT"/ || cp -a "$ISO_MNT"/* "$LVMNT"/
SQUASH=$(find "$LVMNT" -name '*.squashfs' -o -name 'rootfs.sfs' | head -1)
[ -n "$SQUASH" ] || { echo "no squashfs rootfs found in base"; exit 1; }

mkdir -p "$OVMNT"
echo "[remaster] unsquashfs root for overlay merge"
unsquashfs -d "$OVMNT/root" "$SQUASH" >/dev/null 2>&1 \
  || unsquashfs -d "$OVMNT/root" "$SQUASH" | tail -3

echo "[remaster] apply overlay: /toolkit (scripts, manx binaries)"
mkdir -p "$OVMNT/root/toolkit/bin" "$OVMNT/root/toolkit/out"
cp -a "$OVL"/. "$OVMNT/root/toolkit/"
install -m 0755 "$BIND/manx"    "$OVMNT/root/toolkit/bin/manx"
install -m 0755 "$BIND/manx-tui" "$OVMNT/root/toolkit/bin/manx-tui"
chown -R root:root "$OVMNT/root/toolkit"
chmod -R go-w "$OVMNT/root/toolkit"

# startnet: PATH + invoker
cat > "$OVMNT/root/toolkit/setup-env.sh" <<'EOF'
#!/bin/sh
export PATH=/toolkit/bin:$PATH
export OUT=/toolkit/out
alias rescue-menu='/toolkit/bin/rescue-menu.sh'
alias detect-hw='/toolkit/bin/detect-hw.sh'
alias bootstrap-drivers='/toolkit/bin/bootstrap-drivers.sh'
EOF
echo "[remaster] overlay applied"

# 2. re-pack squashfs into the same spot
echo "[remaster] mksquashfs the new root"
gmtime="$(stat -c %Y "$ISO_MNT")"
SOURCE_DATE_EPOCH="$gmtime" mksquashfs "$OVMNT/root" "$SQUASH" -noappend -comp zstd -no-recovery -no-progress \
  || SOURCE_DATE_EPOCH="$gmtime" mksquashfs "$OVMNT/root" "$SQUASH" -noappend -comp gzip -no-progress \
  || { echo "mksquashfs failed"; exit 1; }

# 3. record the CODE inside the ISO for chain-of-custody
echo "[remaster] write SOURCE manifest"
{
  echo "manx version: $VER"
  echo "base iso: $(basename "$BASE") sha256: $(sha256sum "$BASE" | awk '{print $1}')"
  echo "built at epoch: $SOURCE_DATE_EPOCH"
  echo "host tools: $(xorriso --version 2>&1 | head -1)"
  echo "# file inventory:"
  find "$OVMNT/root/toolkit" -type f -exec sha256sum {} + 2>/dev/null | sort
} > "$LVMNT/toolkit/MANIFEST.txt"

# 4. rebuild the iso (isohybrid dual-boot; SystemRescue already UEFI+BIOS capable)
echo "[remaster] xorriso rebuild"
VOLID="manx-$VER"
xorriso -as mkisofs \
  -iso-level 3 \
  -o "$OUTD/out.iso" \
  -V "$VOLID" \
  -isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin \
  -c isolinux/boot.cat \
  -b isolinux/isolinux.bin \
  -no-emul-boot -boot-load-size 4 -boot-info-table \
  -eltorito-alt-boot --efi-boot "efi/boot/bootx64.efi" -no-emul-boot \
  -joliet -r -graft-points \
  "$ISO_MNT/" 2>&1 | tail -5

echo "[remaster] done: $OUTD/out.iso"
