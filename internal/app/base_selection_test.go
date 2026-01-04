package app

import (
	"strings"
	"testing"

	"github.com/chmouel/lazyworktree/internal/config"
)

func TestParseBranchOptions(t *testing.T) {
	raw := strings.Join([]string{
		"main\trefs/heads/main",
		"feature/x\trefs/heads/feature/x",
		"origin/main\trefs/remotes/origin/main",
		"origin/HEAD\trefs/remotes/origin/HEAD",
		"",
	}, "\n")

	got := parseBranchOptions(raw)
	if len(got) != 3 {
		t.Fatalf("expected 3 branch options, got %d", len(got))
	}
	if got[0].name != mainWorktreeName || got[0].isRemote {
		t.Errorf("expected main to be local, got %+v", got[0])
	}
	if got[1].name != "feature/x" || got[1].isRemote {
		t.Errorf("expected feature/x to be local, got %+v", got[1])
	}
	if got[2].name != "origin/main" || !got[2].isRemote {
		t.Errorf("expected origin/main to be remote, got %+v", got[2])
	}
}

func TestPrioritizeBranchOptions(t *testing.T) {
	options := []branchOption{
		{name: "dev"},
		{name: mainWorktreeName},
		{name: "feature"},
	}
	got := prioritizeBranchOptions(options, mainWorktreeName)
	if len(got) != 3 {
		t.Fatalf("expected 3 options, got %d", len(got))
	}
	if got[0].name != mainWorktreeName || got[1].name != "dev" || got[2].name != "feature" {
		t.Errorf("unexpected order: %#v", got)
	}
}

func TestParseCommitOptions(t *testing.T) {
	raw := strings.Join([]string{
		"full1\x1fshort1\x1f2024-01-01\x1fFirst commit",
		"bad-line",
		"full2\x1fshort2\x1f2024-01-02\x1fSecond commit",
	}, "\n")

	got := parseCommitOptions(raw)
	if len(got) != 2 {
		t.Fatalf("expected 2 commit options, got %d", len(got))
	}
	if got[0].fullHash != "full1" || got[0].shortHash != "short1" || got[0].date != "2024-01-01" || got[0].subject != "First commit" {
		t.Errorf("unexpected first commit: %+v", got[0])
	}
	if got[1].fullHash != "full2" || got[1].shortHash != "short2" || got[1].date != "2024-01-02" || got[1].subject != "Second commit" {
		t.Errorf("unexpected second commit: %+v", got[1])
	}
}

func TestSuggestBranchNameWithExisting(t *testing.T) {
	existing := map[string]struct{}{
		"main":   {},
		"main-1": {},
		"dev":    {},
	}

	tests := []struct {
		name     string
		base     string
		expected string
	}{
		{
			name:     "empty base",
			base:     "",
			expected: "",
		},
		{
			name:     "base not taken",
			base:     "feature",
			expected: "feature",
		},
		{
			name:     "base taken uses suffix",
			base:     "dev",
			expected: "dev-1",
		},
		{
			name:     "increments suffix",
			base:     "main",
			expected: "main-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := suggestBranchNameWithExisting(tt.base, existing); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestSanitizeBranchNameFromTitle(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		fallback string
		expected string
	}{
		{
			name:     "basic title",
			title:    "Fix: Add new feature!",
			fallback: "abc123",
			expected: "fix-add-new-feature",
		},
		{
			name:     "limits length",
			title:    "This is a very long commit title that should be truncated to fifty characters",
			fallback: "abc123",
			expected: "this-is-a-very-long-commit-title-that-should-be-tr",
		},
		{
			name:     "empty uses fallback",
			title:    "!!!",
			fallback: "abc123",
			expected: "abc123",
		},
		{
			name:     "empty uses commit",
			title:    "",
			fallback: "",
			expected: "commit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeBranchNameFromTitle(tt.title, tt.fallback); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestBaseRefExists(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")

	if !m.baseRefExists("HEAD") {
		t.Fatal("expected HEAD to exist")
	}
	if m.baseRefExists("refs/does-not-exist") {
		t.Fatal("expected ref to not exist")
	}
}

func TestStripRemotePrefix(t *testing.T) {
	tests := []struct {
		name     string
		branch   string
		expected string
	}{
		{
			name:     "remote branch with origin",
			branch:   "origin/main",
			expected: "main",
		},
		{
			name:     "remote branch with other remote",
			branch:   "upstream/feature",
			expected: "feature",
		},
		{
			name:     "local branch",
			branch:   "main",
			expected: "main",
		},
		{
			name:     "branch with multiple slashes",
			branch:   "origin/feature/test",
			expected: "feature/test",
		},
		{
			name:     "empty string",
			branch:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripRemotePrefix(tt.branch); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
