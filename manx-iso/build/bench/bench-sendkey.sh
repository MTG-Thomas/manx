#!/usr/bin/env bash
# 1. guest has an IP that we can't tell without reading the serial; we can speak to the guest
#    via QEMU monitor (sendkey). This script types a minimal magic line to set up sshd + authorized keys
#    + echo IP on serial (so the laptop can read it) and then finish bench acceptance over ssh.
set -u
VMID=120
MONI='/tmp/120-monitor.txt'

send() {
  # send a single key composition character-by-character
  local keys="$1"
  qm monitor $VMID "sendkey $keys" >/dev/null 2>&1
}

type_line() {
  local line="$1" seqn=''
  local i ch
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
      :) seqn="$seqn shift-semicolon" ;;
      ';') seqn="$seqn semicolon" ;;
      '"') seqn="$seqn shift-apostrophe" ;;
      "'") seqn="$seqn apostrophe" ;;
      '(') seqn="$seqn shift-9" ;;
      ')') seqn="$seqn shift-0" ;;
      '=') seqn="$seqn equal" ;;
      _) seqn="$seqn shift-minus" ;;
      '>') seqn="$seqn shift-dot" ;;
      '$') seqn="$seqn shift-4" ;;
      '#') seqn="$seqn shift-3" ;;
      *) break ;;
    esac
  done
  send "$seqn"
  send kp_enter
  echo "typed: $line"
}

# pre: send Enter to ensure prompt, then clear the tty line with ctrl-u
send kp_enter
sleep 1
send ctrl-u
sleep 1

mkdirkey='mkdir /root/.ssh ;'
send "$mkdirkey"
# compose: echo key > authorized_keys
echokey='echo ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFYFscSweP+p+d4tGfe4UecO73Ttbf+DrbEiKnczJ9KL thomas@midtowntg.com > /root/.ssh/authorized_keys'
send "$echokey"
echo "[+] sendkey sshd bootstrap done"
