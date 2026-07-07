package core

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme holds the complete design token palette for Zero Review.
// All components receive a *Theme and never hardcode colors directly.
type Theme struct {
	// Base
	AppBackground    lipgloss.Style
	TerminalTooSmall lipgloss.Style

	// Header
	HeaderStyle      lipgloss.Style
	HeaderMutedStyle lipgloss.Style
	HeaderAlertStyle lipgloss.Style

	// Sidebar
	SidebarStyle       lipgloss.Style
	SidebarActiveStyle lipgloss.Style
	SidebarItemStyle   lipgloss.Style

	// Panels
	PanelStyle      lipgloss.Style
	PanelTitleStyle lipgloss.Style

	// Typography
	PrimaryText lipgloss.Style
	MutedText   lipgloss.Style
	BoldText    lipgloss.Style
	ErrorText   lipgloss.Style

	// Badges
	SuccessBadge lipgloss.Style
	WarningBadge lipgloss.Style
	DangerBadge  lipgloss.Style
	InfoBadge    lipgloss.Style
	NeutralBadge lipgloss.Style

	// Borders
	BorderNormal  lipgloss.Style
	BorderFocused lipgloss.Style

	// Buttons
	ButtonPrimary   lipgloss.Style
	ButtonSecondary lipgloss.Style
	ButtonDanger    lipgloss.Style

	// Code / Diff
	CodeBlockStyle   lipgloss.Style
	DiffAddedStyle   lipgloss.Style
	DiffRemovedStyle lipgloss.Style
	DiffContextStyle lipgloss.Style
	DiffHunkStyle    lipgloss.Style
	DiffMarkerStyle  lipgloss.Style

	// Toasts
	ToastSuccess lipgloss.Style
	ToastError   lipgloss.Style
	ToastInfo    lipgloss.Style

	// Form / Input
	InputStyle      lipgloss.Style
	InputLabelStyle lipgloss.Style
	HelpStyle       lipgloss.Style
	FooterStyle     lipgloss.Style
	ShortcutKey     lipgloss.Style

	// Chat
	ChatUserBubble lipgloss.Style
	ChatAIBubble   lipgloss.Style
	ChatMetaText   lipgloss.Style

	// Tables
	TableHeaderStyle   lipgloss.Style
	TableCellStyle     lipgloss.Style
	TableSelectedStyle lipgloss.Style
}

