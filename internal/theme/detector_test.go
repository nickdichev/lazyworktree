package theme

import (
	"os"
	"testing"
	"time"

	"golang.org/x/term"
)

func TestParseColorComponent(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "4-bit low", input: "0", want: 0x0000, wantErr: false},
		{name: "4-bit high", input: "f", want: 0xffff, wantErr: false},
		{name: "8-bit low", input: "00", want: 0x0000, wantErr: false},
		{name: "8-bit mid", input: "80", want: 0x8080, wantErr: false},
		{name: "8-bit high", input: "ff", want: 0xffff, wantErr: false},
		{name: "12-bit low", input: "000", want: 0x0000, wantErr: false},
		{name: "12-bit mid", input: "800", want: 0x8008, wantErr: false}, // 0x800 * 16 + 0x800 / 256
		{name: "12-bit high", input: "fff", want: 0xffff, wantErr: false},
		{name: "16-bit low", input: "0000", want: 0x0000, wantErr: false},
		{name: "16-bit high", input: "ffff", want: 0xffff, wantErr: false},
		{name: "invalid hex", input: "xyz", want: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseColorComponent(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseColorComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseColorComponent() = %x, want %x", got, tt.want)
			}
		})
	}
}

// DetectBackground requires a terminal and user interaction, so we skip it in unit tests.
// However, we can mock the behavior if we refactor DetectBackground to take an interface or a reader/writer.
// For now, we just ensure it returns an error or default when not in a terminal (which is the case in CI/tests).
func TestDetectBackgroundNonTerminal(t *testing.T) {
	// Skip if we are running in a real terminal to avoid messing with raw mode
	// during interactive tests (like go test ./...)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		t.Skip("skipping terminal detection test in interactive mode")
	}

	// This test assumes it's running in a non-interactive environment (go test usually is).
	// If it somehow thinks it's a terminal, it might hang waiting for input if we aren't careful.
	// But DetectBackground checks term.IsTerminal(stdin).

	// In a test environment, stdin is typically not a terminal.
	// If it is (e.g. running go test manually), this test might fail or hang.
	// We'll set a very short timeout.

	// We can't easily force IsTerminal to be false without mocking syscalls or 'term' package.
	// But commonly `go test` redirects stdin from /dev/null or a pipe.

	theme, err := DetectBackground(10 * time.Millisecond)

	// We expect either "dracula" (default fallback) and an error, OR a valid theme if it miraculously works.
	// But most likely it returns error "stdin is not a terminal".

	if theme == "" {
		t.Error("DetectBackground returned empty theme")
	}

	// Just ensure it doesn't panic.
	if err == nil {
		// If no error, it means it somehow detected something or timed out gracefully.
		// That's acceptable.
	} else if theme != DefaultDark() {
		// Verify it returns the default dark theme on error
		t.Errorf("DetectBackground returned %q on error, expected default dark %q", theme, DefaultDark())
	}
}
