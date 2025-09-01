package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"cc-launcher/internal/config"
	"cc-launcher/internal/launcher"
	"cc-launcher/internal/ui"
)

func main() {
	// Parse command line flags
	var debugFlag bool
	var localFlag bool
	flag.BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	flag.BoolVar(&localFlag, "local", false, "Only check for local MCP configurations, skip global ones")
	flag.Parse()

	// Set debug mode in config package
	config.SetDebugMode(debugFlag)

	// Find MCP files
	mcpFiles, err := config.FindMCPFiles(localFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error finding MCP files: "+err.Error()))
		os.Exit(1)
	}

	if len(mcpFiles) == 0 {
		// Show styled no-MCP message and launch without MCP
		launcher.ShowNoMCPMessage()

		err := launcher.LaunchClaudeCodeWithoutMCP()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error launching Claude Code: "+err.Error()))
			os.Exit(1)
		}
		return
	}

	// Create UI model and run Bubble Tea program
	m := ui.NewModel(mcpFiles)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("%s\n", ui.RenderError("Error running program: "+err.Error()))
		os.Exit(1)
	}

	// Launch Claude Code after Bubble Tea exits (only if user didn't quit)
	if finalModel, ok := finalModel.(ui.Model); ok && !finalModel.Quitted {
		launcher.ShowLaunchMessage()
		err := launcher.LaunchClaudeCode(finalModel.Selected, finalModel.MCPFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error launching Claude Code: "+err.Error()))
			os.Exit(1)
		}
	}
}