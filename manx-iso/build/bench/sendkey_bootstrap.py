# -*- coding: utf-8 -*-
"""Manx bench: keyboard-inject sshd bootstrap into the pve-t340 bench VM (UMID 120)
via `qm sendkey`, so the laptop can ssh in. Then the serial (ttyS0) gives us the
guest IP. All public-key only (no secrets).

Types these lines on the guest console:
  mkdir /root/.ssh
  echo ssh-ed25519 <pubkey> thomas@midtowntg.com > /root/.ssh/authorized_keys
  echo PermitRootLogin yes >> /etc/ssh/sshd_config
  systemctl restart sshd
  ip -o -4 addr show eth0 > /dev/ttyS0     <-- read on pve via socat
"""
import subprocess, time, sys

PVE = "root@172.16.15.128"
VMID = "120"
KEY = ["-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa"]

CHAR_MAP = {
    ".": "dot", "-": "minus", " ": "spc", "/": "slash", ";": "semicolon",
    "'": "apostrophe", '"': "shift-apostrophe", "\\": "backslash",
    "(": "shift-9", ")": "shift-0", "=": "equal", "_": "shift-minus",
    ">": "shift-dot", "$": "shift-4", "#": "shift-3", "@": "shift-2",
    ":": "shift-semicolon", "+": "shift-equal",
}
for c in "abcdefghijklmnopqrstuvwxyz":
    CHAR_MAP[c] = c
for c in "0123456789":
    CHAR_MAP[c] = c
for c in "ABCDEFGHIJKLMNOPQRSTUVWXYZ":
    CHAR_MAP[c] = "shift-" + c.lower()

def sendkey(seq):
    subprocess.run(["ssh", "-o", "BatchMode=yes"] + KEY + [PVE, "qm", "sendkey", VMID, seq],
                   stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

def type_line(line):
    seq = []
    for ch in line:
        k = CHAR_MAP.get(ch)
        if k is None:
            raise RuntimeError(f"no key mapping for {ch!r}")
        seq.append(k)
    # small batches to stay under monitor command length
    for i in range(0, len(seq), 48):
        sendkey(" ".join(seq[i:i+48]))
        time.sleep(0.02)
    sendkey("ret")  # enter
    time.sleep(0.35)

BANNER = "echo manx-bench-ok > /dev/ttyS0"
LINES = [
    BANNER,  # first: prove we're at a root shell
    "mkdir /root/.ssh",
    "echo ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFYFscSweP+p+d4tGfe4UecO73Ttbf+DrbEiKnczJ9KL "
    "thomas@midtowntg.com > /root/.ssh/authorized_keys",
    "echo PermitRootLogin yes >> /etc/ssh/sshd_config",
    "systemctl restart sshd",
    "ip -o -4 addr show eth0 > /dev/ttyS0",
]

print("typing", len(LINES), "lines...")
for i, line in enumerate(LINES, 1):
    type_line(line)
    print(f"  typed {i}/{len(LINES)}")
time.sleep(1)
print("done; read serial on pve: socat - UNIX-CONNECT:/var/run/qemu-server/120.serial0")
