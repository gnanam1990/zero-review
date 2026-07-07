package review

import "time"

// Severity levels for findings.
type Severity string

const (
	SeverityHigh   Severity = "high"
	SeverityMedium Severity = "medium"
	SeverityLow    Severity = "low"
	SeverityInfo   Severity = "info"
)

// Category labels for findings.
type Category string

const (
	CategorySecurity        Category = "security"
	CategoryBug             Category = "bug"
	CategoryTest            Category = "test"
	CategoryPerformance     Category = "performance"
	CategoryMaintainability Category = "maintainability"
	CategoryDocs            Category = "docs"
)

// FindingStatus tracks user disposition on a finding.
type FindingStatus string

const (
	FindingStatusPending  FindingStatus = "pending"
	FindingStatusApproved FindingStatus = "approved"
	FindingStatusRejected FindingStatus = "rejected"
	FindingStatusEdited   FindingStatus = "edited"
	FindingStatusPosted   FindingStatus = "posted"
)

// PRMetadata is the canonical PR summary used by the UI.
type PRMetadata struct {
	Owner        string
	Repo         string
	Number       int
	Title        string
	Body         string
	Author       string
	SourceBranch string
	TargetBranch string
	ChangedFiles int
	Additions    int
	Deletions    int
	URL          string
}

// Finding is one AI-generated review item validated against the PR diff.
type Finding struct {
	ID               string
	Title            string
	Severity         Severity
	Category         Category
	Confidence       int
	FilePath         string
	LineStart        int
	LineEnd          int
	DiffSide         string
	Explanation      string
	SuggestedComment string
	EditedComment    string
	SuggestedFix     string
	Evidence         string
	Status           FindingStatus
}

// IsPosted reports whether a finding has already been posted to GitHub.
func (f Finding) IsPosted() bool { return f.Status == FindingStatusPosted }

// IsApproved reports whether a finding should be posted.
func (f Finding) IsApproved() bool {
	return f.Status == FindingStatusApproved || f.Status == FindingStatusEdited || f.Status == FindingStatusPosted
}

// PostComment returns the text to post for this finding.
func (f Finding) PostComment() string {
	if f.Status == FindingStatusEdited && f.EditedComment != "" {
		return f.EditedComment
	}
	return f.SuggestedComment
}

// ReviewSession is the UI-facing snapshot of a review run.
type ReviewSession struct {
	PR         PRMetadata
	Findings   []Finding
	Summary    string
	RiskScore  int
	Provider   string
	Model      string
	Mode       string
	ReportPath string
	Posted     bool
	CreatedAt  time.Time
}

// Approved returns findings the user has approved, edited, or already posted.
func (s ReviewSession) Approved() []Finding {
	var out []Finding
	for _, f := range s.Findings {
		if f.IsApproved() {
			out = append(out, f)
		}
	}
	return out
}

// Rejected returns findings the user has rejected.
func (s ReviewSession) Rejected() []Finding {
	var out []Finding
	for _, f := range s.Findings {
		if f.Status == FindingStatusRejected {
			out = append(out, f)
		}
	}
	return out
}

// Pending returns findings awaiting disposition.
func (s ReviewSession) Pending() []Finding {
	var out []Finding
	for _, f := range s.Findings {
		if f.Status == FindingStatusPending {
			out = append(out, f)
		}
	}
	return out
}

// Edited returns findings that have been edited but not yet posted.
func (s ReviewSession) Edited() []Finding {
	var out []Finding
	for _, f := range s.Findings {
		if f.Status == FindingStatusEdited {
			out = append(out, f)
		}
	}
	return out
}

// CountBySeverity returns the number of findings at each severity.
func (s ReviewSession) CountBySeverity() map[Severity]int {
	m := map[Severity]int{SeverityHigh: 0, SeverityMedium: 0, SeverityLow: 0, SeverityInfo: 0}
	for _, f := range s.Findings {
		m[f.Severity]++
	}
	return m
}
