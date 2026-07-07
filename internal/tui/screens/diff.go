package screens

import (
	"fmt"
	"strings"

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

	diff := syntheticDiff(targetFile, finding)
	viewer := components.DiffViewer(theme, diff, targetFile, highlightLine, width, height-2)
	return lipgloss.NewStyle().Width(width).Height(height).Render(viewer)
}

// syntheticDiff generates a plausible diff for the demo when no real diff is loaded.
func syntheticDiff(targetFile string, finding *review.Finding) string {
	if targetFile == "" {
		return "diff --git a/file.go b/file.go\n@@ -1,3 +1,5 @@\n context\n+added line\n+another added line\n context"
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("diff --git a/%s b/%s", targetFile, targetFile))
	lines = append(lines, fmt.Sprintf("--- a/%s", targetFile))
	lines = append(lines, fmt.Sprintf("+++ b/%s", targetFile))
	lines = append(lines, "@@ -10,7 +10,10 @@ func Handler()")
	for i := 0; i < 12; i++ {
		if finding != nil && i == 4 {
			lines = append(lines, fmt.Sprintf("+// %s", finding.Evidence))
		}
		lines = append(lines, fmt.Sprintf(" context line %d", i+1))
	}
	lines = append(lines, fmt.Sprintf("+// added in PR (file: %s)", targetFile))
	return strings.Join(lines, "\n")
}
