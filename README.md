# MANX (rescue toolkit)

> One action surface, three views: `whiptail` menu, `manx <action>` CLI, agent harness over the
> wire. Sponsored by our ops org for offline Windows surgery on real DCs at 2 a.m.

## What this is

`manx` is a *rescue* toolkit - the `bifrost-workspace`-side companion for the Fleetwide
spec shipped at [`docs/rescue/MANX-SPEC.md`](https://github.com/MTG-Thomas/bifrost-workspace/blob/main/docs/rescue/MANX-SPEC.md).
It gives every rescue capability (`collect`, `img-in/out/verify`, `mount-ro`, `hive-inspect/edit`,
`bootstrap-drivers`, `net-up`, `mesh-up`, `channel-up`, `bringup-windows-vm`, `snapshot/rollback`, `status`)
ONE implementation, and a single-name shell surface
(`manx <action>`), a `whiptail`-shaped menu wrapper (`rescue-menu.sh` on the ISO side), and a
Bubble Tea TUI binary (`manx-tui`) as the target interactive interface.

## What it is NOT

- Not a Tailscale replacement or a substitute. The mesh default is
  **Defined Networking + Nebula** with a `tailcat` (static userspace WireGuard over DERP)
  fallback; see the spec.
- Not a driver-carrying ISO. The container image is *self-bootstrapping*:
  it detects unknown hardware at boot (`detect-hw.sh`), then fetches drivers at runtime
  (`bootstrap-drivers.sh` via distro packages / Dell DSU), falling back to PXE.
- Not a portable registry of "known-good" backups. It is a *rescue* tool; actual AD/data
  restore remains Cove/Managed-Nebula-managed and lives outside this repo.

[![ci](https://github.com/MTG-Thomas/manx/actions/workflows/ci.yml/badge.svg)](https://github.com/MTG-Thomas/manx/actions/workflows/ci.yml)

## Layout (codex-swarm parity)

```
cmd/manx/         CLI entry: `manx <action>` (the action surface)
cmd/manx-tui/     Bubble Tea TUI (target UI; degrades to `whiptail`, or plain `dialog`)
internal/         actions (implementation), audit (rows per verb), state (shared w/ menu)
manx-iso/          the ISO overlay: `bin/` scripts + `rescue-menu.sh`
docs/              SECURE-BOOT.md, BENCH-CHECKLIST.md, etc.
```

## Build

```sh
make check  # gofmt + go vet + go test
make build  # compile bin/manx and bin/manx-tui (GOOS/GOARCH aware)
```

Targets are `Makefile`-driven; CI runs `make check`. Bench/vet/go-scan gates are the same
shape as [`MTG-Thomas/codex-swarm`](https://github.com/MTG-Thomas/codex-swarm).

## License

**AGPL-3.0-or-later** (see [`LICENSE`](LICENSE)). Any agent harness or rescue node reachable
over a network falls under AGPL's network-use clause; that is why it is not MIT.
See also [`SECURITY.md`](SECURITY.md) before packaging a release artifact.

---

Derived from the 2026-10-02 WARRENDC incident response, and normalized after the
`codex-swarm` repo-maturity model.
