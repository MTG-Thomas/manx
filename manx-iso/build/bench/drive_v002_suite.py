# -*- coding: utf-8 -*-
"""v0.0.2 ssh driving: login (the guest kept lease 171), check /toolkit state,
run the autorun if it didn't happen, run acceptance verbs, report the motd,
and capture raw output for the bench report."""
import subprocess, sys, time

SSH = ["ssh", "-o", "BatchMode=yes", "-o", "ConnectTimeout=10",
       "-o", "StrictHostKeyChecking=accept-new", "-i", r"C:\Users\ThomasBray\.ssh\id_ed25519",
       "root@172.16.15.171"]

def run(cmd, timeout=90):
    try:
        r = subprocess.run(SSH + [cmd], capture_output=True, text=True, timeout=timeout)
        return (r.stdout + r.stderr).strip()
    except subprocess.TimeoutExpired:
        return "(timeout)"

checks = [
    ("ssh-in", "echo SSH-IN-OK; hostname"),
    ("toolkit-present?", "ls -la /toolkit 2>&1 | head -6"),
    ("autorun banner? (boot log)", "grep -a 'MANX toolkit ready\\|manx-setup\\|REPORT:' /var/log/manx-setup.log 2>/dev/null | tail -8"),
    ("motd on file?", "ls -l /toolkit/MOTD.txt /etc/manx-motd 2>&1 | head -3"),
    ("motd ghost check (console tail)", "cat /etc/manx-motd 2>/dev/null | head -4"),
    ("manx status (real)", "/toolkit/bin/manx status"),
    ("manx detect-hw (real)", "/toolkit/bin/manx detect-hw 2>&1 | tail -6"),
    ("detect-hw report file", "ls -l /toolkit/out/drivers-report.txt 2>&1; head -5 /toolkit/out/drivers-report.txt 2>/dev/null"),
    ("manx img-in (real)", "/toolkit/bin/manx img-in 2>&1 | head -8"),
    ("img-in report file", "ls -l /toolkit/out/img-in-report.txt 2>&1; head -6 /toolkit/out/img-in-report.txt 2>/dev/null"),
    ("audit rows", "wc -l /toolkit/out/audit.log 2>&1; tail -8 /toolkit/out/audit.log 2>&1"),
    ("destructive gate", "/toolkit/bin/manx img-out 2>&1; echo rc=$?"),
    ("gate with i-know", "/toolkit/bin/manx img-out --i-know 2>&1; echo rc=$?"),
]

body = ""
for name, cmd in checks:
    print("=" * 8, name)
    out = run(cmd)
    print(out)
    body += "\n==== %s ====\n%s\n" % (name, out)

open(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\bench-v002-raw.txt", "w", encoding="utf-8").write(body)
print("\nraw captured to bench-v002-raw.txt")
