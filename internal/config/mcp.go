package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var debugMode bool

// SetDebugMode enables or disables debug logging
func SetDebugMode(enabled bool) {
	debugMode = enabled
}

// FindMCPFiles discovers MCP configuration files in both local and global directories.
// It returns a slice of file paths to *.json files found in:
// - .claude/mcp/ (local directory)
// - ~/.claude/mcp/ (global directory)
//
// Returns an error if there are issues accessing the user's home directory
// or if glob operations fail unexpectedly.
func FindMCPFiles() ([]string, error) {
	var mcpFiles []string

	// Scan local directory (.claude/mcp)
	localDir := ".claude/mcp"
	files, err := scanMCPDirectory(localDir)
	if err != nil {
		return nil, fmt.Errorf("failed to scan local MCP directory %s: %w", localDir, err)
	}
	mcpFiles = append(mcpFiles, files...)

	// Scan global directory (~/.claude/mcp)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	globalDir := filepath.Join(homeDir, ".claude", "mcp")
	files, err = scanMCPDirectory(globalDir)
	if err != nil {
		return nil, fmt.Errorf("failed to scan global MCP directory %s: %w", globalDir, err)
	}
	mcpFiles = append(mcpFiles, files...)

	return mcpFiles, nil
}

// scanMCPDirectory scans a specific directory for *.json files.
// It returns an empty slice if the directory doesn't exist or is inaccessible,
// logging appropriate warnings only when debug mode is enabled. Returns an error only for unexpected glob failures.
func scanMCPDirectory(dir string) ([]string, error) {
	// Check if directory exists and is accessible
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			// Directory doesn't exist - this is normal, not an error
			return []string{}, nil
		}
		// Directory exists but is not accessible
		if debugMode {
			log.Printf("Warning: MCP directory %s exists but is not accessible: %v", dir, err)
		}
		return []string{}, nil
	}

	// Directory exists and is accessible, scan for JSON files
	pattern := filepath.Join(dir, "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		// Glob operation failed unexpectedly
		return nil, fmt.Errorf("failed to glob pattern %s: %w", pattern, err)
	}

	if debugMode {
		if len(files) == 0 {
			log.Printf("Info: No MCP configuration files found in %s", dir)
		} else {
			log.Printf("Info: Found %d MCP configuration file(s) in %s", len(files), dir)
		}
	}

	return files, nil
}