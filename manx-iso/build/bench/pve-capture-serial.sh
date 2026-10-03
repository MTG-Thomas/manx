# start the socat serial capture NOW (before the VM resets) so we can capture
# the goto-IP line from the autorun script as it boots.
BENCH_WAIT=${BENCH_WAIT:-180}
end=$(( $(date +%s) + BENCH_WAIT ))
: > /tmp/manx-bench-serial.txt
timeout ${BENCH_WAIT} socat -u UNIX-CONNECT:/var/run/qemu-server/120.serial0 RECV:/tmp/manx-bench-serial.txt 2>/dev/null
echo "=== captured serial (last 30 lines) ==="
tail -30 /tmp/manx-bench-serial.txt
