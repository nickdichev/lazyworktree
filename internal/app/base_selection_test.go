package app

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

func TestParseBranchOptionsWithDate(t *testing.T) {
	now := time.Now().Unix()
	yesterday := now - 86400
	raw := strings.Join([]string{
		fmt.Sprintf("main\trefs/heads/main\t%d", now),
		fmt.Sprintf("feature/x\trefs/heads/feature/x\t%d", yesterday),
		fmt.Sprintf("origin/main\trefs/remotes/origin/main\t%d", now),
		fmt.Sprintf("v1.0.0\trefs/tags/v1.0.0\t%d", yesterday),
		fmt.Sprintf("origin/HEAD\trefs/remotes/origin/HEAD\t%d", now),
		"",
	}, "\n")

	got := parseBranchOptionsWithDate(raw)
	if len(got) != 4 {
		t.Fatalf("expected 4 branch/tag options, got %d", len(got))
	}
	if got[0].name != mainWorktreeName || got[0].isRemote || got[0].isTag {
		t.Errorf("expected main to be local branch, got %+v", got[0])
	}
	if got[0].committerDate.IsZero() {
		t.Errorf("expected main to have a commit date")
	}
	if got[1].name != "feature/x" || got[1].isRemote || got[1].isTag {
		t.Errorf("expected feature/x to be local branch, got %+v", got[1])
	}
	if got[2].name != "origin/main" || !got[2].isRemote || got[2].isTag {
		t.Errorf("expected origin/main to be remote branch, got %+v", got[2])
	}
	if got[3].name != "v1.0.0" || got[3].isRemote || !got[3].isTag {
		t.Errorf("expected v1.0.0 to be tag, got %+v", got[3])
	}
}

func TestSortBranchOptions(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	tests := []struct {
		name     string
		input    []branchOption
		expected []string
	}{
		{
			name: "local main first",
			input: []branchOption{
				{name: "feature", isRemote: false, committerDate: now},
				{name: "main", isRemote: false, committerDate: yesterday},
				{name: "dev", isRemote: false, committerDate: lastWeek},
			},
			expected: []string{"main", "feature", "dev"},
		},
		{
			name: "local master when no main",
			input: []branchOption{
				{name: "feature", isRemote: false, committerDate: now},
				{name: "master", isRemote: false, committerDate: yesterday},
				{name: "dev", isRemote: false, committerDate: lastWeek},
			},
			expected: []string{"master", "feature", "dev"},
		},
		{
			name: "remote origin/main after local main",
			input: []branchOption{
				{name: "feature", isRemote: false, committerDate: now},
				{name: "main", isRemote: false, committerDate: yesterday},
				{name: "origin/main", isRemote: true, committerDate: now},
			},
			expected: []string{"main", "origin/main", "feature"},
		},
		{
			name: "date sorting for others",
			input: []branchOption{
				{name: "old-branch", isRemote: false, committerDate: lastWeek},
				{name: "new-branch", isRemote: false, committerDate: now},
				{name: "mid-branch", isRemote: false, committerDate: yesterday},
			},
			expected: []string{"new-branch", "mid-branch", "old-branch"},
		},
		{
			name: "same date alphabetical tiebreaker",
			input: []branchOption{
				{name: "zebra", isRemote: false, committerDate: now},
				{name: "alpha", isRemote: false, committerDate: now},
				{name: "beta", isRemote: false, committerDate: now},
			},
			expected: []string{"alpha", "beta", "zebra"},
		},
		{
			name:     "empty list",
			input:    []branchOption{},
			expected: []string{},
		},
		{
			name: "all priority branches",
			input: []branchOption{
				{name: "feature", isRemote: false, committerDate: now},
				{name: "origin/master", isRemote: true, committerDate: now},
				{name: "main", isRemote: false, committerDate: yesterday},
				{name: "origin/main", isRemote: true, committerDate: now},
				{name: "master", isRemote: false, committerDate: lastWeek},
			},
			expected: []string{"main", "origin/main", "origin/master", "feature"},
		},
		{
			name: "tags mixed with branches by date",
			input: []branchOption{
				{name: "feature", isRemote: false, isTag: false, committerDate: now},
				{name: "v1.0.0", isRemote: false, isTag: true, committerDate: yesterday},
				{name: "main", isRemote: false, isTag: false, committerDate: lastWeek},
				{name: "v0.9.0", isRemote: false, isTag: true, committerDate: lastWeek.Add(-24 * time.Hour)},
				{name: "dev", isRemote: false, isTag: false, committerDate: yesterday.Add(-12 * time.Hour)},
			},
			expected: []string{"main", "feature", "v1.0.0", "dev", "v0.9.0"},
		},
		{
			name: "tags don't appear in priority positions",
			input: []branchOption{
				{name: "main", isRemote: false, isTag: true, committerDate: now},
				{name: "feature", isRemote: false, isTag: false, committerDate: yesterday},
				{name: "main", isRemote: false, isTag: false, committerDate: lastWeek},
			},
			expected: []string{"main", "main", "feature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortBranchOptions(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d branches, got %d", len(tt.expected), len(got))
			}
			for i, expected := range tt.expected {
				if got[i].name != expected {
					t.Errorf("at index %d: expected %q, got %q", i, expected, got[i].name)
				}
			}
		})
	}
}
