# stream a file to the pve host then into CT 119 (PS-friendly: Get-Content | ssh)
param([string]$Local, [string]$Remote)
Get-Content -Path $Local -Raw -Encoding Byte | ssh -o BatchMode=yes -i C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa root@172.16.15.128 "cat > /tmp/$Remote"
