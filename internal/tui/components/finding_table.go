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
	// Narrow terminals can't fit the full set of columns. Use comfortable
	// widths by default, then shrink category/file columns only when needed
	// so the title column remains readable.
	catW, fileW := 14, 24
	fixed := 4 + 6 + catW + 5 + fileW + 5
	titleW := width - fixed
	if titleW < 8 {
		catW, fileW = 8, 12
		fixed = 4 + 6 + catW + 5 + fileW + 5
		titleW = width - fixed
		if titleW < 0 {
			titleW = 0
		}
	}

	cols := []table.Column{
		{Title: "St", Width: 4},
		{Title: "Sev", Width: 6},
		{Title: "Cat", Width: catW},
		{Title: "Conf", Width: 5},
		{Title: "File", Width: fileW},
		{Title: "Line", Width: 5},
		{Title: "Title", Width: titleW},
	}

	if height < 4 {
		height = 4
	}

	rows := make([]table.Row, 0, len(findings))
	for _, f := range findings {
		rows = append(rows, table.Row{
			statusIcon(string(f.Status)),
			string(f.Severity),
			string(f.Category),
			fmt.Sprintf("%d%%", f.Confidence),
			Truncate(f.FilePath, fileW-2),
			fmt.Sprintf("%d", f.LineStart),
			Truncate(f.Title, titleW-2),
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
