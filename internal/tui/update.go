package tui

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/core"
	"github.com/gnanam1990/zero-review/internal/tui/screens"
)

// tickMsg advances the loading animation.
type tickMsg time.Time

// Update handles all messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Layout = core.NewLayoutFromMsg(msg)
		m.Width = msg.Width
		m.Height = msg.Height
		m.ChatInput.SetWidth(m.Layout.ContentWidth - 6)
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tickMsg:
		return m.advanceLoading()
	}

	// Delegate to active form if present.
	if m.Screen == core.ScreenPRInput && m.PRForm != nil {
		f, cmd := m.PRForm.Update(msg)
		if form, ok := f.(*huh.Form); ok {
			m.PRForm = form
		}
		if m.PRForm.State == huh.StateCompleted {
			m.applyPRInput()
			m.Screen = core.ScreenLoadingReview
			return m, tickCmd()
		}
		return m, cmd
	}

	if m.Screen == core.ScreenSettings && m.SettingsForm != nil {
		f, cmd := m.SettingsForm.Update(msg)
		if form, ok := f.(*huh.Form); ok {
			m.SettingsForm = form
		}
		if m.SettingsForm.State == huh.StateCompleted {
			m.applySettings()
			m.Screen = core.ScreenWelcome
		}
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit
	if key.Matches(msg, m.Keys.Quit) {
		return m, tea.Quit
	}

	// Help toggle
	if key.Matches(msg, m.Keys.Help) {
		m.ShowHelp = !m.ShowHelp
		return m, nil
	}

	// Dismiss modal/toast
	if key.Matches(msg, m.Keys.Back) {
		if m.Confirm != nil {
			m.Confirm = nil
			return m, nil
		}
		if m.ShowHelp {
			m.ShowHelp = false
			return m, nil
		}
		if m.Screen == core.ScreenFindingDetail {
			m.Screen = core.ScreenFindings
			return m, nil
		}
		if m.Screen != core.ScreenWelcome && m.Screen != core.ScreenDashboard {
			m.Screen = core.ScreenDashboard
			return m, nil
		}
	}

	// Navigation shortcuts
	if key.Matches(msg, m.Keys.Overview) {
		m.Screen = core.ScreenDashboard
		return m, nil
	}
	if key.Matches(msg, m.Keys.Findings) {
		m.Screen = core.ScreenFindings
		return m, nil
	}
	if key.Matches(msg, m.Keys.Diff) {
		m.Screen = core.ScreenDiff
		return m, nil
	}
	if key.Matches(msg, m.Keys.Chat) {
		m.Screen = core.ScreenChat
		return m, nil
	}
	if key.Matches(msg, m.Keys.Post) {
		m.Screen = core.ScreenApproval
		return m, nil
	}
	if key.Matches(msg, m.Keys.Save) {
		m.SetToast("Report saved", "success")
		return m, nil
	}

	switch m.Screen {
	case core.ScreenWelcome:
		return m.handleWelcomeKey(msg)
	case core.ScreenFindings:
		return m.handleFindingsKey(msg)
	case core.ScreenFindingDetail:
		return m.handleFindingDetailKey(msg)
	case core.ScreenChat:
		return m.handleChatKey(msg)
	case core.ScreenApproval:
		return m.handleApprovalKey(msg)
	case core.ScreenReport:
		return m.handleReportKey(msg)
	}

	return m, nil
}

func (m Model) handleWelcomeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.PRForm = screens.BuildPRForm(&screens.PRInputValues{
			Provider:   m.Provider,
			Mode:       m.Mode,
			SaveReport: m.SaveReport,
			NoPost:     m.NoPost,
		})
		return m, m.PRForm.Init()
	case "l":
		m.Session = MockSession()
		m.Session.ReportPath = ".zero-review/reports/pr-123-demo.md"
		m.Screen = core.ScreenReport
		return m, nil
	case "comma", "s":
		m.SettingsForm = screens.BuildSettingsForm(&screens.SettingsValues{
			Provider:          m.Provider,
			Mode:              m.Mode,
			Confidence:        strconv.Itoa(m.Config.ConfidenceMin),
			AutoSave:          m.SaveReport,
			Theme:             "auto",
			ShowLowConfidence: false,
			DefaultPostMode:   "comment",
		})
		m.Screen = core.ScreenSettings
		return m, m.SettingsForm.Init()
	}
	return m, nil
}

func (m Model) handleFindingsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.Session == nil || len(m.Session.Findings) == 0 {
		return m, nil
	}

	filtered := screens.ApplyFilter(m.Session.Findings, m.SeverityFilter)
	if len(filtered) == 0 {
		if msg.String() == "x" {
			m.SeverityFilter = ""
			return m, nil
		}
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.FindingsCursor > 0 {
			m.FindingsCursor--
		}
	case "down", "j":
		if m.FindingsCursor < len(filtered)-1 {
			m.FindingsCursor++
		}
	case "enter":
		m.SelectedFinding = &m.FindingsCursor
		m.Screen = core.ScreenFindingDetail
	case "a":
		idx := realIndex(m.Session.Findings, filtered, m.FindingsCursor)
		if idx >= 0 {
			m.Session.Findings[idx].Status = review.FindingStatusApproved
			m.SetToast("Finding approved", "success")
		}
	case "r":
		idx := realIndex(m.Session.Findings, filtered, m.FindingsCursor)
		if idx >= 0 {
			m.Session.Findings[idx].Status = review.FindingStatusRejected
			m.SetToast("Finding rejected", "info")
		}
	case " ":
		idx := realIndex(m.Session.Findings, filtered, m.FindingsCursor)
		if idx >= 0 {
			if m.Session.Findings[idx].Status == review.FindingStatusApproved {
				m.Session.Findings[idx].Status = review.FindingStatusPending
			} else {
				m.Session.Findings[idx].Status = review.FindingStatusApproved
			}
		}
	case "e":
		idx := realIndex(m.Session.Findings, filtered, m.FindingsCursor)
		if idx >= 0 {
			m.SelectedFinding = &idx
			m.Screen = core.ScreenChat
			m.ChatInput.SetValue(m.Session.Findings[idx].SuggestedComment)
		}
	case "c":
		idx := realIndex(m.Session.Findings, filtered, m.FindingsCursor)
		if idx >= 0 {
			m.SelectedFinding = &idx
			m.Screen = core.ScreenChat
		}
	case "d":
		idx := realIndex(m.Session.Findings, filtered, m.FindingsCursor)
		if idx >= 0 {
			m.SelectedFinding = &idx
			m.Screen = core.ScreenDiff
		}
	case "1":
		m.SeverityFilter = "high"
	case "2":
		m.SeverityFilter = "medium"
	case "3":
		m.SeverityFilter = "low"
	case "4":
		m.SeverityFilter = "info"
	case "x":
		m.SeverityFilter = ""
	}
	return m, nil
}

