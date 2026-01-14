package app

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/chmouel/lazyworktree/internal/config"
	"github.com/chmouel/lazyworktree/internal/models"
)

func TestStatusUpdatedMsgUpdatesWorktreeStatus(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")

	wtPath := filepath.Join(cfg.WorktreeDir, "wt1")
	m.worktrees = []*models.WorktreeInfo{{Path: wtPath, Branch: "main"}}
	m.updateTable()

	msg := statusUpdatedMsg{
		statusFiles: []StatusFile{
			{Filename: "staged.txt", Status: "M."},
			{Filename: "modified.txt", Status: ".M"},
			{Filename: "new.txt", Status: " ?", IsUntracked: true},
		},
		path: wtPath,
	}

	_, _ = m.Update(msg)

	wt := m.worktrees[0]
	if !wt.Dirty {
		t.Fatal("expected worktree to be dirty")
	}
	if wt.Staged != 1 {
		t.Fatalf("expected staged count 1, got %d", wt.Staged)
	}
	if wt.Modified != 1 {
		t.Fatalf("expected modified count 1, got %d", wt.Modified)
	}
	if wt.Untracked != 1 {
		t.Fatalf("expected untracked count 1, got %d", wt.Untracked)
	}
}

func TestShouldRefreshGitEventDebounce(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")

	now := time.Now()
	if !m.shouldRefreshGitEvent(now) {
		t.Fatal("expected first refresh to pass")
	}
	if m.shouldRefreshGitEvent(now.Add(gitWatchDebounce / 2)) {
		t.Fatal("expected debounce to block refresh")
	}
	if !m.shouldRefreshGitEvent(now.Add(gitWatchDebounce + time.Millisecond)) {
		t.Fatal("expected refresh after debounce window")
	}
}
