import subprocess, time, zlib, struct, sys, os

PVE_SSH = ["ssh", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa", "root@172.16.15.128"]

def monitor_cmd(cmd):
    r = subprocess.run(["ssh", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa",
                        "root@172.16.15.128", "echo " + cmd + " | qm monitor 120 2>&1 | tail -1"],
                       capture_output=True, text=True)
    return r.stdout.strip()

def screendump(tag):
    monitor_cmd("screendump /tmp/bench-" + tag + ".ppm")
    subprocess.run(["scp", "-o", "BatchMode=yes", "-i",
                    r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa",
                    "root@172.16.15.128:/tmp/bench-" + tag + ".ppm",
                    r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\bench-" + tag + ".ppm"],
                   capture_output=True, text=True)
    ppm = open(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\bench-" + tag + ".ppm", "rb").read()
    parts = ppm.split(b"\n", 3)
    w, h = map(int, parts[1].split())
    pix = parts[3]
    rows = b""
    for row in range(h):
        rows += b"\x00" + pix[row*w*3:(row+1)*w*3]
    comp = zlib.compress(rows, 6)
    png = b"\x89PNG\r\n\x1a\n" + struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0)    # IHDR
    def chunk(tag, data):
        c = struct.pack(">I", len(data)) + tag + data
        c += struct.pack(">I", zlib.crc32(tag + data) & 0xffffffff)
        return c
    png += chunk(b"IDAT", comp)
    png += chunk(b"IEND", b"")
    out = r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\bench-" + tag + ".png"
    open(out, "wb").write(png)
    return out

CM = {".": "dot", "-": "minus", " ": "spc", "/": "slash", ";": "semicolon",
      "'": "apostrophe", '"': "shift-apostrophe", "\\": "backslash",
      "(": "shift-9", ")": "shift-0", "=": "equal", "_": "shift-minus",
      ">": "shift-dot", "$": "shift-4", "#": "shift-3", "@": "shift-2",
      ":": "shift-semicolon", "+": "shift-equal"}

def type_line(line):
    seq = []
    for ch in line:
        if ch in CM:
            seq.append(CM[ch])
        elif ch.isupper():
            seq.append("shift-" + ch.lower())
        elif ch.islower() or ch.isdigit():
            seq.append(ch)
        else:
            raise RuntimeError("unmapped char %r" % ch)
    for i in range(0, len(seq), 40):
        monitor_cmd("sendkey " + " ".join(seq[i:i+40]))
        time.sleep(0.03)
    monitor_cmd("sendkey ret")
    time.sleep(0.6)

# Step 1: type an unmistakable echo, watch it show up on the next screen
type_line("echo manx-DG7000")
time.sleep(2)
print(screendump("t1"))
