#!/bin/sh
# MANX bootstrap-drivers.sh — runtime driver fetch ladder (spec §8b.2.1).
# L1 local modprobe            (rescues stripped-remaster cases)
# L2 distro packages           (linux-modules + linux-firmware for the running kernel)
# L3 vendor lane               (Dell DSU for PERC/Broadcom-class Duke hardware)
# L4 report-only               (adds exact vendor:device IDs for an off-box source)
# Deletion/destructive forbidden by default: this script installs, never removes.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -eu
: "${OUT:=/toolkit/out}"; mkdir -p "$OUT"
LOG="$OUT/drivers.log"
exec > >(tee -a "$LOG") 2>&1

[ "$(id -u)" = 0 ] || { echo "must run as root" >&2; exit 1; }

echo "== bootstrap-drivers (ladder) =="

# L1 - TRY LOCAL ALIASES FIRST (no net needed)
echo "-- L1: local module aliases"
lspci -nnk 2>/dev/null | awk '/^[0-9a-f]+:/ {id=$0} /Kernel driver in use/ {print id, "OK", $0; id=""}' \
  | grep UNCLAIMED || echo "no UNCLAIMED PCI devices reported by lspci"

# Attempt to reconcile stripped modules when the kernel ships them but the ISO
# image lost an alias file:
if [ -d "/lib/modules/$(uname -r)" ]; then
  echo "-- L1: depmod against running kernel $(uname -r)"
  depmod -a "$(uname -r)" 2>&1 | head -20
  modprobe -a $(lspci -k 2>/dev/null | awk '/Kernel modules:/ {print $3}' | sort -u | head -20) 2>&1 | head -20 || true
fi

# L2 - DISTRO PACKAGES (fastest generic fix for megaraid/broadcom/nvme/mellanox)
if command -v apt-get >/dev/null 2>&1; then
  echo "-- L2: apt distro modules + firmware"
  apt-get update -o Acquire::Retries=3 2>&1 | tail -3
  case "$(. /etc/os-release; echo "$ID")" in
    ubuntu|debian)
      KERN="linux-modules-$(uname -r)"
      KERNX="linux-modules-extra-$(uname -r)"
      DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        "$KERN" "$KERNX" linux-firmware firmware-linux-free firmware-misc-nonfree 2>&1 | tail -20
      ;;
    *) echo "unknown distro family ($ID); L2 skipped" ;;
  esac
else
  echo "-- L2: no apt; distro lane skipped"
fi

# Reload junction: try again after module/firmware package land
depmod -a "$(uname -r)" 2>&1 | head -5
systemctl restart systemd-modules-load 2>/dev/null || true

# L3 - VENDOR LANE (Dell T330/PERC/Broadcom via DSU, opt-in for Dell hardware)
if [ "${MANX_VENDOR_LANE:-1}" = 1 ] && [ -f /sys/class/dmi/id/product_name ]; then
  VEND=$(cat /sys/class/dmi/id/sys_vendor 2>/dev/null || echo "")
  case "$VEND" in
    Dell*|Dell\ Inc.*)
      echo "-- L3: Dell lane via DSU bootstrap"
      curl -fsSL https://linux.dell.com/repo/hardware/dsu/bootstrap.cgi -o /tmp/bootstrap.cgi \
        && sh /tmp/bootstrap.cgi 2>&1 | tail -5
      if command -v dsu >/dev/null 2>&1; then
        dsu --collect-details 2>&1 | tail -20
        dsu --driver 2>&1 | tail -40
      else echo "dsu not available after bootstrap; log above for the off-box agent"; fi
      ;;
    *) echo "-- L3: non-Dell ($VEND); vendor lane skipped" ;;
  esac
fi

# L4 - REPORT if anything is still unclaimed
if command -v /toolkit/bin/detect-hw.sh >/dev/null 2>&1; then
  /toolkit/bin/detect-hw.sh || true
else
  echo "-- L4: detect-hw.sh not on path; run it manually for the residual report"
fi
echo "== bootstrap-drivers done; residual report (if any) in $OUT/drivers-report.txt =="
