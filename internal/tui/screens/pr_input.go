package screens

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// PRInputValues holds the form state for a new review.
type PRInputValues struct {
	URL        string
	Provider   string
	Mode       string
	SaveReport bool
	NoPost     bool
}

// BuildPRForm creates the new-review Huh form.
func BuildPRForm(values *PRInputValues) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("GitHub PR URL").
				Placeholder("https://github.com/org/repo/pull/123").
				Value(&values.URL).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("URL is required")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Provider").
				Options(
					huh.NewOption("Fake (demo)", "fake"),
					huh.NewOption("Kimi", "kimi"),
					huh.NewOption("Ollama", "ollama"),
					huh.NewOption("Anthropic", "anthropic"),
					huh.NewOption("OpenAI", "openai"),
				).
				Value(&values.Provider),
			huh.NewSelect[string]().
				Title("Review mode").
				Options(
					huh.NewOption("Balanced", "balanced"),
					huh.NewOption("Strict", "strict"),
					huh.NewOption("Security-only", "security-only"),
					huh.NewOption("Tests-only", "tests-only"),
				).
				Value(&values.Mode),
			huh.NewConfirm().
				Title("Save report automatically?").
				Value(&values.SaveReport).
				Affirmative("Yes").
				Negative("No"),
			huh.NewConfirm().
				Title("Report only, disable posting?").
				Value(&values.NoPost).
				Affirmative("Yes").
				Negative("No"),
		),
	).WithWidth(60).WithShowHelp(true)
}

// PRInput renders the PR input form inside a panel.
func PRInput(theme *core.Theme, form *huh.Form, width, height int) string {
	if form == nil {
		return theme.PanelStyle.Width(width).Height(height).Render(theme.MutedText.Render("Form not initialized."))
	}
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(form.View())
}
