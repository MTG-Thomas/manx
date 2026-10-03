# Contributing to MANX

## Purpose of the repo (read once)

MANX implements [`docs/rescue/MANX-SPEC.md`](https://github.com/MTG-Thomas/bifrost-workspace/blob/main/docs/rescue/MANX-SPEC.md)
(SPEC REPO IS THE AUTHORITY; this repo implements and links back). Design discussions and
spec changes belong in that repo as PRs (they are the lead / governance artifacts), and
then an implementation PR follows in this repo. Neither repo duplicates design docs.

## Contribution flow

1. One branch per action; PR into `main`.
2. `make check` (gofmt + go vet + **go test**) must be clean; `govulncheck` is CI-enforced
   in addition to `make check`.
3. Any *new* tool/verb: add it to the par-tier audit parity contract in the spec:
   - the verb, its `Makefile`-level CLI name, and its audit-row schema must be a
     spec-level decision (so the par-tier contract is not broken silently);
   - an audit row is written for every destructive action;
   - `--i-know` confirmation gates are present in *every* view (menu/CLI/TUI).
4. Open-source code hygiene (codex-swarm parity): docstrings with what + when +
   failure mode; no `panic` on the code paths below the CLI; state extends but never
   breaks in the ISO side.
5. Bench-test the acceptance gate for your PR by running the spec's
   `BENCH-CHECKLIST.md` set (§8b.6 / §10.7 / §10.8) if you touch the ISO layer.

## Naming

The spec picked "MANX" for a few reasons (see the spec for the full list). The rename
to Bilby/Foxhole is a one-commit change if adopted later; the design is name-agnostic.
