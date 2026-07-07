package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Sidebar renders the left navigation rail.
func Sidebar(theme *core.Theme, active core.Screen, height int) string {
	screens := core.NavigableScreens()
	var items []string

	for _, s := range screens {
		label := core.ScreenLabel(s)
		style := theme.SidebarItemStyle
		if s == active {
			style = theme.SidebarActiveStyle
		}
		items = append(items, style.Render(" "+label))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	available := height - 2
	if available > 0 {
		content += "\n" + strings.Repeat("\n", available-1)
	}

	return theme.SidebarStyle.
		Height(height).
		Width(core.SidebarWidth - 1).
		Render(content)
}
