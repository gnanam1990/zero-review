package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gnanam1990/zero-review/internal/tui/components"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
	"github.com/gnanam1990/zero-review/internal/tui/screens"
)

// View renders the entire TUI for the current model state.
func (m Model) View() string {
	if m.Layout.TooSmall {
		return core.TerminalTooSmallMessage(m.Theme, m.Layout.Width, m.Layout.Height)
	}

	var b strings.Builder

	// Header
	header := components.Header(m.Theme, m.Screen, m.Session, m.Layout.Width, m.Layout.ShowSidebar)
	b.WriteString(header + "\n")

	// Sidebar (wide layout)
	if m.Layout.ShowSidebar {
		sidebar := components.Sidebar(m.Theme, m.Screen, m.Layout.ContentHeight)
		b.WriteString(sidebar)
	}

	// Content area
	content := m.renderContent()
	b.WriteString(content)

	// Footer
	footer := components.Footer(m.Theme, m.Keys, m.Screen, m.Layout)
	b.WriteString(footer)

	// Overlays
	out := b.String()
	if m.ShowHelp {
		out = overlay(out, components.HelpOverlay(m.Theme, m.Screen, m.Layout.Width, m.Layout.Height))
	}
	if m.Confirm != nil {
		out = overlay(out, components.ConfirmModal(m.Theme, m.Confirm.Title, m.Confirm.Body, m.Confirm.YesText, m.Confirm.NoText, 50, m.Layout.Height))
	}
	if m.Toast != nil {
		out = overlayBottom(out, components.Toast(m.Theme, m.Toast, m.Layout.Width))
	}
	if m.CommandOpen {
		out = overlayBottom(out, components.CommandBar(m.Theme, "", m.Layout.Width-4))
	}

	return out
}

func (m Model) renderContent() string {
	w := m.Layout.ContentWidth
	h := m.Layout.ContentHeight

	switch m.Screen {
	case core.ScreenWelcome:
		return screens.Welcome(m.Theme, w, h)
	case core.ScreenPRInput:
		return screens.PRInput(m.Theme, m.PRForm, w, h)
	case core.ScreenLoadingReview:
		url := ""
		if m.Session != nil {
			url = m.Session.PR.URL
		}
		return screens.LoadingReview(m.Theme, url, m.LoadingSteps, m.LoadingTip, w, h)
	case core.ScreenDashboard:
		return screens.Dashboard(m.Theme, m.Session, w, h)
	case core.ScreenFindings:
		return screens.Findings(m.Theme, m.Session, m.FindingsCursor, m.SelectedFinding, m.SeverityFilter, w, h)
	case core.ScreenFindingDetail:
		return screens.FindingDetail(m.Theme, m.SelectedFindingPtr(), w, h)
	case core.ScreenDiff:
		return screens.Diff(m.Theme, m.Session, m.SelectedFindingPtr(), w, h)
	case core.ScreenChat:
		return screens.Chat(m.Theme, m.Session, m.ChatMessages, m.SelectedFindingPtr(), w, h)
	case core.ScreenApproval:
		mode := screens.PostModeComment
		if m.NoPost {
			mode = screens.PostModeReportOnly
		}
		return screens.Approval(m.Theme, m.Session, mode, m.NoPost, w, h)
	case core.ScreenReport:
		return screens.Report(m.Theme, m.Session, m.Session.ReportPath, 0, w, h)
	case core.ScreenSettings:
		return screens.Settings(m.Theme, m.SettingsForm, w, h)
	default:
		return components.EmptyState(m.Theme, "?", "Unknown screen", "Press esc to go back.", w, h)
	}
}

func overlay(background, foreground string) string {
	return lipgloss.JoinVertical(lipgloss.Left, foreground)
}

func overlayBottom(background, foreground string) string {
	return background + "\n" + foreground
}