func realIndex(all, filtered []review.Finding, filteredCursor int) int {
	if filteredCursor < 0 || filteredCursor >= len(filtered) {
		return -1
	}
	target := filtered[filteredCursor]
	for i, f := range all {
		if f.ID == target.ID {
			return i
		}
	}
	return -1
}

func (m Model) handleFindingDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.SelectedFindingPtr()
	if f == nil {
		return m, nil
	}
	switch msg.String() {
	case "a":
		f.Status = review.FindingStatusApproved
		m.SetToast("Finding approved", "success")
	case "r":
		f.Status = review.FindingStatusRejected
		m.SetToast("Finding rejected", "info")
	case "e":
		m.Screen = core.ScreenChat
		m.ChatInput.SetValue(f.SuggestedComment)
	case "c":
		m.Screen = core.ScreenChat
	case "d":
		m.Screen = core.ScreenDiff
	}
	return m, nil
}

func (m Model) handleChatKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.Screen = core.ScreenFindings
		m.ChatInput.Blur()
		return m, nil
	case "ctrl+l":
		m.ChatMessages = nil
		return m, nil
	case "enter":
		msgText := m.ChatInput.Value()
		if msgText == "" {
			return m, nil
		}
		m.ChatMessages = append(m.ChatMessages, core.ChatMessage{Role: "user", Text: msgText, Timestamp: time.Now()})
		reply := screens.SuggestedReply(msgText, m.SelectedFindingPtr())
		m.ChatMessages = append(m.ChatMessages, core.ChatMessage{Role: "ai", Text: reply, Timestamp: time.Now()})
		m.ChatInput.SetValue("")
		return m, nil
	}

	var cmd tea.Cmd
	m.ChatInput, cmd = m.ChatInput.Update(msg)
	return m, cmd
}

func (m Model) handleApprovalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "p":
		if m.Config.NoPost {
			m.SetToast("Posting disabled", "info")
			return m, nil
		}
		m.Confirm = &ConfirmPrompt{
			Title:   "Confirm Post",
			Body:    "Post approved comments to GitHub? This cannot be undone by Zero Review.",
			YesText: "Post Review",
			NoText:  "Cancel",
			Action:  "post",
		}
	case "m":
		// Cycle modes for demo
		return m, nil
	case "s":
		m.Session.ReportPath = ".zero-review/reports/pr-123-demo.md"
		m.SetToast("Report saved", "success")
	case "y":
		if m.Confirm != nil {
			m.Confirm = nil
			m.Session.Posted = true
			for i := range m.Session.Findings {
				if m.Session.Findings[i].IsApproved() {
					m.Session.Findings[i].Status = review.FindingStatusPosted
				}
			}
			m.SetToast("Review posted", "success")
			m.Screen = core.ScreenReport
		}
	case "n":
		m.Confirm = nil
	}
	return m, nil
}

func (m Model) handleReportKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "o":
		m.Screen = core.ScreenDashboard
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) applyPRInput() {
	if m.PRForm == nil {
		return
	}
	values := m.PRForm.Get("").(*screens.PRInputValues)
	m.Config.PRURL = values.URL
	m.Provider = values.Provider
	m.Mode = values.Mode
	m.SaveReport = values.SaveReport
	m.NoPost = values.NoPost
	m.Config.NoPost = values.NoPost
}

func (m *Model) applySettings() {
	if m.SettingsForm == nil {
		return
	}
	values := m.SettingsForm.Get("").(*screens.SettingsValues)
	m.Provider = values.Provider
	m.Mode = values.Mode
	if n, err := strconv.Atoi(values.Confidence); err == nil {
		m.Config.ConfidenceMin = n
	}
	m.SaveReport = values.AutoSave
	m.Config.NoPost = values.DefaultPostMode == "report_only"
	m.NoPost = m.Config.NoPost
}

func (m Model) advanceLoading() (tea.Model, tea.Cmd) {
	done := 0
	for i := range m.LoadingSteps {
		if m.LoadingSteps[i].Done {
			done++
		}
	}
	if done >= len(m.LoadingSteps) {
		m.Session = MockSession()
		m.Screen = core.ScreenDashboard
		return m, nil
	}

	for i := range m.LoadingSteps {
		if !m.LoadingSteps[i].Done {
			m.LoadingSteps[i].Done = true
			break
		}
	}
	return m, tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func tickCmd() tea.Cmd {
	return tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}
