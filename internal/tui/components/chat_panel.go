package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// ChatPanel renders the chat history and input area.
func ChatPanel(theme *core.Theme, messages []core.ChatMessage, contextLabel string, width, height int) string {
	var b strings.Builder

	if contextLabel != "" {
		b.WriteString(theme.ChatMetaText.Render("Context: "+contextLabel) + "\n")
	}

	for _, msg := range messages {
		if msg.Role == "user" {
			b.WriteString(theme.ChatUserBubble.Render("you: "+msg.Text) + "\n")
		} else {
			b.WriteString(theme.ChatAIBubble.Render("ai: "+msg.Text) + "\n")
		}
	}

	historyHeight := height - 6
	if historyHeight < 1 {
		historyHeight = 1
	}
	history := theme.PanelStyle.
		Width(width - 2).
		Height(historyHeight).
		Render(b.String())

	inputBox := theme.BorderFocused.Width(width - 2).Render(
		theme.MutedText.Render("Type a message. Enter to send, esc to back."),
	)

	return lipgloss.JoinVertical(lipgloss.Left, history, inputBox)
}

// QuickPrompts returns suggested chat prompts.
func QuickPrompts() []string {
	return []string{
		"Explain this PR",
		"What is risky here?",
		"Make comment shorter",
		"Generate test checklist",
		"Only show security issues",
		"Is this worth posting?",
		"Improve this review comment",
	}
}

// RenderQuickPrompts draws prompt chips.
func RenderQuickPrompts(theme *core.Theme, width int) string {
	var chips []string
	for _, p := range QuickPrompts() {
		chips = append(chips, theme.ButtonSecondary.Render(fmt.Sprintf(" %s ", p)))
	}
	return theme.HelpStyle.Width(width).Render(strings.Join(chips, " "))
}
