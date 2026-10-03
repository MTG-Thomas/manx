#!/usr/bin/env bash
# capture the last console lines via the serial socket (readable dump)
timeout 6 socat - UNIX-CONNECT:/var/run/qemu-server/120.serial0 2>/dev/null | tail -40
