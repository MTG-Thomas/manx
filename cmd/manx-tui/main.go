// MANX TUI: Bubble Tea view over the same action surface as `manx <action>`.
// It is a VIEW, not a second implementation: any capability not present in the
// action machinery must not appear here. Enter runs the selected verb through
// GateRunner.RunView("tui", ...) — the same gate + audit row as the CLI — with
// the verb's stdout captured into the output viewport. The header strip
// re-renders actions.StatusData on a ticker (display refresh is a render,
// not an action: strip and audit-tail refresh write no audit rows).
//
// Degrades on a dumb/serial tty (spec §10.8): this binary detects the tty and
// exits non-zero if the target terminal is incapable; the wrapper then drops to
// `whiptail`/`dialog` (Tier 1).
//
// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MTG-Thomas/manx/internal/actions"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	stripStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	truthStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("136"))
	verbStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	dangerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("144"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ruleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
)

// viewTUI is the audit-row view tag for every verb run from this surface.
const viewTUI = "tui"

// maxOutLines keeps the output viewport inside one 80x24 screen.
const maxOutLines = 5

// maxAuditLines is the audit-tail pane depth.
const maxAuditLines = 2

// stripTick / truthTick cadences (seconds).
const (
	stripTick = 2
	truthTick = 14
)

// truths are the motd's incident-response + creature-care lines, one at a
// time (rotating): the same truths the boot motd teaches, said where the
// work happens.
var truths = []string{
	"breathe -> snapshot -> act (in that order)",
	"`manx status` BEFORE anything else: know the box",
	"collect evidence BEFORE repairing the patient",
	"write down what you did BEFORE you do it",
	"the backup you have not restore-tested is a rumor",
	"two changes at once = two suspects",
	"destructive verb? read it twice, then --i-know",
	"correlated failure on a cadence? look for a ticker, not a ghost",
	"you are not a machine: water + stretch",
	"tired eyes reboot zero servers; take the break",
	"the snarky one-liner can wait until the practice is back online",
}

// viewOf mirrors one spec-contract verb for display.
type viewOf struct {
	verb   string
	help   string
	danger bool
}

type (
	stripMsg struct {
		snap  actions.StatusData
		audit []string
	}
	truthMsg struct{ idx int }
	tickMsg  struct{}
)

type model struct {
	verbs  []viewOf
	cursor int
	gr     *actions.GateRunner
	last   string

	snap    actions.StatusData
	truth   int
	auditLn []string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Second*time.Duration(stripTick), func(time.Time) tea.Msg { return tickMsg{} }),
		tea.Tick(time.Second*time.Duration(truthTick), func(time.Time) tea.Msg { return truthMsg{m.truth + 1} }),
	)
}

// runVerb runs the verb through the shared GateRunner while capturing stdout
// into a string. This is the contract behavior, not the TUI; the capture is
// presentation-only.
func (m *model) runVerb(verb string) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", m.gr.RunView(viewTUI, verb, nil)
	}
	saved := os.Stdout
	os.Stdout = w
	out := make(chan []byte, 1)
	go func() {
		b, _ := io.ReadAll(r)
		out <- b
	}()
	runErr := m.gr.RunView(viewTUI, verb, nil)
	_ = w.Close() // unblocks the reader regardless of run outcome
	b := <-out
	_ = r.Close()
	os.Stdout = saved
	return string(b), runErr
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.verbs)-1 {
				m.cursor++
			}
		case "enter":
			verb := m.verbs[m.cursor].verb
			out, err := m.runVerb(verb)
			switch {
			case errors.Is(err, actions.ErrGateRefused):
				m.last = fmt.Sprintf("\u2717 %s\n  refused: %v\n  destructive verbs are CLI-only:\n  run `manx %s --i-know` from the rescue shell",
					verb, err, verb)
			case err != nil:
				m.last = fmt.Sprintf("\u2717 %s\n%s", verb, strings.TrimSpace(out+"\n"+err.Error()))
			default:
				m.last = fmt.Sprintf("\u2713 %s\n%s", verb, strings.TrimSpace(out))
			}
			m.last = trimLines(m.last, maxOutLines)
		case "t":
			return m, func() tea.Msg { return truthMsg{m.truth + 1} }
		}
	case stripMsg:
		m.snap = msg.snap
		m.auditLn = msg.audit
	case tickMsg:
		return m, tea.Batch(
			tea.Tick(time.Second*time.Duration(stripTick), func(time.Time) tea.Msg { return tickMsg{} }),
			snapshotCmd(),
		)
	case truthMsg:
		m.truth = msg.idx % len(truths)
		return m, tea.Tick(time.Second*time.Duration(truthTick), func(time.Time) tea.Msg { return truthMsg{m.truth + 1} })
	}
	return m, nil
}

func snapshotCmd() tea.Cmd {
	return func() tea.Msg {
		return stripMsg{snap: actions.StatusSnapshot(), audit: auditTail(maxAuditLines)}
	}
}

