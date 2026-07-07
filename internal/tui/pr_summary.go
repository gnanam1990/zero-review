package tui

import (
	"fmt"
	"strings"
)

func handleSummaryKey(m Model, key string) (teaModel, teaCmd) {
	switch key {
	case "enter", "s":
		m.Screen = ScreenFindings
	case "d":
		m.Screen = ScreenDiff
	case "q":
		return m, quitCmd()
	}
	return m, nil
}

func renderSummary(m Model) string {
	pr := m.Result.PRData
	var b strings.Builder
	b.WriteString(m.Styles.Header.Render(" zero-review ") + "\n\n")
	b.WriteString(m.Styles.Title.Render(pr.Title) + "\n")
	b.WriteString(m.Styles.Subtitle.Render(fmt.Sprintf("#%d by %s · %s → %s · %d files (+ %d / - %d)",
		pr.PRNumber, pr.Author, pr.Branch, pr.Base, pr.ChangedFiles, pr.Additions, pr.Deletions)) + "\n\n")

	risk := reviewRiskScore(m.Result.Findings)
	riskStyle := m.Styles.Highlight
	if risk > 50 {
		riskStyle = m.Styles.Danger
	} else if risk > 25 {
		riskStyle = m.Styles.Warning
	}
	b.WriteString(m.Styles.Normal.Render("Risk score: ") + riskStyle.Render(fmt.Sprintf("%d/100", risk)) + "\n\n")

	b.WriteString(m.Styles.Panel.Render(
		m.Styles.Title.Render("Actions")+"\n"+
			"[Enter / s] Start review\n"+
			"[d] View diff\n"+
			"[q] Exit",
	) + "\n")

	b.WriteString("\n" + m.Styles.Help.Render(" AI: "+m.Result.Provider+" / "+m.Result.Model+" · Findings: "+itoa(len(m.Result.Findings))))
	return b.String()
}
