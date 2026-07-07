package review

import "github.com/gnanam1990/zero-review/internal/config"

// ApplyFilters applies optional security/test filters after validation.
func ApplyFilters(findings []Finding, opts config.Options) []Finding {
	if !opts.SecurityOnly && !opts.TestsOnly {
		return findings
	}
	var out []Finding
	for _, f := range findings {
		if opts.SecurityOnly && f.Category != CategorySecurity {
			continue
		}
		if opts.TestsOnly && f.Category != CategoryTest {
			continue
		}
		out = append(out, f)
	}
	return out
}

// RiskScore computes a simple numeric risk score from approved findings.
func RiskScore(findings []Finding) int {
	score := 0
	for _, f := range findings {
		if !f.IsApproved() {
			continue
		}
		switch f.Severity {
		case SeverityHigh:
			score += 25
		case SeverityMedium:
			score += 15
		case SeverityLow:
			score += 5
		default:
			score += 1
		}
	}
	if score > 100 {
		return 100
	}
	return score
}