func auditTail(n int) []string {
	b, err := os.ReadFile(actions.DefaultAuditPath())
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	if n > len(lines) {
		n = len(lines)
	}
	return lines[len(lines)-n:]
}

// compactAuditRow renders one audit JSON line as "HH:MM view verb result".
func compactAuditRow(line string) string {
	var row struct {
		Ts   string `json:"ts"`
		View string `json:"view"`
		Verb string `json:"verb"`
		Dang string `json:"danger"`
		Res  string `json:"result"`
	}
	if err := json.Unmarshal([]byte(line), &row); err != nil {
		return "bad audit row"
	}
	stamp := ""
	if len(row.Ts) >= 16 {
		stamp = row.Ts[11:16] // HH:MM of the UTC ts
	}
	danger := ""
	if row.Dang != "" && row.Dang != "i-know" {
		danger = " (" + row.Dang + ")"
	}
	return fmt.Sprintf("%s \u00b7 %-4s \u00b7 %-16s \u00b7 %s%s", stamp, row.View, row.Verb, row.Res, danger)
}

func trimLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = append([]string{"..."}, lines[len(lines)-n:]...)
	}
	return strings.Join(lines, "\n")
}

func compactUptime(secondsField string) string {
	secs, err := strconv.ParseFloat(secondsField, 64)
	if err != nil {
		return "?" + secondsField + "s"
	}
	return fmt.Sprintf("up %dm%02ds", int(secs)/60, int(secs)%60)
}

func (m model) stripLine() string {
	up := compactUptime(strings.Fields(m.snap.Uptime)[0])
	net := m.snap.NetShort
	if net == "" {
		net = "no net"
	}
	kit := "kit up"
	if !m.snap.Toolkit {
		kit = "no kit"
	}
	return fmt.Sprintf("%s %s \u00b7 %s \u00b7 %s \u00b7 %s \u00b7 %s \u00b7 audit %d",
		time.Now().UTC().Format("15:04:05"), m.snap.Host, m.snap.Kernel, up, kit, net, m.snap.AuditRows)
}

func (m model) View() string {
	s := titleStyle.Render("   /\\_/\\     MANX \u00b7 one tool, three views, zero panic") + "\n"
	s += stripStyle.Render("  ( o.o )   "+m.stripLine()) + "\n"
	s += truthStyle.Render(" meow<  >\u2500 \"") + truthStyle.Render(truths[m.truth]) + "\"\n"
	s += ruleStyle.Render(strings.Repeat("\u2500", 64)) + "\n"
	for i, v := range m.verbs {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		style := verbStyle
		if v.danger {
			style = dangerStyle
		}
		s += fmt.Sprintf("%s%-18s %s\n", cursor, v.verb, style.Render(v.help))
	}
	if m.last != "" {
		s += ruleStyle.Render(strings.Repeat("\u2500", 64)) + "\n"
		s += okStyle.Render(m.last) + "\n"
	}
	if len(m.auditLn) > 0 {
		s += ruleStyle.Render(strings.Repeat("\u2500", 64)) + "\n"
		for _, ln := range m.auditLn {
			s += helpStyle.Render("audit \u00b7 "+compactAuditRow(ln)) + "\n"
		}
	}
	s += ruleStyle.Render(strings.Repeat("\u2500", 64)) + "\n"
	s += helpStyle.Render("\u2191/\u2193 select \u00b7 Enter run (GateRunner-audited, view=tui) \u00b7 t truth \u00b7 q quit")
	return s
}

func initialVerbs() []viewOf {
	// Mirrored from the parity contract (internal/actions/parity.go — SPEC_VERBS).
	return []viewOf{
		{verb: "status", help: "one-line runtime summary"},
		{verb: "detect-hw", help: "unclaimed/missing-driver hardware"},
		{verb: "collect", help: "harvest logs/evtx/minidumps (read-only)"},
		{verb: "setup", help: "stage /toolkit from the boot media"},
		{verb: "img-in", help: "disk inventory (read-only)"},
		{verb: "img-out", help: "image a volume (destructive)", danger: true},
		{verb: "hive-edit", help: "SAFE offline hive edit (destructive)", danger: true},
		{verb: "bootstrap-drivers", help: "runtime driver fetch ladder (mutating)", danger: true},
		{verb: "bringup-windows-vm", help: "qemu/OVMF bringup (destructive)", danger: true},
	}
}

func main() {
	if !isCapableTTY(os.Stdout) {
		fmt.Fprintln(os.Stderr, "not a capable tty: use `manx <action>`, or the whiptail fallback")
		os.Exit(2)
	}
	gate := actions.Default().NewAuditor(actions.DefaultAuditPath())
	m := model{verbs: initialVerbs(), gr: gate}
	m.snap = actions.StatusSnapshot() // prime the strip before the first tick lands
	m.auditLn = auditTail(maxAuditLines)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tui error:", err)
		os.Exit(1)
	}
}
