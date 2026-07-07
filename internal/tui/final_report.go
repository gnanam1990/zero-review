package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gnanam1990/zero-review/internal/report"
	"github.com/gnanam1990/zero-review/internal/review"
)

// reportsSavedMsg is emitted after local reports are written.
type reportsSavedMsg struct {
	Paths []string
	Err   error
}

func renderFinal(m Model) string {
	var b strings.Builder
	b.WriteString(m.Styles.Header.Render(" Final Report ") + "\n\n")

	risk := reviewRiskScore(m.Findings)
	riskStyle := m.Styles.Highlight
	if risk > 50 {
		riskStyle = m.Styles.Danger
	} else if risk > 25 {
		riskStyle = m.Styles.Warning
	}

	b.WriteString(m.Styles.Title.Render(m.Result.PRData.Title) + "\n")
	b.WriteString(m.Styles.Subtitle.Render(fmt.Sprintf("#%d · %s/%s", m.Result.PRData.PRNumber, m.Result.PRData.Owner, m.Result.PRData.Repo)) + "\n\n")
	b.WriteString(m.Styles.Normal.Render("Risk score: " + riskStyle.Render(fmt.Sprintf("%d/100", risk)) + "\n"))
	b.WriteString(m.Styles.Normal.Render("Approved: " + m.Styles.Highlight.Render(fmt.Sprintf("%d", len(review.Approved(m.Findings)))) + "  " +
		"Rejected: " + m.Styles.Danger.Render(fmt.Sprintf("%d", len(review.Rejected(m.Findings)))) + "  " +
		"Pending: " + m.Styles.Muted.Render(fmt.Sprintf("%d", len(review.Pending(m.Findings)))) + "\n\n"))

	if m.Summary != "" {
		b.WriteString(m.Styles.Panel.Render(m.Summary) + "\n\n")
	}

	if len(m.ReportPaths) > 0 {
		b.WriteString(m.Styles.Subtitle.Render("Saved reports") + "\n")
		for _, p := range m.ReportPaths {
			b.WriteString("  • " + p + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(m.Styles.Help.Render("[q] quit"))
	return b.String()
}

// saveReportsCmd writes markdown and JSON reports to ~/.zero-review/reports.
func saveReportsCmd(result review.Result) teaCmd {
	return func() tea.Msg {
		dir := filepath.Join(os.Getenv("HOME"), ".zero-review", "reports")
		mdPath, err := report.SaveMarkdown(result, dir)
		if err != nil {
			return reportsSavedMsg{Err: err}
		}
		jsonPath, err := report.SaveJSON(result, dir)
		if err != nil {
			return reportsSavedMsg{Err: err}
		}
		return reportsSavedMsg{Paths: []string{mdPath, jsonPath}}
	}
}

// teaMsg is a tiny alias so this file can reference the type without importing tea everywhere.
type teaMsg = tea.Msg
