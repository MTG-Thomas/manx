# ISO overlay (rescue-medium side of MANX)

Scripts here are the *ISO view* of the same action surface. They are baked into
`manx-iso` (`/toolkit/bin/`), and are the Tier-1 (`whiptail`) view plus the shell
helpers. The Go binary (`manx`, `manx-tui`) is compiled into the ISO at build time
and is the Tier-2 view; both read/write the same state file (`/toolkit/out/`)
and the same audit log (`/toolkit/out/audit.log`).

The authoritative behavior lives in `bifrost-workspace/docs/rescue/MANX-SPEC.md`
§2.3, §8b, and §10.8. Scripts here must not re-implement logic that exists in
`internal/actions/` — if a script grows beyond orchestration, port it.

Baked dependencies (spec §2.1/§8b.2): `whiptail`, `dialog`, `jq`, `rg`, `7z`,
`curl`, `tmux`, `rsync`, `chntpw`, `hivexsh`, `libguestfs-tools`, `ntfs-3g`,
`mdadm`, `lvm2`, `partclone`, `ddrescue`, `smartmontools`, `nvme-cli`.

Nothing here may enqueue a secret. The mesh/OAuth/API key that bootstrap needs is
injected at runtime by the operator or fetched from a writable persist partition.
