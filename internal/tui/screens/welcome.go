package screens

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Welcome renders the welcome screen. focus is 0=start, 1=last, 2=settings.
func Welcome(theme *core.Theme, focus, width, height int) string {
	logo := theme.PrimaryText.
		Bold(true).
		Render("Zero Review")

	subtitle := theme.MutedText.Render("Interactive AI PR review inside your terminal.")
	body := theme.MutedText.Render("Paste a GitHub PR link. Review AI findings. Approve only what matters. Post clean PR comments.")

	start := renderWelcomeButton(theme, focus == 0, " [ Start Review ] ")
	last := renderWelcomeButton(theme, focus == 1, " [ Open Last Report ] ")
	settings := renderWelcomeButton(theme, focus == 2, " [ Settings ] ")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, start, "  ", last, "  ", settings)
	footer := theme.MutedText.Render("No auto-posting. You stay in control.")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		subtitle,
		"",
		body,
		"",
		"",
		buttons,
		"",
		"",
		footer,
	)

	box := theme.PanelStyle.
		Width(width - 4).
		Height(height - 4).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)

	return components.Center(box, width, height)
}

func renderWelcomeButton(theme *core.Theme, active bool, label string) string {
	if active {
		return theme.ButtonPrimary.Render(label)
	}
	return theme.ButtonSecondary.Render(label)
}
