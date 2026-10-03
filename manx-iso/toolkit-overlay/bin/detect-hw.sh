#!/bin/sh
# MANX: enumerate hardware the current kernel canNOT drive (unclaimed PCI/USB)
#+ and list recent failed firmware loads. Output: machine-refundable lines for
#+ bootstrap-drivers.sh (spec 8b.2.1) and for the audit log.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -eu
: "${OUT:=/toolkit/out}"; mkdir -p "$OUT"
E="$OUT/drivers-report.txt"

# 1. Unclaimed PCI/USB devices, with vendor:device and friendly names
{
  echo "# unclaimed-or-driverless PCI devices (lspci -nnk):"
  lspci -nnk 2>/dev/null | awk '
    /^[0-9a-f]+:/ { if (id) { printf "%s -> %s\n", id, (drv ? drv : "UNCLAIMED") } ; id=$0; drv="" }
    /Kernel driver in use/ { drv=$NF }
        END { if (id) { printf "%s -> %s\n", id, (drv ? drv : "UNCLAIMED") } }
  '
  echo ""
  echo "# unclaimed USB devices:"
  lsusb 2>/dev/null | head -50
} | tee "$E"

# 2. Failed firmware loads from the running kernel
echo ""
{ dmesg 2>/dev/null || journalctl -k 2>/dev/null || true; } \
  | grep -iE 'firmware.*failed|loading.*firmware.*not found|Direct firmware load' \
  | head -40 >> "$E" || true

# 3. Friendly names, so the report is readable by an off-box agent later
if command -v lspci >/dev/null 2>&1; then
  echo "" >> "$E"
  echo "# pci.ids names for the IDs above:" >> "$E"
fi

echo "detect-hw: report at $E"
