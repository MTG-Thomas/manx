import subprocess, time

SSHBASE = ["ssh", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa", "root@172.16.15.128"]
CM = {".": "dot", " ": "spc", "/": "slash", "-": "minus", "b": "b"}
for c in "acdefghijklmnopqrstuvwxyz0123456789": CM[c] = c

def sendkey(name):
    subprocess.run(SSHBASE + ["qm", "sendkey", "120", name], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

def shot(name):
    subprocess.run(SSHBASE + ["printf 'screendump /tmp/v003-" + name + ".ppm\\n' | qm monitor 120"], capture_output=True)
    subprocess.run(["scp", "-q", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa",
                    "root@172.16.15.128:/tmp/v003-" + name + ".ppm",
                    r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\v003-" + name + ".ppm"], capture_output=True)
    from PIL import Image
    Image.open(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\v003-" + name + ".ppm",
               formats=["PPM"]).save(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\v003-" + name + ".png")
    print("saved v003-" + name + ".png")

def type_line(line):
    for ch in line:
        sendkey(CM[ch]); time.sleep(0.09)
    sendkey("ret"); time.sleep(2)

# launch the TUI on the console
type_line("/toolkit/bin/manx-tui")
time.sleep(4)
# run 'status' (cursor starts there): Enter, wait, screenshot with output captured
sendkey("ret"); time.sleep(3)
shot("tui-status")
# navigate down to img-out (index 5): 5x down, enter -> refusal + CLI hint
for _ in range(5):
    sendkey("down"); time.sleep(0.25)
sendkey("ret"); time.sleep(1.5)
shot("tui-refusal")
sendkey("q")
time.sleep(1)
# audit rows: tui view rows must exist now
type_line("tail -4 /toolkit/out/audit.log")
time.sleep(1)
shot("tui-audit")
