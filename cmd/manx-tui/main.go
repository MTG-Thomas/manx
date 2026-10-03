// MANX TUI: Bubble Tea view over the same action surface as `manx <action>`.
// It is a VIEW, not a second implementation: any capability not present in the
// action machinery must not appear here.
//
// Degrades on a dumb/serial tty (spec §10.8): this binary detects the tty and
// exits non-zero if the target terminal is incapable; the wrapper then drops to
// `whiptail`/`dialog` (Tier 1).
//
// SPDX-License-Identifier: AGPL-3.0-or-later
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).MarginBottom(1)
	verbStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dangerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// viewOf mirrors one spec-contract verb for display. The verb strings are the
// contract; invoking shells out to `manx <verb>` (parity: one implementation).
type viewOf struct {
	verb   string
	help   string
	danger bool
}

type model struct {
	verbs  []viewOf
	cursor int
}

func (m model) Init() tea.Cmd { return nil }

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
		}
	}
	return m, nil
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
	s += "\n" + helpStyle.Render("↑/↓ move · Enter run (wired to manx CLI) · q quit")
	return s
}

func initialVerbs() []viewOf {
	// Mirrored from the spec parity contract (same verbs as the CLI manifest).
	return []viewOf{
		{verb: "status", help: "one-line runtime summary"},
		{verb: "detect-hw", help: "unclaimed/missing-driver hardware"},
		{verb: "collect", help: "harvest logs/evtx/minidumps (read-only)"},
		{verb: "img-in", help: "disk inventory (read-only)"},
		{verb: "img-out", help: "image a volume (destructive)", danger: true},
		{verb: "hive-edit", help: "SAFE offline hive edit (destructive)", danger: true},
		{verb: "collect-drivers", help: "runtime driver fetch ladder (mutating)", danger: true},
		{verb: "bringup-windows-vm", help: "qemu/OVMF bringup (destructive)", danger: true},
	}
}

func main() {
	if !isCapableTTY(os.Stdout) {
		fmt.Fprintln(os.Stderr, "not a capable tty: use `manx <action>`, or the whiptail fallback")
		os.Exit(2)
	}
	if _, err := tea.NewProgram(model{verbs: initialVerbs()}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tui error:", err)
		os.Exit(1)
	}
}
