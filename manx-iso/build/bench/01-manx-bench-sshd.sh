#!/bin/bash
# MANX BENCH-ONLY autorun (never in the public release): start sshd with
# MTG operator's public key and print a bench-ready banner on the serial console.
# Public keys are not secrets - this is the least-interactive bench injection.
set -u
KEY="ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFYFscSweP+p+d4tGfe4UecO73Ttbf+DrbEiKnczJ9KL thomas@midtowntg.com"
{
  echo "bench-sshd init $(date -u +%FT%TZ)"
  mkdir -p /root/.ssh
  printf '%s\n' "$KEY" > /root/.ssh/authorized_keys
  chmod 700 /root/.ssh; chmod 600 /root/.ssh/authorized_keys
  if [ -f /etc/ssh/sshd_config ]; then
    sed -i 's/^#\?PermitRootLogin.*/PermitRootLogin prohibit-password/' /etc/ssh/sshd_config
    systemctl restart sshd 2>/dev/null || service ssh restart 2>/dev/null || true
  fi
  systemctl start sshd 2>/dev/null || service ssh start 2>/dev/null || true
  ip -o -4 addr show 2>/dev/null | awk '{print "bench-sshd iface:", $2, $4}'
  echo "bench-sshd ready; operator key installed (not a secret: public half only)"
} 2>&1 | tee /dev/console | tee -a /var/log/manx-setup.log
