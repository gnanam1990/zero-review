package components

import "strings"

// Wrap wraps text to a maximum width.
func Wrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	var lines []string
	var line string
	for _, word := range strings.Fields(s) {
		if len(line)+len(word)+1 > width && line != "" {
			lines = append(lines, line)
			line = word
		} else {
			if line == "" {
				line = word
			} else {
				line += " " + word
			}
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// Truncate shortens a string with ellipsis.
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
