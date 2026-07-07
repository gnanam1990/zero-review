package screens

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// SettingsValues holds the settings form state.
type SettingsValues struct {
	Provider          string
	Mode              string
	Confidence        string
	AutoSave          bool
	Theme             string
	ShowLowConfidence bool
	DefaultPostMode   string
}

// BuildSettingsForm creates the settings Huh form.
func BuildSettingsForm(values *SettingsValues) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Default provider").
				Options(
					huh.NewOption("Fake", "fake"),
					huh.NewOption("Kimi", "kimi"),
					huh.NewOption("Ollama", "ollama"),
					huh.NewOption("Anthropic", "anthropic"),
					huh.NewOption("OpenAI", "openai"),
				).
				Value(&values.Provider),
			huh.NewSelect[string]().
				Title("Default review mode").
				Options(
					huh.NewOption("Balanced", "balanced"),
					huh.NewOption("Strict", "strict"),
					huh.NewOption("Security-only", "security-only"),
					huh.NewOption("Tests-only", "tests-only"),
				).
				Value(&values.Mode),
			huh.NewSelect[string]().
				Title("Confidence threshold").
				Options(
					huh.NewOption("50%", "50"),
					huh.NewOption("60%", "60"),
					huh.NewOption("70%", "70"),
					huh.NewOption("75%", "75"),
					huh.NewOption("80%", "80"),
					huh.NewOption("90%", "90"),
					huh.NewOption("95%", "95"),
				).
				Value(&values.Confidence),
			huh.NewConfirm().
				Title("Auto-save report?").
				Value(&values.AutoSave).
				Affirmative("Yes").
				Negative("No"),
			huh.NewSelect[string]().
				Title("Theme").
				Options(
					huh.NewOption("Auto", "auto"),
					huh.NewOption("Dark", "dark"),
					huh.NewOption("Light", "light"),
				).
				Value(&values.Theme),
			huh.NewConfirm().
				Title("Show low-confidence findings?").
				Value(&values.ShowLowConfidence).
				Affirmative("Yes").
				Negative("No"),
			huh.NewSelect[string]().
				Title("Default posting mode").
				Options(
					huh.NewOption("Comment only", "comment"),
					huh.NewOption("Request changes", "request_changes"),
					huh.NewOption("Approve PR", "approve"),
					huh.NewOption("Save report only", "report_only"),
				).
				Value(&values.DefaultPostMode),
		),
	).WithWidth(60).WithShowHelp(true)
}

// Settings renders the settings screen.
func Settings(theme *core.Theme, form *huh.Form, width, height int) string {
	if form == nil {
		return theme.PanelStyle.Width(width).Height(height).Render(theme.MutedText.Render("Form not initialized."))
	}
	// Huh forms manage their own height; top-align so the whole form is usable
	// even when it is taller than the panel.
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		Render(form.View())
}
