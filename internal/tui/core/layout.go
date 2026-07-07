package core

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	minWidth     = 60
	minHeight    = 20
	SidebarWidth = 18
	headerHeight = 1
	footerHeight = 1
)

// Layout holds the computed geometry for the current frame.
type Layout struct {
	Width         int
	Height        int
	ShowSidebar   bool
	TooSmall      bool
	ContentX      int
	ContentY      int
	ContentWidth  int
	ContentHeight int
}

// NewLayout computes layout bounds from a window size message.
func NewLayout(width, height int) Layout {
	showSidebar := width >= 100
	tooSmall := width < minWidth || height < minHeight

	contentX := 0
	contentWidth := width
	if showSidebar {
		contentX = SidebarWidth
		contentWidth = width - SidebarWidth
	}

	contentY := headerHeight
	contentHeight := height - headerHeight - footerHeight

	return Layout{
		Width:         width,
		Height:        height,
		ShowSidebar:   showSidebar,
		TooSmall:      tooSmall,
		ContentX:      contentX,
		ContentY:      contentY,
		ContentWidth:  contentWidth,
		ContentHeight: contentHeight,
	}
}

// NewLayoutFromMsg is a convenience wrapper.
func NewLayoutFromMsg(msg tea.WindowSizeMsg) Layout {
	return NewLayout(msg.Width, msg.Height)
}

// TerminalTooSmallMessage returns the full-screen "resize" message.
func TerminalTooSmallMessage(theme *Theme, width, height int) string {
	return theme.TerminalTooSmall.
		Width(width).
		Height(height).
		Render("Terminal too small.\nPlease resize to at least 60×20.")
}
