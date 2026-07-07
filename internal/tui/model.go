package tui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gnanam1990/zero-review/internal/config"
	"github.com/gnanam1990/zero-review/internal/github"
	"github.com/gnanam1990/zero-review/internal/review"
)

// Screen is the current TUI screen.
type Screen int

const (
	ScreenSummary Screen = iota
	ScreenFindings
	ScreenDetail
	ScreenDiff
	ScreenChat
	ScreenApproval
	ScreenFinal
)

// Model is the Bubble Tea model for zero-review.
type Model struct {
	Styles    Styles
	Screen    Screen
	Result    review.Result
	Findings  []review.Finding
	Cursor    int
	Selected  *review.Finding
	Chat      []ChatMessage
	ChatInput string
	ChatMode  string // "chat" | "edit"
	Width     int
	Height    int
	Err       error

	GitHub  github.Client
	Options config.Options

	// Final state
	Summary     string
	Posted      int
	ReportPaths []string

	// Approval state
	ApprovedCount int
	RejectedCount int
	EditedCount   int
	PostMode      string
	NoPost        bool

	table table.Model
}

// ChatMessage is one chat turn.
type ChatMessage struct {
	Role string // user | ai
	Text string
}

// NewModel creates the initial TUI model.
func NewModel(result review.Result, gh github.Client, opts config.Options) Model {
	m := Model{
		Styles:   NewStyles(),
		Screen:   ScreenSummary,
		Result:   result,
		Findings: result.Findings,
		GitHub:   gh,
		Options:  opts,
		NoPost:   opts.NoPost,
		PostMode: "comment",
		Chat: []ChatMessage{
			{Role: "ai", Text: "Review complete. I found " + countString(len(result.Findings), "finding") + ". Use Enter to inspect, A to approve, R to reject, E to edit."},
		},
	}
	m.table = newFindingsTable(m.Styles, m.Findings)
	return m
}

func countString(n int, word string) string {
	if n == 1 {
		return "1 " + word
	}
	return strconv.Itoa(n) + " " + word + "s"
}

func itoa(n int) string { return strconv.Itoa(n) }

// Init is the initial command.
func (m Model) Init() tea.Cmd { return nil }

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.table.SetWidth(msg.Width - 6)
		m.table.SetHeight(msg.Height - 12)
		return m, nil

	case reportsSavedMsg:
		if msg.Err != nil {
			m.Err = msg.Err
		} else {
			m.ReportPaths = msg.Paths
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	if m.Screen == ScreenFindings {
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		m.Cursor = m.table.Cursor()
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	switch m.Screen {
	case ScreenSummary:
		return handleSummaryKey(m, key)
	case ScreenFindings:
		return handleFindingsKey(m, key)
	case ScreenDetail:
		return handleDetailKey(m, key)
	case ScreenChat:
		return handleChatKey(m, key)
	case ScreenApproval:
		return handleApprovalKey(m, key)
	case ScreenFinal:
		if key == "enter" || key == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the active screen.
func (m Model) View() string {
	if m.Err != nil {
		return m.Styles.Danger.Render("Error: "+m.Err.Error()) + "\n\n" + m.Styles.Help.Render("Press q to quit")
	}

	switch m.Screen {
	case ScreenSummary:
		return renderSummary(m)
	case ScreenFindings:
		return renderFindings(m)
	case ScreenDetail:
		return renderDetail(m)
	case ScreenDiff:
		return renderDiff(m)
	case ScreenChat:
		return renderChat(m)
	case ScreenApproval:
		return renderApproval(m)
	case ScreenFinal:
		return renderFinal(m)
	}
	return ""
}
