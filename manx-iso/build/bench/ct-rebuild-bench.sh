#!/usr/bin/env bash
set -euo pipefail
# manx v0.0.1 bench-ISO rebuild (uses the sysrescue-customize official lane)
cd /opt/manx-build
echo "== fetch sysrescue-customize if missing =="
[ -f sysrescue-customize ] || curl -fsSL 'https://gitlab.com/systemrescue/systemrescue-sources/-/raw/main/airootfs/usr/share/sysrescue/bin/sysrescue-customize?inline=false' -o sysrescue-customize
chmod +x sysrescue-customize
echo "== rebuild bench iso ==" 
rm -rf work
bash sysrescue-customize \
  --auto \
  --source=/opt/manx-build/sysrescue-13.01-amd64.iso \
  --dest=/opt/manx-build/manx-bench-iso-v0.0.1.iso \
  --recipe-dir=/opt/manx-build/recipe-bench \
  --work-dir=/opt/manx-build/work \
  --overwrite 2>&1 | tail -12
echo "== verify ==" 
sha256sum manx-bench-iso-v0.0.1.iso
echo "== BENCH ISO (v2 autobuild) done =="
