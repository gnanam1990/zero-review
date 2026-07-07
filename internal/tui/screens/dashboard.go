package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Dashboard renders the overview screen.
func Dashboard(theme *core.Theme, session *review.ReviewSession, width, height int) string {
	if session == nil {
		return components.EmptyState(theme, "", "No review loaded", "Press s to start a new review.", width, height)
	}

	leftWidth := width/2 - 2
	rightWidth := width - leftWidth - 4

	// PR Summary panel
	prSummary := theme.PanelStyle.Width(leftWidth).Render(strings.Join([]string{
		theme.PanelTitleStyle.Render("PR Summary"),
		theme.BoldText.Render(session.PR.Title),
		theme.MutedText.Render("Author: @" + session.PR.Author),
		theme.MutedText.Render(fmt.Sprintf("Branch: %s → %s", session.PR.SourceBranch, session.PR.TargetBranch)),
		theme.MutedText.Render(fmt.Sprintf("Files: %d  +%d / -%d", session.PR.ChangedFiles, session.PR.Additions, session.PR.Deletions)),
	}, "\n"))

	// Stats panel
	counts := session.CountBySeverity()
	stats := theme.PanelStyle.Width(rightWidth).Render(strings.Join([]string{
		theme.PanelTitleStyle.Render("Review Stats"),
		fmt.Sprintf("Findings: %d", len(session.Findings)),
		fmt.Sprintf("High:   %d", counts[review.SeverityHigh]),
		fmt.Sprintf("Medium: %d", counts[review.SeverityMedium]),
		fmt.Sprintf("Low:    %d", counts[review.SeverityLow]),
		fmt.Sprintf("Info:   %d", counts[review.SeverityInfo]),
	}, "\n"))

	top := lipgloss.JoinHorizontal(lipgloss.Top, prSummary, "  ", stats)

	// AI Summary
	aiSummary := theme.PanelStyle.Width(width - 4).Render(strings.Join([]string{
		theme.PanelTitleStyle.Render("AI Summary"),
		theme.PrimaryText.Render(session.Summary),
	}, "\n"))

	// Next steps
	steps := theme.PanelStyle.Width(width - 4).Render(strings.Join([]string{
		theme.PanelTitleStyle.Render("Suggested next steps"),
		"1. Review the high-severity security finding first.",
		"2. Approve only comments that are useful to the PR author.",
		"3. Save a report before posting.",
	}, "\n"))

	content := lipgloss.JoinVertical(lipgloss.Left, top, "", aiSummary, "", steps)
	return theme.PanelStyle.Width(width).Height(height).Render(content)
}
