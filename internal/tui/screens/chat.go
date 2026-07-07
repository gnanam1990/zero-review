package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// Chat renders the chatbot screen.
func Chat(theme *core.Theme, session *review.ReviewSession, messages []core.ChatMessage, finding *review.Finding, width, height int) string {
	contextLabel := ""
	if session != nil {
		contextLabel = fmt.Sprintf("PR #%d", session.PR.Number)
	}
	if finding != nil {
		contextLabel += fmt.Sprintf(" · %s:%d", finding.FilePath, finding.LineStart)
	}

	historyHeight := height - 8
	if historyHeight < 1 {
		historyHeight = 1
	}

	chatPanel := components.ChatPanel(theme, messages, contextLabel, width, historyHeight)
	prompts := components.RenderQuickPrompts(theme, width)

	return lipgloss.JoinVertical(lipgloss.Left,
		theme.PanelTitleStyle.Render("Chat with Zero"),
		chatPanel,
		prompts,
	)
}

// SuggestedReply returns a canned AI response for a user message.
func SuggestedReply(userMsg string, finding *review.Finding) string {
	lower := strings.ToLower(userMsg)
	switch {
	case strings.Contains(lower, "shorter"):
		return "Shorter version: \"Verify the webhook signature before updating payment state.\""
	case strings.Contains(lower, "explain"):
		if finding != nil {
			return finding.Explanation
		}
		return "This PR introduces payment webhook handling. The main concern is ordering of signature verification versus state mutation."
	case strings.Contains(lower, "risk"):
		return "The highest risk is the high-severity finding: signature verification happens after payment state mutation."
	case strings.Contains(lower, "test"):
		return "Consider adding tests for invalid signatures, replay attacks, and oversized request bodies."
	case strings.Contains(lower, "security"):
		return "There are two security findings: signature-before-mutation order and raw error responses leaking details."
	default:
		return "I can explain the finding, suggest a shorter comment, or help draft tests. What would you like?"
	}
}
