# Security policy

## Supported versions

| Version | Supported          |
| ------- | ------------------ |
| v0.x.y  | latest tag only    |

## Reporting a vulnerability

GitHub is the reporting channel: [open a private security advisory](https://github.com/MTG-Thomas/manx/security/advisories/new).
Please do NOT open a public issue for unpatched vulnerabilities.

## Residual-trust model for release artifacts

Every artifact produced by this repo must not contain secrets. Membership on the
verification side (what CI checks before a release ships):

- `SHA256SUMS` is produced over every release artifact, and is itself signed in CI
  (`cosign`/`gpg`), so a corrupted download is detectable.
- The `manx-iso` build never bakes into the image:
  - no enrollment tokens, API tokens, auth JSON, or SSH private keys;
  - no per-box secrets (mesh certs are issued at boot, not baked);
  - no password files of any kind.
- CI scans release artifacts (the ISO + toolkit overlay) for secret-shaped strings
  (high-entropy `bfen_…`-style bearer prefixes, PEM private-key blocks, `.ssh/`
  directory markers) and fails the release job if any suspicion surfaces.

## Secure Boot expectations

The `manx-iso` rescue medium inherits Secure Boot from the base distro's
**Microsoft-signed shim + grub + kernel** chain (see `docs/SECURE-BOOT.md`).
Nothing in this repository needs to ship a Microsoft-side key. If a future release
needs a custom-signed kernel module, it will require an explicit, documented MOK
enrollment flow, and the CI chain-of-custody documentation must be updated first.

## Disclosure policy

We follow coordinated disclosure: a confirmed report gets a reproducing-configuration
commit (necessarily sometimes a redacted one), a CVE after we are sure the downstream
deploys we know about are patched, and then a public advisory after the fix ships.
