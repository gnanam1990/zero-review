package review

import (
	"fmt"
	"time"

	"github.com/gnanam1990/zero-review/internal/config"
)

// ToSession converts the engine Result into a UI-facing ReviewSession.
func (r Result) ToSession(opts config.Options) ReviewSession {
	mode := "balanced"
	switch {
	case opts.Strict:
		mode = "strict"
	case opts.SecurityOnly:
		mode = "security-only"
	case opts.TestsOnly:
		mode = "tests-only"
	}

	findings := make([]Finding, 0, len(r.Findings))
	for _, old := range r.Findings {
		findings = append(findings, Finding{
			ID:               old.ID,
			Title:            old.Title,
			Severity:         Severity(old.Severity),
			Category:         Category(old.Category),
			Confidence:       old.Confidence,
			FilePath:         old.FilePath,
			LineStart:        old.LineStart,
			LineEnd:          old.LineEnd,
			DiffSide:         old.DiffSide,
			Explanation:      old.Explanation,
			SuggestedComment: old.SuggestedComment,
			EditedComment:    old.EditedComment,
			SuggestedFix:     old.SuggestedFix,
			Evidence:         old.Evidence,
			Status:           FindingStatus(old.Status),
		})
	}

	return ReviewSession{
		PR: PRMetadata{
			Owner:        r.PRData.Owner,
			Repo:         r.PRData.Repo,
			Number:       r.PRData.PRNumber,
			Title:        r.PRData.Title,
			Body:         r.PRData.Body,
			Author:       r.PRData.Author,
			SourceBranch: r.PRData.Branch,
			TargetBranch: r.PRData.Base,
			ChangedFiles: r.PRData.ChangedFiles,
			Additions:    r.PRData.Additions,
			Deletions:    r.PRData.Deletions,
			URL:          fmt.Sprintf("https://github.com/%s/%s/pull/%d", r.PRData.Owner, r.PRData.Repo, r.PRData.PRNumber),
		},
		Findings:  findings,
		Summary:   "",
		RiskScore: RiskScore(findings),
		Provider:  r.Provider,
		Model:     r.Model,
		Mode:      mode,
		Posted:    false,
		CreatedAt: time.Now().UTC(),
	}
}
