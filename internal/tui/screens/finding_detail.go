package screens

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// FindingDetail renders the deep-dive finding view.
func FindingDetail(theme *core.Theme, finding *review.Finding, width, height int) string {
	card := components.FindingCard(theme, finding, width)
	return lipgloss.NewStyle().Width(width).Height(height).Render(card)
}
