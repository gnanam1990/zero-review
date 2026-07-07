package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gnanam1990/zero-review/internal/tui/core"
	"github.com/gnanam1990/zero-review/internal/tui/screens"
)

func TestAllScreensRenderWithoutPanic(t *testing.T) {
	theme := core.NewTheme()
	session := MockSession()

	cases := []struct {
		name   string
		setup  func(*Model)
		needle string
	}{
		{
			name:   "welcome",
			setup:  func(m *Model) { m.Screen = core.ScreenWelcome },
			needle: "Zero Review",
		},
		{
			name: "pr-input",
			setup: func(m *Model) {
				m.Screen = core.ScreenPRInput
				m.PRForm = screens.BuildPRForm(&screens.PRInputValues{})
				_ = m.PRForm.Init()
			},
			needle: "GitHub PR URL",
		},
		{
			name: "loading",
			setup: func(m *Model) {
				m.Screen = core.ScreenLoadingReview
				m.Layout = core.NewLayout(100, 30)
			},
			needle: "Reviewing PR",
		},
		{
			name:   "dashboard",
			setup:  func(m *Model) { m.Screen = core.ScreenDashboard; m.Session = session },
			needle: "PR Summary",
		},
		{
			name: "findings",
			setup: func(m *Model) {
				m.Screen = core.ScreenFindings
				m.Session = session
				m.Layout = core.NewLayout(120, 40)
			},
			needle: "Findings",
		},
		{
			name: "finding-detail",
			setup: func(m *Model) {
				m.Screen = core.ScreenFindingDetail
				m.Session = session
				m.SelectedFinding = intPtr(0)
			},
			needle: session.Findings[0].Title,
		},
		{
			name: "diff",
			setup: func(m *Model) {
				m.Screen = core.ScreenDiff
				m.Session = session
				m.SelectedFinding = intPtr(0)
			},
			needle: "diff --git",
		},
		{
			name: "chat",
			setup: func(m *Model) {
				m.Screen = core.ScreenChat
				m.Session = session
				m.SelectedFinding = intPtr(0)
				m.ChatInput.Focus()
			},
			needle: "Chat with Zero",
		},
		{
			name:   "approval",
			setup:  func(m *Model) { m.Screen = core.ScreenApproval; m.Session = session },
			needle: "Approval",
		},
		{
			name:   "report",
			setup:  func(m *Model) { m.Screen = core.ScreenReport; m.Session = session },
			needle: "Review Report",
		},
		{
			name: "settings",
			setup: func(m *Model) {
				m.Screen = core.ScreenSettings
				m.SettingsForm = screens.BuildSettingsForm(&screens.SettingsValues{})
				_ = m.SettingsForm.Init()
			},
			needle: "Default provider",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel(theme)
			m.Layout = core.NewLayout(100, 30)
			m.Width = 100
			m.Height = 30
			tc.setup(&m)

			view := m.View()
			if view == "" {
				t.Fatal("rendered empty view")
			}
			if !strings.Contains(stripANSI(view), tc.needle) {
				t.Fatalf("expected view to contain %q", tc.needle)
			}
		})
	}
}

func TestHelpOverlayRenders(t *testing.T) {
	m := NewModel(core.NewTheme())
	m.Layout = core.NewLayout(100, 30)
	m.ShowHelp = true

	view := m.View()
	if !strings.Contains(stripANSI(view), "Keyboard shortcuts") {
		t.Fatal("help overlay not rendered")
	}
}

func TestConfirmModalRenders(t *testing.T) {
	m := NewModel(core.NewTheme())
	m.Layout = core.NewLayout(100, 30)
	m.Session = MockSession()
	m.Screen = core.ScreenApproval
	m.Confirm = &ConfirmPrompt{
		Title:   "Post?",
		Body:    "really post?",
		YesText: "Yes",
		NoText:  "No",
		Action:  "post",
	}

	view := m.View()
	if !strings.Contains(stripANSI(view), "Post?") {
		t.Fatal("confirm modal not rendered")
	}
}

func TestCommandPaletteRenders(t *testing.T) {
	m := NewModel(core.NewTheme())
	m.Layout = core.NewLayout(100, 30)
	m.CommandOpen = true
	m.CommandInput.Focus()

	view := m.View()
	if !strings.Contains(stripANSI(view), "Type a command") {
		t.Fatal("command palette not rendered")
	}
}

func TestWindowSizeDoesNotPanic(t *testing.T) {
	m := NewModel(core.NewTheme())
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	um := updated.(Model)
	if um.Layout.Width != 120 {
		t.Fatalf("expected width 120, got %d", um.Layout.Width)
	}
}

func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '\x1b' {
			j := i + 1
			for j < len(s) && (s[j] == '[' || (s[j] >= 0x30 && s[j] <= 0x3f)) {
				j++
			}
			for j < len(s) && (s[j] >= 0x20 && s[j] <= 0x2f) {
				j++
			}
			if j < len(s) {
				j++
			}
			i = j
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
