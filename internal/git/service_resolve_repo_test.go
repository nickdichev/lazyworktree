package git

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveRepoName(t *testing.T) {
	// Note: These tests modify the process working directory, so they cannot run in parallel.

	t.Run("resolve from github remote url with .git", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := exec.Command("git", "init")
		cmd.Dir = tmpDir
		require.NoError(t, cmd.Run())

		cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/owner/repo.git")
		cmd.Dir = tmpDir
		require.NoError(t, cmd.Run())

		oldWd, err := os.Getwd()
		require.NoError(t, err)
		defer func() { _ = os.Chdir(oldWd) }()

		require.NoError(t, os.Chdir(tmpDir))

		notify := func(_ string, _ string) {}
		notifyOnce := func(_ string, _ string, _ string) {}
		service := NewService(notify, notifyOnce)

		name := service.ResolveRepoName(context.Background())
		assert.Equal(t, "owner/repo", name)
	})

	t.Run("resolve from gitlab remote url with .git", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := exec.Command("git", "init")
		cmd.Dir = tmpDir
		require.NoError(t, cmd.Run())

		cmd = exec.Command("git", "remote", "add", "origin", "https://gitlab.com/group/subgroup/project.git")
		cmd.Dir = tmpDir
		require.NoError(t, cmd.Run())

		oldWd, err := os.Getwd()
		require.NoError(t, err)
		defer func() { _ = os.Chdir(oldWd) }()
		require.NoError(t, os.Chdir(tmpDir))

		notify := func(_ string, _ string) {}
		notifyOnce := func(_ string, _ string, _ string) {}
		service := NewService(notify, notifyOnce)

		name := service.ResolveRepoName(context.Background())
		// The regex matches the part after the first slash in the path after domain
		// For gitlab.com/group/subgroup/project.git, it matches group/subgroup/project.git
		// And then trims .git
		assert.Equal(t, "group/subgroup/project", name)
	})

	t.Run("resolve from remote url without .git", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := exec.Command("git", "init")
		cmd.Dir = tmpDir
		require.NoError(t, cmd.Run())

		cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/owner/repo")
		cmd.Dir = tmpDir
		require.NoError(t, cmd.Run())

		oldWd, err := os.Getwd()
		require.NoError(t, err)
		defer func() { _ = os.Chdir(oldWd) }()
		require.NoError(t, os.Chdir(tmpDir))

		notify := func(_ string, _ string) {}
		notifyOnce := func(_ string, _ string, _ string) {}
		service := NewService(notify, notifyOnce)

		name := service.ResolveRepoName(context.Background())
		assert.Equal(t, "owner/repo", name)
	})
}
