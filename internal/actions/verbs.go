// SPDX-License-Identifier: AGPL-3.0-or-later
package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ExecCommand is the process runner used by the Linux-side verbs. It is a
// package variable purely so tests can inject canned process output;
// production code must not replace it.
var ExecCommand = exec.Command

// DefaultAuditPath resolves the shared append-only audit log location:
// 1. $MANX_AUDIT if set,
// 2. /toolkit/out/audit.log when that directory exists (ISO/systemrescue),
// 3. ./out/audit.log otherwise (developer benches).
func DefaultAuditPath() string {
	if p := os.Getenv("MANX_AUDIT"); p != "" {
		return p
	}
	if st, err := os.Stat("/toolkit/out"); err == nil && st.IsDir() {
		return "/toolkit/out/audit.log"
	}
	_ = os.MkdirAll("out", 0o755)
	return filepath.Join("out", "audit.log")
}

// toolkitOut resolves the exam-report output dir (same fallback rules).
func toolkitOut() string {
	if p := os.Getenv("MANX_OUT"); p != "" {
		return p
	}
	if st, err := os.Stat("/toolkit/out"); err == nil && st.IsDir() {
		return "/toolkit/out"
	}
	_ = os.MkdirAll("out", 0o755)
	return "out"
}

func onLinux() bool { return runtime.GOOS == "linux" }

func binExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// unclaimedLines turns `lspci -nnk` output into the report lines the
// spec's detect-hw.sh contract produces (device line + "-> UNCLAIMED | driver").
// Pure function; tested with canned lspci output.
func unclaimedLines(lspciOut string) []string {
	var lines []string
	id := ""
	drv := ""
	flush := func() {
		if id != "" {
			if drv == "" {
				lines = append(lines, id+" -> UNCLAIMED")
			} else {
				lines = append(lines, id+" -> "+drv)
			}
		}
		drv = "" // drivers belong to one device; never leak to the next flush
	}
	for _, l := range strings.Split(lspciOut, "\n") {
		trimmed := strings.TrimSpace(l)
		switch {
		case strings.HasPrefix(trimmed, "Kernel driver in use:"):
			drv = strings.TrimSpace(strings.TrimPrefix(trimmed, "Kernel driver in use:"))
		case len(trimmed) > 0 && trimmed[0] >= '0' && trimmed[0] <= '9' && strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "Kernel"):
			// device line: leading hex bus id like "00:00.0" (with tabs after)
			flush()
			id = trimmed
		}
	}
	flush()
	// keep only UNCLAIMED for the report (mirrors detect-hw.sh)
	var out []string
	for _, l := range lines {
		if strings.HasSuffix(l, "UNCLAIMED") {
			out = append(out, l)
		}
	}
	return out
}

// composeStatus is the `manx status` renderer. Pure function; tested.
func composeStatus(host, kernel, netLine, uptime string, auditRows int, toolkit bool) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("MANX status\n"))
	b.WriteString(fmt.Sprintf("  host:      %s\n", host))
	b.WriteString(fmt.Sprintf("  kernel:    %s\n", kernel))
	b.WriteString("  toolkit:   ")
	if toolkit {
		b.WriteString("present (/toolkit)\n")
	} else {
		b.WriteString("absent (run from the ISO, or set MANX_OUT)\n")
	}
	b.WriteString(fmt.Sprintf("  net:       %s\n", netLine))
	b.WriteString(fmt.Sprintf("  uptime:    %s\n", uptime))
	b.WriteString(fmt.Sprintf("  audit rows: %d\n", auditRows))
	return b.String()
}

// runStatus implements `manx status` with real data (Linux-tested; degrades on
// other OSes rather than lying about a rescue environment).
func runStatus(args []string) error {
	s := StatusSnapshot()
	fmt.Print(composeStatus(s.Host, s.Kernel, s.Net, s.Uptime, s.AuditRows, s.Toolkit))
	return nil
}

// StatusSnapshot gathers the observed values the `status` verb renders.
// The verb prints them and exits; the TUI's live strip re-reads and re-renders
// them (display refresh is a render, not an action — no audit row).
// Pure data gathering; each field is best-effort per platform.
type StatusData struct {
	Host      string
	Kernel    string
	Net       string // full first-INET line exactly as `ip -o -4` printed it
	NetShort  string // "iface cidr" for compact strips
	Uptime    string
	AuditRows int
	Toolkit   bool
}

