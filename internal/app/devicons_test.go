package app

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIconFileInfoName(t *testing.T) {
	info := iconFileInfo{name: "example.go", isDir: false}
	assert.Equal(t, "example.go", info.Name())
}

func TestIconFileInfoSize(t *testing.T) {
	info := iconFileInfo{name: "file.txt", isDir: false}
	assert.Equal(t, int64(0), info.Size())
}

func TestIconFileInfoModeFile(t *testing.T) {
	info := iconFileInfo{name: "file.txt", isDir: false}
	assert.Equal(t, os.FileMode(0), info.Mode())
}

func TestIconFileInfoModeDirectory(t *testing.T) {
	info := iconFileInfo{name: "dir", isDir: true}
	expectedMode := os.ModeDir | 0o755
	assert.Equal(t, expectedMode, info.Mode())
}

func TestIconFileInfoModTime(t *testing.T) {
	info := iconFileInfo{name: "file.txt", isDir: false}
	assert.Equal(t, time.Time{}, info.ModTime())
}

func TestIconFileInfoIsDir(t *testing.T) {
	tests := []struct {
		name   string
		isDir  bool
		expect bool
	}{
		{"file", false, false},
		{"directory", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := iconFileInfo{name: tt.name, isDir: tt.isDir}
			assert.Equal(t, tt.expect, info.IsDir())
		})
	}
}

func TestIconFileInfoSys(t *testing.T) {
	info := iconFileInfo{name: "file.txt", isDir: false}
	assert.Nil(t, info.Sys())
}

func TestDeviconForNameEmpty(t *testing.T) {
	result := deviconForName("", false)
	assert.Empty(t, result)
}

func TestDeviconForNameFile(t *testing.T) {
	tests := []struct {
		name     string
		isDir    bool
		fileName string
	}{
		{"go file", false, "main.go"},
		{"python file", false, "script.py"},
		{"markdown file", false, "README.md"},
		{"directory", true, "src"},
		{"json file", false, "config.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deviconForName(tt.fileName, tt.isDir)
			// Result should not be empty for valid files
			// (actual icons depend on devicons library)
			assert.NotEmpty(t, result)
		})
	}
}

func TestCIIconForConclusionSuccess(t *testing.T) {
	result := ciIconForConclusion("success")
	assert.Equal(t, iconCISuccess, result)
}

func TestCIIconForConclusionFailure(t *testing.T) {
	result := ciIconForConclusion("failure")
	assert.Equal(t, iconCIFailure, result)
}

func TestCIIconForConclusionSkipped(t *testing.T) {
	result := ciIconForConclusion("skipped")
	assert.Equal(t, iconCISkipped, result)
}

func TestCIIconForConclusionCancelled(t *testing.T) {
	result := ciIconForConclusion("cancelled")
	assert.Equal(t, iconCICancelled, result)
}

func TestCIIconForConclusionPending(t *testing.T) {
	result := ciIconForConclusion("pending")
	assert.Equal(t, iconCIPending, result)
}

func TestCIIconForConclusionEmpty(t *testing.T) {
	result := ciIconForConclusion("")
	assert.Equal(t, iconCIPending, result)
}

func TestCIIconForConclusionUnknown(t *testing.T) {
	result := ciIconForConclusion("unknown_status")
	assert.Equal(t, iconCIUnknown, result)
}

func TestCIIconForConclusionAllStates(t *testing.T) {
	tests := []struct {
		conclusion string
		expected   string
	}{
		{"success", iconCISuccess},
		{"failure", iconCIFailure},
		{"skipped", iconCISkipped},
		{"cancelled", iconCICancelled},
		{"pending", iconCIPending},
		{"", iconCIPending},
		{"unknown", iconCIUnknown},
		{"random_value", iconCIUnknown},
	}

	for _, tt := range tests {
		t.Run("conclusion_"+tt.conclusion, func(t *testing.T) {
			result := ciIconForConclusion(tt.conclusion)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIconWithSpaceEmpty(t *testing.T) {
	result := iconWithSpace("")
	assert.Empty(t, result)
}

func TestIconWithSpaceWithIcon(t *testing.T) {
	// Test with a non-empty icon (use any non-empty string)
	result := iconWithSpace("test")
	assert.Equal(t, "test ", result)
}

func TestIconWithSpaceMultipleIcons(t *testing.T) {
	tests := []struct {
		icon     string
		expected string
	}{
		{"", ""},
		{"", " "},
		{"󰄱", "󰄱 "},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("icon_%d", i), func(t *testing.T) {
			result := iconWithSpace(tt.icon)
			// Empty icon returns empty string, non-empty returns icon with space
			if tt.icon == "" {
				assert.Empty(t, result)
			} else {
				assert.Equal(t, tt.icon+" ", result)
			}
		})
	}
}
