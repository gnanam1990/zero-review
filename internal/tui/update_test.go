package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gnanam1990/zero-review/internal/tui/core"
	"github.com/gnanam1990/zero-review/internal/tui/screens"
)

func TestQDoesNotQuitInFormOrChat(t *testing.T) {
	m := NewModel(core.NewTheme())
	m.Screen = core.ScreenPRInput
	m.PRForm = screens.BuildPRForm(&screens.PRInputValues{})

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd != nil {
		if _, ok := cmd().(tea.QuitMsg); ok {
			t.Fatal("typing 'q' inside a Huh form must not quit")
		}
	}
	um := updated.(Model)
	if um.Screen != core.ScreenPRInput {
		t.Fatalf("expected to stay on PRInput, got %v", um.Screen)
	}

	// Chat should also let 'q' through to the textarea.
	m2 := NewModelWithSession(core.NewTheme(), MockSession())
	m2.Screen = core.ScreenChat
	m2.ChatInput.Focus()
	updated2, _ := m2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	um2 := updated2.(Model)
	if um2.Screen != core.ScreenChat {
		t.Fatalf("typing 'q' in chat changed screen to %v", um2.Screen)
	}
	if um2.ChatInput.Value() != "q" {
		t.Fatalf("expected 'q' in chat input, got %q", um2.ChatInput.Value())
	}
}

func TestEscInChatReturnsToDashboardWhenNoFinding(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())
	m.Screen = core.ScreenChat
	m.ChatInput.Focus()
	m.ChatInput.SetValue("hello")

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	um := updated.(Model)
	if um.Screen != core.ScreenDashboard {
		t.Fatalf("expected dashboard, got %v", um.Screen)
	}
	if um.ChatInput.Value() != "" {
		t.Fatalf("expected chat input cleared, got %q", um.ChatInput.Value())
	}
}

func TestFindingsEnterUsesRealIndexWithFilter(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())
	m.Screen = core.ScreenFindings
	m.SeverityFilter = "high"
	// With the mock session, only the high finding is visible when filtered.
	m.FindingsCursor = 0

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	um := updated.(Model)
	if um.Screen != core.ScreenFindingDetail {
		t.Fatalf("expected finding detail, got %v", um.Screen)
	}
	if um.SelectedFinding == nil || *um.SelectedFinding != 0 {
		t.Fatalf("expected selected finding real index 0 (high), got %v", um.SelectedFinding)
	}
}

func TestFindingsDiffSetsSelectedFinding(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())
	m.Screen = core.ScreenFindings
	m.FindingsCursor = 2

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	um := updated.(Model)
	if um.Screen != core.ScreenDiff {
		t.Fatalf("expected diff screen, got %v", um.Screen)
	}
	if um.SelectedFinding == nil || *um.SelectedFinding != 2 {
		t.Fatalf("expected selected finding index 2, got %v", um.SelectedFinding)
	}
}

func TestApprovalPOpensConfirm(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())
	m.Screen = core.ScreenApproval

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	um := updated.(Model)
	if um.Confirm == nil {
		t.Fatal("expected confirm prompt to open")
	}
	if um.Confirm.Action != "post" {
		t.Fatalf("expected post action, got %q", um.Confirm.Action)
	}
}

func TestCommandPaletteOpensAndExecutes(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())
	m.Screen = core.ScreenDashboard

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	um := updated.(Model)
	if !um.CommandOpen {
		t.Fatal("expected command palette to open")
	}

	// Type "findings".
	for _, r := range "findings" {
		updated, _ = um.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		um = updated.(Model)
	}
	if um.CommandInput.Value() != "findings" {
		t.Fatalf("expected command query 'findings', got %q", um.CommandInput.Value())
	}

	// Execute.
	updated, _ = um.Update(tea.KeyMsg{Type: tea.KeyEnter})
	um = updated.(Model)
	if um.Screen != core.ScreenFindings {
		t.Fatalf("expected findings screen, got %v", um.Screen)
	}
	if um.CommandOpen {
		t.Fatal("expected command palette to close")
	}
}

func TestBInChatIsTypedNotBack(t *testing.T) {
	m := NewModelWithSession(core.NewTheme(), MockSession())
	m.Screen = core.ScreenChat
	m.ChatInput.Focus()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	um := updated.(Model)
	if um.Screen != core.ScreenChat {
		t.Fatalf("typing 'b' in chat changed screen to %v", um.Screen)
	}
	if um.ChatInput.Value() != "b" {
		t.Fatalf("expected 'b' in chat input, got %q", um.ChatInput.Value())
	}
}
