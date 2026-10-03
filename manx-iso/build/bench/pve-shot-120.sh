#!/usr/bin/env bash
set -u
qm screenshot $VMID --out-file /tmp/120-screenshot.png 2>&1 | tail -1
base64 -w0 /tmp/120-screenshot.png > /tmp/120-screenshot.png.b64
echo "screenshot: /tmp/120-screenshot.png"
