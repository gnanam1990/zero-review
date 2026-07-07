package components

import "github.com/charmbracelet/lipgloss"

// Center places content in the center of the given terminal bounds.
func Center(content string, width, height int) string {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)
}
