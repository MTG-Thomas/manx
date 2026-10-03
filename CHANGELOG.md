# Changelog

All notable changes to MANX are documented here (Keep a Changelog, semver from v0).

## [Unreleased]

## [0.0.1] - 2026-10-03

### Added
- bootstrap: spec, LICENSE (AGPL-3.0-or-later), SECURITY.md, CONTRIBUTING.md, AGENTS.md,
  CODEOWNERS, CI + govulncheck + dependabot, Makefile with codex-swarm-parity targets.
- CLI skeleton: `manx --list-actions` + `manx status` (one action surface, verb registry).
- TUI skeleton: `manx-tui` (Bubble Tea status/verb view, degrades to `whiptail`).
- ISO overlay: `rescue-menu.sh` + `detect-hw.sh` + `bootstrap-drivers.sh` + `net-up.sh`
  + `mesh-up.sh` + `channel-up.sh` + `collect.sh` + `mount-ro.sh` + `hive-inspect.sh`
  + `hive-edit.sh` + `img-*.sh` placeholder scripts.
- SECURE-BOOT.md (signed-distro chain-of-custody model: Microsoft-signed shim/grub
  chain from the base distro + cosign/gpg-signed release manifests; no MOK enrollment
  required if no custom kernel is introduced; bench-only MOK enrollment only when we
  introduce one).
