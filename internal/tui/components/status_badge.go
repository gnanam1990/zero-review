package components

import (
	"fmt"

	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// SeverityBadge renders a severity badge.
func SeverityBadge(theme *core.Theme, severity string) string {
	return theme.SeverityBadge(severity).Render(fmt.Sprintf(" %s ", severity))
}

// StatusBadge renders a finding status badge.
func StatusBadge(theme *core.Theme, status string) string {
	icon := ""
	switch status {
	case "approved":
		icon = "✓"
	case "rejected":
		icon = "✕"
	case "edited":
		icon = "✎"
	case "posted":
		icon = "↗"
	default:
		icon = "○"
	}
	return theme.StatusBadge(status).Render(fmt.Sprintf(" %s %s ", icon, status))
}
