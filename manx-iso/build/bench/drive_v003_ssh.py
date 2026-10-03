# -*- coding: utf-8 -*-
"""v0.0.3 acceptance over SSH (laptop -> guest 172.16.15.171)."""
import subprocess

SSH = ["ssh", "-o", "BatchMode=yes", "-o", "ConnectTimeout=10",
       "-o", "StrictHostKeyChecking=accept-new", "-i",
       r"C:\Users\ThomasBray\.ssh\id_ed25519", "root@172.16.15.171"]

def run(label, cmd, timeout=90):
    try:
        r = subprocess.run(SSH + [cmd], capture_output=True, text=True, timeout=timeout)
        out = (r.stdout + r.stderr).strip()
    except subprocess.TimeoutExpired:
        out = "(timeout)"
    print("=" * 8, label)
    print(out[:1400])
    return out

body = []
body.append(run("ssh-in+marker", "echo SSH-IN-OK; hostname; cat /toolkit/.setup-complete"))
body.append(run("status (real)", "/toolkit/bin/manx status | head -8"))
body.append(run("collect (real)", "/toolkit/bin/manx collect | head -6"))
body.append(run("collect artifact", "ls -la /toolkit/out/ | head -8"))
body.append(run("list-actions", "/toolkit/bin/manx --list-actions"))
body.append(run("audit log tail", "tail -6 /toolkit/out/audit.log"))
body.append(run("gate refusal (img-out)", "/toolkit/bin/manx img-out 2>&1; echo rc=$?"))
body.append(run("setup verb idempotent", "/toolkit/bin/manx setup | head -3"))
body.append(run("audit row count", "wc -l /toolkit/out/audit.log"))
open(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\bench-v003-ssh.txt", "w",
     encoding="utf-8").write("\n".join(body))
print("captured to bench-v003-ssh.txt")
