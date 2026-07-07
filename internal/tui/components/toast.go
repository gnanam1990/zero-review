package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Toast renders a transient toast message anchored to the bottom right.
func Toast(theme *core.Theme, toast *core.Toast, width int) string {
	if toast == nil {
		return ""
	}

	var style lipgloss.Style
	switch toast.Kind {
	case "success":
		style = theme.ToastSuccess
	case "error":
		style = theme.ToastError
	default:
		style = theme.ToastInfo
	}

	msg := style.Render(" " + toast.Message + " ")
	pad := width - lipgloss.Width(msg) - 2
	if pad < 0 {
		pad = 0
	}
	return strings.Repeat("\n", 0) + strings.Repeat(" ", pad) + msg
}
