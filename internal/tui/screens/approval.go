package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// PostMode represents the GitHub review event to submit.
type PostMode string

const (
	PostModeComment        PostMode = "comment"
	PostModeRequestChanges PostMode = "request_changes"
	PostModeApprove        PostMode = "approve"
	PostModeReportOnly     PostMode = "report_only"
)

// Approval renders the final approval screen.
func Approval(theme *core.Theme, session *review.ReviewSession, mode PostMode, noPost bool, width, height int) string {
	if session == nil {
		return components.EmptyState(theme, "", "Nothing to approve", "Start a review first.", width, height)
	}

	approved := session.Approved()
	rejected := session.Rejected()
	edited := session.Edited()

	summary := strings.Join([]string{
		fmt.Sprintf("Approved comments: %d", len(approved)),
		fmt.Sprintf("Rejected comments: %d", len(rejected)),
		fmt.Sprintf("Edited comments:   %d", len(edited)),
		fmt.Sprintf("Summary comment:   %s", "Yes"),
	}, "\n")

	modePanel := theme.PanelStyle.Width(width / 2).Render(strings.Join([]string{
		theme.PanelTitleStyle.Render("Posting mode (m to cycle)"),
		fmt.Sprintf("> %s", ModeLabel(mode)),
		"",
		"  " + ModeLabel(PostModeComment),
		"  " + ModeLabel(PostModeRequestChanges),
		"  " + ModeLabel(PostModeApprove),
		"  " + ModeLabel(PostModeReportOnly),
	}, "\n"))

	var findingLines []string
	findingLines = append(findingLines, theme.PanelTitleStyle.Render("Approved findings"))
	for _, f := range approved {
		findingLines = append(findingLines, fmt.Sprintf("✓ %-7s %s:%-3d %s", f.Severity, f.FilePath, f.LineStart, f.Title))
	}
	findingsPanel := theme.PanelStyle.Width(width - width/2 - 3).Render(strings.Join(findingLines, "\n"))

	top := lipgloss.JoinHorizontal(lipgloss.Top, modePanel, "  ", findingsPanel)

	postButton := theme.ButtonPrimary.Render(" [ Post Review ] ")
	if noPost || mode == PostModeReportOnly {
		postButton = theme.ButtonSecondary.Render(" [ Posting Disabled ] ")
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Center, postButton, "  ", theme.ButtonSecondary.Render(" [ Save Report ] "), "  ", theme.ButtonSecondary.Render(" [ Back ] "))

	content := lipgloss.JoinVertical(lipgloss.Left,
		theme.PanelTitleStyle.Render("Approval"),
		theme.MutedText.Render("Review is ready. Nothing has been posted yet."),
		"",
		summary,
		"",
		top,
		"",
		buttons,
	)

	return theme.PanelStyle.Width(width).Height(height).Render(content)
}

func ModeLabel(m PostMode) string {
	switch m {
	case PostModeComment:
		return "Comment only"
	case PostModeRequestChanges:
		return "Request changes"
	case PostModeApprove:
		return "Approve PR"
	case PostModeReportOnly:
		return "Save report only"
	default:
		return string(m)
	}
}
