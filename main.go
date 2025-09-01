package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices     []string
	cursor      int
	selected    map[int]struct{}
	mcpFiles    []string
	multiSelect bool
}

func initialModel() model {
	mcpFiles := findMCPFiles()
	choices := []string{"No mcp servers"}

	for _, file := range mcpFiles {
		baseName := strings.TrimSuffix(filepath.Base(file), ".json")
		choices = append(choices, baseName)
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
	mcpDir := ".claude/mcp"

	if _, err := os.Stat(mcpDir); os.IsNotExist(err) {
		return mcpFiles
	}

	files, err := filepath.Glob(filepath.Join(mcpDir, "*.json"))
	if err != nil {
		return mcpFiles
	}

	return files
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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
	s := "Choose Claude Code launch configuration:\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checkbox := " "
		if _, ok := m.selected[i]; ok {
			checkbox = "✓"
		} else {
			checkbox = "○"
		}

		s += fmt.Sprintf("%s %s %s\n", cursor, checkbox, choice)
	}

	helpText := "\nPress enter to launch, q to quit"
	if m.multiSelect {
		helpText += ", space to select/deselect"
	}
	s += helpText + ".\n"

	return s
}

func main() {
	m := initialModel()

	if len(m.mcpFiles) == 0 {
		fmt.Println("No MCP configuration files found in .claude/mcp/")
		fmt.Println("Launching Claude Code without MCP servers...")

		cmd := exec.Command("claude")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error launching Claude Code: %v\n", err)
			os.Exit(1)
		}
		return
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}

	// Launch Claude Code after Bubble Tea exits
	if finalModel, ok := finalModel.(model); ok {
		fmt.Print("🚀 Launching Claude Code...\n\n")
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
		fmt.Fprintf(os.Stderr, "Error launching Claude Code: %v\n", err)
		os.Exit(1)
	}
}
