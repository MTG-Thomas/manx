#!/usr/bin/env bash
set -euo pipefail
# v0.0.3 bench deploy: updated ISO, long serial capture, reset, then a few
# checks INCLUDING a login-free look: we deliberately send NO keystrokes.
qm set 120 --ide2 local:iso/manx-bench-iso-v0.0.3.iso,media=cdrom
nohup timeout 420 socat -u UNIX-CONNECT:/var/run/qemu-server/120.serial0 RECV:/tmp/manx-serial-v003.txt > /tmp/socat-v003.out 2>&1 &
qm reset 120
echo reset-sent
sleep 135
printf 'screendump /tmp/v003-boot.ppm\n' | qm monitor 120 >/dev/null
echo "--- pve -> guest ssh (auto sshd + keys should just work) ---"
ssh -o BatchMode=yes -o ConnectTimeout=10 -o StrictHostKeyChecking=accept-new root@172.16.15.171 'hostname; ls /toolkit/bin | head -3; cat /toolkit/.setup-complete 2>/dev/null; manx status 2>/dev/null | head -3' || echo "guest ssh failed"
echo "--- serial capture grep ---"
grep -a 'MANX toolkit\|setup ok\|REPORT' /tmp/manx-serial-v003.txt 2>/dev/null | head -6 || true
ls -l /tmp/v003-boot.ppm
