# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based launcher utility for Claude Code that provides an interactive TUI interface for selecting MCP (Model Context Protocol) server configurations. The launcher automatically detects MCP configuration files and allows users to choose which servers to launch with Claude Code.

## Architecture

- **Single binary application** built with Go
- **Bubble Tea TUI framework** for interactive terminal interface
- **MCP configuration discovery**: Scans `.claude/mcp/` directory for `*.json` configuration files
- **Claude Code integration**: Launches the `claude` command with appropriate `--mcp-config` flags

## Key Components

- `main.go`: Contains the entire application logic
  - `model` struct: Bubble Tea model managing UI state and MCP file selection
  - `findMCPFiles()`: Discovers MCP configuration files in `.claude/mcp/`
  - `launchClaudeCodeFromSelection()`: Executes Claude Code with selected MCP configurations

## Development Commands

```bash
# Build the application
go build -o claude-launcher

# Run directly
go run main.go

# Clean build artifacts
rm claude-launcher
```

## Usage Flow

1. Application scans for MCP configuration files in `.claude/mcp/*.json`
2. If no MCP files found: launches Claude Code directly without MCP servers
3. If MCP files found: presents interactive selection interface
4. User selects desired MCP configurations or "No mcp servers"
5. Launches Claude Code with `--mcp-config` flags for selected configurations

## Dependencies

- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/bubbles`: UI components
- `github.com/charmbracelet/lipgloss`: Styling library
- Standard Go libraries: `os`, `exec`, `path/filepath`, `strings`