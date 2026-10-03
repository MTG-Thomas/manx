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

def screenshot(shotname):
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

typed_line("which sshd; systemctl list-unit-files | grep -i ssh; systemctl start sshd; ss -tln | head -8; ps aux | head -2")
time.sleep(2)
screenshot("sshd-diag")
