package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Choices     []string
	Cursor      int
	Selected    map[int]struct{}
	MCPFiles    []string
	MultiSelect bool
	Quitted     bool
}

func NewModel(mcpFiles []string) Model {
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

	return Model{
		Choices:     choices,
		Selected:    selected,
		MCPFiles:    mcpFiles,
		MultiSelect: len(mcpFiles) > 1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitted = true
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}

		case "enter":
			return m, tea.Quit

		case " ":
			if m.MultiSelect {
				if m.Cursor == 0 {
					// Selecting "No mcp servers" - clear all other selections
					m.Selected = make(map[int]struct{})
					m.Selected[0] = struct{}{}
				} else {
					// Selecting an MCP server - clear "No mcp servers" first
					delete(m.Selected, 0)

					// Toggle the current MCP server selection
					if _, ok := m.Selected[m.Cursor]; ok {
						delete(m.Selected, m.Cursor)

						// If no MCP servers selected, re-select "No mcp servers"
						if len(m.Selected) == 0 {
							m.Selected[0] = struct{}{}
						}
					} else {
						m.Selected[m.Cursor] = struct{}{}
					}
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var s strings.Builder

	// Title with gradient
	title := CreateGradientText("⚡ Claude Code Launcher", PurpleGradientStart, PurpleGradientEnd)
	s.WriteString(title + "\n\n")

	// Header
	header := HeaderStyle.Render("🚀 Choose your MCP configuration:")
	s.WriteString(header + "\n")

	// Menu items
	for i, choice := range m.Choices {
		var cursor, checkbox, item string

		// Cursor
		if m.Cursor == i {
			cursor = CursorStyle.Render("❯")
		} else {
			cursor = " "
		}

		// Checkbox
		if _, ok := m.Selected[i]; ok {
			checkbox = CheckboxSelectedStyle.Render("⬢")
		} else {
			checkbox = CheckboxUnselectedStyle.Render("⬡")
		}

		// Item styling
		if m.Cursor == i {
			// Special handling for different choice types
			if i == 0 {
				// "No mcp servers" option
				item = SelectedItemStyle.Render("🚫 " + choice)
			} else {
				// Parse choice to separate name and location
				parts := strings.Split(choice, " (")
				name := parts[0]
				location := ""
				if len(parts) > 1 {
					location = " (" + parts[1]
				}

				styledName := SelectedItemStyle.Render(name)
				styledLocation := LocationStyle.Render(location)
				item = styledName + styledLocation
			}
		} else {
			// Unselected item
			if i == 0 {
				item = UnselectedItemStyle.Render("🚫 " + choice)
			} else {
				parts := strings.Split(choice, " (")
				name := parts[0]
				location := ""
				if len(parts) > 1 {
					location = " (" + parts[1]
				}

				styledName := UnselectedItemStyle.Render(name)
				styledLocation := LocationStyle.Render(location)
				item = styledName + styledLocation
			}
		}

		s.WriteString(fmt.Sprintf(" %s %s %s\n", cursor, checkbox, item))
	}

	// Help text
	helpText := "💡 Controls: "
	if m.MultiSelect {
		helpText += "↑/↓ navigate • space select • enter launch • q quit"
	} else {
		helpText += "↑/↓ navigate • enter launch • q quit"
	}

	help := HelpStyle.Render(helpText)
	s.WriteString("\n" + help + "\n")

	return s.String()
}