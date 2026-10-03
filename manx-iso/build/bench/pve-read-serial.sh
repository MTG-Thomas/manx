#!/usr/bin/env bash
# read the bench VM serial console for the manx bench banner (DHCP IP + sshd-ready)
# serial0 socket: /var/run/qemu-server/120.serial0 on the pve host
timeout=${BENCH_WAIT:-120}
end=$(( $(date +%s) + timeout ))
echo "waiting up to ${timeout}s for bench banner..."
found=""
while [ $(date +%s) -lt $end ]; do
  out=$(timeout 5 socat -u UNIX-CONNECT:/var/run/qemu-server/120.serial0 - 2>/dev/null | grep -aE 'bench-sshd (ready|iface)|MANX toolkit ready|manx-setup' | head -5)
  if [ -n "$out" ]; then found="$out"; break; fi
  sleep 5
done
if [ -n "$found" ]; then
  echo "=== BANNER ==="
  echo "$found"
else
  echo "no banner captured within ${timeout}s; last console follows:"
  timeout 5 socat -u UNIX-CONNECT:/var/run/qemu-server/120.serial0 - 2>/dev/null | tail -30
fi
