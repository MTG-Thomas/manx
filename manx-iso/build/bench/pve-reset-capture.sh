#!/usr/bin/env bash
set -euo pipefail
# pve-level: fresh bench ISO into local iso storage, then start serial capture, then reset VM 120
pct exec 119 -- cat /opt/manx-build/manx-bench-iso-v0.0.1.iso > /var/lib/vz/template/iso/manx-bench-iso-v0.0.1.iso
sha256sum /var/lib/vz/template/iso/manx-bench-iso-v0.0.1.iso
# start the serial capture in the background for 240 s
nohup timeout 240 socat -u UNIX-CONNECT:/var/run/qemu-server/120.serial0 RECV:/tmp/manx-bench-serial.txt > /tmp/socat-serial.out 2>&1 &
echo socat-started
qm set 120 --ide2 local:iso/manx-bench-iso-v0.0.1.iso,media=cdrom
qm reset 120
echo reset-sent
