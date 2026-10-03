// MANX TUI: Bubble Tea view over the same action surface as `manx <action>`.
// It is a VIEW, not a second implementation: any capability not present in the
// action machinery must not appear here. Enter runs the selected verb through
// GateRunner.RunView("tui", ...) — the same gate + audit row as the CLI — with
// the verb's stdout captured into the viewport (verbs print to stdout; behind
// the alt screen their output must be captured, not leaked).
//
// Degrades on a dumb/serial tty (spec §10.8): this binary detects the tty and
// exits non-zero if the target terminal is incapable; the wrapper then drops to
// `whiptail`/`dialog` (Tier 1).
//
// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MTG-Thomas/manx/internal/actions"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).MarginBottom(1)
	verbStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dangerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// viewTUI is the audit-row view tag for every verb run from this surface.
const viewTUI = "tui"

// maxOutLines keeps the viewport inside one 80x24 screen.
const maxOutLines = 12

// viewOf mirrors one spec-contract verb for display.
type viewOf struct {
	verb   string
	help   string
	danger bool
}

type model struct {
	verbs  []viewOf
	cursor int
	gr     *actions.GateRunner
	last   string
}

func (m model) Init() tea.Cmd { return nil }

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
				m.last = fmt.Sprintf("✗ %s\n  refused: %v\n  destructive verbs are CLI-only:\n  run `manx %s --i-know` from the rescue shell",
					verb, err, verb)
			case err != nil:
				m.last = fmt.Sprintf("✗ %s\n%s", verb, strings.TrimSpace(out+"\n"+err.Error()))
			default:
				m.last = fmt.Sprintf("✓ %s\n%s", verb, strings.TrimSpace(out))
			}
			m.last = trimLines(m.last, maxOutLines)
		}
	}
	return m, nil
}

func trimLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = append([]string{"…"}, lines[len(lines)-n:]...)
	}
	return strings.Join(lines, "\n")
}

func (m model) View() string {
	s := titleStyle.Render("MANX - one action surface, three views")
	s += "\n\n"
	for i, v := range m.verbs {
		cursor := "  "
		if m.cursor == i {
			cursor = "❯ "
		}
		style := verbStyle
		if v.danger {
			style = dangerStyle
		}
		s += fmt.Sprintf("%s%-20s %s\n", cursor, v.verb, style.Render(v.help))
	}
	s += "\n"
	if m.last != "" {
		s += verbStyle.Render(m.last)
	}
	s += "\n" + helpStyle.Render("↑/↓ select · Enter run (GateRunner-audited, view=tui) · q quit")
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
	if _, err := tea.NewProgram(model{verbs: initialVerbs(), gr: gate}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tui error:", err)
		os.Exit(1)
	}
}
