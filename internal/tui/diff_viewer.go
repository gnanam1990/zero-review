package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func handleDiffKey(m Model, key string) (teaModel, teaCmd) {
	switch key {
	case "esc", "b":
		if m.Selected != nil {
			m.Screen = ScreenDetail
		} else {
			m.Screen = ScreenFindings
		}
	}
	return m, nil
}

func renderDiff(m Model) string {
	var b strings.Builder
	b.WriteString(m.Styles.Header.Render(" Diff ") + "\n\n")

	target := ""
	if m.Selected != nil {
		target = m.Selected.FilePath
	}

	lines := strings.Split(m.Result.PRData.Diff, "\n")
	var out []string
	inFile := target == ""
	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			inFile = target == "" || strings.Contains(line, target)
		}
		if !inFile {
			continue
		}
		styled := m.Styles.Normal.Render(line)
		if strings.HasPrefix(line, "+") {
			styled = lipgloss.NewStyle().Foreground(lipgloss.Color("#2CDA9D")).Render(line)
		} else if strings.HasPrefix(line, "-") {
			styled = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render(line)
		} else if strings.HasPrefix(line, "@@") {
			styled = lipgloss.NewStyle().Foreground(lipgloss.Color("#F7B801")).Render(line)
		}
		out = append(out, styled)
		if len(out) > m.Height-8 {
			break
		}
	}

	if len(out) == 0 {
		b.WriteString(m.Styles.Muted.Render("No diff context available."))
	} else {
		b.WriteString(strings.Join(out, "\n"))
	}

	b.WriteString("\n\n" + m.Styles.Help.Render("[esc] back"))
	return b.String()
}
