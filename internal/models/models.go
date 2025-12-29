package models

// PRInfo represents information about a Pull Request or Merge Request
type PRInfo struct {
	Number int
	State  string
	Title  string
	URL    string
}

// WorktreeInfo represents information about a Git worktree
type WorktreeInfo struct {
	Path         string
	Branch       string
	IsMain       bool
	Dirty        bool
	Ahead        int
	Behind       int
	LastActive   string
	LastActiveTS int64
	PR           *PRInfo
	Untracked    int
	Modified     int
	Staged       int
	Divergence   string
}

// Constants for cache and state files
const (
	LastSelectedFilename = ".last-selected"
	CacheFilename        = ".worktree-cache.json"
)
