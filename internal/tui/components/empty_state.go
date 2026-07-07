package components

import (
	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// EmptyState renders a helpful empty screen.
func EmptyState(theme *core.Theme, icon, title, subtitle string, width, height int) string {
	if icon == "" {
		icon = "∅"
	}
	body := lipgloss.JoinVertical(
		lipgloss.Center,
		theme.PrimaryText.Render(icon),
		"",
		theme.BoldText.Render(title),
		theme.MutedText.Render(subtitle),
	)
	return theme.PanelStyle.
		Width(width).
		Height(height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(body)
}
