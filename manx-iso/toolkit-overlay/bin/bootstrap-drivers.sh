#!/bin/sh
# MANX bootstrap-drivers.sh - runtime driver fetch ladder (spec 8b.2.1).
# L1 local modprobe    (rescues stripped-remaster cases)
# L2 distro packages   (linux-modules + linux-firmware for the running kernel)
# L3 vendor lane       (Dell DSU for PERC/Broadcom-class hardware)
# L4 report-only       (adds exact vendor:device IDs for an off-box source)
# Deletion/destructive forbidden by default: this script installs, never removes.
# SPDX-License-Identifier: AGPL-3.0-or-later
# shellcheck shell=sh

set -eu
: "${OUT:=/toolkit/out}"
mkdir -p "$OUT"
LOG="$OUT/drivers.log"

runlog() {
  echo "+ $*" | tee -a "$LOG"
  "$@" 2>&1 | tee -a "$LOG"
}

show_unclaimed() {
  lspci -k 2>/dev/null | awk '
    /^[0-9a-f]+:/ { if (id) { printf "%s -> %s\n", id, (drv ? drv : "UNCLAIMED") } id=$0; drv="" }
    /Kernel driver in use/ { drv=$NF }
    END { if (id) { printf "%s -> %s\n", id, (drv ? drv : "UNCLAIMED") } }
  ' | grep UNCLAIMED || echo "no UNCLAIMED PCI devices reported by lspci"
}

KRN="$(uname -r)"

# ---- L1: local module aliases, no network needed -------------------------------
runlog echo "-- L1: local module aliases"
runlog show_unclaimed
if [ -d "/lib/modules/$KRN" ]; then
  runlog depmod -a "$KRN"
  MODS="$(lspci -k 2>/dev/null | awk '/Kernel modules:/ {print $3}' | sort -u | head -20)"
  if [ -n "$MODS" ]; then
    # shellcheck disable=SC2086
    # (intended word-splitting: one module name per modprobe argument)
    runlog modprobe -a $MODS
  fi
fi

# ---- L2: distro packages (fastest generic fix) --------------------------------
if command -v apt-get >/dev/null 2>&1; then
  runlog echo "-- L2: apt distro modules + firmware"
  runlog apt-get update
  # shellcheck disable=SC1091
  # (os-release not present on all boot media)
  DIST_ID="$(sed -n 's/^ID=//p' /etc/os-release 2>/dev/null | tr -d '\"' | head -1)"
  case "$DIST_ID" in
    ubuntu|debian)
      runlog env DEBIAN_FRONTEND=noninteractive apt-get install -y \
        --no-install-recommends \
        "linux-modules-$KRN" "linux-modules-extra-$KRN" \
        linux-firmware firmware-linux-free firmware-misc-nonfree
      ;;
    *)
      echo "L2 skipped: distro family is not ubuntu/debian, but $DIST_ID" | tee -a "$LOG"
      ;;
  esac
else
  echo "L2 skipped: no apt-get on media" | tee -a "$LOG"
fi

# ---- L2.5: reload modules from the freshly installed set ----------------------
runlog depmod -a "$KRN"
systemctl restart systemd-modules-load 2>/dev/null || true

# ---- L3: vendor lane (Dell DSU) -----------------------------------------------
if [ "${MANX_VENDOR_LANE:-1}" = 1 ] && [ -f /sys/class/dmi/id/product_name ]; then
  VEND="$(cat /sys/class/dmi/id/sys_vendor 2>/dev/null || true)"
  case "$VEND" in
    Dell*)
      runlog echo "-- L3: Dell lane via DSU bootstrap"
      runlog curl -fsSL https://linux.dell.com/repo/hardware/dsu/bootstrap.cgi -o /tmp/bootstrap.cgi
      runlog sh /tmp/bootstrap.cgi
      if command -v dsu >/dev/null 2>&1; then
        runlog dsu --collect-details
        runlog dsu --driver
      else
        echo "L3: dsu not available after DSU bootstrap; the DSU bootstrap run is in the log" | tee -a "$LOG"
      fi
      ;;
    *)
      echo "L3 skipped: vendor lane for $VEND not implemented" | tee -a "$LOG"
      ;;
  esac
fi

# ---- L4: report ---------------------------------------------------------------
echo "-- L4: report" | tee -a "$LOG"
if command -v /toolkit/manx-iso-overlay/detect-hw.sh >/dev/null 2>&1; then
  runlog /toolkit/manx-iso-overlay/detect-hw.sh
else
  echo "L4: detect-hw.sh not yet staged; run it manually if available" | tee -a "$LOG"
fi
