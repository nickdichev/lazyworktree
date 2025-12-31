// Package models defines the data objects shared across lazyworktree packages.
package models

// PRInfo captures the relevant metadata for a pull request.
type PRInfo struct {
	Number int
	State  string
	Title  string
	URL    string
}

// WorktreeInfo summarizes the information for a git worktree.
type WorktreeInfo struct {
	Path         string
	Branch       string
	IsMain       bool
	Dirty        bool
	Ahead        int
	Behind       int
	HasUpstream  bool
	LastActive   string
	LastActiveTS int64
	PR           *PRInfo
	Untracked    int
	Modified     int
	Staged       int
	Divergence   string
}

const (
	// LastSelectedFilename stores the last worktree selection for a repo.
	LastSelectedFilename = ".last-selected"
	// CacheFilename stores cached worktree metadata for faster loads.
	CacheFilename = ".worktree-cache.json"
)
