package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"cc-launcher/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

// LaunchClaudeCode launches Claude Code with the specified MCP configuration files
func LaunchClaudeCode(selected map[int]struct{}, mcpFiles []string, yolo bool, happy bool, resume bool, continueFlag bool) error {
	var executablePath, executableName string
	
	// Check if happy flag is set and happy is available
	if happy {
		happyPath, err := exec.LookPath("happy")
		if err != nil {
			// Happy not found, show warning and fall back to claude
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Warning: 'happy' command not found in PATH. Falling back to 'claude'."))
			claudePath, err := exec.LookPath("claude")
			if err != nil {
				return fmt.Errorf("claude executable not found: %w", err)
			}
			executablePath = claudePath
			executableName = "claude"
		} else {
			executablePath = happyPath
			executableName = "happy"
		}
	} else {
		// Find the full path to claude executable
		claudePath, err := exec.LookPath("claude")
		if err != nil {
			return fmt.Errorf("claude executable not found: %w", err)
		}
		executablePath = claudePath
		executableName = "claude"
	}

	// Build arguments array
	args := []string{executableName}

	// Add --dangerously-skip-permissions if yolo flag is set
	if yolo {
		args = append(args, "--dangerously-skip-permissions")
	}

	// Add --resume if resume flag is set
	if resume {
		args = append(args, "--resume")
	}

	// Add --continue if continue flag is set
	if continueFlag {
		args = append(args, "--continue")
	}

	// Always add --strict-mcp-config to ensure only specified MCP servers are used
	args = append(args, "--strict-mcp-config")

	// If "No mcp servers" is selected, only add --strict-mcp-config (no --mcp-config)
	if _, noMcpSelected := selected[0]; !noMcpSelected {
		args = append(args, "--mcp-config")

		// Add all selected MCP files
		for i := range selected {
			if i > 0 && i-1 < len(mcpFiles) {
				args = append(args, mcpFiles[i-1])
			}
		}
	}

	// Use syscall.Exec to replace current process with Claude Code
	return syscall.Exec(executablePath, args, os.Environ())
}

// LaunchClaudeCodeWithoutMCP launches Claude Code without any MCP servers
func LaunchClaudeCodeWithoutMCP(yolo bool, happy bool, resume bool, continueFlag bool) error {
	var executablePath, executableName string
	
	// Check if happy flag is set and happy is available
	if happy {
		happyPath, err := exec.LookPath("happy")
		if err != nil {
			// Happy not found, show warning and fall back to claude
			fmt.Fprintf(os.Stderr, "%s\n", ui.RenderError("Warning: 'happy' command not found in PATH. Falling back to 'claude'."))
			claudePath, err := exec.LookPath("claude")
			if err != nil {
				return fmt.Errorf("claude executable not found: %w", err)
			}
			executablePath = claudePath
			executableName = "claude"
		} else {
			executablePath = happyPath
			executableName = "happy"
		}
	} else {
		// Find the full path to claude executable
		claudePath, err := exec.LookPath("claude")
		if err != nil {
			return fmt.Errorf("claude executable not found: %w", err)
		}
		executablePath = claudePath
		executableName = "claude"
	}

	// Build arguments array
	args := []string{executableName}

	// Add --dangerously-skip-permissions if yolo flag is set
	if yolo {
		args = append(args, "--dangerously-skip-permissions")
	}

	// Add --resume if resume flag is set
	if resume {
		args = append(args, "--resume")
	}

	// Add --continue if continue flag is set
	if continueFlag {
		args = append(args, "--continue")
	}

	// Always add --strict-mcp-config to ensure no MCP servers are loaded
	args = append(args, "--strict-mcp-config")

	// Use syscall.Exec to replace current process with Claude Code
	return syscall.Exec(executablePath, args, os.Environ())
}

// ShowNoMCPMessage displays a styled message when no MCP files are found
func ShowNoMCPMessage(happy bool) {
	title := ui.CreateGradientText("⚡ Claude Code Launcher", ui.PurpleGradientStart, ui.PurpleGradientEnd)
	fmt.Println(title)

	// Show confirmation if happy flag is set
	if happy {
		happyStyle := lipgloss.NewStyle().Foreground(ui.MutedColor)
		fmt.Println(happyStyle.Render("🎉 Happy mode enabled - will use 'happy' command instead of 'claude'"))
	}
	fmt.Println()
	noMcpStyle := lipgloss.NewStyle().
		Foreground(ui.MutedColor).
		Italic(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.MutedColor)

	launchMsg := ui.CreateGradientText("🚀 Launching Claude Code without MCP servers...", ui.PurpleGradientStart, ui.PurpleGradientEnd)

	fmt.Println(noMcpStyle.Render("📁 No MCP configuration files found in .claude/mcp/ or ~/.claude/mcp/"))
	fmt.Println()
	fmt.Println(launchMsg)
	fmt.Println()
}

// ShowLaunchMessage displays a styled message when launching Claude Code
func ShowLaunchMessage(happy bool) {
	var launchMsg string
	if happy {
		launchMsg = ui.CreateGradientText("🚀 Launching Claude Code (with happy)...", ui.PurpleGradientStart, ui.PurpleGradientEnd)
	} else {
		launchMsg = ui.CreateGradientText("🚀 Launching Claude Code...", ui.PurpleGradientStart, ui.PurpleGradientEnd)
	}
	fmt.Println()
	fmt.Println(launchMsg)
	fmt.Println()
}