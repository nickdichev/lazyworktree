package app

import (
	"strings"
	"testing"

	"github.com/chmouel/lazyworktree/internal/models"
)

func TestGeneratePRWorktreeName(t *testing.T) {
	tests := []struct {
		name     string
		pr       *models.PRInfo
		expected string
	}{
		{
			name: "simple title",
			pr: &models.PRInfo{
				Number: 123,
				Title:  "Add feature",
			},
			expected: "pr123-add-feature",
		},
		{
			name: "title with special characters",
			pr: &models.PRInfo{
				Number: 2367,
				Title:  "Feat: Add one-per-pipeline comment strategy!",
			},
			expected: "pr2367-feat-add-one-per-pipeline-comment-strategy",
		},
		{
			name: "long title gets truncated",
			pr: &models.PRInfo{
				Number: 999,
				Title:  "This is a very long title that should be truncated because it exceeds the maximum length limit of one hundred characters total including the pr prefix",
			},
			// Result should be exactly 100 chars or less, and not end with hyphen
			expected: "pr999-this-is-a-very-long-title-that-should-be-truncated-because-it-exceeds-the-maximum-length-limit",
		},
		{
			name: "title with multiple spaces",
			pr: &models.PRInfo{
				Number: 456,
				Title:  "Fix   multiple    spaces",
			},
			expected: "pr456-fix-multiple-spaces",
		},
		{
			name: "title with numbers and symbols",
			pr: &models.PRInfo{
				Number: 789,
				Title:  "Update v2.0 API (breaking changes)",
			},
			expected: "pr789-update-v2-0-api-breaking-changes",
		},
		{
			name: "empty title",
			pr: &models.PRInfo{
				Number: 100,
				Title:  "",
			},
			expected: "pr100",
		},
		{
			name: "title with only special characters",
			pr: &models.PRInfo{
				Number: 200,
				Title:  "!!!@@@###$$$",
			},
			expected: "pr200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generatePRWorktreeName(tt.pr)
			// For the long title test, just verify it's <= 100 chars and doesn't end with hyphen
			if tt.name == "long title gets truncated" {
				if len(result) > 100 {
					t.Errorf("generatePRWorktreeName() result length = %d, want <= 100", len(result))
				}
				if strings.HasSuffix(result, "-") {
					t.Errorf("generatePRWorktreeName() result ends with hyphen: %q", result)
				}
			} else if result != tt.expected {
				t.Errorf("generatePRWorktreeName() = %q, want %q", result, tt.expected)
			}
			// Ensure result is max 100 chars
			if len(result) > 100 {
				t.Errorf("generatePRWorktreeName() result length = %d, want <= 100", len(result))
			}
		})
	}
}
