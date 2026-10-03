#!/usr/bin/env bash
set -euo pipefail
# pve-level driver for the v0.0.2 ISO build cycle (base + bench variants)
# CT 119 = manx-build (Debian 13 LXC). Payloads arrive from the operator
# workstation over scp; this script is Idempotent per step.
echo "=== 1. push latest binaries into CT ==="
mkdir -p /tmp/manx-payload/bin /tmp/manx-payload
scp_from_host() { echo "(scped from workstation)"; }
pct push 119 /tmp/manx-payload/bin/manx /opt/manx-build/bin/manx --perms 0755
pct push 119 /tmp/manx-payload/bin/manx-tui /opt/manx-build/bin/manx-tui --perms 0755
pct push 119 /tmp/manx-payload/autorun/00-manx-setup.sh /opt/manx-build/recipe/iso_add/autorun/00-manx-setup.sh --perms 0755
pct push 119 /tmp/manx-payload/autorun/00-manx-setup.sh /opt/manx-build/recipe-bench/iso_add/autorun/00-manx-setup.sh --perms 0755
pct push 119 /tmp/manx-payload/MOTD.txt /opt/manx-build/recipe/iso_add/manx/toolkit/MOTD.txt
pct push 119 /tmp/manx-payload/MOTD.txt /opt/manx-build/recipe-bench/iso_add/manx/toolkit/MOTD.txt
pct exec 119 -- bash -c 'ls -l /opt/manx-build/recipe/iso_add/manx/toolkit/MOTD.txt /opt/manx-build/bin/manx'

echo "=== 2. build base release ISO (v0.0.2) ==="
pct exec 119 -- bash -c 'cd /opt/manx-build && rm -rf work && bash sysrescue-customize --auto --source=/opt/manx-build/sysrescue-13.01-amd64.iso --dest=/opt/manx-build/manx-iso-v0.0.2.iso --recipe-dir=/opt/manx-build/recipe --work-dir=/opt/manx-build/work --overwrite 2>&1 | tail -3'
pct exec 119 -- bash -c 'sha256sum /opt/manx-build/manx-iso-v0.0.2.iso'

echo "=== 3. build bench ISO (v0.0.2) ==="
pct exec 119 -- bash -c 'cd /opt/manx-build && rm -rf work && bash sysrescue-customize --auto --source=/opt/manx-build/sysrescue-13.01-amd64.iso --dest=/opt/manx-build/manx-bench-iso-v0.0.2.iso --recipe-dir=/opt/manx-build/recipe-bench --work-dir=/opt/manx-build/work --overwrite 2>&1 | tail -3'
pct exec 119 -- bash -c 'sha256sum /opt/manx-build/manx-bench-iso-v0.0.2.iso'
