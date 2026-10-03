// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// FirstBootMarker is the file the autorun banner leaves behind to prove the
// toolkit setup ran successfully on this boot. Its presence lets `manx status`
// and the bench tell "setup succeeded" from "setup was never invoked".
const FirstBootMarker = "/toolkit/.setup-complete"

// MarkFirstBoot writes the first-boot marker (idempotent; best-effort).
// Returns an error string for the caller to surface, never panics.
func MarkFirstBoot() error {
	if err := os.MkdirAll("/toolkit", 0o755); err != nil {
		return fmt.Errorf("cannot create /toolkit: %v", err)
	}
	stamp := fmt.Sprintf("setup ok at %s (manx status for state)\n", stampNow())
	// duplicate into both the ISO-side tracker and the shell history
	if err := os.WriteFile(FirstBootMarker, []byte(stamp), 0o644); err != nil {
		return fmt.Errorf("marker write failed: %v", err)
	}
	return nil
}

// FirstBootDone reports whether the marker exists (and is readable).
func FirstBootDone() bool {
	st, err := os.Stat(FirstBootMarker)
	return err == nil && st.Size() > 0
}

// RunIfFirstBoot executes the fn only if the first-boot marker is absent.
// Used by runSetup to make the autorun path idempotent.
func RunFirstBootOnce(fn func() error) error {
	if FirstBootDone() {
		return nil
	}
	if err := fn(); err != nil {
		return err
	}
	return MarkFirstBoot()
}

// runSetup is the in-Go equivalent of the shell autorun banner: verifies the
// toolkit payload on the boot media and stages it into /toolkit, setting PATH
// and the motd. It is the verb `manx setup`.
func runSetup(args []string) error {
	if !onLinux() {
		return fmt.Errorf("setup is ISO-native (running on %s); boot the rescue media", runtime.GOOS)
	}
	bootMnt := "/run/archiso/bootmnt"
	src := bootMnt + "/manx/toolkit"
	dst := "/toolkit"
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("toolkit payload not found on boot media: %s", src)
	}
	// shell-out-style staging (cp -a semantics preserved by per-file walk)
	for _, pair := range [][2]string{
		{filepath.Join(src, "bin"), filepath.Join(dst, "bin")},
		{filepath.Join(src, "manx-iso-overlay"), filepath.Join(dst, "manx-iso-overlay")},
		{filepath.Join(src, "MOTD.txt"), filepath.Join(dst, "MOTD.txt")},
	} {
		if err := copyAll(pair[0], pair[1]); err != nil {
			return fmt.Errorf("staging %s: %v", pair[0], err)
		}
	}
	// motd persistence (mirrors the bash autorun semantics)
	if b, err := os.ReadFile(filepath.Join(dst, "MOTD.txt")); err == nil {
		_ = os.WriteFile("/etc/manx-motd", b, 0o644)
		bashrc := "/root/.bashrc"
		if b2, err2 := os.ReadFile(bashrc); err2 == nil && !strings.Contains(string(b2), "manx-motd") {
			f, err2 := os.OpenFile(bashrc, os.O_APPEND|os.O_WRONLY, 0o644)
			if err2 == nil {
				_, _ = f.WriteString("test -f /etc/manx-motd && cat /etc/manx-motd\n")
				_ = f.Close()
			}
		}
	}
	if err := MarkFirstBoot(); err != nil {
		return err
	}
	fmt.Println("MANX toolkit ready at /toolkit (setup marker written; motd active for new shells)")
	return nil
}

func stampNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// collect walks the well-known log/evidence locations and tars what exists
// into <out>/collect-<ts>.tar.gz, printing everything it found. This is the
// verb an operator runs FIRST (spec: collect > diagnose > repair).
func runCollect(args []string) error {
	if !onLinux() {
		return fmt.Errorf("collect is ISO-native (running on %s); boot the rescue media", runtime.GOOS)
	}
	outDir := toolkitOut()
	targets := []string{
		"/var/log",
		"/etc/fstab",
		"/proc/mounts",
		"/etc/hosts",
	}
	// SystemRescue keeps Windows-side evidence here when the OS volume is mounted; if
	// the mounted OS root is visible we also harvest its CBS/dism/winevt evidence.
	osVol := findOsVolume()
	if osVol != "" {
		osTargets := []string{
			filepath.Join(osVol, "Windows", "Logs", "CBS"),
			filepath.Join(osVol, "Windows", "Logs", "DISM"),
			filepath.Join(osVol, "Windows", "System32", "winevt", "Logs"),
			filepath.Join(osVol, "Windows", "Minidump"),
			filepath.Join(osVol, "Windows", "Memory.DMP"),
		}
		targets = append(targets, osTargets...)
	}
	var found []string
	for _, t := range targets {
		if _, err := os.Stat(t); err == nil {
			found = append(found, t)
		}
	}
	if len(found) == 0 {
		return fmt.Errorf("no collectable targets found (mount the OS volume first; try `manx status`)")
	}
	archive := filepath.Join(outDir, fmt.Sprintf("collect-%s.tar.gz", stampNow()))
	tarArgs := append([]string{"-czf", archive}, found...)
	cmd := ExecCommand("tar", tarArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tar collect failed: %v (%s)", err, strings.TrimSpace(string(out)))
	}
	fmt.Printf("collected %d targets into %s\n", len(found), archive)
	fmt.Println(strings.Join(found, "\n  "))
	return nil
}

func findOsVolume() string {
	// scan /proc/mounts for an NTFS volume hosting or Windows/
	b, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return ""
	}
	for _, l := range strings.Split(string(b), "\n") {
		f := strings.Fields(l)
		if len(f) < 3 {
			continue
		}
		mnt := strings.ReplaceAll(f[1], "\\040", " ")
		if strings.HasSuffix(strings.ToLower(f[2]), "ntfs") || strings.Contains(strings.ToLower(l), "ntfs") {
			if _, err := os.Stat(filepath.Join(mnt, "Windows", "System32", "config", "SOFTWARE")); err == nil {
				return mnt
			}
		}
	}
	return ""
}
