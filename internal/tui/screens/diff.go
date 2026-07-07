package screens

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Diff renders the diff viewer screen.
func Diff(theme *core.Theme, session *review.ReviewSession, finding *review.Finding, width, height int) string {
	if session == nil {
		return components.EmptyState(theme, "", "No diff available", "Start a review to load diff context.", width, height)
	}

	targetFile := ""
	highlightLine := 0
	if finding != nil {
		targetFile = finding.FilePath
		highlightLine = finding.LineStart
	}

	viewer := components.DiffViewer(theme, "", targetFile, highlightLine, width, height-2)
	return lipgloss.NewStyle().Width(width).Height(height).Render(viewer)
}
