#!/usr/bin/env bash
set -euo pipefail
# re-stage the freshly built binaries into BOTH recipes (they were pushed to
# /opt/manx-build/bin in the v0.0.2 cycle but the recipe dirs kept v0.0.1 copies)
pct exec 119 -- bash -c '
cp -f /opt/manx-build/bin/manx    /opt/manx-build/recipe/iso_add/manx/toolkit/bin/manx
cp -f /opt/manx-build/bin/manx-tui /opt/manx-build/recipe/iso_add/manx/toolkit/bin/manx-tui
cp -f /opt/manx-build/bin/manx    /opt/manx-build/recipe-bench/iso_add/manx/toolkit/bin/manx
cp -f /opt/manx-build/bin/manx-tui /opt/manx-build/recipe-bench/iso_add/manx/toolkit/bin/manx-tui
chmod -R 755 /opt/manx-build/recipe /opt/manx-build/recipe-bench
sha256sum /opt/manx-build/recipe/iso_add/manx/toolkit/bin/manx /opt/manx-build/recipe-bench/iso_add/manx/toolkit/bin/manx
'
pct exec 119 -- bash -c 'cd /opt/manx-build && rm -rf work && bash sysrescue-customize --auto --source=/opt/manx-build/sysrescue-13.01-amd64.iso --dest=/opt/manx-build/manx-bench-iso-v0.0.2.iso --recipe-dir=/opt/manx-build/recipe-bench --work-dir=/opt/manx-build/work --overwrite 2>&1 | tail -3'
pct exec 119 -- bash -c 'sha256sum /opt/manx-build/manx-bench-iso-v0.0.2.iso'
