#!/usr/bin/env bash
set -euo pipefail
# manx bench VM (MTG-owned; never customer equipment) on the ISO-builder host
qm create 120 \
  --name manx-bench-vm \
  --memory 4096 --cores 2 \
  --machine q35 --bios ovmf \
  --efidisk0 local-lvm:1 \
  --scsi0 local-lvm:10,iothread=1 \
  --net0 virtio,bridge=vmbr0,firewall=0 \
  --ide2 local:iso/manx-iso-v0.0.1.iso,media=cdrom \
  --boot order='ide2;scsi0' --onboot 1 \
  --ostype l26 \
  --tags manx \
  --description "MANX rescue-ISO acceptance bench (MTG-owned). ISO in cdrom (local iso storage). Boot UEFI. Never use for customer data/equipment."
qm set 120 --serial0 socket
qm start 120 2>&1 | tail -1 || true
sleep 3
qm status 120
