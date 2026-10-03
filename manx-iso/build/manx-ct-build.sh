#!/usr/bin/env bash
set -euo pipefail
# manx v0.0.1 ISO build inside the manx-build LXC (unprivileged).
# xorriso + squashfs-tools + rsync + curl are already installed.
echo "== base iso (pinned 13.01, expected sha256 56289b69...) =="
cd /opt/manx-build
[ -f sysrescue-13.01-amd64.iso ] || \
  curl -fsSL --retry 3 -o sysrescue-13.01-amd64.iso \
    https://fastly-cdn.system-rescue.org/releases/13.01/systemrescue-13.01-amd64.iso
sha256sum sysrescue-13.01-amd64.iso

echo "== fetch sysrescue-customize (official) =="
if [ ! -f sysrescue-customize ]; then
  curl -fsSL 'https://gitlab.com/systemrescue/systemrescue-sources/-/raw/main/airootfs/usr/share/sysrescue/bin/sysrescue-customize?inline=false' \
    -o sysrescue-customize
  chmod +x sysrescue-customize
fi
head -3 sysrescue-customize

echo "== recipe staged =="
find recipe -type f | sort

echo "== run sysrescue-customize (auto) =="
bash sysrescue-customize \
  --auto \
  --source=/opt/manx-build/sysrescue-13.01-amd64.iso \
  --dest=/opt/manx-build/manx-iso-v0.0.1.iso \
  --recipe-dir=/opt/manx-build/recipe \
  --work-dir=/opt/manx-build/work \
  --overwrite \
  2>&1 | tee /opt/manx-build/customize.log | tail -25

echo "== verify ==" # sha + xorriso boot report
ls -l manx-iso-v0.0.1.iso
sha256sum manx-iso-v0.0.1.iso
xorriso -indev manx-iso-v0.0.1.iso -report_el_torito plain 2>/dev/null | head -20
echo "== MANX ISO v0.0.1 build finished =="
