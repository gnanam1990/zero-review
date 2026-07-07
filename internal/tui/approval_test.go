package tui

import (
	"testing"

	"github.com/gnanam1990/zero-review/internal/config"
	"github.com/gnanam1990/zero-review/internal/github"
	"github.com/gnanam1990/zero-review/internal/review"
)

func TestApprovalFlow(t *testing.T) {
	gh := github.NewFakeClient()
	result := review.Result{
		PRData: review.PRData{
			Owner:    "gnanam1990",
			Repo:     "zero-review",
			PRNumber: 7,
			Title:    "Add auth middleware",
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
				SuggestedComment: "Check nil.",
				Status:           review.StatusPending,
			},
			{
				ID:         "F002",
				Title:      "Missing error handling",
				Severity:   review.SeverityMedium,
				Confidence: 80,
				Category:   review.CategoryMaintainability,
				FilePath:   "internal/review/engine.go",
				LineStart:  88,
				Status:     review.StatusPending,
			},
		},
	}
	m := NewModel(result, gh, config.Options{})

	// Approve all.
	m2, _ := handleApprovalKey(m, "a")
	mm := m2.(Model)
	if len(review.Approved(mm.Findings)) != 2 {
		t.Fatalf("expected 2 approved findings, got %d", len(review.Approved(mm.Findings)))
	}

	// Confirm post.
	m3, _ := handleApprovalKey(mm, "enter")
	mm2 := m3.(Model)
	if mm2.Posted != 2 {
		t.Fatalf("expected 2 posted comments, got %d", mm2.Posted)
	}
	if gh.LastReview == nil {
		t.Fatalf("expected a review submission")
	}
	if len(gh.LastReview.Comments) != 2 {
		t.Fatalf("expected 2 review comments, got %d", len(gh.LastReview.Comments))
	}
}

func TestNoPost(t *testing.T) {
	gh := github.NewFakeClient()
	result := review.Result{
		PRData: review.PRData{
			Owner:    "gnanam1990",
			Repo:     "zero-review",
			PRNumber: 7,
			Title:    "Add auth middleware",
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
				SuggestedComment: "Check nil.",
				Status:           review.StatusApproved,
			},
		},
	}
	m := NewModel(result, gh, config.Options{NoPost: true})
	m2, _ := handleApprovalKey(m, "enter")
	mm := m2.(Model)
	if mm.Posted != 0 {
		t.Fatalf("expected 0 posted comments with NoPost, got %d", mm.Posted)
	}
	if gh.LastReview != nil {
		t.Fatalf("expected no review submission with NoPost")
	}
}
