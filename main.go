package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// createGradientText creates a gradient effect for text
func createGradientText(text string, startColor, endColor lipgloss.Color) string {
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

var (
	// Color palette
	primaryColor        = lipgloss.Color("#FF6B9D")
	secondaryColor      = lipgloss.Color("#4ECDC4")
	accentColor         = lipgloss.Color("#45B7D1")
	successColor        = lipgloss.Color("#96CEB4")
	mutedColor          = lipgloss.Color("#6C7B7F")
	purpleGradientStart = lipgloss.Color("#8B5CF6")
	purpleGradientEnd   = lipgloss.Color("#A855F7")
	textColor           = lipgloss.Color("#ECF0F1")

	// Title styles with gradient effect
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(1, 2)

	headerStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// Item styles
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Padding(0, 1)

	unselectedItemStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Padding(0, 1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	checkboxSelectedStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	checkboxUnselectedStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	locationStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			MarginTop(1)

	launchStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(0, 1)
)

type model struct {
	choices     []string
	cursor      int
	selected    map[int]struct{}
	mcpFiles    []string
	multiSelect bool
	quitted     bool
}

func initialModel() model {
	mcpFiles := findMCPFiles()
	choices := []string{"No mcp servers"}

	for _, file := range mcpFiles {
		baseName := strings.TrimSuffix(filepath.Base(file), ".json")

		// Determine if file is from local or global directory
		var location string
		if strings.HasPrefix(file, ".claude/mcp/") {
			location = "local"
		} else {
			location = "global"
		}

		choices = append(choices, fmt.Sprintf("%s (%s)", baseName, location))
	}

	selected := make(map[int]struct{})
	selected[0] = struct{}{} // Pre-select "No mcp servers"

	return model{
		choices:     choices,
		selected:    selected,
		mcpFiles:    mcpFiles,
		multiSelect: len(mcpFiles) > 1,
	}
}

func findMCPFiles() []string {
	var mcpFiles []string

	// Scan local directory (.claude/mcp)
	localDir := ".claude/mcp"
	if _, err := os.Stat(localDir); !os.IsNotExist(err) {
		if files, err := filepath.Glob(filepath.Join(localDir, "*.json")); err == nil {
			mcpFiles = append(mcpFiles, files...)
		}
	}

	// Scan global directory (~/.claude/mcp)
	if homeDir, err := os.UserHomeDir(); err == nil {
		globalDir := filepath.Join(homeDir, ".claude", "mcp")
		if _, err := os.Stat(globalDir); !os.IsNotExist(err) {
			if files, err := filepath.Glob(filepath.Join(globalDir, "*.json")); err == nil {
				mcpFiles = append(mcpFiles, files...)
			}
		}
	}

	return mcpFiles
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitted = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			return m, tea.Quit

		case " ":
			if m.multiSelect {
				if m.cursor == 0 {
					// Selecting "No mcp servers" - clear all other selections
					m.selected = make(map[int]struct{})
					m.selected[0] = struct{}{}
				} else {
					// Selecting an MCP server - clear "No mcp servers" first
					delete(m.selected, 0)

					// Toggle the current MCP server selection
					if _, ok := m.selected[m.cursor]; ok {
						delete(m.selected, m.cursor)

						// If no MCP servers selected, re-select "No mcp servers"
						if len(m.selected) == 0 {
							m.selected[0] = struct{}{}
						}
					} else {
						m.selected[m.cursor] = struct{}{}
					}
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	// Title with gradient
	title := createGradientText("⚡ Claude Code Launcher", purpleGradientStart, purpleGradientEnd)
	s.WriteString(title + "\n\n")

	// Header
	header := headerStyle.Render("🚀 Choose your MCP configuration:")
	s.WriteString(header + "\n")

	// Menu items
	for i, choice := range m.choices {
		var cursor, checkbox, item string

		// Cursor
		if m.cursor == i {
			cursor = cursorStyle.Render("❯")
		} else {
			cursor = " "
		}

		// Checkbox
		if _, ok := m.selected[i]; ok {
			checkbox = checkboxSelectedStyle.Render("⬢")
		} else {
			checkbox = checkboxUnselectedStyle.Render("⬡")
		}

		// Item styling
		if m.cursor == i {
			// Special handling for different choice types
			if i == 0 {
				// "No mcp servers" option
				item = selectedItemStyle.Render("🚫 " + choice)
			} else {
				// Parse choice to separate name and location
				parts := strings.Split(choice, " (")
				name := parts[0]
				location := ""
				if len(parts) > 1 {
					location = " (" + parts[1]
				}

				styledName := selectedItemStyle.Render(name)
				styledLocation := locationStyle.Render(location)
				item = styledName + styledLocation
			}
		} else {
			// Unselected item
			if i == 0 {
				item = unselectedItemStyle.Render("🚫 " + choice)
			} else {
				parts := strings.Split(choice, " (")
				name := parts[0]
				location := ""
				if len(parts) > 1 {
					location = " (" + parts[1]
				}

				styledName := unselectedItemStyle.Render(name)
				styledLocation := locationStyle.Render(location)
				item = styledName + styledLocation
			}
		}

		s.WriteString(fmt.Sprintf(" %s %s %s\n", cursor, checkbox, item))
	}

	// Help text
	helpText := "💡 Controls: "
	if m.multiSelect {
		helpText += "↑/↓ navigate • space select • enter launch • q quit"
	} else {
		helpText += "↑/↓ navigate • enter launch • q quit"
	}

	help := helpStyle.Render(helpText)
	s.WriteString("\n" + help + "\n")

	return s.String()
}

func main() {
	m := initialModel()

	if len(m.mcpFiles) == 0 {
		// Styled no-MCP message
		title := createGradientText("⚡ Claude Code Launcher", purpleGradientStart, purpleGradientEnd)
		noMcpStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor)

		launchMsg := createGradientText("🚀 Launching Claude Code without MCP servers...", purpleGradientStart, purpleGradientEnd)

		fmt.Println(title)
		fmt.Println()
		fmt.Println(noMcpStyle.Render("📁 No MCP configuration files found in .claude/mcp/ or ~/.claude/mcp/"))
		fmt.Println()
		fmt.Println(launchMsg)
		fmt.Println()

		cmd := exec.Command("claude")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E74C3C")).
				Bold(true).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#E74C3C")).
				Padding(0, 1)
			fmt.Fprintf(os.Stderr, "%s\n", errorStyle.Render("❌ Error launching Claude Code: "+err.Error()))
			os.Exit(1)
		}
		return
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E74C3C")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#E74C3C")).
			Padding(0, 1)
		fmt.Printf("%s\n", errorStyle.Render("❌ Error running program: "+err.Error()))
		os.Exit(1)
	}

	// Launch Claude Code after Bubble Tea exits (only if user didn't quit)
	if finalModel, ok := finalModel.(model); ok && !finalModel.quitted {
		launchMsg := createGradientText("🚀 Launching Claude Code...", purpleGradientStart, purpleGradientEnd)
		fmt.Println()
		fmt.Println(launchMsg)
		fmt.Println()
		launchClaudeCodeFromSelection(finalModel.selected, finalModel.mcpFiles)
	}
}

func launchClaudeCodeFromSelection(selected map[int]struct{}, mcpFiles []string) {
	var args []string

	// If "No mcp servers" is selected
	if _, noMcpSelected := selected[0]; noMcpSelected {
		args = []string{}
	} else {
		args = []string{"--mcp-config"}

		// Add all selected MCP files
		for i := range selected {
			if i > 0 && i-1 < len(mcpFiles) {
				args = append(args, mcpFiles[i-1])
			}
		}
	}

	cmd := exec.Command("claude", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E74C3C")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#E74C3C")).
			Padding(0, 1)
		fmt.Fprintf(os.Stderr, "%s\n", errorStyle.Render("❌ Error launching Claude Code: "+err.Error()))
		os.Exit(1)
	}
}