// NewTheme builds the default theme using adaptive colors.
func NewTheme() *Theme {
	// Colors
	primary := lipgloss.Color("#7D56F4")
	primaryLight := lipgloss.Color("#A78BFA")
	bgDark := lipgloss.Color("#1E1E2E")
	bgLight := lipgloss.Color("#F5F5F7")
	textDark := lipgloss.Color("#FAFAFA")
	textLight := lipgloss.Color("#1D1D1F")
	mutedDark := lipgloss.Color("#6C6C7D")
	mutedLight := lipgloss.Color("#6E6E73")
	panelDark := lipgloss.Color("#252536")
	panelLight := lipgloss.Color("#FFFFFF")
	sidebarDark := lipgloss.Color("#1A1A2A")
	sidebarLight := lipgloss.Color("#EBEBF0")

	adaptiveText := lipgloss.AdaptiveColor{Dark: string(textDark), Light: string(textLight)}
	adaptiveMuted := lipgloss.AdaptiveColor{Dark: string(mutedDark), Light: string(mutedLight)}
	adaptivePanel := lipgloss.AdaptiveColor{Dark: string(panelDark), Light: string(panelLight)}
	adaptiveSidebar := lipgloss.AdaptiveColor{Dark: string(sidebarDark), Light: string(sidebarLight)}

	t := &Theme{}

	t.AppBackground = lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Dark: string(bgDark), Light: string(bgLight)})

	t.TerminalTooSmall = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.AdaptiveColor{Dark: string(textDark), Light: string(textLight)}).
		Background(lipgloss.AdaptiveColor{Dark: string(bgDark), Light: string(bgLight)}).
		Padding(2).
		Align(lipgloss.Center)

	t.HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(primary).
		Padding(0, 1)

	t.HeaderMutedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D0D0E0")).
		Background(primary).
		Padding(0, 1)

	t.HeaderAlertStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B6B")).
		Background(primary).
		Padding(0, 1)

	t.SidebarStyle = lipgloss.NewStyle().
		Background(adaptiveSidebar).
		Foreground(adaptiveText).
		Padding(1, 0)

	t.SidebarItemStyle = lipgloss.NewStyle().
		Padding(0, 1).
		MarginLeft(1).
		MarginRight(1).
		Foreground(adaptiveText)

	t.SidebarActiveStyle = lipgloss.NewStyle().
		Padding(0, 1).
		MarginLeft(1).
		MarginRight(1).
		Bold(true).
		Background(primary).
		Foreground(lipgloss.Color("#FAFAFA"))

	t.PanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Dark: "#2A2A3C", Light: "#D1D1D6"}).
		Background(adaptivePanel).
		Foreground(adaptiveText).
		Padding(1, 2)

	t.PanelTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryLight).
		MarginBottom(1)

	t.PrimaryText = lipgloss.NewStyle().Foreground(adaptiveText)
	t.MutedText = lipgloss.NewStyle().Foreground(adaptiveMuted)
	t.BoldText = lipgloss.NewStyle().Bold(true).Foreground(adaptiveText)
	t.ErrorText = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))

	badge := func(fg, bg string) lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(fg)).
			Background(lipgloss.Color(bg)).
			Padding(0, 1).
			Bold(true)
	}
	t.SuccessBadge = badge("#11111B", "#2CDA9D")
	t.WarningBadge = badge("#11111B", "#F7B801")
	t.DangerBadge = badge("#11111B", "#FF6B6B")
	t.InfoBadge = badge("#11111B", "#58A6FF")
	t.NeutralBadge = badge("#FAFAFA", "#4A4A5A")

	t.BorderNormal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Dark: "#2A2A3C", Light: "#D1D1D6"})

	t.BorderFocused = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary)

	button := func(fg, bg string) lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(fg)).
			Background(lipgloss.Color(bg)).
			Padding(0, 2).
			Bold(true)
	}
	t.ButtonPrimary = button("#FAFAFA", "#7D56F4")
	t.ButtonSecondary = button("#FAFAFA", "#4A4A5A")
	t.ButtonDanger = button("#FAFAFA", "#FF6B6B")

	t.CodeBlockStyle = lipgloss.NewStyle().
		Foreground(adaptiveText).
		Background(lipgloss.AdaptiveColor{Dark: "#16161E", Light: "#F0F0F5"}).
		Padding(1, 2)

	t.DiffAddedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2CDA9D"))
	t.DiffRemovedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
	t.DiffContextStyle = lipgloss.NewStyle().Foreground(adaptiveMuted)
	t.DiffHunkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F7B801")).Bold(true)
	t.DiffMarkerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6C6C7D"))

	toast := func(fg, bg string) lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(fg)).
			Background(lipgloss.Color(bg)).
			Padding(0, 2).
			Bold(true)
	}
	t.ToastSuccess = toast("#11111B", "#2CDA9D")
	t.ToastError = toast("#FAFAFA", "#FF6B6B")
	t.ToastInfo = toast("#11111B", "#58A6FF")

	t.InputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Dark: "#4A4A5A", Light: "#A1A1AA"}).
		Padding(0, 1).
		Width(50)

	t.InputLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(adaptiveText).
		MarginBottom(1)

	t.HelpStyle = lipgloss.NewStyle().Foreground(adaptiveMuted)
	t.FooterStyle = lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Dark: "#1A1A2A", Light: "#E5E5EA"}).
		Foreground(adaptiveText).
		Padding(0, 1)

	t.ShortcutKey = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryLight).
		Background(lipgloss.AdaptiveColor{Dark: "#2A2A3C", Light: "#E0E0E8"}).
		Padding(0, 1)

	t.ChatUserBubble = lipgloss.NewStyle().
		Foreground(adaptiveText).
		Background(lipgloss.AdaptiveColor{Dark: "#2A2A3C", Light: "#E0E0E8"}).
		Padding(0, 1).
		MarginLeft(2)

	t.ChatAIBubble = lipgloss.NewStyle().
		Foreground(adaptiveText).
		Background(lipgloss.AdaptiveColor{Dark: "#252536", Light: "#FFFFFF"}).
		Padding(0, 1).
		MarginRight(2)

	t.ChatMetaText = lipgloss.NewStyle().Foreground(adaptiveMuted)

	t.TableHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.AdaptiveColor{Dark: "#2A2A3C", Light: "#6C6C7D"})

	t.TableCellStyle = lipgloss.NewStyle().Foreground(adaptiveText)
	t.TableSelectedStyle = lipgloss.NewStyle().
		Bold(true).
		Background(primary).
		Foreground(lipgloss.Color("#FAFAFA"))

	return t
}

// SeverityBadge returns the correct badge style for a severity string.
func (t *Theme) SeverityBadge(severity string) lipgloss.Style {
	switch severity {
	case "high":
		return t.DangerBadge
	case "medium":
		return t.WarningBadge
	case "low":
		return t.InfoBadge
	case "info":
		return t.NeutralBadge
	default:
		return t.NeutralBadge
	}
}

// StatusBadge returns the correct badge style for a finding status.
func (t *Theme) StatusBadge(status string) lipgloss.Style {
	switch status {
	case "approved", "posted":
		return t.SuccessBadge
	case "rejected":
		return t.DangerBadge
	case "edited":
		return t.WarningBadge
	default:
		return t.NeutralBadge
	}
}
