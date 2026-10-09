#!/usr/bin/env bash
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
# Run inside a disposable Manx VM, never a customer recovery session.
phase=${1:?register, offline, or sync}
endpoint=${2:?public lab Worker HTTPS origin}
evidence=${3:?operator-owned persistent acceptance directory}
manx_binary=${MANX_BINARY:-manx}
umask 077
mkdir -p "$evidence/state"
case "$phase" in
  register)
    "$manx_binary" session register --endpoint "$endpoint" --state-dir "$evidence/state" > "$evidence/registered.json"
    jq -e '.session_id and .pairing_code and .state == "awaiting_claim"' "$evidence/registered.json" >/dev/null
    cat "$evidence/registered.json"
    ;;
  offline)
    "$manx_binary" session record --endpoint "$endpoint" --state-dir "$evidence/state" --diagnostic hardware_inventory > "$evidence/offline.json"
    if HTTPS_PROXY=http://127.0.0.1:9 NO_PROXY= ALL_PROXY= "$manx_binary" session sync --endpoint "$endpoint" --state-dir "$evidence/state"; then
      echo 'Expected simulated network interruption did not occur.' >&2; exit 1
    fi
    "$manx_binary" session status --state-dir "$evidence/state" > "$evidence/interrupted.json"
    jq -e '.pending_events == 1' "$evidence/interrupted.json" >/dev/null
    jq -r '.pending[0].event_id' "$evidence/state/checkpoint.json" > "$evidence/diagnostic-event-id.txt"
    ;;
  sync)
    "$manx_binary" session sync --endpoint "$endpoint" --state-dir "$evidence/state" > "$evidence/synced.json"
    "$manx_binary" session sync --endpoint "$endpoint" --state-dir "$evidence/state" > "$evidence/repeated.json"
    jq -e '.state == "approved" and .pending_events == 0 and .acknowledged_sequence == 1' "$evidence/synced.json" >/dev/null
    jq -e '.acknowledged_sequence == 1' "$evidence/repeated.json" >/dev/null
    test "$(jq -r .session_id "$evidence/registered.json")" = "$(jq -r .session_id "$evidence/synced.json")"
    echo 'Client interruption checks passed. Verify the recorded event ID appears exactly once in Workspace and save its evidence.'
    ;;
  *) echo 'Expected register, offline, or sync.' >&2; exit 2 ;;
esac
