package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/gnanam1990/zero-review/internal/config"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/core"
	"github.com/gnanam1990/zero-review/internal/tui/screens"
)

// Model is the top-level Bubble Tea model for Zero Review.
type Model struct {
	// Core
	Screen  core.Screen
	Theme   *core.Theme
	Layout  core.Layout
	Keys    core.KeyMap
	Session *review.ReviewSession
	Config  config.Options

	// Navigation
	SidebarIndex int
	ShowHelp     bool
	WelcomeFocus int // 0=start, 1=last report, 2=settings

	// Findings
	FindingsCursor  int
	SelectedFinding *int // index into Session.Findings
	SeverityFilter  string

	// Chat
	ChatMessages []core.ChatMessage
	ChatInput    textarea.Model
	ChatViewport viewport.Model
	ChatContext  string

	// Command palette
	CommandOpen  bool
	CommandInput textinput.Model

	// PR Input form (Huh)
	PRForm     *huh.Form
	PRURLInput textinput.Model
	Provider   string
	Mode       string
	SaveReport bool
	NoPost     bool
	PostMode   screens.PostMode

	// Settings form (Huh)
	SettingsForm *huh.Form

	// Loading
	LoadingSteps []core.LoadingStep
	LoadingTip   string

	// Modal / Toast
	Confirm *ConfirmPrompt
	Toast   *core.Toast

	// Misc
	Width  int
	Height int
}

// ConfirmPrompt is a reusable yes/no modal.
type ConfirmPrompt struct {
	Title   string
	Body    string
	YesText string
	NoText  string
	Action  string
}

// NewModel creates the initial TUI model in the welcome screen.
func NewModel(theme *core.Theme) Model {
	m := Model{
		Screen:       core.ScreenWelcome,
		Theme:        theme,
		Layout:       core.NewLayout(80, 24),
		Keys:         core.DefaultKeyMap(),
		Config:       config.DefaultOptions(),
		Provider:     "fake",
		Mode:         "balanced",
		SaveReport:   true,
		NoPost:       false,
		PostMode:     screens.PostModeComment,
		WelcomeFocus: 0,
		LoadingSteps: defaultLoadingSteps(),
		LoadingTip:   "Zero Review never posts comments without your approval.",
	}

	m.ChatInput = textarea.New()
	m.ChatInput.Placeholder = "Ask about this PR..."
	m.ChatInput.SetWidth(40)
	m.ChatInput.SetHeight(3)

	m.CommandInput = textinput.New()
	m.CommandInput.Placeholder = "Type a command..."
	m.CommandInput.Width = 40

	m.PRURLInput = textinput.New()
	m.PRURLInput.Placeholder = "https://github.com/org/repo/pull/123"
	m.PRURLInput.Focus()

	return m
}

// NewModelWithSession creates the TUI model already populated with a review session.
func NewModelWithSession(theme *core.Theme, session *review.ReviewSession) Model {
	m := NewModel(theme)
	m.Session = session
	m.Screen = core.ScreenDashboard
	return m
}

func defaultLoadingSteps() []core.LoadingStep {
	return []core.LoadingStep{
		{Label: "Parse PR link"},
		{Label: "Fetch PR metadata"},
		{Label: "Load changed files"},
		{Label: "Build diff context"},
		{Label: "Ask AI reviewer"},
		{Label: "Validate line anchors"},
		{Label: "Remove noisy findings"},
		{Label: "Generate report"},
	}
}

// Init is the initial command.
func (m Model) Init() tea.Cmd {
	return nil
}

// SelectedFindingPtr returns the selected finding or nil.
func (m Model) SelectedFindingPtr() *review.Finding {
	if m.Session == nil || m.SelectedFinding == nil {
		return nil
	}
	idx := *m.SelectedFinding
	if idx < 0 || idx >= len(m.Session.Findings) {
		return nil
	}
	return &m.Session.Findings[idx]
}

// SetToast sets a transient toast message.
func (m *Model) SetToast(message, kind string) {
	m.Toast = &core.Toast{Message: message, Kind: kind, Until: time.Now().Add(3 * time.Second)}
}

// ClearToast removes the active toast.
func (m *Model) ClearToast() {
	m.Toast = nil
}
