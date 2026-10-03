import subprocess, textwrap

SSH = ["ssh", "-o", "BatchMode=yes", "-o", "StrictHostKeyChecking=accept-new",
       "-i", r"C:\Users\ThomasBray\.ssh\id_ed25519", "root@172.16.15.171"]

CMDS = [
    ("hostname", "hostname"),
    ("kernel", "uname -r"),
    ("net (ip -o -4 addr)", "ip -o -4 addr show | awk '{print $2, $4}'"),
    ("manx --list-actions", "/toolkit/bin/manx --list-actions"),
    ("manx status", "/toolkit/bin/manx status"),
    ("manx detect-hw", "/toolkit/bin/manx detect-hw"),
    ("detect-hw.sh report", "bash /toolkit/manx-iso-overlay/detect-hw.sh 2>&1 | head -12"),
    ("manx img-in (read-only)", "/toolkit/bin/manx img-in 2>&1 | head -6"),
    ("audit log rows", "ls -l /toolkit/out/audit.log; wc -l /toolkit/out/audit.log; tail -6 /toolkit/out/audit.log"),
    ("toolkit layout", "ls -la /toolkit/bin /toolkit/manx-iso-overlay; ls -l /manx/MANIFEST.txt 2>/dev/null; head -12 /manx/MANIFEST.txt 2>/dev/null"),
    ("destructive verb gate (refuses without --i-know)", "/toolkit/bin/manx img-out; echo \"exit-code=$?\""),
    ("rescue-menu non-tty (should bail, not hang)", "timeout 5 bash /toolkit/manx-iso-overlay/rescue-menu.sh < /dev/null > /dev/null 2>&1 ; echo \"menu-non-tty-exit=$?\""),
    ("sshd state", "ss -tln | grep 22 | head -2"),
    ("nomnially (nested list test)", "/toolkit/bin/manx --list-actions 2>&1 | tail -6"),
]
out = {}
for k, cmd in CMDS:
    r = subprocess.run(SSH + [cmd], capture_output=True, text=True, timeout=180)
    out[k] = (r.stdout or r.stderr).strip()

for k, v in out.items():
    print("=" * 8, k)
    print(v)
