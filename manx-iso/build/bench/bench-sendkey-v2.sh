#!/usr/bin/env bash
# manx bench: send key sequences over the QEMU console via QEMU monitor, using
# the unmaintained `sendkey` mechanism. Line-by-line.
# Also useful as documentation of what was done.
set -u
VMID=120
send() {
  local keys="$1"
  qm monitor $VMID "sendkey $keys" >/dev/null 2>&1
}
type_line() {
  local line="$1" seqn='' i ch
  for i in $(seq 1 ${#line}); do
    ch="${line:i-1:1}"
    case "$ch" in
      [a-z]) seqn="$seqn $ch" ;;
      [A-Z]) seqn="$seqn shift-$ch" ;;
      [0-9]) seqn="$seqn $ch" ;;
      ' ') seqn="$seqn spc" ;;
      -) seqn="$seqn minus" ;;
      .) seqn="$seqn dot" ;;
      /) seqn="$seqn slash" ;;
      ;) seqn="$seqn semicolon" ;;
      '"') seqn="$seqn shift-apostrophe" ;;
      "'") seqn="$seqn apostrophe" ;;
      '(') seqn="$seqn shift-9" ;;
      ')') seqn="$seqn shift-0" ;;
      '=') seqn="$seqn equal" ;;
      _) seqn="$seqn shift-minus" ;;
      '>') seqn="$seqn shift-dot" ;;
      '$') seqn="$seqn shift-4" ;;
      '#') seqn="$seqn shift-3" ;;
      *) return 1 ;;
    esac
  done
  send "$seqn"
  send kp_enter
  echo "[sendkey] typed: ${line:0:60}..."
}

# --- commands to type on the guest console ---
type_line "mkdir /root/.ssh"
type_line "echo 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFYFscSweP+p+d4tGfe4UecO73Ttbf+DrbEiKnczJ9KL thomas@midtowntg.com' > /root/.ssh/authorized_keys"
type_line "sed -i 's/#\?PermitRootLogin.*/PermitRootLogin yes/' /etc/ssh/sshd_config"
type_line "systemctl start sshd"
type_line "ip -o -4 addr show eth0 | tail -1 > /dev/ttyS0"
sleep 1
