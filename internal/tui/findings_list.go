package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
)

func newFindingsTable(styles Styles, findings []review.Finding) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "Severity", Width: 10},
		{Title: "Category", Width: 18},
		{Title: "File", Width: 24},
		{Title: "Line", Width: 7},
		{Title: "Conf", Width: 6},
		{Title: "Status", Width: 10},
	}

	rows := make([]table.Row, 0, len(findings))
	for _, f := range findings {
		rows = append(rows, table.Row{
			f.ID,
			string(f.Severity),
			string(f.Category),
			truncate(f.FilePath, 22),
			fmt.Sprintf("%d", f.LineStart),
			fmt.Sprintf("%d%%", f.Confidence),
			string(f.Status),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t.SetStyles(table.Styles{
		Header:   styles.Title.Bold(true),
		Selected: styles.Selected,
		Cell:     styles.Normal,
	})
	return t
}

func handleFindingsKey(m Model, key string) (teaModel, teaCmd) {
	switch key {
	case "enter":
		if m.Cursor >= 0 && m.Cursor < len(m.Findings) {
			m.Selected = &m.Findings[m.Cursor]
			m.Screen = ScreenDetail
		}
	case "a":
		if m.Cursor >= 0 && m.Cursor < len(m.Findings) {
			m.Findings[m.Cursor].Status = review.StatusApproved
			m.Findings[m.Cursor].EditedComment = ""
			refreshTable(&m)
		}
	case "r":
		if m.Cursor >= 0 && m.Cursor < len(m.Findings) {
			m.Findings[m.Cursor].Status = review.StatusRejected
			m.Findings[m.Cursor].EditedComment = ""
			refreshTable(&m)
		}
	case "e":
		if m.Cursor >= 0 && m.Cursor < len(m.Findings) {
			m.Selected = &m.Findings[m.Cursor]
			m.ChatInput = m.Findings[m.Cursor].SuggestedComment
			m.Screen = ScreenChat
		}
	case "c":
		if m.Cursor >= 0 && m.Cursor < len(m.Findings) {
			m.Selected = &m.Findings[m.Cursor]
			m.ChatInput = ""
			m.Screen = ScreenChat
		}
	case "p":
		m.Screen = ScreenApproval
	case "s":
		return m, saveReportsCmd(m.Result)
	case "q":
		return m, quitCmd()
	case "d":
		if m.Cursor >= 0 && m.Cursor < len(m.Findings) {
			m.Selected = &m.Findings[m.Cursor]
			m.Screen = ScreenDiff
		}
	}
	return m, nil
}

func refreshTable(m *Model) {
	m.table = newFindingsTable(m.Styles, m.Findings)
}

func renderFindings(m Model) string {
	var b strings.Builder
	b.WriteString(m.Styles.Header.Render(" Findings ") + "\n\n")
	if len(m.Findings) == 0 {
		b.WriteString(m.Styles.Muted.Render("No findings to display."))
		return b.String()
	}
	b.WriteString(m.table.View() + "\n")
	b.WriteString(m.Styles.Help.Render(
		"[Enter] detail  [a] approve  [r] reject  [e] edit  [c] chat  [d] diff  [p] post  [s] save  [q] quit"))
	return b.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func statusStyle(styles Styles, status review.Status) lipgloss.Style {
	switch status {
	case review.StatusApproved:
		return styles.Success
	case review.StatusRejected:
		return styles.Danger
	case review.StatusEdited:
		return styles.Warning
	default:
		return styles.Muted
	}
}
