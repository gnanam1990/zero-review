package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// HelpOverlay renders a full-screen help sheet.
func HelpOverlay(theme *core.Theme, screen core.Screen, width, height int) string {
	help := []struct{ key, desc string }{
		{"q / ctrl+c", "quit"},
		{"esc / b", "back or cancel"},
		{"?", "toggle this help"},
		{"tab", "next panel"},
		{"shift+tab", "previous panel"},
		{"/", "command palette"},
		{"o", "overview"},
		{"f", "findings"},
		{"d", "diff"},
		{"c", "chat"},
		{"p", "approval / post"},
		{"s", "save report"},
	}

	var lines []string
	lines = append(lines, theme.PanelTitleStyle.Render("Keyboard shortcuts"))
	lines = append(lines, "")
	for _, h := range help {
		lines = append(lines, fmt.Sprintf("%s  %s",
			theme.ShortcutKey.Render(" "+h.key+" "),
			theme.MutedText.Render(h.desc),
		))
	}

	if screen == core.ScreenFindings {
		lines = append(lines, "", theme.BoldText.Render("Findings screen"))
		for _, h := range []struct{ key, desc string }{
			{"↑/↓ or j/k", "move"},
			{"enter", "open selected"},
			{"a", "approve selected"},
			{"r", "reject selected"},
			{"e", "edit comment"},
			{"space", "toggle approval"},
			{"1-4", "filter severity"},
			{"x", "clear filters"},
		} {
			lines = append(lines, fmt.Sprintf("%s  %s",
				theme.ShortcutKey.Render(" "+h.key+" "),
				theme.MutedText.Render(h.desc),
			))
		}
	}

	content := theme.PanelStyle.Width(width - 8).Padding(2).Render(strings.Join(lines, "\n"))
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)
}
