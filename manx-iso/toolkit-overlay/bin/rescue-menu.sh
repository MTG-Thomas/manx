#!/bin/sh
# MANX rescue menu — Tier-1 view (whiptail). Same action surface as `manx <verb>`.
# Every menu entry calls the SAME implementation as `manx <verb>` (spec §10.8 parity).
# Destructive verbs keep the --i-know gate here, and write audit rows via manx.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -eu
OUT=/toolkit/out; mkdir -p "$OUT"

banner() {
  echo "MANX rescue menu"
  echo "scratch: $(df -h /toolkit 2>/dev/null | tail -1 | awk '{print $4}') available"
}

ask_i_know() {
  # $1 = label
  whiptail --title "destructive: $1" \
    --yesno "This is a DESTRUCTIVE verb.\nIt requires explicit confirmation (spec §10.8).\nProceed?" 12 60
}

run_verb() {
  # wraps `manx <verb>`; passes --i-know only after human confirmation
  verb="$1"; shift
  if grep -q ",danger" <<EOF
$(manx --actions 2>/dev/null)
EOF
  then :; fi  # danger detection is authoritative in the Go binary
  if manx "$verb" --check-danger >/dev/null 2>&1 && ! ask_i_know "$verb"; then
    manx audit denied --verb "$verb" 2>/dev/null || true
    return
  fi
  set +e
  manx "$verb" "$@" --i-know 2>&1 | tee -a "$OUT/last-run.log"
  rc=$?
  set -e
  [ $rc -eq 0 ] || whiptail --msgbox "verb '$verb' exited with rc=$rc (see $OUT/last-run.log)" 12 60
  return $rc
}

MENU_LIST=$(cat <<NAV
status            one-line runtime summary (net/mesh/drivers/scratch)
detect-hw         inventory unclaimed / missing-driver hardware
collect           harvest logs, evtx, minidumps, catalog listings (read-only)
img-in            disk/volume inventory (read-only)
img-out           image a volume/partition to a file (destructive)
hive-inspect      read-only hive dump (hivexsh / reglookup / chntpw view)
hive-edit         SAFE offline hive edit (snapshot first, destructive)
bootstrap-drivers runtime driver fetch ladder L1 local / L2 distro / L3 vendor (mutating)
bringup-vm        qemu/OVMF bringup of the Windows guest (destructive)
net-up            bring up DHCP + drivers for the first NIC
mesh-up           join the DN/Nebula overlay with a runtime-issued cert
channel-up        tailcat fallback channel (no control plane needed)
status-tail       tail the audit log (watch what happened)
shell             drop to a rescue bash under the supervisor
NAV
)

# build the whiptail menu args from the verb list, one entry at a time
# (deliberately word-split-tagged: the loop-controlled construction is the point)
MENU_ARGS=""
while IFS= read -r line_entry; do
  verb_name=$(printf '%s' "$line_entry" | awk '{print $1}')
  verb_help=$(printf '%s' "$line_entry" | cut -d' ' -f2-)
  MENU_ARGS="$MENU_ARGS \"$(printf '%s' "$verb_name")\" \"$(printf '%s' "$verb_help")\""
done <<EOF
$MENU_LIST
EOF
# shellcheck disable=SC2086  # we control the quoting of MENU_ARGS construction
while true; do
  banner
  SEL=$(eval whiptail --title '"MANX rescue menu"' --menu '"Pick an action (spec §10.8)"' 24 78 16 $MENU_ARGS 3>&1 1>&2 2>&3)
  rc=$?
  [ $rc -ne 0 ] && whiptail --msgbox "Quitting the menu drops to the rescue shell. Type 'menu' to return." 10 50 && break
  case "$SEL" in
    status|detect-hw|collect|img-in)                       run_verb "$SEL" ;;
    img-out)                                              run_verb img-out ;;
    hive-inspect)                                         run_verb hive-inspect ;;
    hive-edit)                                            run_verb hive-edit ;;
    bootstrap-drivers)                                    run_verb bootstrap-drivers ;;
    bringup-vm)                                           run_verb bringup-windows-vm ;;
    net-up)                                               run_verb net-up ;;
    mesh-up)                                              run_verb mesh-up ;;
    channel-up)                                           run_verb channel-up ;;
    status-tail)                                          manx audit tail ;;
    shell)                                                bash ;;
    *)                                                    whiptail --msgbox "unknown selection: $SEL" 10 40 ;;
  esac
done
