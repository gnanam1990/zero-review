package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/review"
)

func handleDetailKey(m Model, key string) (teaModel, teaCmd) {
	if m.Selected == nil {
		m.Screen = ScreenFindings
		return m, nil
	}
	switch key {
	case "a":
		m.Selected.Status = review.StatusApproved
		m.Selected.EditedComment = ""
	case "r":
		m.Selected.Status = review.StatusRejected
		m.Selected.EditedComment = ""
	case "e":
		m.ChatInput = m.Selected.SuggestedComment
		m.ChatMode = "edit"
		m.Screen = ScreenChat
		return m, nil
	case "c":
		m.ChatInput = ""
		m.ChatMode = "chat"
		m.Screen = ScreenChat
		return m, nil
	case "d":
		m.Screen = ScreenDiff
		return m, nil
	case "esc", "b":
		m.Screen = ScreenFindings
	}
	return m, nil
}

func renderDetail(m Model) string {
	if m.Selected == nil {
		return m.Styles.Muted.Render("No finding selected. Press esc.")
	}
	f := m.Selected
	var b strings.Builder
	b.WriteString(m.Styles.Header.Render(" Finding Detail ") + "\n\n")
	b.WriteString(m.Styles.Title.Render(f.ID+": "+f.Title) + "\n")
	b.WriteString(fmt.Sprintf("%s · %s · %d%% confidence · %s\n\n",
		severityBadge(m.Styles, string(f.Severity)), f.Category, f.Confidence, statusBadge(m.Styles, string(f.Status))))

	b.WriteString(m.Styles.Subtitle.Render("File") + "\n")
	b.WriteString(m.Styles.Normal.Render(fmt.Sprintf("%s:%d-%d (%s side)\n\n", f.FilePath, f.LineStart, f.LineEnd, f.DiffSide)))

	b.WriteString(m.Styles.Subtitle.Render("Explanation") + "\n")
	b.WriteString(wrap(f.Explanation, m.Width-4) + "\n\n")

	b.WriteString(m.Styles.Subtitle.Render("Suggested comment") + "\n")
	b.WriteString(m.Styles.Panel.Render(wrap(f.SuggestedComment, m.Width-8)) + "\n\n")

	if f.SuggestedFix != "" {
		b.WriteString(m.Styles.Subtitle.Render("Suggested fix") + "\n")
		b.WriteString(m.Styles.Code.Render(wrap(f.SuggestedFix, m.Width-8)) + "\n\n")
	}

	if f.Evidence != "" {
		b.WriteString(m.Styles.Subtitle.Render("Evidence") + "\n")
		b.WriteString(wrap(f.Evidence, m.Width-4) + "\n\n")
	}

	b.WriteString(m.Styles.Help.Render("[a] approve  [r] reject  [e] edit  [c] chat  [d] diff  [esc] back"))
	return b.String()
}

func wrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	return wrapText(s, width)
}

func wrapText(s string, width int) string {
	var lines []string
	var line string
	for _, word := range strings.Fields(s) {
		if len(line)+len(word)+1 > width && line != "" {
			lines = append(lines, line)
			line = word
		} else {
			if line == "" {
				line = word
			} else {
				line += " " + word
			}
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func severityBadge(styles Styles, sev string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#11111B")).
		Background(severityColor(sev)).
		Padding(0, 1).
		Render(" " + sev + " ")
}

func statusBadge(styles Styles, status string) string {
	var style lipgloss.Style
	switch status {
	case "approved":
		style = lipgloss.NewStyle().Foreground(styles.Success.GetForeground())
	case "rejected":
		style = lipgloss.NewStyle().Foreground(styles.Danger.GetForeground())
	case "edited":
		style = lipgloss.NewStyle().Foreground(styles.Warning.GetForeground())
	default:
		style = lipgloss.NewStyle().Foreground(styles.Muted.GetForeground())
	}
	return style.Render(status)
}
