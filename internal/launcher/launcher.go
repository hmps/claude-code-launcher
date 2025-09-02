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
func LaunchClaudeCode(selected map[int]struct{}, mcpFiles []string, yolo bool) error {
	// Find the full path to claude executable
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude executable not found: %w", err)
	}

	// Build arguments array
	args := []string{"claude"}

	// Add --dangerously-skip-permissions if yolo flag is set
	if yolo {
		args = append(args, "--dangerously-skip-permissions")
	}

	// If "No mcp servers" is selected
	if _, noMcpSelected := selected[0]; noMcpSelected {
		// args already contains yolo flag if needed, no need to add anything else
	} else {
		args = append(args, "--mcp-config")

		// Add all selected MCP files
		for i := range selected {
			if i > 0 && i-1 < len(mcpFiles) {
				args = append(args, mcpFiles[i-1])
			}
		}
	}

	// Use syscall.Exec to replace current process with Claude Code
	return syscall.Exec(claudePath, args, os.Environ())
}

// LaunchClaudeCodeWithoutMCP launches Claude Code without any MCP servers
func LaunchClaudeCodeWithoutMCP(yolo bool) error {
	// Find the full path to claude executable
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude executable not found: %w", err)
	}

	// Build arguments array
	args := []string{"claude"}

	// Add --dangerously-skip-permissions if yolo flag is set
	if yolo {
		args = append(args, "--dangerously-skip-permissions")
	}

	// Use syscall.Exec to replace current process with Claude Code
	return syscall.Exec(claudePath, args, os.Environ())
}

// ShowNoMCPMessage displays a styled message when no MCP files are found
func ShowNoMCPMessage() {
	title := ui.CreateGradientText("⚡ Claude Code Launcher", ui.PurpleGradientStart, ui.PurpleGradientEnd)
	noMcpStyle := lipgloss.NewStyle().
		Foreground(ui.MutedColor).
		Italic(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.MutedColor)

	launchMsg := ui.CreateGradientText("🚀 Launching Claude Code without MCP servers...", ui.PurpleGradientStart, ui.PurpleGradientEnd)

	fmt.Println(title)
	fmt.Println()
	fmt.Println(noMcpStyle.Render("📁 No MCP configuration files found in .claude/mcp/ or ~/.claude/mcp/"))
	fmt.Println()
	fmt.Println(launchMsg)
	fmt.Println()
}

// ShowLaunchMessage displays a styled message when launching Claude Code
func ShowLaunchMessage() {
	launchMsg := ui.CreateGradientText("🚀 Launching Claude Code...", ui.PurpleGradientStart, ui.PurpleGradientEnd)
	fmt.Println()
	fmt.Println(launchMsg)
	fmt.Println()
}