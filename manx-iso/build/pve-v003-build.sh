#!/usr/bin/env bash
set -euo pipefail
# pve-level v0.0.3 build cycle: stage fresh binaries + overlays + recipe YAML
# into CT 119's recipe dirs, build base + bench ISOs, verify hashes, publish to
# pve iso storage. Payloads (binaries, MOTD, autorun scripts, recipe YAMLs)
# arrive from the operator workstation over scp into /tmp/manx-v003/.
set -x
RECIPE=/opt/manx-build/recipe
BREC=/opt/manx-build/recipe-bench

# 1. binaries (rebuilt from merged main)
pct push 119 /tmp/manx-v003/manx /opt/manx-build/bin/manx --perms 0755
pct push 119 /tmp/manx-v003/manx-tui /opt/manx-build/bin/manx-tui --perms 0755
pct exec 119 -- bash -c '
set -e
for R in '"$RECIPE"' '"$BREC"'; do
  cp -f /opt/manx-build/bin/manx     "$R"/iso_add/manx/toolkit/bin/manx
  cp -f /opt/manx-build/bin/manx-tui "$R"/iso_add/manx/toolkit/bin/manx-tui
  chmod 755 "$R"/iso_add/manx/toolkit/bin/*
done
sha256sum '"$RECIPE"'/iso_add/manx/toolkit/bin/manx '"$BREC"'/iso_add/manx/toolkit/bin/manx
'

# 2. overlays + recipe YAMLs (canonical trees now live in-repo)
pct push 119 /tmp/manx-v003/MOTD.txt /opt/manx-build/recipe/iso_add/manx/toolkit/MOTD.txt
pct push 119 /tmp/manx-v003/MOTD.txt /opt/manx-build/recipe-bench/iso_add/manx/toolkit/MOTD.txt
pct push 119 /tmp/manx-v003/00-manx-setup.sh /opt/manx-build/recipe/iso_add/autorun/00-manx-setup.sh --perms 0755
pct push 119 /tmp/manx-v003/00-manx-setup.sh /opt/manx-build/recipe-bench/iso_add/autorun/00-manx-setup.sh --perms 0755
pct push 119 /tmp/manx-v003/00-manx-setup.sh /opt/manx-build/recipe-bench/iso_add/manx/toolkit/autorun/00-manx-setup.sh --perms 0755
pct push 119 /tmp/manx-v003/00-manx-bench-sshd.sh /opt/manx-build/recipe-bench/iso_add/autorun/00-manx-bench-sshd.sh --perms 0755
pct push 119 /tmp/manx-v003/10-base.yaml /opt/manx-build/recipe/iso_add/sysrescue.d/10-manx-autorun.yaml
pct push 119 /tmp/manx-v003/10-bench.yaml /opt/manx-build/recipe-bench/iso_add/sysrescue.d/10-manx-autorun.yaml
pct push 119 /tmp/manx-v003/20-bench.yaml /opt/manx-build/recipe-bench/iso_add/sysrescue.d/20-manx-bench.yaml

# 3. build both ISOs
pct exec 119 -- bash -c 'cd /opt/manx-build && rm -rf work && bash sysrescue-customize --auto \
  --source=/opt/manx-build/sysrescue-13.01-amd64.iso --dest=/opt/manx-build/manx-iso-v0.0.3.iso \
  --recipe-dir=/opt/manx-build/recipe --work-dir=/opt/manx-build/work --overwrite 2>&1 | tail -3'
pct exec 119 -- bash -c 'cd /opt/manx-build && rm -rf work && bash sysrescue-customize --auto \
  --source=/opt/manx-build/sysrescue-13.01-amd64.iso --dest=/opt/manx-build/manx-bench-iso-v0.0.3.iso \
  --recipe-dir=/opt/manx-build/recipe-bench --work-dir=/opt/manx-build/work --overwrite 2>&1 | tail -3'
pct exec 119 -- bash -c 'sha256sum /opt/manx-build/manx-iso-v0.0.3.iso /opt/manx-build/manx-bench-iso-v0.0.3.iso'

# 4. place into pve iso storage
pct exec 119 -- cat /opt/manx-build/manx-iso-v0.0.3.iso > /var/lib/vz/template/iso/manx-iso-v0.0.3.iso
pct exec 119 -- cat /opt/manx-build/manx-bench-iso-v0.0.3.iso > /var/lib/vz/template/iso/manx-bench-iso-v0.0.3.iso
sha256sum /var/lib/vz/template/iso/manx-iso-v0.0.3.iso /var/lib/vz/template/iso/manx-bench-iso-v0.0.3.iso
echo "== v0.0.3 build cycle finished =="
