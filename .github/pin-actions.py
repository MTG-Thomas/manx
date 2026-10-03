"""Pin all GitHub Action refs to immutable SHAs (codex-swarm security parity).

Fetches the latest commit SHA for the ref used in each workflow `uses:` line
(vX tag converted to latest commit), and rewrites the files in place.
Prints an audit line per pin.
"""
import json, re, subprocess, sys
from pathlib import Path

REPO_ROOT = Path(r"C:\Users\ThomasBray\src\manx")
WF_DIR = REPO_ROOT / ".github" / "workflows"

REFS = {}


def resolve_sha(owner_repo: str, ref: str) -> str:
    key = f"{owner_repo}@{ref}"
    if key not in REFS:
        out = subprocess.check_output(
            ["gh", "api", f"repos/{owner_repo}/commits/{ref}", "--jq", ".sha"],
            text=True)
        REFS[key] = out.strip()
    return REFS[key]


def main():
    changed = 0
    for f in WF_DIR.glob("*.yml"):
        text = original = f.read_text(encoding="utf-8")
        def sub(m):
            uses, ref = m.group(1), m.group(2)
            if "@" in uses:  # already sha-pinned (uses something@sha / composite? not here)
                return m.group(0)
            repo = uses.split("@")[0]
            try:
                sha = resolve_sha(repo, ref)
            except Exception as e:
                print(f"skip {repo}@{ref}: {e}")
                return m.group(0)
            print(f"pin {uses}@{ref} -> {repo}@{sha[:12]}")
            return f"{repo}@{sha}"
        text = re.sub(rf"uses:\s*([A-Za-z0-9_.\-/]+)@(\S+)", sub, text)
        if text != original:
            f.write_text(text, encoding="utf-8")
            changed += 1
    print(f"pinned workflows: {changed}")


if __name__ == "__main__":
    sys.exit(main())
