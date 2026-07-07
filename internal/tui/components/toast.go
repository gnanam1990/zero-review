package components

import (
	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Toast renders a transient toast message.
// The caller is responsible for positioning it inside the layout.
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

	_ = width
	return style.Render(" " + toast.Message + " ")
}
