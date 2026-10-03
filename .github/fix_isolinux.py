# -*- coding: utf-8 -*-
p = r"C:\Users\ThomasBray\src\manx\.github\workflows\iso-linux.yml"
s = open(p, encoding="utf-8").read()
s = s.replace("url=$(cat manx-iso/build/BASE-URL.sha256 | head -1)",
              "url=$(python3 -c \"import sys; print(open('manx-iso/build/BASE-URL.sha256').read().splitlines()[0])\")")
s = s.replace("expected=$(cat manx-iso/build/BASE-URL.sha256 | tail -1)",
              "expected=$(python3 -c \"import sys; print(open('manx-iso/build/BASE-URL.sha256').read().splitlines()[-1])\")")
open(p, "w", encoding="utf-8", newline="\n").write(s)
print("patched")
import subprocess
r = subprocess.run(["grep", "-n", "BASE-URL.sha256", p], capture_output=True, text=True)
print(r.stdout or r.stderr)
