import subprocess, time

"""v0.0.2 build: relay the SSH acceptance through the pve host.
The laptop cannot reach the guest over 172.16.15.x (its own adapter doesn't
sit in that DHCP range consistently); but pve-t340 host CAN reach the guest.
Host-side helper gets the same acceptance suite with the proxmox key.
"""
def main():
    out = subprocess.run(
        ["ssh", "-o", "BatchMode=yes", "-i", r"C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa",
         "root@172.16.15.128",
         "ssh -o BatchMode=yes -o ConnectTimeout=10 -o StrictHostKeyChecking=accept-new "
         "root@172.16.15.171 "
         "'hostname; which manx; ls /toolkit' 2>&1 | head -10"],
        capture_output=True, text=True, timeout=90)
    print(out.stdout.strip())
    print("stderr:", out.stderr.strip())

main()
