package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Choices     []string
	Cursor      int
	Selected    map[int]struct{}
	MCPFiles    []string
	MultiSelect bool
	Quitted     bool
	Happy       bool
	// Flag states
	HappyFlag    bool
	ContinueFlag bool
	ResumeFlag   bool
	YoloFlag     bool
	// UI state
	ShowingMCPSelection bool
	FlagCursor          int
}

func NewModel(mcpFiles []string, happy bool) Model {
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
		Happy:       happy,
		// Initialize flags to false by default
		HappyFlag:    false,
		ContinueFlag: false,
		ResumeFlag:   false,
		YoloFlag:     false,
		// Start with showing MCP selection
		ShowingMCPSelection: true,
		FlagCursor:          0,
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

		case "tab":
			// Toggle between flag selection and MCP selection
			m.ShowingMCPSelection = !m.ShowingMCPSelection
			if m.ShowingMCPSelection {
				m.Cursor = 0
			} else {
				m.FlagCursor = 0
			}

		case "up", "k":
			if m.ShowingMCPSelection {
				if m.Cursor > 0 {
					m.Cursor--
				}
			} else {
				if m.FlagCursor > 0 {
					m.FlagCursor--
				}
			}

		case "down", "j":
			if m.ShowingMCPSelection {
				if m.Cursor < len(m.Choices)-1 {
					m.Cursor++
				}
			} else {
				if m.FlagCursor < 3 { // 4 flags total (0-3)
					m.FlagCursor++
				}
			}

		case "enter":
			return m, tea.Quit

		// Number key shortcuts for MCP server selection
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if m.MultiSelect {
				key := msg.String()
				index := int(key[0] - '0') // Convert string digit to int

				if index < len(m.Choices) {
					if index == 0 {
						// Selecting "No mcp servers" - clear all other selections
						m.Selected = make(map[int]struct{})
						m.Selected[0] = struct{}{}
					} else {
						// Selecting an MCP server - clear "No mcp servers" first
						delete(m.Selected, 0)

						// Toggle the MCP server selection
						if _, ok := m.Selected[index]; ok {
							delete(m.Selected, index)

							// If no MCP servers selected, re-select "No mcp servers"
							if len(m.Selected) == 0 {
								m.Selected[0] = struct{}{}
							}
						} else {
							m.Selected[index] = struct{}{}
						}
					}
				}
			}

		// Letter key shortcuts for flag toggles
		case "h":
			m.HappyFlag = !m.HappyFlag
		case "c":
			m.ContinueFlag = !m.ContinueFlag
			// If both continue and resume are selected, resume takes priority
			if m.ContinueFlag && m.ResumeFlag {
				m.ResumeFlag = false
			}
		case "r":
			m.ResumeFlag = !m.ResumeFlag
			// If both continue and resume are selected, resume takes priority
			if m.ResumeFlag && m.ContinueFlag {
				m.ContinueFlag = false
			}
		case "y":
			m.YoloFlag = !m.YoloFlag

		case " ":
			if m.ShowingMCPSelection {
				// Handle MCP selection
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
			} else {
				// Handle flag selection
				switch m.FlagCursor {
				case 0:
					m.HappyFlag = !m.HappyFlag
				case 1:
					m.ContinueFlag = !m.ContinueFlag
					// If both continue and resume are selected, resume takes priority
					if m.ContinueFlag && m.ResumeFlag {
						m.ContinueFlag = false
					}
				case 2:
					m.ResumeFlag = !m.ResumeFlag
					// If both continue and resume are selected, resume takes priority
					if m.ResumeFlag && m.ContinueFlag {
						m.ContinueFlag = false
					}
				case 3:
					m.YoloFlag = !m.YoloFlag
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
	s.WriteString(title + "\n")

	// Show confirmation if happy flag is set
	if m.Happy {
		happyStyle := lipgloss.NewStyle().Foreground(MutedColor)
		s.WriteString(happyStyle.Render("🦦 Happy mode enabled") + "\n")
	}
	s.WriteString("\n")

	// MCP section header
	mcpHeaderStyle := HeaderStyle
	if m.ShowingMCPSelection {
		mcpHeaderStyle = mcpHeaderStyle.Foreground(SelectedItemStyle.GetForeground())
	}
	mcpHeader := mcpHeaderStyle.Render("🚀 Choose your MCP configuration:")
	s.WriteString(mcpHeader + "\n")

	// MCP menu items
	for i, choice := range m.Choices {
		var cursor, checkbox, item string

		// Cursor (only show when in MCP selection mode)
		if m.ShowingMCPSelection && m.Cursor == i {
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

		// Item styling with number shortcut
		if m.ShowingMCPSelection && m.Cursor == i {
			// Special handling for different choice types
			if i == 0 {
				// "No mcp servers" option
				item = SelectedItemStyle.Render("🚫 " + choice + " [" + fmt.Sprintf("%d", i) + "]")
			} else {
				// Parse choice to separate name and location
				parts := strings.Split(choice, " (")
				name := parts[0]
				location := ""
				if len(parts) > 1 {
					location = " (" + parts[1]
				}

				// Add number shortcut to the name
				nameWithShortcut := name + " [" + fmt.Sprintf("%d", i) + "]"
				styledName := SelectedItemStyle.Render(nameWithShortcut)
				styledLocation := LocationStyle.Render(location)
				item = styledName + styledLocation
			}
		} else {
			// Unselected item
			if i == 0 {
				item = UnselectedItemStyle.Render("🚫 " + choice + " [" + fmt.Sprintf("%d", i) + "]")
			} else {
				parts := strings.Split(choice, " (")
				name := parts[0]
				location := ""
				if len(parts) > 1 {
					location = " (" + parts[1]
				}

				// Add number shortcut to the name
				nameWithShortcut := name + " [" + fmt.Sprintf("%d", i) + "]"
				styledName := UnselectedItemStyle.Render(nameWithShortcut)
				styledLocation := LocationStyle.Render(location)
				item = styledName + styledLocation
			}
		}

		s.WriteString(fmt.Sprintf(" %s %s %s\n", cursor, checkbox, item))
	}

	s.WriteString("\n")

	// Flags section
	flagHeaderStyle := HeaderStyle
	if !m.ShowingMCPSelection {
		flagHeaderStyle = flagHeaderStyle.Foreground(SelectedItemStyle.GetForeground())
	}
	flagHeader := flagHeaderStyle.Render("⚙️ Configuration Flags:")
	s.WriteString(flagHeader + "\n")

	// Flag choices
	flagChoices := []struct {
		name     string
		label    string
		value    bool
		shortcut string
	}{
		{"happy", "🦦 Use happy [h]", m.HappyFlag, "h"},
		{"continue", "🔄 Continue previous session [c]", m.ContinueFlag, "c"},
		{"resume", "📂 Resume previous session [r]", m.ResumeFlag, "r"},
		{"yolo", "⚠️ Skip permissions check [y]", m.YoloFlag, "y"},
	}

	for i, flag := range flagChoices {
		var cursor, checkbox, item string

		// Cursor (only show when in flag selection mode)
		if !m.ShowingMCPSelection && m.FlagCursor == i {
			cursor = CursorStyle.Render("❯")
		} else {
			cursor = " "
		}

		// Checkbox
		if flag.value {
			checkbox = CheckboxSelectedStyle.Render("⬢")
		} else {
			checkbox = CheckboxUnselectedStyle.Render("⬡")
		}

		// Item styling
		if !m.ShowingMCPSelection && m.FlagCursor == i {
			item = SelectedItemStyle.Render(flag.label)
		} else {
			item = UnselectedItemStyle.Render(flag.label)
		}

		s.WriteString(fmt.Sprintf(" %s %s %s\n", cursor, checkbox, item))
	}

	// Help text
	helpText := "💡 Controls: "
	if m.MultiSelect {
		helpText += "tab switch sections • ↑/↓ navigate • space select • enter launch • q quit"
	} else {
		helpText += "tab switch sections • ↑/↓ navigate • enter launch • q quit"
	}

	help := HelpStyle.Render(helpText)
	s.WriteString("\n" + help + "\n")

	return s.String()
}
