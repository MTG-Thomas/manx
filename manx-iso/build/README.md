# manx-iso/build: remaster inputs

This directory is the *build contract* for the Linux ISO remaster
(the inputs the CI container consumes; see `docs/BUILD-PIPELINE.md`).

- `BASE-URL.sha256` — two lines: (1) the pinned SystemRescue ISO URL, (2) its
  sha256, matching the sha the cache key hashes. CI refuses to remaster from an
  upstream whose hash drifts.
- `build-remaster.sh` — the one script that turns (base ISO + `toolkit-overlay/`
  + `dist/bin/*`) into `out.iso`. Containerized in the CI lane; no human steps.

Long-term these merge into a `Dockerfile` (digest-pinned) behind a `deps.lock`
(apt version pinning) — until then the honest claim is "reproducible inputs,
verified manifest" (see the pipeline doc §Reproducibility).
