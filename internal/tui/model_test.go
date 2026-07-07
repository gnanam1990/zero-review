package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gnanam1990/zero-review/internal/review"
	"github.com/gnanam1990/zero-review/internal/tui/core"
)

func TestNewModelStartsAtWelcome(t *testing.T) {
	m := NewModel(core.NewTheme())
	if m.Screen != core.ScreenWelcome {
		t.Fatalf("expected ScreenWelcome, got %v", m.Screen)
	}
	if m.Session != nil {
		t.Fatal("expected no session on fresh model")
	}
}

func TestNewModelWithSessionStartsAtDashboard(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())
	if m.Screen != core.ScreenDashboard {
		t.Fatalf("expected ScreenDashboard, got %v", m.Screen)
	}
	if m.Session == nil {
		t.Fatal("expected session to be set")
	}
}

func TestQuitKeyQuits(t *testing.T) {
	m := NewModel(core.NewTheme())
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected a quit command")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected tea.QuitMsg, got %T", msg)
	}
	_ = updated
}

func TestNavigationKeys(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())

	cases := []struct {
		key    rune
		screen core.Screen
	}{
		{'o', core.ScreenDashboard},
		{'f', core.ScreenFindings},
		{'d', core.ScreenDiff},
		{'c', core.ScreenChat},
	}

	for _, tc := range cases {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tc.key}})
		if um, ok := updated.(Model); ok {
			if um.Screen != tc.screen {
				t.Errorf("key %q: expected %v, got %v", string(tc.key), tc.screen, um.Screen)
			}
		} else {
			t.Errorf("key %q: expected Model, got %T", string(tc.key), updated)
		}
	}
}

func TestApproveFinding(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())
	m.Screen = core.ScreenFindings
	m.FindingsCursor = 0
	m.SelectedFinding = intPtr(0)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	um := updated.(Model)
	if um.Session.Findings[0].Status != review.FindingStatusApproved {
		t.Fatalf("expected approved status, got %v", um.Session.Findings[0].Status)
	}
}

func TestRejectFinding(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())
	m.Screen = core.ScreenFindings
	m.FindingsCursor = 0
	m.SelectedFinding = intPtr(0)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	um := updated.(Model)
	if um.Session.Findings[0].Status != review.FindingStatusRejected {
		t.Fatalf("expected rejected status, got %v", um.Session.Findings[0].Status)
	}
}

func TestMockSessionHasExpectedFindings(t *testing.T) {
	s := MockSession()
	if len(s.Findings) != 7 {
		t.Fatalf("expected 7 mock findings, got %d", len(s.Findings))
	}
	if s.PR.Title == "" {
		t.Fatal("expected PR title")
	}
}

func intPtr(i int) *int {
	return &i
}
