// SPDX-License-Identifier: AGPL-3.0-or-later

package actions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDriverBootstrapUsesPrivateTemporaryHandoff(t *testing.T) {
	path := filepath.Join("..", "..", "manx-iso", "toolkit-overlay", "bin", "bootstrap-drivers.sh")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read bootstrap script: %v", err)
	}
	script := string(content)
	for _, required := range []string{
		`mktemp -d "${TMPDIR:-/tmp}/manx-dsu.XXXXXXXX"`,
		`chmod 700 "$DSU_TMP"`,
		`[ ! -f "$DSU_BOOTSTRAP" ]`,
		`[ -L "$DSU_BOOTSTRAP" ]`,
		`stat -c %u "$DSU_BOOTSTRAP"`,
		`trap cleanup_dsu 0`,
		`runlog sh "$DSU_BOOTSTRAP"`,
	} {
		if !strings.Contains(script, required) {
			t.Errorf("bootstrap script missing security control %q", required)
		}
	}
	if strings.Contains(script, "/tmp/bootstrap.cgi") {
		t.Fatal("bootstrap script still uses the predictable shared path")
	}
}
