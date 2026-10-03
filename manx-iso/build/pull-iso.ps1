# stream the ISO out of CT 119 via pct exec to a local file (no host path)
ssh -o BatchMode=yes -i C:\Users\ThomasBray\.ssh\proxmox-root-id_rsa root@172.16.15.128 "pct exec 119 -- cat /opt/manx-build/manx-iso-v0.0.1.iso" | Set-Content -Path C:\Users\ThomasBray\src\manx\manx-iso-v0.0.1.iso -Encoding Byte -AsByteStream
# verify locally
Get-Item C:\Users\ThomasBray\src\manx\manx-iso-v0.0.1.iso | Select-Object Length
(Get-FileHash C:\Users\ThomasBray\src\manx\manx-iso-v0.0.1.iso -Algorithm SHA256).Hash
