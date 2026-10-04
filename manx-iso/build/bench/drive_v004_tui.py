import subprocess, time

SSHBASE = ["ssh", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa", "root@172.16.15.128"]
CM = {".": "dot", " ": "spc", "/": "slash", "-": "minus"}
for c in "abcdefghijklmnopqrstuvwxyz0123456789": CM[c] = c

def sendkey(name):
    subprocess.run(SSHBASE + ["qm", "sendkey", "120", name], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

def shot(name):
    subprocess.run(SSHBASE + ["printf 'screendump /tmp/v004-" + name + ".ppm\\n' | qm monitor 120"], capture_output=True)
    subprocess.run(["scp", "-q", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa",
                    "root@172.16.15.128:/tmp/v004-" + name + ".ppm",
                    r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\v004-" + name + ".ppm"], capture_output=True)
    from PIL import Image
    Image.open(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\v004-" + name + ".ppm",
               formats=["PPM"]).save(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\v004-" + name + ".png")
    print("saved v004-" + name + ".png")

def type_line(line):
    for ch in line:
        sendkey(CM[ch]); time.sleep(0.09)
    sendkey("ret"); time.sleep(2)

# fresh term so old altscreen remnants don't confuse: run tui
type_line("reset")
time.sleep(1)
type_line("/toolkit/bin/manx-tui")
time.sleep(2)
# let the 2s strip tick land once + truth rotate a few words in
time.sleep(4)
shot("tui-head")
# run status
sendkey("ret"); time.sleep(3)
shot("tui-status")
# destructive refusal view: navigate 5 down to img-out
for _ in range(5):
    sendkey("down"); time.sleep(0.2)
sendkey("ret"); time.sleep(2)
shot("tui-refusal")
sendkey("q"); time.sleep(1)
shot("tui-after")
