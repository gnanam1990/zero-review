package components

import (
	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// ConfirmModal renders a centered confirmation dialog.
func ConfirmModal(theme *core.Theme, title, body, yesText, noText string, width, height int) string {
	titleBox := theme.PanelTitleStyle.Render(title)
	bodyBox := theme.PrimaryText.Render(Wrap(body, width-10))

	if yesText == "" {
		yesText = "Yes"
	}
	if noText == "" {
		noText = "Cancel"
	}

	yes := theme.ButtonPrimary.Render(" " + yesText + " ")
	no := theme.ButtonSecondary.Render(" " + noText + " ")
	buttons := lipgloss.JoinHorizontal(lipgloss.Center, yes, "  ", no)

	content := lipgloss.JoinVertical(lipgloss.Left, titleBox, "", bodyBox, "", buttons)
	box := theme.BorderFocused.
		Width(width).
		Padding(2, 4).
		Render(content)

	// Center in terminal.
	return lipgloss.NewStyle().
		Width(width + 8).
		Height(height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(box)
}
