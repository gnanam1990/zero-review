package tui

import "github.com/charmbracelet/lipgloss"

// Styles holds the design tokens for the zero-review TUI.
type Styles struct {
	Header     lipgloss.Style
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Normal     lipgloss.Style
	Muted      lipgloss.Style
	Selected   lipgloss.Style
	Highlight  lipgloss.Style
	Success    lipgloss.Style
	Warning    lipgloss.Style
	Danger     lipgloss.Style
	Info       lipgloss.Style
	Help       lipgloss.Style
	Panel      lipgloss.Style
	ChatUser   lipgloss.Style
	ChatSystem lipgloss.Style
	Code       lipgloss.Style
}

// NewStyles returns the default theme.
func NewStyles() Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")),
		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0")),
		Normal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")),
		Muted: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C6C7D")),
		Selected: lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FAFAFA")),
		Highlight: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2CDA9D")),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2CDA9D")),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7B801")),
		Danger: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")),
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#58A6FF")),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C6C7D")),
		Panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#2A2A3C")).
			Padding(1, 2),
		ChatUser: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#2A2A3C")).
			Padding(0, 1).
			MarginLeft(2),
		ChatSystem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#1E1E2E")).
			Padding(0, 1).
			MarginRight(2),
		Code: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#1E1E2E")),
	}
}

func severityColor(sev string) lipgloss.Color {
	switch sev {
	case "high":
		return lipgloss.Color("#FF6B6B")
	case "medium":
		return lipgloss.Color("#F7B801")
	case "low":
		return lipgloss.Color("#58A6FF")
	default:
		return lipgloss.Color("#A0A0A0")
	}
}
