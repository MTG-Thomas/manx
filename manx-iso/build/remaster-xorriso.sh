#!/bin/bash
# manx v0.0.1 ISO remaster, pure xorriso (no root, no squashfs repack).
# base ISO already downloaded + hash-verified to the pinned manifest.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
B=~/manx-build
BASE="$B/sysrescue.iso"
TREE="$B/tree"
OUT="$B/manx-iso-v0.0.1.iso"
OVERLAY_SRC=/mnt/c/Users/ThomasBray/src/manx/manx-iso/toolkit-overlay
BIND=/mnt/c/Users/ThomasBray/src/manx/out
cd "$B"

echo "[1/5] extract base iso tree (xorriso osirrox, no root)"
rm -rf "$TREE"
xorriso -osirrox on -indev "$BASE" -extract / "$TREE" >/dev/null && echo "  extracted: $(find "$TREE" | wc -l) paths"

echo "[2/5] stage toolkit overlay + go binaries into tree"
mkdir -p "$TREE/manx/toolkit/bin" "$TREE/manx/toolkit/manx-iso-overlay" "$TREE/autorun"
cp /mnt/c/Users/ThomasBray/src/manx/out/manx "$TREE/manx/toolkit/bin/manx"
cp /mnt/c/Users/ThomasBray/src/manx/out/manx-tui "$TREE/manx/toolkit/bin/manx-tui"
chmod 755 "$TREE"/manx/toolkit/bin/*
cp /mnt/c/Users/ThomasBray/src/manx/manx-iso/toolkit-overlay/bin/*.sh "$TREE/manx/toolkit/manx-iso-overlay/"
chmod 755 "$TREE"/manx/toolkit/manx-iso-overlay/*.sh
mkdir -p "$TREE/manx/toolkit/manx-iso-overlay/autorun"
cp /mnt/c/Users/ThomasBray/src/manx/manx-iso/toolkit-overlay/autorun/00-manx-setup.sh "$TREE/manx/toolkit/manx-iso-overlay/autorun/"
chmod 755 "$TREE/manx/toolkit/manx-iso-overlay/autorun/00-manx-setup.sh"
cp /mnt/c/Users/ThomasBray/src/manx/manx-iso/toolkit-overlay/README.ISO.md "$TREE/manx/toolkit/"
# autorun at ISO root (SystemRescue looks here after boot -> sets up /toolkit)
cp /mnt/c/Users/ThomasBray/src/manx/manx-iso/toolkit-overlay/autorun/00-manx-setup.sh "$TREE/autorun/"
chmod 755 "$TREE/autorun/00-manx-setup.sh"

echo "[2b/5] audit-row for staged files (server-side manifest)"
find "$TREE/manx" "$TREE/autorun" -type f -exec sha256sum {} + | sort > "$TREE/manx/MANIFEST.txt"

echo "[3/5] capture preserved el-torito args from base"
xorriso -indev "$BASE" -report_el_torito as_mkisofs > "$B/eltorito-args.txt"
cat "$B/eltorito-args.txt"

echo "[4/5] rebuild iso (xorriso as mkisofs with preserved boot config) into out.iso"
cd "$TREE"
# capture args but drop the '-indev/-outdev/volid/relaxed' lines containing dev paths, then reapply volid
xorriso -as mkisofs $(sed 's/`-indev [^ ]*//' "$B/eltorito-args.txt" | sed 's/-outdev .*-o//' | tr -d '\n' | sed 's/\x27close-on-error.*//' ) \
  -o "$OUT" . >> "$B/build.log" 2>&1 || { echo "as_mkisofs call failed; showing args:"; cat "$B/eltorito-args.txt"; exit 1; }

echo "[5/5] verify and report"
ls -l "$OUT"
sha256sum "$OUT"
xorriso -indev "$OUT" -report_el_torito plain | head -20
