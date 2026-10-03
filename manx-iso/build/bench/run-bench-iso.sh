#!/usr/bin/env bash
cd /opt/manx-build
[ -f sysrescue-customize ] || curl -fsSL 'https://gitlab.com/systemrescue/systemrescue-sources/-/raw/main/airootfs/usr/share/sysrescue/bin/sysrescue-customize?inline=false' -o sysrescue-customize
chmod +x sysrescue-customize
test -f recipe-bench/iso_add/autorun/00-manx-bench-sshd.sh
rm -rf work
bash sysrescue-customize \
  --auto \
  --source=/opt/manx-build/sysrescue-13.01-amd64.iso \
  --dest=/opt/manx-build/manx-bench-iso-v0.0.1.iso \
  --recipe-dir=/opt/manx-build/recipe-bench \
  --work-dir=/opt/manx-build/work \
  --overwrite 2>&1 | tail -15
sha256sum manx-bench-iso-v0.0.1.iso
echo "== BENCH ISO DONE =="
