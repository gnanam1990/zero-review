package components

import (
	"fmt"
	"strings"

	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// ProgressTimeline renders the loading timeline.
func ProgressTimeline(theme *core.Theme, steps []core.LoadingStep, width int) string {
	var lines []string
	for _, step := range steps {
		icon := "·"
		if step.Done {
			icon = "✓"
		} else if step.Active {
			icon = "◐"
		}
		line := fmt.Sprintf("%s %s", icon, step.Label)
		if step.Active {
			lines = append(lines, theme.PrimaryText.Render(line))
		} else if step.Done {
			lines = append(lines, theme.MutedText.Render(line))
		} else {
			lines = append(lines, theme.MutedText.Render(line))
		}
	}
	return theme.PanelStyle.Width(width).Render(strings.Join(lines, "\n"))
}
