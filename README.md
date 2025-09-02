# Claude Code Launcher

A terminal-based launcher for [Claude Code](https://claude.ai/code) that provides an interactive interface for selecting and launching MCP (Model Context Protocol) server configurations.

## Overview

Claude Code Launcher is a Go-based TUI (Terminal User Interface) application that streamlines the process of launching Claude Code with different MCP server configurations. Instead of manually specifying `--mcp-config` flags, this tool automatically discovers your MCP configuration files and presents them in an interactive selection menu.

## Features

- 🔍 **Automatic Discovery**: Scans `.claude/mcp/` directory for JSON configuration files
- 🎯 **Interactive Selection**: Choose which MCP servers to launch using a clean TUI interface
- ⚡ **Smart Defaults**: Launches Claude Code directly if no MCP configurations are found
- 🎨 **Beautiful Interface**: Built with Bubble Tea for a smooth terminal experience

## Platform Support

**This tool currently supports macOS only.** 

The launcher uses Unix system calls (`syscall.Exec`) to completely replace the launcher process with Claude Code, ensuring clean process management and proper terminal handling. This approach prevents issues with process suspension (Ctrl+C) and provides the smoothest user experience, but is not available on Windows.

## Installation

### Prerequisites

- **macOS** (required)
- Go 1.18 or higher
- Claude Code CLI installed (`claude` command available)

### Install from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/claude-code-launcher.git
cd claude-code-launcher

# Install the binary to $GOPATH/bin
go install

# The command 'cc-launcher' should now be available globally
```

Make sure `$(go env GOPATH)/bin` is in your PATH. You can verify this with:

```bash
echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "✓ GOPATH/bin is in PATH" || echo "✗ Add $(go env GOPATH)/bin to PATH"
```

If it's not in your PATH, add this to your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Usage

Simply run the launcher from any directory:

```bash
cc-launcher
```

### Workflow

1. The launcher scans for MCP configuration files in `.claude/mcp/*.json`
2. If MCP files are found, you'll see an interactive menu:
   - Use arrow keys to navigate
   - Press `Space` to select/deselect configurations
   - Press `Enter` to launch Claude Code with selected configurations
   - Press `q` or `Ctrl+C` to quit
3. If no MCP files are found, Claude Code launches directly

### Example

```bash
$ cc-launcher

Select MCP configurations to launch:
  [ ] taskmaster.json
  [x] shadcn-ui.json
  [ ] context7.json
  [ ] No mcp servers

Press Space to select, Enter to launch, q to quit
```

## MCP Configuration

Place your MCP server configuration files in the `.claude/mcp/` directory relative to your current working directory. Each configuration should be a valid JSON file.

Example structure:
```
.claude/
└── mcp/
    ├── taskmaster.json
    ├── shadcn-ui.json
    └── context7.json
```

## Development

### Building from Source

```bash
# Build the binary
make build

# Run directly without building
make run

# Clean build artifacts
make clean
```

Available make targets:
- `make build` - Build the application binary (`cc-launcher`)
- `make clean` - Remove build artifacts
- `make run` - Run the application directly without building

### Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
