# MANX repo maturity â€” implemented state (2026-10-03)

| Maturity piece | Where | State |
|---|---|---|
| Public repo `MTG-Thomas/manx` | github | done, AGPL-3.0-or-later |
| codex-swarm Makefile parity (fmt-check/vet/test/build) | `Makefile` | done (all targets green in CI) |
| CI: `go` + `ci` + `govulncheck` | `.github/workflows/{go,c,govulncheck}.yml` | done, pinned SHAs |
| CodeQL (security-extended, weekly) | `.github/workflows/codeql.yml` | done, pinned SHAs |
| Security parity: govulncheck | workflow | done |
| Actions SHA-pinning tool + applied | `.github/pin-actions.py` | done |
| Branch protection `main`: status checks go/govulncheck/CodeQL required on PRs, force-push/deletions blocked; `enforce_admins` intentionally off (solo maintainer) and PR-review requirement intentionally absent (no second reviewer; documented tradeoff) | branch ruleset | done |
| Release pipeline (tag -> dry-run -> publish) | `dry-run-release.yml` + `publish-github-release.yml` | done, two modes: CI lane & manual lane |
| ISO builder CI lane | `iso-linux.yml` | wired (containerized Debian; needs 1st bench confirmation to tag) |
| Bench lane (self-hosted runner) | `bench-e2e.yml` | wired, offline until an MTG bench runner is registered |
| Release v0.0.1 (bootstrap) | `v0.0.1` release | done, prerelease, assets verified, dry-run reports manual lane and exits 0 |
| Build provenance doc | `manx-iso/build/BUILD-PROVENANCE.md` | done |

Note: dry-run-release's mode detection now treats a manually-published tag as valid
(release digest verified server-side), so bootstrap tags no longer turn the release
pipeline red. This is intentional parity with how codex-swarm treated the private
repo (where release lanes were manual); we keep the CI lane as *the* path for all
future releases after v0.0.1.
