package tui

import (
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

	// Header
	header := components.Header(m.Theme, m.Screen, m.Session, m.Layout.Width, m.Layout.ShowSidebar)

	// Body: either the active screen or a full-body modal.
	body := m.renderContent()
	if m.Layout.ShowSidebar && !m.modalActive() {
		sidebar := components.Sidebar(m.Theme, m.Screen, m.Layout.ContentHeight)
		body = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, body)
	}
	if m.ShowHelp {
		body = components.HelpOverlay(m.Theme, m.Screen, m.Layout.Width, m.Layout.ContentHeight)
	} else if m.Confirm != nil {
		body = components.ConfirmModal(m.Theme, m.Confirm.Title, m.Confirm.Body, m.Confirm.YesText, m.Confirm.NoText, 50, m.Layout.ContentHeight)
	}

	// Footer, possibly with a toast notification.
	footer := m.renderFooter()

	// Command palette expands the footer area when open.
	bodyHeight := m.Layout.ContentHeight
	if m.CommandOpen {
		cmdBar := components.CommandBar(m.Theme, m.CommandInput.Value(), m.Layout.Width-4)
		suggestions := components.RenderSuggestions(m.Theme, m.CommandInput.Value(), m.Layout.Width-4)
		cmdBlock := lipgloss.JoinVertical(lipgloss.Left, cmdBar, suggestions)
		cmdH := lipgloss.Height(cmdBlock)
		bodyHeight = m.Layout.ContentHeight - cmdH
		if bodyHeight < 1 {
			bodyHeight = 1
		}
		body = lipgloss.NewStyle().Height(bodyHeight).Render(body)
		footer = lipgloss.JoinVertical(lipgloss.Left, cmdBlock, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// modalActive returns true when a screen-level modal should hide the sidebar.
func (m Model) modalActive() bool {
	return m.ShowHelp || m.Confirm != nil
}

// renderFooter renders the footer bar, overlaying a transient toast when active.
func (m Model) renderFooter() string {
	text := m.Keys.FooterHelp(m.Screen, m.Layout)
	if m.Toast == nil || m.Toast.Expired() {
		return m.Theme.FooterStyle.Width(m.Layout.Width).Render(text)
	}

	toast := components.Toast(m.Theme, m.Toast, m.Layout.Width)
	toastW := lipgloss.Width(toast)
	if toastW > m.Layout.Width {
		toastW = m.Layout.Width
	}

	helpW := m.Layout.Width - toastW - 1
	if helpW < 0 {
		helpW = 0
	}
	help := m.Theme.FooterStyle.Width(helpW).Render(text)
	return lipgloss.JoinHorizontal(lipgloss.Center, help, " ", toast)
}

func (m Model) renderContent() string {
	w := m.Layout.ContentWidth
	h := m.Layout.ContentHeight

	switch m.Screen {
	case core.ScreenWelcome:
		return screens.Welcome(m.Theme, m.WelcomeFocus, w, h)
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
		return screens.Chat(m.Theme, m.Session, m.ChatMessages, m.ChatInput, m.SelectedFindingPtr(), w, h)
	case core.ScreenApproval:
		return screens.Approval(m.Theme, m.Session, m.PostMode, m.NoPost, w, h)
	case core.ScreenReport:
		return screens.Report(m.Theme, m.Session, m.Session.ReportPath, 0, w, h)
	case core.ScreenSettings:
		return screens.Settings(m.Theme, m.SettingsForm, w, h)
	default:
		return components.EmptyState(m.Theme, "?", "Unknown screen", "Press esc to go back.", w, h)
	}
}
