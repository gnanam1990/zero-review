package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Header renders the top app bar.
func Header(theme *core.Theme, screen core.Screen, session *review.ReviewSession, width int, showSidebar bool) string {
	var parts []string

	parts = append(parts, theme.HeaderStyle.Render(" Zero Review "))

	if showSidebar {
		parts = append(parts, theme.HeaderMutedStyle.Render("  "+core.ScreenName(screen)))
	}

	if session != nil {
		prLabel := fmt.Sprintf("#%d %s", session.PR.Number, Truncate(session.PR.Title, 35))
		parts = append(parts, theme.HeaderMutedStyle.Render(prLabel))

		risk := riskBadge(theme, session.RiskScore)
		parts = append(parts, risk)

		if session.Provider != "" {
			parts = append(parts, theme.HeaderMutedStyle.Render(session.Provider+"/"+session.Model))
		}
	}

	// Push auth status to the right.
	auth := theme.HeaderMutedStyle.Render("● GitHub ready")
	if session == nil {
		auth = theme.HeaderAlertStyle.Render("● No PR loaded")
	}

	left := lipgloss.JoinHorizontal(lipgloss.Center, parts...)
	padding := width - lipgloss.Width(left) - lipgloss.Width(auth)
	if padding < 0 {
		padding = 0
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, left, strings.Repeat(" ", padding), auth)
}

func riskBadge(theme *core.Theme, score int) string {
	label := fmt.Sprintf("Risk %d/100", score)
	if score >= 50 {
		return theme.DangerBadge.Render(label)
	}
	if score >= 25 {
		return theme.WarningBadge.Render(label)
	}
	return theme.SuccessBadge.Render(label)
}
