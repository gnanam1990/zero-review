package tui

import (
	"strings"

	"github.com/gnanam1990/zero-review/internal/review"
)

func handleChatKey(m Model, key string) (teaModel, teaCmd) {
	switch key {
	case "esc", "b":
		m.Screen = ScreenFindings
		m.ChatInput = ""
		m.ChatMode = ""
	case "enter":
		return submitChat(m)
	}
	return m, nil
}

func submitChat(m Model) (teaModel, teaCmd) {
	if m.ChatInput == "" {
		return m, nil
	}
	m.Chat = append(m.Chat, ChatMessage{Role: "user", Text: m.ChatInput})
	reply := aiReply(m)
	m.Chat = append(m.Chat, ChatMessage{Role: "ai", Text: reply})

	if m.Selected != nil && m.ChatMode == "edit" && strings.TrimSpace(m.ChatInput) != "" {
		m.Selected.EditedComment = strings.TrimSpace(m.ChatInput)
		m.Selected.Status = review.StatusEdited
	}
	m.ChatInput = ""
	return m, nil
}

func aiReply(m Model) string {
	if m.Selected == nil {
		return "Select a finding first to chat about it."
	}
	q := strings.ToLower(m.ChatInput)
	switch {
	case strings.Contains(q, "explain"):
		return m.Selected.Explanation
	case strings.Contains(q, "security"):
		if m.Selected.Category == review.CategorySecurity {
			return "Yes, this is flagged as a security issue."
		}
		return "This finding is not categorized as security."
	case strings.Contains(q, "shorter") || strings.Contains(q, "shorter comment"):
		return shorten(m.Selected.SuggestedComment)
	case strings.Contains(q, "test") || strings.Contains(q, "tests"):
		return "Consider adding a test case that exercises this path and asserts the expected behavior."
	default:
		return "I can explain the issue, suggest a shorter comment, or help draft a test. What would you like?"
	}
}

func shorten(s string) string {
	if len(s) < 80 {
		return s
	}
	words := strings.Fields(s)
	var out []string
	total := 0
	for _, w := range words {
		if total+len(w)+1 > 80 {
			break
		}
		out = append(out, w)
		total += len(w) + 1
	}
	if len(out) == 0 {
		return s[:80] + "..."
	}
	return strings.Join(out, " ") + "."
}

func renderChat(m Model) string {
	var b strings.Builder
	b.WriteString(m.Styles.Header.Render(" Chat ") + "\n\n")
	for _, msg := range m.Chat {
		if msg.Role == "user" {
			b.WriteString(m.Styles.ChatUser.Render("you: "+msg.Text) + "\n")
		} else {
			b.WriteString(m.Styles.ChatSystem.Render("ai: "+msg.Text) + "\n")
		}
	}
	b.WriteString("\n" + m.Styles.Panel.Render("Input: "+m.ChatInput+"_") + "\n")
	b.WriteString(m.Styles.Help.Render("[Enter] send  [esc] back"))
	return b.String()
}
