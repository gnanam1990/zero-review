package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/gnanam1990/zero-review/internal/review"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// FindingTable builds the findings table component.
type FindingTable struct {
	Table table.Model
}

// NewFindingTable creates a table for the given findings.
func NewFindingTable(theme *core.Theme, findings []review.Finding, width, height int) FindingTable {
	cols := []table.Column{
		{Title: "St", Width: 4},
		{Title: "Sev", Width: 6},
		{Title: "Cat", Width: 14},
		{Title: "Conf", Width: 5},
		{Title: "File", Width: 24},
		{Title: "Line", Width: 5},
		{Title: "Title", Width: width - 70},
	}

	rows := make([]table.Row, 0, len(findings))
	for _, f := range findings {
		rows = append(rows, table.Row{
			statusIcon(string(f.Status)),
			string(f.Severity),
			string(f.Category),
			fmt.Sprintf("%d%%", f.Confidence),
			Truncate(f.FilePath, 22),
			fmt.Sprintf("%d", f.LineStart),
			Truncate(f.Title, width-75),
		})
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height-2),
	)

	t.SetStyles(table.Styles{
		Header:   theme.TableHeaderStyle,
		Selected: theme.TableSelectedStyle,
		Cell:     theme.TableCellStyle,
	})

	return FindingTable{Table: t}
}

func statusIcon(status string) string {
	switch status {
	case "approved":
		return "✓"
	case "rejected":
		return "✕"
	case "edited":
		return "✎"
	case "posted":
		return "↗"
	default:
		return "○"
	}
}
