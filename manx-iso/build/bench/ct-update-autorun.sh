#!/usr/bin/env bash
# bench-build: update 00-manx-setup.sh with serial reporting (no fancy quoting), then rebuild.
set -euo pipefail
echo "== bench build: update 00-manx-setup.sh, rebuild bench ISO =="
# 1. patched autorun with IP-to-serial report (uses tee to both console and serial)
cat > /opt/manx-build/recipe-bench/iso_add/autorun/00-manx-setup.sh <<'EOF'
#!/bin/bash
# MANX autorun: copy the toolkit from the boot media into the live system and
# print a one-line banner on BOTH the console and the serial (ttyS0)
# SPDX-License-Identifier: AGPL-3.0-or-later
set -u
BOOTMNT=/run/archiso/bootmnt
SRC="$BOOTMNT/manx/toolkit"
DST=/toolkit
{
  echo "manx-setup $(date -u +%FT%TZ)"
  ip -o -4 addr show 2>/dev/null | awk '{print "iface:", $2, $4}' | sed 's/^/REPORT:/'
  if [ -d "$SRC" ]; then
    mkdir -p "$DST/bin" "$DST/out" "$DST/manx-iso-overlay"
    cp -a "$SRC/bin/." "$DST/bin/"
    cp -a "$SRC/manx-iso-overlay/." "$DST/manx-iso-overlay/" 2>/dev/null || true
    chmod +x "$DST/bin/"*
    echo "MANX toolkit ready at /toolkit (menu: whiptail $DST/manx-iso-overlay/rescue-menu.sh, CLI: manx <verb>)"
    echo "REPORT: toolkit-ready"
  else
    echo "manx dir not found on boot media at $SRC"
    echo "REPORT: missing-toolkit-source"
  fi
  echo "REPORT: end-of-banner"
} 2>&1 | tee /dev/console 2>/dev/null | tee /dev/ttyS0 2>/dev/null
exit 0
EOF
chmod 755 /opt/manx-build/recipe-bench/iso_add/autorun/00-manx-setup.sh
echo "autorun updated (reports to console + serial)"
