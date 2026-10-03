#!/usr/bin/env bash
set -euo pipefail
# final v0.0.3 bench redeploy + zero-key boot proof
qm set 120 --ide2 local:iso/manx-bench-iso-v0.0.3.iso,media=cdrom
qm reset 120
echo reset-sent
sleep 135
printf 'screendump /tmp/v003-final-boot.ppm\n' | qm monitor 120 >/dev/null
echo done
