package theme

import "testing"

func TestGetTheme(t *testing.T) {
	for _, name := range AvailableThemes() {
		got := GetTheme(name)
		if got == nil {
			t.Fatalf("expected theme for %q", name)
		}
	}

	fallback := GetTheme("unknown")
	if fallback.Background != Dracula().Background {
		t.Fatalf("expected Dracula fallback, got %q", fallback.Background)
	}
}

func TestAvailableThemesIncludesDefaults(t *testing.T) {
	themes := AvailableThemes()
	required := map[string]bool{
		"dracula":          false,
		"narna":            false,
		"clean-light":      false,
		"solarized-dark":   false,
		"solarized-light":  false,
		"gruvbox-dark":     false,
		"gruvbox-light":    false,
		"nord":             false,
		"monokai":          false,
		"catppuccin-mocha": false,
	}

	for _, name := range themes {
		if _, ok := required[name]; ok {
			required[name] = true
		}
	}

	for name, seen := range required {
		if !seen {
			t.Fatalf("expected theme %q to be available", name)
		}
	}
}
