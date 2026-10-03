#!/usr/bin/env bash
set -euo pipefail
for s in bootstrap-drivers.sh rescue-menu.sh detect-hw.sh 00-manx-setup.sh; do
  pct push 119 /tmp/sc-overlay/$s /shellcheck/overlay/$s --perms 0755
done
pct exec 119 -- bash -c 'shellcheck --severity=style /shellcheck/overlay/00-manx-setup.sh /shellcheck/overlay/*.sh; echo RC=$?'
