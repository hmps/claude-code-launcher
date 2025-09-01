package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// RenderError creates a styled error message with consistent formatting
// Uses red color scheme with rounded border and bold text
func RenderError(message string) string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E74C3C")).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E74C3C")).
		Padding(0, 1)

	return errorStyle.Render("❌ " + message)
}