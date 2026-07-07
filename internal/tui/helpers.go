package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gnanam1990/zero-review/internal/review"
)

type teaModel = tea.Model
type teaCmd = tea.Cmd

func quitCmd() tea.Cmd {
	return tea.Quit
}

func reviewRiskScore(findings []review.Finding) int {
	return review.RiskScore(findings)
}
