package review

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

// Finding is one AI-generated review item validated against the PR diff.
type Finding struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Severity         Severity `json:"severity"`
	Confidence       int      `json:"confidence"`
	Category         Category `json:"category"`
	FilePath         string   `json:"file_path"`
	LineStart        int      `json:"line_start"`
	LineEnd          int      `json:"line_end"`
	DiffSide         string   `json:"diff_side"`
	Explanation      string   `json:"explanation"`
	SuggestedComment string   `json:"suggested_comment"`
	SuggestedFix     string   `json:"suggested_fix"`
	Evidence         string   `json:"evidence"`
	Status           Status   `json:"status"`
	EditedComment    string   `json:"edited_comment,omitempty"`
}

// Status tracks user disposition on a finding.
type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
	StatusEdited   Status = "edited"
)

// IsApproved reports whether a finding should be posted.
func (f Finding) IsApproved() bool {
	return f.Status == StatusApproved || (f.Status == StatusEdited && f.EditedComment != "")
}

// PostComment returns the text to post for this finding.
func (f Finding) PostComment() string {
	if f.Status == StatusEdited && f.EditedComment != "" {
		return f.EditedComment
	}
	return f.SuggestedComment
}

// Approved returns findings the user has approved or edited.
func Approved(findings []Finding) []Finding {
	var out []Finding
	for _, f := range findings {
		if f.IsApproved() {
			out = append(out, f)
		}
	}
	return out
}

// Rejected returns findings the user has rejected.
func Rejected(findings []Finding) []Finding {
	var out []Finding
	for _, f := range findings {
		if f.Status == StatusRejected {
			out = append(out, f)
		}
	}
	return out
}

// Pending returns findings awaiting disposition.
func Pending(findings []Finding) []Finding {
	var out []Finding
	for _, f := range findings {
		if f.Status == StatusPending {
			out = append(out, f)
		}
	}
	return out
}
