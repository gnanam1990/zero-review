package components

import (
	"github.com/charmbracelet/glamour"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// MarkdownPanel renders markdown text with Glamour, with a plain fallback.
func MarkdownPanel(theme *core.Theme, markdown string, width int) string {
	if width <= 0 {
		width = 60
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-4),
	)
	if err != nil {
		return theme.PanelStyle.Width(width).Render(theme.MutedText.Render(markdown))
	}

	out, err := r.Render(markdown)
	if err != nil {
		return theme.PanelStyle.Width(width).Render(theme.MutedText.Render(markdown))
	}

	return theme.PanelStyle.Width(width).Render(out)
}
