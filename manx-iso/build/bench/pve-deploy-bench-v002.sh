#!/usr/bin/env bash
set -euo pipefail
# pve-level: place v0.0.2 bench ISO into local iso storage, capture serial, reset VM 120
pct exec 119 -- cat /opt/manx-build/manx-bench-iso-v0.0.2.iso > /var/lib/vz/template/iso/manx-bench-iso-v0.0.2.iso
sha256sum /var/lib/vz/template/iso/manx-bench-iso-v0.0.2.iso
nohup timeout 300 socat -u UNIX-CONNECT:/var/run/qemu-server/120.serial0 RECV:/tmp/manx-bench-serial-v002.txt > /tmp/socat-serial-v002.out 2>&1 &
echo socat-started
qm set 120 --ide2 local:iso/manx-bench-iso-v0.0.2.iso,media=cdrom
qm reset 120
echo reset-sent
