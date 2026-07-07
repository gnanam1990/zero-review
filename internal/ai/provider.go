package ai

import "context"

// Provider is the AI review backend.
type Provider interface {
	Review(ctx context.Context, prompt string) ([]ProviderFinding, error)
	Model() string
	Name() string
}

// ProviderFinding matches the JSON shape the AI returns.
type ProviderFinding struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Severity         string `json:"severity"`
	Confidence       int    `json:"confidence"`
	Category         string `json:"category"`
	FilePath         string `json:"file_path"`
	LineStart        int    `json:"line_start"`
	LineEnd          int    `json:"line_end"`
	DiffSide         string `json:"diff_side"`
	Explanation      string `json:"explanation"`
	SuggestedComment string `json:"suggested_comment"`
	SuggestedFix     string `json:"suggested_fix"`
	Evidence         string `json:"evidence"`
}
