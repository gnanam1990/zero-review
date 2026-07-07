package screens

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// LoadingReview renders the progress timeline.
func LoadingReview(theme *core.Theme, prURL string, steps []core.LoadingStep, tip string, width, height int) string {
	title := theme.PanelTitleStyle.Render("Reviewing PR")
	url := theme.MutedText.Render(prURL)

	timeline := components.ProgressTimeline(theme, steps, width-6)
	footer := theme.MutedText.Render("Tip: " + tip)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		url,
		"",
		timeline,
		"",
		footer,
	)

	box := theme.PanelStyle.Width(width - 4).Height(height - 4).Render(content)
	return components.Center(box, width, height)
}
