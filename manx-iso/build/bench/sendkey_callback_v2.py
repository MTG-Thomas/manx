import subprocess, time

PVE = "root@172.16.15.128"
VMID = "120"
KEY = ["-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa"]
CHAR_MAP = {
    ".": "dot", "-": "minus", " ": "spc", "/": "slash", ";": "semicolon",
    "'": "apostrophe", '"': "shift-apostrophe", "\\": "backslash",
    "(": "shift-9", ")": "shift-0", "=": "equal", "_": "shift-minus",
    ">": "shift-dot", "$": "shift-4", "#": "shift-3", "@": "shift-2",
    ":": "shift-semicolon", "+": "shift-equal", "?": "shift-slash",
}
for c in "abcdefghijklmnopqrstuvwxyz": CHAR_MAP[c] = c
for c in "0123456789": CHAR_MAP[c] = c
for c in "ABCDEFGHIJKLMNOPQRSTUVWXYZ": CHAR_MAP[c] = "shift-" + c.lower()

def sendkey(seq):
    subprocess.run("ssh -o BatchMode=yes -i C:/Users/ThomasBray/.ssh/proxmox-root-id_rsa "
                   "root@172.16.15.128 qm sendkey 120 " + seq,
                   stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, shell=False)

def type_line(line):
    seq = []
    for ch in line:
        k = CHAR_MAP.get(ch)
        if k is None:
            raise RuntimeError("no key for %r in %r" % (ch, line))
        seq.append(k)
    for i in range(0, len(seq), 48):
        sendkey(" ".join(seq[i:i + 48]))
        time.sleep(0.02)
    sendkey("ret")
    time.sleep(0.4)

# callback URL: no quotes needed — shell expands $(...) inside double free-form
cmd = "curl -s http://172.16.15.128:4097/callback-$(hostname)"
for cmd_line in [cmd,
                 "ip -o -4 addr show eth0 > /dev/ttyS0"]:
    type_line(cmd_line)
    print("typed:", cmd_line[:60])
time.sleep(1)
print("read /tmp/listen.log on pve for the callback")
