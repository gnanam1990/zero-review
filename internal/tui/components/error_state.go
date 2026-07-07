package components

import (
	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// ErrorState renders a friendly error with next action.
func ErrorState(theme *core.Theme, title, message, action string, width, height int) string {
	if title == "" {
		title = "Something went wrong"
	}
	if action == "" {
		action = "Press esc to go back."
	}
	body := lipgloss.JoinVertical(
		lipgloss.Center,
		theme.ErrorText.Render("⚠ "+title),
		"",
		theme.PrimaryText.Render(message),
		"",
		theme.MutedText.Render(action),
	)
	return theme.PanelStyle.
		Width(width).
		Height(height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(body)
}
