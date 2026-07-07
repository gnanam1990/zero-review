package components

import core "github.com/gnanam1990/zero-review/internal/tui/core"

// Footer renders the bottom command/help bar.
func Footer(theme *core.Theme, keys core.KeyMap, screen core.Screen, layout core.Layout) string {
	text := keys.FooterHelp(screen, layout)
	return theme.FooterStyle.
		Width(layout.Width).
		Render(text)
}
