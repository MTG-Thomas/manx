// Runbook parity lock: every doc in docs/runbooks must declare flat
// front-matter (---title/triggers/first-verbs/spec---) and whatever it lists
// under first-verbs must be verbs this registry actually has. Runbooks name
// capabilities; the action surface owns them. One drifts, CI fails.
// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const runbooksDir = "../../docs/runbooks"

func parseFrontMatter(t *testing.T, path string) map[string][]string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	s := string(b)
	if !strings.HasPrefix(s, "---\n") {
		t.Fatalf("%s: missing front-matter block (must start with ---\\n)", path)
	}
	end := strings.Index(s[3:], "\n---")
	if end < 0 {
		t.Fatalf("%s: front-matter block not terminated", path)
	}
	out := map[string][]string{}
	for _, l := range strings.Split(s[4:end+3], "\n") {
		l = strings.TrimSpace(l)
		if l == "" || !strings.Contains(l, ":") {
			continue
		}
		parts := strings.SplitN(l, ":", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		for _, item := range strings.Split(val, ",") {
			item = strings.TrimSpace(item)
			if item != "" {
				out[key] = append(out[key], item)
			}
		}
	}
	return out
}

func TestRunbookFrontMatter(t *testing.T) {
	entries, err := os.ReadDir(runbooksDir)
	if err != nil {
		t.Fatalf("walk %s: %v", runbooksDir, err)
	}
	found := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		found++
		path := filepath.Join(runbooksDir, e.Name())
		fm := parseFrontMatter(t, path)
		for _, key := range []string{"title", "triggers", "first-verbs", "spec"} {
			if len(fm[key]) == 0 {
				t.Errorf("%s: front-matter key %q missing/empty", e.Name(), key)
			}
		}
		reg := Default()
		for _, v := range fm["first-verbs"] {
			if _, ok := reg.Get(v); !ok {
				t.Errorf("%s: first-verb %q is not a registry verb (runbooks must only lead with real capabilities)", e.Name(), v)
			}
		}
	}
	if found < 1 {
		t.Fatalf("no runbooks found under %s", runbooksDir)
	}
}
