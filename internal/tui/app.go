package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gnanam1990/zero-review/internal/review"
	core "github.com/gnanam1990/zero-review/internal/tui/core"
)

// NewApp creates a new Bubble Tea program for the mock TUI.
func NewApp() *tea.Program {
	theme := core.NewTheme()
	m := NewModel(theme)
	return tea.NewProgram(m, tea.WithAltScreen())
}

// NewAppWithSession creates a program preloaded with a review session.
func NewAppWithSession(session *review.ReviewSession) *tea.Program {
	theme := core.NewTheme()
	m := NewModelWithSession(theme, session)
	return tea.NewProgram(m, tea.WithAltScreen())
}
