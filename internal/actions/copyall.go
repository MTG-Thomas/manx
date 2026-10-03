// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

import (
	"io"
	"os"
	"path/filepath"
)

// copyAll is a dependency-free recursive copy (cp -a for our limited case:
// no xattrs, no devices, no symlinks inside the toolkit payload).
func copyAll(src, dst string) error {
	running := runningExe()
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		// skip overwriting the binary we are currently running (`manx setup`
		// executed from /toolkit/bin is a re-stage onto itself; "text file busy"
		// would abort the whole staging otherwise)
		if running != "" && sameFile(target, running) {
			return nil // skip overwriting the binary we are running from
		}
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	})
}

// runningExe returns the absolute path of the currently running executable
// (empty string when unavailable).
func runningExe() string {
	p, err := os.Executable()
	if err != nil {
		return ""
	}
	if r, err := filepath.EvalSymlinks(p); err == nil {
		return r
	}
	return p
}

// sameFile compares cleaned absolute paths lexically (no syscall needed for
// the toolkit's flat layout; different devices still compare equal by path).
func sameFile(a, b string) bool {
	aa, _ := filepath.Abs(a)
	bb, _ := filepath.Abs(b)
	return filepath.Clean(aa) == filepath.Clean(bb)
}
