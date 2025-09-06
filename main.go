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
	var yoloFlag bool
	var happyFlag bool
	var resumeFlag bool
	var continueFlag bool
	var blankFlag bool
	flag.BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	flag.BoolVar(&localFlag, "local", false, "Only check for local MCP configurations, skip global ones")
	flag.BoolVar(&yoloFlag, "yolo", false, "Launch Claude Code with --dangerously-skip-permissions")
	flag.BoolVar(&happyFlag, "happy", false, "Use happy instead of claude command")
	flag.BoolVar(&resumeFlag, "r", false, "Launch Claude Code with --resume flag (-r, --resume)")
	flag.BoolVar(&resumeFlag, "resume", false, "Launch Claude Code with --resume flag (-r, --resume)")
	flag.BoolVar(&continueFlag, "c", false, "Launch Claude Code with --continue flag (-c, --continue)")
	flag.BoolVar(&continueFlag, "continue", false, "Launch Claude Code with --continue flag (-c, --continue)")
	flag.BoolVar(&blankFlag, "b", false, "Launch Claude Code without MCP servers (skip TUI) (-b, --blank)")
	flag.BoolVar(&blankFlag, "blank", false, "Launch Claude Code without MCP servers (skip TUI) (-b, --blank)")
	
	// Custom usage function to show double dashes for all flags except -r, -c, and -b
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  -b, --blank\n        Launch Claude Code without MCP servers (skip TUI)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -c, --continue\n        Launch Claude Code with --continue flag\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --debug\n        Enable debug logging\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --happy\n        Use happy instead of claude command\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --local\n        Only check for local MCP configurations, skip global ones\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -r, --resume\n        Launch Claude Code with --resume flag\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --yolo\n        Launch Claude Code with --dangerously-skip-permissions\n")
	}
	
	flag.Parse()

	// Set debug mode in config package
	config.SetDebugMode(debugFlag)

	// Check if any flags were provided
	anyFlagProvided := debugFlag || localFlag || yoloFlag || happyFlag || resumeFlag || continueFlag || blankFlag
	
	// If any flag is provided, bypass TUI and launch directly with defaults
	if anyFlagProvided {
		// Apply default values: blank=true (skip TUI), all others keep their parsed values (defaulting to false)
		err := launcher.LaunchClaudeCodeWithoutMCP(yoloFlag, happyFlag, resumeFlag, continueFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error launching Claude Code: "+err.Error()))
			os.Exit(1)
		}
		return
	}

	// Find MCP files (command line flag takes precedence over TUI setting)
	mcpFiles, err := config.FindMCPFiles(localFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error finding MCP files: "+err.Error()))
		os.Exit(1)
	}

	if len(mcpFiles) == 0 {
		// Show styled no-MCP message and launch without MCP
		launcher.ShowNoMCPMessage(happyFlag)

		err := launcher.LaunchClaudeCodeWithoutMCP(yoloFlag, happyFlag, resumeFlag, continueFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error launching Claude Code: "+err.Error()))
			os.Exit(1)
		}
		return
	}

	// Create UI model and run Bubble Tea program
	m := ui.NewModel(mcpFiles, happyFlag)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("%s\n", ui.RenderError("Error running program: "+err.Error()))
		os.Exit(1)
	}

	// Launch Claude Code after Bubble Tea exits (only if user didn't quit)
	if finalModel, ok := finalModel.(ui.Model); ok && !finalModel.Quitted {
		// Use flags from TUI model, prioritizing resume over continue as specified
		effectiveResumeFlag := finalModel.ResumeFlag
		effectiveContinueFlag := finalModel.ContinueFlag
		if finalModel.ResumeFlag && finalModel.ContinueFlag {
			effectiveContinueFlag = false // resume takes priority
		}

		launcher.ShowLaunchMessage(finalModel.HappyFlag || happyFlag)
		err := launcher.LaunchClaudeCode(
			finalModel.Selected, 
			finalModel.MCPFiles, 
			finalModel.YoloFlag, 
			finalModel.HappyFlag || happyFlag, 
			effectiveResumeFlag, 
			effectiveContinueFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error launching Claude Code: "+err.Error()))
			os.Exit(1)
		}
	}
}