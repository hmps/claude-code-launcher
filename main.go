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
	var configFlag bool
	var zaiFlag bool
	flag.BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	flag.BoolVar(&localFlag, "local", false, "Only check for local MCP configurations, skip global ones")
	flag.BoolVar(&yoloFlag, "yolo", false, "Launch Claude Code with --dangerously-skip-permissions")
	flag.BoolVar(&happyFlag, "happy", false, "Use happy instead of claude command")
	flag.BoolVar(&resumeFlag, "r", false, "Launch Claude Code with --resume flag (-r, --resume)")
	flag.BoolVar(&resumeFlag, "resume", false, "Launch Claude Code with --resume flag (-r, --resume)")
	flag.BoolVar(&continueFlag, "continue", false, "Launch Claude Code with --continue flag (--continue)")
	flag.BoolVar(&configFlag, "c", false, "Always show TUI config interface (overrides other flags) (-c, --config)")
	flag.BoolVar(&configFlag, "config", false, "Always show TUI config interface (overrides other flags) (-c, --config)")
	flag.BoolVar(&blankFlag, "b", false, "Launch Claude Code without MCP servers (skip TUI) (-b, --blank)")
	flag.BoolVar(&blankFlag, "blank", false, "Launch Claude Code without MCP servers (skip TUI) (-b, --blank)")
	flag.BoolVar(&zaiFlag, "zai", false, "Use z.ai coding plan (requires Z_AI_API_KEY environment variable)")
	
	// Custom usage function to show double dashes for all flags except -r, -c, and -b
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  -b, --blank\n        Launch Claude Code without MCP servers (skip TUI)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -c, --config\n        Always show TUI config interface (overrides other flags)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --continue\n        Launch Claude Code with --continue flag\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --debug\n        Enable debug logging\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --happy\n        Use happy instead of claude command\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --local\n        Only check for local MCP configurations, skip global ones\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  -r, --resume\n        Launch Claude Code with --resume flag\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --yolo\n        Launch Claude Code with --dangerously-skip-permissions\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  --zai\n        Use z.ai coding plan (requires Z_AI_API_KEY environment variable)\n")
	}
	
	flag.Parse()

	// Validate zai flag requires Z_AI_API_KEY
	if zaiFlag {
		zaiAPIKey := os.Getenv("Z_AI_API_KEY")
		if zaiAPIKey == "" {
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error: --zai flag requires Z_AI_API_KEY environment variable to be set"))
			os.Exit(1)
		}
	}

	// Check if Z_AI_API_KEY is available for TUI
	zaiAvailable := os.Getenv("Z_AI_API_KEY") != ""

	// Set debug mode in config package
	config.SetDebugMode(debugFlag)

	// Check if any flags were provided (excluding configFlag which forces TUI)
	anyFlagProvided := debugFlag || localFlag || yoloFlag || happyFlag || resumeFlag || continueFlag || blankFlag || zaiFlag
	
	// If any flag is provided but config flag is NOT set, bypass TUI and launch directly with defaults
	if anyFlagProvided && !configFlag {
		// Apply default values: blank=true (skip TUI), all others keep their parsed values (defaulting to false)
		err := launcher.LaunchClaudeCodeWithoutMCP(yoloFlag, happyFlag, resumeFlag, continueFlag, zaiFlag)
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

	// If config flag is set, always show TUI even with no MCP files
	if configFlag && len(mcpFiles) == 0 {
		// Create empty MCP files list to force TUI
		mcpFiles = []string{}
	}

	if len(mcpFiles) == 0 && !configFlag {
		// Show styled no-MCP message and launch without MCP
		launcher.ShowNoMCPMessage(happyFlag)

		err := launcher.LaunchClaudeCodeWithoutMCP(yoloFlag, happyFlag, resumeFlag, continueFlag, zaiFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error launching Claude Code: "+err.Error()))
			os.Exit(1)
		}
		return
	}

	// Create UI model and run Bubble Tea program
	// Pass command line flags as defaults when config flag is used
	m := ui.NewModelWithDefaults(mcpFiles, happyFlag, yoloFlag, continueFlag, resumeFlag, blankFlag, zaiFlag, zaiAvailable)
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
			effectiveContinueFlag,
			finalModel.ZaiFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Error launching Claude Code: "+err.Error()))
			os.Exit(1)
		}
	}
}