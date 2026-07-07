package review

import (
	"testing"

	"github.com/gnanam1990/zero-review/internal/ai"
	"github.com/gnanam1990/zero-review/internal/github"
)

func TestValidateAll(t *testing.T) {
	data := PRData{
		Files: []github.PRFile{
			{Filename: "cmd/app/main.go", Additions: 10},
		},
		Diff: `diff --git a/cmd/app/main.go b/cmd/app/main.go
new file mode 100644
+++ b/cmd/app/main.go
@@ -0,0 +1,10 @@
+package main
`,
	}

	raw := []Finding{
		{
			ID:         "F001",
			Title:      "Nil pointer",
			Severity:   SeverityHigh,
			Confidence: 90,
			Category:   CategoryBug,
			FilePath:   "cmd/app/main.go",
			LineStart:  5,
			LineEnd:    5,
			DiffSide:   "RIGHT",
		},
		{
			ID:         "F002",
			Title:      "Unknown file",
			Severity:   SeverityMedium,
			Confidence: 80,
			Category:   CategoryBug,
			FilePath:   "does/not/exist.go",
			LineStart:  1,
		},
		{
			ID:         "F003",
			Title:      "Low confidence",
			Severity:   SeverityLow,
			Confidence: 30,
			Category:   CategoryMaintainability,
			FilePath:   "cmd/app/main.go",
			LineStart:  3,
		},
		{
			ID:        "F004",
			Title:     "Bad severity",
			Severity:  Severity("critical"),
			Category:  CategoryBug,
			FilePath:  "cmd/app/main.go",
			LineStart: 1,
		},
	}

	got := ValidateAll(raw, data, 75)
	if len(got) != 1 {
		t.Fatalf("expected 1 validated finding, got %d", len(got))
	}
	if got[0].ID != "F001" {
		t.Fatalf("expected F001, got %s", got[0].ID)
	}
	if got[0].Status != StatusPending {
		t.Fatalf("expected status pending, got %s", got[0].Status)
	}

	// A finding outside the changed hunk range should be rejected.
	outOfRange := []Finding{raw[0]}
	outOfRange[0].LineStart = 999
	if len(ValidateAll(outOfRange, data, 75)) != 0 {
		t.Fatalf("expected out-of-range finding to be rejected")
	}
}

func TestFromProviderFindings(t *testing.T) {
	raw := []ai.ProviderFinding{
		{
			ID:       "A1",
			Title:    "Issue",
			Severity: "medium",
			Category: "test",
		},
	}
	findings := FromProviderFindings(raw)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != SeverityMedium {
		t.Fatalf("expected severity medium, got %s", findings[0].Severity)
	}
}
