package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gnanam1990/zero-review/internal/review"
)

func TestSaveMarkdown(t *testing.T) {
	dir := t.TempDir()
	result := review.Result{
		PRData: review.PRData{
			Owner:    "gnanam1990",
			Repo:     "zero-review",
			PRNumber: 7,
			Title:    "Add auth middleware",
			Author:   "alice",
			Branch:   "feature/auth",
			Base:     "main",
		},
		Findings: []review.Finding{
			{
				ID:               "F001",
				Title:            "Nil pointer",
				Severity:         review.SeverityHigh,
				Confidence:       90,
				Category:         review.CategoryBug,
				FilePath:         "cmd/app/main.go",
				LineStart:        42,
				LineEnd:          42,
				Status:           review.FindingStatusApproved,
				SuggestedComment: "Check nil before use.",
			},
		},
		Provider: "fake",
		Model:    "fake",
	}

	path, err := SaveMarkdown(result, dir)
	if err != nil {
		t.Fatalf("SaveMarkdown error: %v", err)
	}
	if !strings.HasPrefix(filepath.Base(path), "pr-7-") {
		t.Fatalf("unexpected report filename: %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Add auth middleware") {
		t.Fatalf("report missing PR title")
	}
	if !strings.Contains(content, "Nil pointer") {
		t.Fatalf("report missing finding title")
	}
	if !strings.Contains(content, "Check nil before use.") {
		t.Fatalf("report missing approved comment")
	}
}
