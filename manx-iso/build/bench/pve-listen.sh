#!/usr/bin/env bash
# run a tiny HTTP listener on the pve host that records every request (to see the guest's callback)
nohup python3 -m http.server 4097 --bind 0.0.0.0 --directory /tmp/listen > /tmp/listen.log 2>&1 &
echo "listener 4097 started"
sleep 1
curl -s http://127.0.0.1:4097/pve-local-check >/dev/null && echo "listener reachable locally"
