// Package theme provides theme definitions and management for the TUI.
package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines all colors used in the application UI.
type Theme struct {
	Background lipgloss.Color
	Accent     lipgloss.Color
	AccentDim  lipgloss.Color
	Border     lipgloss.Color
	BorderDim  lipgloss.Color
	MutedFg    lipgloss.Color
	TextFg     lipgloss.Color
	SuccessFg  lipgloss.Color
	WarnFg     lipgloss.Color
	ErrorFg    lipgloss.Color
	Cyan       lipgloss.Color
	Pink       lipgloss.Color
	Yellow     lipgloss.Color
}

// Dracula returns the Dracula theme (dark background, vibrant colors).
func Dracula() *Theme {
	return &Theme{
		Background: lipgloss.Color("#282A36"), // Background
		Accent:     lipgloss.Color("#BD93F9"), // Purple (primary accent)
		AccentDim:  lipgloss.Color("#44475A"), // Current Line / Selection
		Border:     lipgloss.Color("#6272A4"), // Comment (subtle borders)
		BorderDim:  lipgloss.Color("#44475A"), // Darker borders
		MutedFg:    lipgloss.Color("#6272A4"), // Comment (muted text)
		TextFg:     lipgloss.Color("#F8F8F2"), // Foreground (primary text)
		SuccessFg:  lipgloss.Color("#50FA7B"), // Green (success)
		WarnFg:     lipgloss.Color("#FFB86C"), // Orange (warning)
		ErrorFg:    lipgloss.Color("#FF5555"), // Red (error)
		Cyan:       lipgloss.Color("#8BE9FD"), // Cyan (info/secondary)
		Pink:       lipgloss.Color("#FF79C6"), // Pink (alternative accent)
		Yellow:     lipgloss.Color("#F1FA8C"), // Yellow (alternative highlight)
	}
}

// LazyGit returns a balanced LazyGit-inspired dark theme with blue accents.
func LazyGit() *Theme {
	return &Theme{
		Background: lipgloss.Color("#0D1117"), // Charcoal background
		Accent:     lipgloss.Color("#41ADFF"), // LazyGit blue accent
		AccentDim:  lipgloss.Color("#1A2230"), // Selected rows / panels
		Border:     lipgloss.Color("#30363D"), // Subtle borders
		BorderDim:  lipgloss.Color("#20252D"), // Dim borders
		MutedFg:    lipgloss.Color("#8B949E"), // Muted text
		TextFg:     lipgloss.Color("#E6EDF3"), // Primary text
		SuccessFg:  lipgloss.Color("#3FB950"), // Success green
		WarnFg:     lipgloss.Color("#E3B341"), // Warning amber
		ErrorFg:    lipgloss.Color("#F47067"), // Soft red
		Cyan:       lipgloss.Color("#7CE0F3"), // Cyan highlights
		Pink:       lipgloss.Color("#D2A8FF"), // Accent purple/pink
		Yellow:     lipgloss.Color("#F2CC60"), // Highlight yellow
	}
}

// Light returns a theme optimized for light terminal backgrounds.
func Light() *Theme {
	return &Theme{
		Background: lipgloss.Color("#FFFFFF"), // White
		Accent:     lipgloss.Color("#0066CC"), // Dark blue
		AccentDim:  lipgloss.Color("#E6E6E6"), // Light gray for selection
		Border:     lipgloss.Color("#CCCCCC"), // Light gray
		BorderDim:  lipgloss.Color("#D9D9D9"), // Slightly darker gray
		MutedFg:    lipgloss.Color("#666666"), // Medium gray (muted text)
		TextFg:     lipgloss.Color("#000000"), // Black (primary text)
		SuccessFg:  lipgloss.Color("#009900"), // Dark green (success)
		WarnFg:     lipgloss.Color("#FF8800"), // Dark orange (warning)
		ErrorFg:    lipgloss.Color("#CC0000"), // Dark red (error)
		Cyan:       lipgloss.Color("#0099CC"), // Dark cyan (info/secondary)
		Pink:       lipgloss.Color("#CC0066"), // Dark pink
		Yellow:     lipgloss.Color("#CC9900"), // Dark yellow
	}
}

// GetTheme returns a theme by name, or Dracula if not found.
func GetTheme(name string) *Theme {
	switch name {
	case "lazygit":
		return LazyGit()
	case "light":
		return Light()
	default:
		return Dracula()
	}
}

// AvailableThemes returns a list of available theme names.
func AvailableThemes() []string {
	return []string{"dracula", "lazygit", "light"}
}
