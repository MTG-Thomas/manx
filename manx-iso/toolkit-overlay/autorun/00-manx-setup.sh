#!/bin/bash
# MANX autorun: copy the toolkit from the boot media into the live system and
# print a one-line banner. Runs automatically at boot (SystemRescue autorun).
# Sets up /toolkit from /run/archiso/bootmnt/manx/toolkit so the manx CLI, the
# Bubble Tea TUI and the Tier-1 whiptail menu are all available immediately.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -u
BOOTMNT=/run/archiso/bootmnt
SRC="$BOOTMNT/manx/toolkit"
DST=/toolkit
LOG=/var/log/manx-setup.log
{
  echo "manx-setup $(date -u +%FT%TZ)"
  if [ -d "$SRC" ]; then
    mkdir -p "$DST/bin" "$DST/out" "/toolkit/manx-iso-overlay"
    cp -a "$SRC/bin/." "$DST/bin/" 2>/dev/null
    ip -o -4 addr show 2>/dev/null | awk '{print "  iface:", $2, $4}' | head -5
    cp -a "$SRC/manx-iso-overlay/." "/toolkit/manx-iso-overlay/" 2>/dev/null || true
    cp -a "$SRC/MOTD.txt" "$DST/MOTD.txt" 2>/dev/null || true
    if [ -x "$DST/bin/manx" ]; then
      export PATH="$DST/bin:$PATH"
      chmod +x "$DST/bin/"*
      echo "MANX toolkit ready at /toolkit (menu: whiptail $DST/manx-iso-overlay/rescue-menu.sh, CLI: manx <verb>)"
      # motd: incident-response truths + creature care, on console and /etc/motd
      if [ -f "$DST/MOTD.txt" ]; then
        cp "$DST/MOTD.txt" /etc/manx-motd 2>/dev/null || true
        cat "$DST/MOTD.txt" | tee /dev/ttyS0 2>/dev/null || cat "$DST/MOTD.txt"
        grep -q manx-motd /root/.bashrc 2>/dev/null || \
          echo "test -f /etc/manx-motd && cat /etc/manx-motd" >> /root/.bashrc
      fi
    else
      echo "manx binaary missing from live layer; menu script still at $DST/manx-iso-overlay/"
    fi
  else
    echo "manx dir not found on boot media at $SRC"
  fi
  echo "audit log: $DST/out/audit.log"
} 2>&1 | tee -a "$LOG" | tee /dev/console
# Provide the menu even when the banner is skipped: PATH must be set for interactive shells.
echo "export PATH=/toolkit/bin:\$PATH; alias rescue-menu='/toolkit/manx-iso-overlay/rescue-menu.sh'; alias manx='/toolkit/bin/manx'; alias manx-tui='/toolkit/bin/manx-tui;'" >> /root/.bashrc
echo "  (added manx paths to /root/.bashrc)" | tee -a "$LOG"
exit 0
