package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// CommandBar renders a simple command palette input line.
func CommandBar(theme *core.Theme, query string, width int) string {
	label := theme.PrimaryText.Render("> ")
	display := query
	if display == "" {
		display = "Type a command..."
	}
	input := theme.InputStyle.Width(width - 6).Render(display + "_")
	box := theme.BorderFocused.Width(width).Padding(1).Render(label + input)
	return lipgloss.NewStyle().Render(box)
}

// CommandSuggestions returns common commands for the palette.
func CommandSuggestions() []string {
	return []string{
		"overview",
		"findings",
		"diff",
		"chat",
		"approval",
		"report",
		"settings",
		"save report",
		"post review",
		"quit",
	}
}

// RenderSuggestions draws the filtered suggestion list.
func RenderSuggestions(theme *core.Theme, query string, width int) string {
	var lines []string
	for _, cmd := range CommandSuggestions() {
		if query == "" || strings.Contains(cmd, query) {
			lines = append(lines, theme.MutedText.Render(fmt.Sprintf("  %s", cmd)))
		}
	}
	if len(lines) == 0 {
		lines = append(lines, theme.MutedText.Render("  no matches"))
	}
	return theme.PanelStyle.Width(width).Render(strings.Join(lines, "\n"))
}
