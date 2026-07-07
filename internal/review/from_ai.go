package review

import "github.com/gnanam1990/zero-review/internal/ai"

// FromProviderFindings converts AI provider findings into the review model.
func FromProviderFindings(raw []ai.ProviderFinding) []Finding {
	out := make([]Finding, 0, len(raw))
	for _, r := range raw {
		out = append(out, Finding{
			ID:               r.ID,
			Title:            r.Title,
			Severity:         Severity(r.Severity),
			Confidence:       r.Confidence,
			Category:         Category(r.Category),
			FilePath:         r.FilePath,
			LineStart:        r.LineStart,
			LineEnd:          r.LineEnd,
			DiffSide:         r.DiffSide,
			Explanation:      r.Explanation,
			SuggestedComment: r.SuggestedComment,
			SuggestedFix:     r.SuggestedFix,
			Evidence:         r.Evidence,
			Status:           StatusPending,
		})
	}
	return out
}
