package main

import (
	"flag"
	"os"
	"testing"
)

func TestHandleWtCreateValidation(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "missing both flags",
			args:        []string{},
			expectError: true,
			errorMsg:    "must specify either --from-branch or --from-pr",
		},
		{
			name:        "both flags specified",
			args:        []string{"--from-branch", "main", "--from-pr", "123"},
			expectError: true,
			errorMsg:    "mutually exclusive",
		},
		{
			name:        "with-change with from-pr",
			args:        []string{"--from-pr", "123", "--with-change"},
			expectError: true,
			errorMsg:    "--with-change can only be used with --from-branch",
		},
		{
			name:        "valid from-branch",
			args:        []string{"--from-branch", "main"},
			expectError: false,
		},
		{
			name:        "valid from-pr",
			args:        []string{"--from-pr", "123"},
			expectError: false,
		},
		{
			name:        "valid from-branch with with-change",
			args:        []string{"--from-branch", "main", "--with-change"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a flag set to test validation logic
			fs := flag.NewFlagSet("wt-create", flag.ContinueOnError)
			fromBranch := fs.String("from-branch", "", "")
			fromPR := fs.Int("from-pr", 0, "")
			withChange := fs.Bool("with-change", false, "")
			_ = fs.Bool("silent", false, "")

			// Capture stderr
			oldStderr := os.Stderr
			_, w, _ := os.Pipe()
			os.Stderr = w

			err := fs.Parse(tt.args)
			if err != nil {
				_ = w.Close()
				os.Stderr = oldStderr
				if !tt.expectError {
					t.Errorf("unexpected parse error: %v", err)
				}
				return
			}

			// Test validation logic
			hasError := false
			switch {
			case *fromBranch != "" && *fromPR > 0:
				hasError = true
			case *fromBranch == "" && *fromPR == 0:
				hasError = true
			case *withChange && *fromPR > 0:
				hasError = true
			}

			_ = w.Close()
			os.Stderr = oldStderr

			if tt.expectError && !hasError {
				t.Error("expected validation error but got none")
			}
			if !tt.expectError && hasError {
				t.Error("unexpected validation error")
			}
		})
	}
}

func TestHandleWtDeleteFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		noBranch bool
		silent   bool
		worktree string
	}{
		{
			name:     "default flags",
			args:     []string{},
			noBranch: false,
			silent:   false,
		},
		{
			name:     "no-branch flag",
			args:     []string{"--no-branch"},
			noBranch: true,
			silent:   false,
		},
		{
			name:     "silent flag",
			args:     []string{"--silent"},
			noBranch: false,
			silent:   true,
		},
		{
			name:     "worktree path",
			args:     []string{"/path/to/worktree"},
			noBranch: false,
			silent:   false,
			worktree: "/path/to/worktree",
		},
		{
			name:     "all flags and path",
			args:     []string{"--no-branch", "--silent", "/path/to/worktree"},
			noBranch: true,
			silent:   true,
			worktree: "/path/to/worktree",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("wt-delete", flag.ContinueOnError)
			noBranch := fs.Bool("no-branch", false, "")
			silent := fs.Bool("silent", false, "")

			if err := fs.Parse(tt.args); err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}

			if *noBranch != tt.noBranch {
				t.Errorf("noBranch = %v, want %v", *noBranch, tt.noBranch)
			}
			if *silent != tt.silent {
				t.Errorf("silent = %v, want %v", *silent, tt.silent)
			}

			var worktreePath string
			if len(fs.Args()) > 0 {
				worktreePath = fs.Args()[0]
			}
			if worktreePath != tt.worktree {
				t.Errorf("worktreePath = %q, want %q", worktreePath, tt.worktree)
			}
		})
	}
}
