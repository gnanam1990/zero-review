package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Findings renders the findings list screen.
func Findings(theme *core.Theme, session *review.ReviewSession, cursor int, selected *int, filter string, width, height int) string {
	if session == nil || len(session.Findings) == 0 {
		return components.EmptyState(theme, "∅", "No findings", "Start a review to see AI findings here.", width, height)
	}

	findings := ApplyFilter(session.Findings, filter)
	if len(findings) == 0 {
		return components.EmptyState(theme, "∅", "No findings matched this filter", "Press x to clear filters.", width, height)
	}

	approved := len(session.Approved())
	rejected := len(session.Rejected())
	pending := len(session.Pending())
	filterText := "All"
	if filter != "" {
		filterText = filter
	}

	header := fmt.Sprintf("Filter: %s · Approved: %d · Rejected: %d · Pending: %d", filterText, approved, rejected, pending)

	showPreview := width >= 100
	tableWidth := width - 4
	if showPreview {
		tableWidth = width*2/3 - 2
	}

	tableHeight := height - 4
	ft := components.NewFindingTable(theme, findings, tableWidth, tableHeight)
	if cursor >= 0 && cursor < len(findings) {
		ft.Table.SetCursor(cursor)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		theme.PanelTitleStyle.Render("Findings"),
		theme.MutedText.Render(header),
		ft.Table.View(),
	)

	if !showPreview {
		return theme.PanelStyle.Width(width).Height(height).Render(content)
	}

	previewWidth := width - tableWidth - 4
	var preview *review.Finding
	if selected != nil && *selected >= 0 && *selected < len(findings) {
		preview = &findings[*selected]
	}
	card := components.FindingCard(theme, preview, previewWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		theme.PanelStyle.Width(tableWidth).Height(height).Render(content),
		" ",
		card,
	)
}

// ApplyFilter filters findings by severity (exported for update logic).
func ApplyFilter(findings []review.Finding, filter string) []review.Finding {
	if filter == "" {
		return findings
	}
	var out []review.Finding
	for _, f := range findings {
		if strings.EqualFold(string(f.Severity), filter) {
			out = append(out, f)
		}
	}
	return out
}