func StatusSnapshot() StatusData {
	s := StatusData{
		Kernel: runtime.GOOS,
		Net:    "no IPv4 address found",
		Uptime: "unknown",
	}
	if h, err := os.Hostname(); err == nil {
		s.Host = h
	}
	if onLinux() {
		if o, e := ExecCommand("ip", "-o", "-4", "addr", "show").Output(); e == nil {
			firstInet := ""
			for _, l := range strings.Split(string(o), "\n") {
				if strings.Contains(l, "inet") {
					if firstInet == "" {
						firstInet = strings.TrimSpace(l)
					}
					// prefer a non-loopback INET: "lo" is never the useful answer
					if !strings.HasPrefix(strings.TrimSpace(strings.SplitN(l, " ", 2)[1]), "lo ") {
						s.Net = strings.TrimSpace(l)
						if f := strings.Fields(l); len(f) > 3 && f[2] == "inet" {
							s.NetShort = f[1] + " " + f[3]
						}
						break
					}
				}
			}
			if firstInet != "" && s.Net == "no IPv4 address found" {
				s.Net = firstInet // only loopback exists; show it honestly
			}
		}
		if o, e := ExecCommand("uname", "-r").Output(); e == nil {
			s.Kernel = strings.TrimSpace(string(o))
		}
	}
	if o, err := os.ReadFile("/proc/uptime"); err == nil {
		s.Uptime = strings.TrimSpace(strings.SplitN(string(o), " ", 2)[0]) + " seconds"
	}
	if b, err := os.ReadFile(DefaultAuditPath()); err == nil {
		s.AuditRows = len(strings.Split(strings.TrimRight(string(b), "\n"), "\n"))
	}
	if st, err := os.Stat("/toolkit"); err == nil && st.IsDir() {
		s.Toolkit = true
	}
	return s
}

// runDetectHW implements `manx detect-hw` (read-only): mirrors detect-hw.sh by
// parsing lspci -nnk and writing the UNCLAIMED report to the output dir.
func runDetectHW(args []string) error {
	if !onLinux() {
		return fmt.Errorf("detect-hw is Linux-native (running on %s); use the ISO", runtime.GOOS)
	}
	if !binExists("lspci") {
		return fmt.Errorf("lspci not available; install pciutils on the rescue media")
	}
	out, err := ExecCommand("lspci", "-nnk").Output()
	if err != nil {
		return fmt.Errorf("lspci failed: %v", err)
	}
	report := strings.Join(unclaimedLines(string(out)), "\n")
	path := filepath.Join(toolkitOut(), "drivers-report.txt")
	writeErr := os.WriteFile(path, []byte(report+"\n"), 0o600)
	if report == "" {
		fmt.Println("detect-hw: no UNCLAIMED PCI devices found")
	} else {
		fmt.Println(report)
	}
	if writeErr != nil {
		fmt.Printf("detect-hw: report could not be written (%s): %v\n", path, writeErr)
	} else {
		fmt.Printf("detect-hw: report at %s\n", path)
	}
	return nil
}

// runImgIn implements `manx img-in` (read-only): disk inventory via lsblk/blkid.
func runImgIn(args []string) error {
	if !onLinux() {
		return fmt.Errorf("img-in is Linux-native (running on %s); use the ISO", runtime.GOOS)
	}
	var report strings.Builder
	missing := []string{}
	for _, pair := range [][2]string{
		{"lsblk", "lsblk -o NAME,SIZE,TYPE,FSTYPE,MOUNTPOINT"},
		{"blkid", "blkid"},
	} {
		if !binExists(pair[0]) {
			missing = append(missing, pair[0])
			continue
		}
		fields := strings.Fields(pair[1])
		out, err := ExecCommand(fields[0], fields[1:]...).Output()
		if err != nil {
			return fmt.Errorf("%s failed: %v", pair[0], err)
		}
		report.WriteString("-- " + pair[0] + " --\n")
		report.Write(out)
		report.WriteString("\n")
	}
	if len(missing) > 0 {
		fmt.Printf("img-in: skipped (not on the rescue media): %s\n", strings.Join(missing, ", "))
	}
	path := filepath.Join(toolkitOut(), "img-in-report.txt")
	writeErr := os.WriteFile(path, []byte(report.String()), 0o600)
	fmt.Print(report.String())
	if writeErr != nil {
		fmt.Printf("img-in: report could not be written (%s): %v\n", path, writeErr)
	} else {
		fmt.Printf("img-in: report at %s\n", path)
	}
	return nil
}
