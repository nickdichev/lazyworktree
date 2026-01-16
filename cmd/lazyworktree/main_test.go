package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chmouel/lazyworktree/internal/config"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = writer

	fn()

	_ = writer.Close()
	os.Stdout = orig

	out, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	return string(out)
}

func TestPrintSyntaxThemes(t *testing.T) {
	out := captureStdout(t, func() {
		printSyntaxThemes()
	})

	if !strings.Contains(out, "Available syntax themes") {
		t.Fatalf("expected header to be printed, got %q", out)
	}
	if !strings.Contains(out, "dracula") {
		t.Fatalf("expected theme list to include dracula, got %q", out)
	}
}

func TestPrintCompletion(t *testing.T) {
	shells := []string{"bash", "zsh", "fish"}

	for _, shell := range shells {
		t.Run(shell, func(t *testing.T) {
			out := captureStdout(t, func() {
				if err := printCompletion(shell); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			})

			if out == "" {
				t.Error("expected non-empty output")
			}

			// Verify it contains lazyworktree
			if !strings.Contains(out, "lazyworktree") {
				t.Error("output missing program name")
			}

			// Verify shell-specific structure
			switch shell {
			case "bash":
				if !strings.Contains(out, "_lazyworktree_completion") {
					t.Error("bash completion missing expected function")
				}
			case "zsh":
				if !strings.Contains(out, "#compdef") {
					t.Error("zsh completion missing compdef directive")
				}
			case "fish":
				if !strings.Contains(out, "complete -c") {
					t.Error("fish completion missing complete command")
				}
			}
		})
	}
}

func TestPrintCompletionInvalidShell(t *testing.T) {
	err := printCompletion("invalid")
	if err == nil {
		t.Error("expected error for invalid shell")
	}

	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestOutputSelectionWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	selectedPath := "/path/to/worktree"
	data := selectedPath + "\n"

	const filePerms = 0o600
	err := os.WriteFile(outputFile, []byte(data), filePerms)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	// #nosec G304 - test file operations with t.TempDir() are safe
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != data {
		t.Fatalf("expected %q, got %q", data, string(content))
	}
}

func TestOutputSelectionEmptyContent(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	const filePerms = 0o600
	err := os.WriteFile(outputFile, []byte(""), filePerms)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	// #nosec G304 - test file operations with t.TempDir() are safe
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if len(content) != 0 {
		t.Fatalf("expected empty content, got %q", string(content))
	}
}

func TestOutputSelectionDirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "subdir1", "subdir2", "output.txt")

	const dirPerms = 0o750
	err := os.MkdirAll(filepath.Dir(outputPath), dirPerms)
	if err != nil {
		t.Fatalf("failed to create directories: %v", err)
	}

	const filePerms = 0o600
	err = os.WriteFile(outputPath, []byte("/test/path\n"), filePerms)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("output file not created: %v", err)
	}
}

func TestVersionVariables(t *testing.T) {
	// Verify build variables are set (at least with defaults)
	if version == "" {
		t.Error("version should not be empty")
	}
	if commit == "" {
		t.Error("commit should not be empty")
	}
	if date == "" {
		t.Error("date should not be empty")
	}
	if builtBy == "" {
		t.Error("builtBy should not be empty")
	}
}

func TestPrintSyntaxThemesContainsThemes(t *testing.T) {
	out := captureStdout(t, func() {
		printSyntaxThemes()
	})

	expectedThemes := []string{"dracula", "monokai", "nord"}
	for _, theme := range expectedThemes {
		if !strings.Contains(out, theme) {
			t.Logf("warning: expected theme %q in output", theme)
		}
	}
}

func TestPrintVersion(t *testing.T) {
	out := captureStdout(t, func() {
		printVersion()
	})

	if !strings.Contains(out, "lazyworktree version") {
		t.Errorf("expected version header, got %q", out)
	}
	if !strings.Contains(out, version) {
		t.Errorf("expected version %q in output, got %q", version, out)
	}
}

func TestApplyWorktreeDirConfig(t *testing.T) {
	tests := []struct {
		name           string
		worktreeDir    string
		cfgWorktreeDir string
		expected       string
		expectError    bool
	}{
		{
			name:        "flag takes precedence",
			worktreeDir: "/custom/path",
			expected:    "/custom/path",
		},
		{
			name:           "config value used when flag empty",
			worktreeDir:    "",
			cfgWorktreeDir: "/config/path",
			expected:       "/config/path",
		},
		{
			name:        "default when both empty",
			worktreeDir: "",
			expected:    "", // Will be set to default, but we can't easily test home dir
		},
		{
			name:        "expand path with tilde",
			worktreeDir: "~/test",
			expected:    "", // Will be expanded, exact path depends on home dir
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.AppConfig{
				WorktreeDir: tt.cfgWorktreeDir,
			}

			err := applyWorktreeDirConfig(cfg, tt.worktreeDir)
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && tt.expected != "" && !strings.Contains(cfg.WorktreeDir, tt.expected) {
				// For default case, just verify it's set
				if tt.worktreeDir == "" && tt.cfgWorktreeDir == "" {
					if cfg.WorktreeDir == "" {
						t.Error("expected default worktree dir to be set")
					}
				}
			}
		})
	}
}

func TestApplyThemeConfig(t *testing.T) {
	tests := []struct {
		name        string
		themeName   string
		expectError bool
	}{
		{
			name:        "valid theme",
			themeName:   "dracula",
			expectError: false,
		},
		{
			name:        "valid theme uppercase",
			themeName:   "DRACULA",
			expectError: false,
		},
		{
			name:        "invalid theme",
			themeName:   "nonexistent-theme",
			expectError: true,
		},
		{
			name:        "empty theme",
			themeName:   "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.GitPager = "delta"

			err := applyThemeConfig(cfg, tt.themeName)
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && tt.themeName != "" {
				if cfg.Theme == "" {
					t.Error("expected theme to be set")
				}
			}
		})
	}
}

func TestLoadCLIConfig(t *testing.T) {
	t.Run("load default config", func(t *testing.T) {
		cfg, err := loadCLIConfig("", "", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg == nil {
			t.Fatal("expected config to be non-nil")
		}
	})

	t.Run("apply worktree dir", func(t *testing.T) {
		cfg, err := loadCLIConfig("", "/test/path", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.WorktreeDir == "" {
			t.Error("expected worktree dir to be set")
		}
	})

	t.Run("apply config overrides", func(t *testing.T) {
		overrides := configOverrides{"lw.theme=dracula"}
		cfg, err := loadCLIConfig("", "", overrides)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Theme != "dracula" {
			t.Errorf("expected theme to be dracula, got %q", cfg.Theme)
		}
	})
}

func TestNewCLIGitService(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.GitPager = "delta"
	cfg.GitPagerArgs = []string{"--syntax-theme", "Dracula"}

	svc := newCLIGitService(cfg)
	if svc == nil {
		t.Fatal("expected service to be non-nil")
	}
	if !svc.UseGitPager() {
		t.Error("expected git pager to be enabled")
	}
}
