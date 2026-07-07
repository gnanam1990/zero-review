package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Report renders the saved report screen.
func Report(theme *core.Theme, session *review.ReviewSession, reportPath string, posted int, width, height int) string {
	if session == nil {
		return components.EmptyState(theme, "", "No report available", "Start and save a review to see the report.", width, height)
	}

	status := "Saved"
	if reportPath == "" {
		status = "Not saved yet"
	}

	summary := strings.Join([]string{
		fmt.Sprintf("Status: %s", status),
		fmt.Sprintf("Path:   %s", reportPath),
		"",
		"Summary",
		fmt.Sprintf("- %d findings found", len(session.Findings)),
		fmt.Sprintf("- %d approved", len(session.Approved())),
		fmt.Sprintf("- %d rejected", len(session.Rejected())),
		fmt.Sprintf("- %d post failures", 0),
	}, "\n")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center,
		theme.ButtonPrimary.Render(" [ Open Report ] "),
		"  ",
		theme.ButtonSecondary.Render(" [ Copy Path ] "),
		"  ",
		theme.ButtonSecondary.Render(" [ Back to Dashboard ] "),
		"  ",
		theme.ButtonSecondary.Render(" [ Exit ] "),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		theme.PanelTitleStyle.Render("Review Report"),
		summary,
		"",
		buttons,
	)

	return theme.PanelStyle.Width(width).Height(height).Render(content)
}
