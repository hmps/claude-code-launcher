package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CreateGradientText creates a gradient effect for text
func CreateGradientText(text string, startColor, endColor lipgloss.Color) string {
	// Create a gradient effect by alternating between colors
	// Split the text and apply different background colors to create gradient illusion
	textRunes := []rune(text)
	if len(textRunes) == 0 {
		return ""
	}

	// For a visual gradient effect, we'll create segments with different shades
	var result strings.Builder

	// Create base style
	baseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	// Split text into chunks and alternate colors to simulate gradient
	for i, r := range textRunes {
		var bgColor lipgloss.Color
		if i < len(textRunes)/2 {
			bgColor = startColor
		} else {
			bgColor = endColor
		}

		charStyle := baseStyle.Background(bgColor)
		result.WriteString(charStyle.Render(string(r)))
	}

	// Add padding around the entire gradient
	finalStyle := lipgloss.NewStyle().Padding(1, 2)
	return finalStyle.Render(result.String())
}

// Color palette
var (
	PrimaryColor        = lipgloss.Color("#FF6B9D")
	SecondaryColor      = lipgloss.Color("#4ECDC4")
	AccentColor         = lipgloss.Color("#45B7D1")
	SuccessColor        = lipgloss.Color("#96CEB4")
	MutedColor          = lipgloss.Color("#6C7B7F")
	PurpleGradientStart = lipgloss.Color("#8B5CF6")
	PurpleGradientEnd   = lipgloss.Color("#A855F7")
	TextColor           = lipgloss.Color("#ECF0F1")
)

// Title styles with gradient effect
var (
	TitleStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
		Foreground(TextColor).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	// Item styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true).
				Padding(0, 1)

	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Padding(0, 1)

	CursorStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	CheckboxSelectedStyle = lipgloss.NewStyle().
				Foreground(SuccessColor).
				Bold(true)

	CheckboxUnselectedStyle = lipgloss.NewStyle().
				Foreground(MutedColor)

	LocationStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Italic(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			MarginTop(1)

	LaunchStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			Padding(0, 1)
)