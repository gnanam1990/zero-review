package components

import (
	"fmt"
	"strings"

	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// DiffViewer renders unified diff context around a file or finding.
func DiffViewer(theme *core.Theme, diff, targetFile string, highlightLine int, width, height int) string {
	if diff == "" {
		return theme.PanelStyle.Width(width).Height(height).Render(theme.MutedText.Render("No diff context available."))
	}

	lines := strings.Split(diff, "\n")
	var out []string
	inFile := targetFile == ""
	rendered := 0
	maxLines := height - 2

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			inFile = targetFile == "" || strings.Contains(line, targetFile)
		}
		if !inFile {
			continue
		}
		if rendered >= maxLines {
			break
		}

		styled := styleDiffLine(theme, line, highlightLine)
		out = append(out, styled)
		rendered++
	}

	if len(out) == 0 {
		return theme.PanelStyle.Width(width).Render(theme.MutedText.Render("No diff context for selected file."))
	}

	header := theme.PanelTitleStyle.Render(fmt.Sprintf("Diff: %s", targetFile))
	body := theme.CodeBlockStyle.Width(width - 2).Height(height - 3).Render(strings.Join(out, "\n"))
	return theme.PanelStyle.Width(width).Height(height).Render(header + "\n" + body)
}

func styleDiffLine(theme *core.Theme, line string, highlightLine int) string {
	prefix := " "
	if len(line) > 0 {
		prefix = line[:1]
	}

	switch prefix {
	case "+":
		return theme.DiffAddedStyle.Render(line)
	case "-":
		return theme.DiffRemovedStyle.Render(line)
	case "@":
		return theme.DiffHunkStyle.Render(line)
	default:
		return theme.DiffContextStyle.Render(line)
	}
}
