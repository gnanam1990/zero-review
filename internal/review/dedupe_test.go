package review

import "testing"

func TestDeduplicate(t *testing.T) {
	findings := []Finding{
		{
			ID:         "F001",
			Title:      "Nil pointer dereference",
			Severity:   SeverityHigh,
			Confidence: 80,
			Category:   CategoryBug,
			FilePath:   "cmd/app/main.go",
			LineStart:  42,
		},
		{
			ID:         "F002",
			Title:      "Nil pointer dereference",
			Severity:   SeverityHigh,
			Confidence: 95,
			Category:   CategoryBug,
			FilePath:   "cmd/app/main.go",
			LineStart:  42,
		},
		{
			ID:         "F003",
			Title:      "Missing error handling",
			Severity:   SeverityMedium,
			Confidence: 85,
			Category:   CategoryMaintainability,
			FilePath:   "internal/review/engine.go",
			LineStart:  88,
		},
	}

	got := Deduplicate(findings)
	if len(got) != 2 {
		t.Fatalf("expected 2 findings after dedupe, got %d", len(got))
	}
	if got[0].ID != "F002" {
		t.Fatalf("expected higher-confidence F002 to win, got %s", got[0].ID)
	}
}
