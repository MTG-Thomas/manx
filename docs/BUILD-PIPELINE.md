# MANX build pipeline (ISO + release)

ISOs are heavy artifacts, not special ones. The pipeline is tag-driven CI with an
explicit human gate between "artifact exists" and "artifact is a release", plus an
opt-in bench lane that runs the §8b.6 acceptance set on real MTG infra.

## Jobs and triggers

| Job | Trigger | Runner | Output |
|---|---|---|---|
| `go` | every push/PR | `ubuntu-latest` | `bin/manx`, `bin/manx-tui` (workflow artifact) |
| `iso-linux` | `workflow_dispatch` (bench candidate), `push: tags v*` | `ubuntu-latest`, **container: pinned Debian digest** | `manx-iso-<ver>.iso` + `SHA256SUMS.parts` (workflow artifact) |
| `iso-winpe` (optional lane) | tag or manual; `wimlib` builds the `boot.wim` **from Linux** (spec §11.8 — no Windows runner) | `ubuntu-latest` | `manx-winpe-<ver>.img` |
| `dry-run-release` | tag | `ubuntu-latest` | secrets-scan report + **`SHA256SUMS` + cosign (keyless OIDC)** |
| `publish-github-release` | tag, **after dry-run greens** | `ubuntu-latest` | GitHub Release with the artifacts |
| **`bench-e2e`** | **manual only** (`workflow_dispatch`, inputs = build run id) | `self-hosted, manx-bench` (our infra) | bench booted/verified log; result recorded but never auto-promotes |

## The three hard rules

1. **Public repo ⇒ self-hosted runners never see PR code.** The bench lane can only
   be triggered by `workflow_dispatch` on `main`/tags by an org member, and the
   `bench-e2e` workflow is gated behind the `bench` environment (required reviewer).
   No `pull_request` trigger ever reaches the self-hosted runner. If a fork wants
   bench coverage, they must run their own bench per the spec.
2. **Never silently re-release.** Everything is immutable per tag: assets upload once
   from the `publish` job; a defect fix means a new tag (`v0.0.1-a.1` style), never a
   re-upload under the same name (the KB5123099 corrupted-download lesson).
3. **Reproducibility: attempted, honestly caveated.** The builder container is pinned
   by digest; `SOURCE_DATE_EPOCH` is fixed per tag; paths are relative and sorted;
   build output prints a diff-relevant manifest. Squashfs/xorriso block-order may still
   vary between rebuilds, so the claim in the docs is "deterministic *inputs*, stable
   layout, verified manifest" — not "byte-identical ISO" (that claim mathematically
   requires full reproducible-build infrastructure this repo does not yet staff).

## Cache and cost

- Base SystemRescue ISO (~800–900 MB) fetched once per builder version; cached with
  `actions/cache` keyed on the pinned upstream URL + its sha256.
- Docker layers: we do not build a Docker *image* for the ISO; we run a Debian
  container in the job (`container:`) and `apt-get install` the remaster suite
  (`xorriso`, `live-build`, `wimlib`, `squashfs-tools`, `dosfstools`, `mtools`,
  `erofs-utils`). Installed-dep manifest (`dpkg -l`) ships in the artifact so we can
  prove what built it. Digest pinning of the container is a first-release gate; until
  then the version is substantively pinned by the `debian:bookworm-slim` tag plus
  the recorded `dpkg -l` snapshot.
- ISO build should land in the 20–40 min band on a shared runner. Cost is acceptable
  on tags/dispatch only.

## Publishing (only `v*` tags)

`dry-run-release` runs the §12.2 release gates (secrets scan, manifest generation),
cosigns `SHA256SUMS` **keyless via GitHub OIDC** (no long-lived key to leak; the
verify experience is documented in `docs/SECURE-BOOT.md` §4), and only after that
does `publish-github-release` upload assets under the tag. Bench coverage is
*expected* but is not a technical gate inside this pipeline — it is the human's
responsibility per §10.7/§8b.6 acceptance sets, and the bench lane is opt-in rather
than binding because "the bench Dell was free today" is not a promise we can encode.

## Bench lane on our infra

`bench-e2e.yml` is the opt-in path for e2e testing on **MTG-owned** hardware:

- **Rule zero — never customer hardware.** Bench testing is performed on
  equipment MTG owns and operates, on networks MTG controls. A customer's
  server, customer LAN, or anything at a customer site is out of bounds for
  this lane, regardless of convenience: a public repo's opt-in self-hosted
  job must never place agent code on someone else's compute.
- **Runner identity is org-scoped.** The `manx-bench` runner is registered from
  a MTG-owned device/VM (e.g. a bench box on the office LAN or a lab VM), with
  a label `self-hosted, manx-bench`, and a runner registration scoped to this
  repository only (`actions:read` on MTG-Thomas/manx for sibling artifacts).
- All customer-confidential identifiers (specific sites, IPs, hostnames,
  accounts) never appear in this repository — they live in the runner's local
  profile, and the acceptance checksums they produce are generic artifacts.
- The runner is audited per §10.5 security-model row: what a bench run
  executed, what disks it touched, and whether the boot result was
  1. bootable, 2. net-up, 3. detect-hw clean — written back as a bench
  summary artifact (no customer data, no secrets).


- **Trigger**: `workflow_dispatch` with the build run id (`inputs.run_id`) of the
  artifact to test, chosen by the operator. Dispatching is available only to
  maintainers, and the environment `bench` requires a reviewer confirmation before
  the runner starts.
- **Runner**: repo-scoped `manx-bench` runner, labeled `self-hosted` + `bench`
  (a Dell + KVM + a USB stick + a test storage controller on our bench).
- **Egress**: the bench only needs the LAN and, if disk imaging is exercised, our
  SMB/Nebula endpoints — no arbitrary internet. If the runner's machine also
  holds credentials, they must live in the build's scan-excluded profile, and the
  audit contract from the spec applies to anything the runner executes (the runner
  writes its output to `/toolkit/out/` on the target drive and posts artifacts).
- **Post-run behavior**: the runner posts the `§8b.6` outcome into the workflow
  summary; a failing bench never deletes the artifact, it annotates the dispatch
  run and closes the loop for the operator; the publish is untouched (bench is
  advisory, publishing is the human's action).
