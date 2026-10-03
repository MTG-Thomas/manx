import subprocess, hashlib, sys, os
dest = r"C:\Users\ThomasBray\src\manx\manx-iso-v0.0.1.iso"
h = hashlib.sha256()
n = 0
proc = subprocess.Popen(
    ["ssh", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa",
     "root@172.16.15.128", "pct exec 119 -- tar -cf - -C /opt/manx-build manx-iso-v0.0.1.iso"],
    stdout=subprocess.PIPE, bufsize=1024*1024)
with open(dest, "wb") as out:
    while True:
        chunk = proc.stdout.read(4 * 1024 * 1024)
        if not chunk:
            break
        out.write(chunk)
        h.update(chunk)
        n += len(chunk)
rc = proc.wait()
print("bytes:", n, "rc:", rc)
print("sha256:", h.hexdigest())
