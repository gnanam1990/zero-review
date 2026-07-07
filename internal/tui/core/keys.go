package core

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the global and context-aware key bindings.
type KeyMap struct {
	Quit     key.Binding
	Back     key.Binding
	Help     key.Binding
	NextPane key.Binding
	PrevPane key.Binding
	Command  key.Binding
	Save     key.Binding
	Post     key.Binding
	Chat     key.Binding
	Diff     key.Binding
	Findings key.Binding
	Overview key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "b"),
			key.WithHelp("esc", "back"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		NextPane: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		PrevPane: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev"),
		),
		Command: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "command"),
		),
		Save: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "save report"),
		),
		Post: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "post"),
		),
		Chat: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "chat"),
		),
		Diff: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "diff"),
		),
		Findings: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "findings"),
		),
		Overview: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "overview"),
		),
	}
}

// FooterHelp returns a compact help string for the footer.
func (k KeyMap) FooterHelp(screen Screen, layout Layout) string {
	base := "q quit · esc back · ? help"
	if layout.ShowSidebar {
		base = "↑/↓ nav · enter select · " + base
	}

	switch screen {
	case ScreenFindings:
		return "↑/↓ move · enter open · a approve · r reject · e edit · c chat · d diff · p post · x clear filters · " + base
	case ScreenFindingDetail:
		return "a approve · r reject · e edit · c chat · d diff · b/esc back · " + base
	case ScreenChat:
		return "enter send · shift+enter newline · ctrl+l clear · ctrl+r regen · esc back · " + base
	case ScreenApproval:
		return "p post · m mode · s save · esc back · " + base
	case ScreenPRInput:
		return "tab next field · shift+tab prev · enter submit · esc cancel · " + base
	default:
		return base
	}
}
