import subprocess, time

SSHBASE = ["ssh", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa", "root@172.16.15.128"]
CM = {".": "dot", " ": "spc", "/": "slash", "-": "minus",
      "'": "apostrophe", '"': "shift-apostrophe", "\\": "backslash",
      "(": "shift-9", ")": "shift-0", "=": "equal", "_": "shift-minus",
      ">": "shift-dot", "$": "shift-4", "#": "shift-3", "@": "shift-2",
      ":": "shift-semicolon", "+": "shift-equal", "?": "shift-slash",
      ";": "semicolon", "&": "shift-7", "|": "shift-backslash",
      "<": "shift-comma"}
for c in "abcdefghijklmnopqrstuvwxyz": CM[c] = c
for c in "0123456789": CM[c] = c
for c in "ABCDEFGHIJKLMNOPQRSTUVWXYZ": CM[c] = "shift-" + c.lower()

def sendkey(name):
    subprocess.run(SSHBASE + ["qm", "sendkey", "120", name], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

def screendump_return_png(shotname):
    subprocess.run(SSHBASE + ["printf 'screendump /tmp/'" + shotname + "'.ppm\\n' | qm monitor 120"], capture_output=True)
    subprocess.run(["scp", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa",
                    "root@172.16.15.128:/tmp/" + shotname + ".ppm",
                    r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\bench-" + shotname + ".ppm"], capture_output=True)
    from PIL import Image
    Image.open(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\bench-" + shotname + ".ppm",
               formats=["PPM"]).save(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\bench-" + shotname + ".png")
    print("saved bench-" + shotname + ".png")

def typed_line(line):
    for ch in line:
        sendkey(CM[ch])
        time.sleep(0.10)
    sendkey("ret")
    time.sleep(1.5)

# Solid path: run the entire acceptance suite ON the guest console and see the
# output as image—no ssh needed. This is the *reliable* bench lane.
SUITE = [
    # (label, command); use head to limit console use
    ("motd on console already verified visually"),
    ("host", "echo HOST: $(hostname)"),
    ("status", "/toolkit/bin/manx status"),
    ("detect-hw", "/toolkit/bin/manx detect-hw; echo DETECT-HW-DONE"),
    ("img-in", "/toolkit/bin/manx img-in; echo IMG-IN-DONE"),
    ("audit rows", "wc -l /toolkit/out/audit.log"),
    ("audit log tail", "tail -6 /toolkit/out/audit.log"),
    ("gate without i-know", "/toolkit/bin/manx img-out 2>&1; echo rc=$?"),
    ("gate with i-know", "/toolkit/bin/manx img-out --i-know 2>&1; echo rc=$?"),
    ("reports on disk", "ls -l /toolkit/out/"),
    ("motd file", "head -4 /etc/manx-motd 2>/dev/null || cat /toolkit/MOTD.txt 2>/dev/null | head -4"),
]
print("=== typing suite, one command per line, then screenshot ===")
# note: first run setup script once to make /toolkit real
typed_line("bash /run/archiso/bootmnt/autorun/00-manx-setup.sh")
time.sleep(3)
for cmd in SUITE[1:]:
    typed_line(cmd[1])
    time.sleep(1)
screendump_return_png("b002-suite-1")
