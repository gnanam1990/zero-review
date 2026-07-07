package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/gnanam1990/zero-review/internal/github"
	"github.com/gnanam1990/zero-review/internal/review"
)

func handleApprovalKey(m Model, key string) (teaModel, teaCmd) {
	switch key {
	case "esc":
		m.Screen = ScreenFindings
	case "a":
		for i := range m.Findings {
			m.Findings[i].Status = review.StatusApproved
		}
		refreshTable(&m)
	case "r":
		for i := range m.Findings {
			m.Findings[i].Status = review.StatusRejected
		}
		refreshTable(&m)
	case "c":
		return resetApproval(m)
	case "enter", " ":
		return confirmPost(m)
	}
	return m, nil
}

func resetApproval(m Model) (teaModel, teaCmd) {
	for i := range m.Findings {
		m.Findings[i].Status = review.StatusPending
	}
	refreshTable(&m)
	return m, nil
}

func confirmPost(m Model) (teaModel, teaCmd) {
	approved := review.Approved(m.Findings)
	if len(approved) == 0 {
		m.Screen = ScreenFinal
		m.Summary = "No approved findings. Nothing posted to GitHub."
		m.Posted = 0
		return m, nil
	}
	if m.Options.NoPost {
		m.Screen = ScreenFinal
		m.Summary = fmt.Sprintf("%d approved. Posting disabled by --no-post.", len(approved))
		m.Posted = 0
		return m, nil
	}

	comments := make([]github.ReviewComment, 0, len(approved))
	for _, f := range approved {
		comments = append(comments, github.ReviewComment{
			Path: f.FilePath,
			Line: f.LineStart,
			Side: f.DiffSide,
			Body: f.PostComment(),
		})
	}

	m.Screen = ScreenFinal
	ctx := context.Background()
	err := m.GitHub.PostReview(ctx, m.Result.PRData.Owner, m.Result.PRData.Repo, m.Result.PRData.PRNumber, github.ReviewSubmission{
		Body:     fmt.Sprintf("zero-review: %d approved finding(s)", len(comments)),
		Event:    "COMMENT",
		Comments: comments,
	})
	if err != nil {
		m.Summary = "Failed to post review: " + err.Error()
		m.Posted = 0
	} else {
		m.Summary = fmt.Sprintf("Posted %d approved finding(s) to PR #%d.", len(comments), m.Result.PRData.PRNumber)
		m.Posted = len(comments)
	}
	return m, nil
}

func renderApproval(m Model) string {
	var b strings.Builder
	b.WriteString(m.Styles.Header.Render(" Approval ") + "\n\n")

	approved := review.Approved(m.Findings)
	rejected := review.Rejected(m.Findings)
	pending := review.Pending(m.Findings)
	b.WriteString(fmt.Sprintf("Approved: %s  Rejected: %s  Pending: %s\n\n",
		m.Styles.Highlight.Render(fmt.Sprintf("%d", len(approved))),
		m.Styles.Warning.Render(fmt.Sprintf("%d", len(rejected))),
		m.Styles.Muted.Render(fmt.Sprintf("%d", len(pending))),
	))

	for _, f := range m.Findings {
		icon := statusIcon(f.Status)
		b.WriteString(fmt.Sprintf("%s %s\n", icon, f.Title))
	}

	b.WriteString("\n" + m.Styles.Help.Render("[a] approve all  [r] reject all  [c] clear  [enter] post approved  [esc] back"))
	return b.String()
}

func statusIcon(s review.Status) string {
	switch s {
	case review.StatusApproved:
		return "✅"
	case review.StatusRejected:
		return "❌"
	case review.StatusEdited:
		return "✏️"
	default:
		return "⏳"
	}
}
