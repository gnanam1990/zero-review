package components

import (
	"fmt"
	"strings"

	"github.com/gnanam1990/zero-review/internal/review"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// FindingCard renders a compact preview card for a finding.
func FindingCard(theme *core.Theme, f *review.Finding, width int) string {
	if f == nil {
		return theme.PanelStyle.Width(width).Render(theme.MutedText.Render("No finding selected."))
	}

	var b strings.Builder
	b.WriteString(theme.PanelTitleStyle.Render(f.ID+": "+f.Title) + "\n")
	b.WriteString(fmt.Sprintf("%s · %s · %d%% confidence · %s\n\n",
		SeverityBadge(theme, string(f.Severity)),
		StatusBadge(theme, string(f.Status)),
		f.Confidence,
		f.Category,
	))

	b.WriteString(theme.MutedText.Render(fmt.Sprintf("%s:%d-%d (%s side)\n\n", f.FilePath, f.LineStart, f.LineEnd, f.DiffSide)))

	b.WriteString(theme.BoldText.Render("Why it matters") + "\n")
	b.WriteString(Wrap(f.Explanation, width-6) + "\n\n")

	b.WriteString(theme.BoldText.Render("Suggested PR comment") + "\n")
	b.WriteString(theme.CodeBlockStyle.Width(width-6).Render(Wrap(f.PostComment(), width-10)) + "\n")

	if f.SuggestedFix != "" {
		b.WriteString("\n" + theme.BoldText.Render("Suggested fix") + "\n")
		b.WriteString(theme.CodeBlockStyle.Width(width-6).Render(Wrap(f.SuggestedFix, width-10)) + "\n")
	}

	actions := theme.HelpStyle.Render("[a] approve · [r] reject · [e] edit · [c] chat · [d] diff")
	b.WriteString("\n" + actions)

	return theme.PanelStyle.Width(width).Render(b.String())
}
