package ai

import "context"

// FakeProvider returns canned findings for tests and demos.
type FakeProvider struct {
	Findings  []ProviderFinding
	ModelName string
}

// NewFakeProvider creates a provider that returns two canned findings.
func NewFakeProvider() *FakeProvider {
	return &FakeProvider{
		ModelName: "fake",
		Findings: []ProviderFinding{
			{
				ID:               "F001",
				Title:            "Possible nil pointer dereference",
				Severity:         "high",
				Confidence:       90,
				Category:         "bug",
				FilePath:         "cmd/app/main.go",
				LineStart:        42,
				LineEnd:          42,
				DiffSide:         "RIGHT",
				Explanation:      "The result of getClient() is used without a nil check.",
				SuggestedComment: "Consider checking if `client` is nil before calling methods on it.",
				SuggestedFix:     "if client == nil { return err }",
				Evidence:         "client.Do(req) on line 42 where client comes from getClient().",
			},
			{
				ID:               "F002",
				Title:            "Missing error handling",
				Severity:         "medium",
				Confidence:       80,
				Category:         "maintainability",
				FilePath:         "internal/review/engine.go",
				LineStart:        88,
				LineEnd:          88,
				DiffSide:         "RIGHT",
				Explanation:      "The error from process() is ignored.",
				SuggestedComment: "Please handle the returned error or explicitly ignore it with `_`.",
				SuggestedFix:     "_, err := process()",
				Evidence:         "`_ = process(...)` on line 88.",
			},
		},
	}
}

// Review returns the canned findings.
func (f *FakeProvider) Review(ctx context.Context, prompt string) ([]ProviderFinding, error) {
	return append([]ProviderFinding(nil), f.Findings...), nil
}

// Model returns the fake model name.
func (f *FakeProvider) Model() string {
	if f.ModelName != "" {
		return f.ModelName
	}
	return "fake"
}

// Name returns the provider name.
func (f *FakeProvider) Name() string {
	return "fake"
}
