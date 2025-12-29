package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chmouel/lazyworktree/internal/app"
	"github.com/chmouel/lazyworktree/internal/config"
)

func main() {
	var worktreeDir string
	var debugLog string

	flag.StringVar(&worktreeDir, "worktree-dir", "", "Override the default worktree root directory")
	flag.StringVar(&debugLog, "debug-log", "", "Path to debug log file")
	flag.Parse()

	// Get initial filter from remaining args
	initialFilter := strings.Join(flag.Args(), " ")

	// Load config
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		cfg = config.DefaultConfig()
	}

	// Override worktree_dir if provided via flag
	if worktreeDir != "" {
		expanded, err := expandPath(worktreeDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error expanding worktree-dir: %v\n", err)
			os.Exit(1)
		}
		cfg.WorktreeDir = expanded
	} else if cfg.WorktreeDir != "" {
		expanded, err := expandPath(cfg.WorktreeDir)
		if err == nil {
			cfg.WorktreeDir = expanded
		}
	} else {
		// Default fallback
		home, _ := os.UserHomeDir()
		cfg.WorktreeDir = filepath.Join(home, ".local", "share", "worktrees")
	}

	// Set debug log if provided
	if debugLog != "" {
		expanded, err := expandPath(debugLog)
		if err == nil {
			cfg.DebugLog = expanded
		} else {
			cfg.DebugLog = debugLog
		}
	}

	// Create and run the app
	model := app.NewAppModel(cfg, initialFilter)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}

	// Output selected path for shell integration
	selectedPath := model.GetSelectedPath()
	if selectedPath != "" {
		fmt.Println(selectedPath)
	}
}

// expandPath expands ~ and environment variables in a path
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[1:])
	}
	return os.ExpandEnv(path), nil
}
