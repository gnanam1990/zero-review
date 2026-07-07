package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gnanam1990/zero-review/internal/review"
)

// SaveMarkdown writes a markdown report for the review result.
func SaveMarkdown(result review.Result, dir string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil { //nolint:gosec
		return "", err
	}
	ts := time.Now().UTC().Format("20060102-150405")
	path := filepath.Join(dir, fmt.Sprintf("pr-%d-%s.md", result.PRData.PRNumber, ts))

	var b strings.Builder
	fmt.Fprintf(&b, "# Review Report: %s\n\n", result.PRData.Title)
	fmt.Fprintf(&b, "- **Repository:** %s/%s\n", result.PRData.Owner, result.PRData.Repo)
	fmt.Fprintf(&b, "- **PR:** #%d\n", result.PRData.PRNumber)
	fmt.Fprintf(&b, "- **Author:** %s\n", result.PRData.Author)
	fmt.Fprintf(&b, "- **Branch:** %s -> %s\n", result.PRData.Branch, result.PRData.Base)
	fmt.Fprintf(&b, "- **Changed files:** %d (+ %d / - %d)\n", result.PRData.ChangedFiles, result.PRData.Additions, result.PRData.Deletions)
	fmt.Fprintf(&b, "- **Provider:** %s / %s\n", result.Provider, result.Model)
	fmt.Fprintf(&b, "- **Risk score:** %d/100\n", review.RiskScore(result.Findings))
	fmt.Fprintf(&b, "- **Timestamp:** %s\n\n", time.Now().UTC().Format(time.RFC3339))

	b.WriteString("## Summary\n\n")
	fmt.Fprintf(&b, "Total findings: %d | Approved: %d | Rejected: %d | Edited: %d\n\n",
		len(result.Findings), countStatus(result.Findings, review.StatusApproved), countStatus(result.Findings, review.StatusRejected), countStatus(result.Findings, review.StatusEdited))

	writeSection(&b, "Approved Findings", filterStatus(result.Findings, review.StatusApproved, review.StatusEdited))
	writeSection(&b, "Rejected Findings", filterStatus(result.Findings, review.StatusRejected))

	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil { //nolint:gosec
		return "", err
	}
	return path, nil
}

func writeSection(b *strings.Builder, title string, findings []review.Finding) {
	fmt.Fprintf(b, "## %s\n\n", title)
	if len(findings) == 0 {
		b.WriteString("None.\n\n")
		return
	}
	for _, f := range findings {
		fmt.Fprintf(b, "### %s (%s, %d%% confidence)\n", f.Title, f.Severity, f.Confidence)
		fmt.Fprintf(b, "- **File:** `%s:%d-%d`\n", f.FilePath, f.LineStart, f.LineEnd)
		fmt.Fprintf(b, "- **Category:** %s\n", f.Category)
		fmt.Fprintf(b, "- **Status:** %s\n", f.Status)
		fmt.Fprintf(b, "- **Comment:** %s\n", f.PostComment())
		if f.SuggestedFix != "" {
			fmt.Fprintf(b, "- **Suggested fix:** %s\n", f.SuggestedFix)
		}
		if f.Evidence != "" {
			fmt.Fprintf(b, "- **Evidence:** %s\n", f.Evidence)
		}
		b.WriteString("\n")
	}
}

func countStatus(findings []review.Finding, status review.Status) int {
	c := 0
	for _, f := range findings {
		if f.Status == status {
			c++
		}
	}
	return c
}

func filterStatus(findings []review.Finding, statuses ...review.Status) []review.Finding {
	var out []review.Finding
	for _, f := range findings {
		for _, s := range statuses {
			if f.Status == s {
				out = append(out, f)
				break
			}
		}
	}
	return out
}
