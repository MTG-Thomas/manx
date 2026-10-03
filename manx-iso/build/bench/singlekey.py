import subprocess, time

def main():
    # type one 'q', then screendump and pull the ppm
    def monitor(cmd):
        r = subprocess.run(["ssh", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa",
                            "root@172.16.15.128", "echo '" + cmd + "' | qm monitor 120 2>&1 | tail -1"],
                           capture_output=True, text=True)
        print("monitor:", r.stdout.strip())
        return r.stdout.strip()
    monitor("sendkey q")
    time.sleep(1.2)
    monitor("screendump /tmp/b-q.ppm")
    subprocess.run(["scp", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa",
                    "root@172.16.15.128:/tmp/b-q.ppm",
                    r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\b-q.ppm"],
                   capture_output=True, text=True)
    from PIL import Image
    Image.open(r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\b-q.ppm", formats=["PPM"]).save(
        r"C:\Users\ThomasBray\src\manx\manx-iso\build\bench\b-q.png")
    print("converted; ready to view")

main()
