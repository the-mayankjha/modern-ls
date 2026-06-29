// Package terminal provides terminal dimension detection.
package terminal

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const defaultWidth = 80

// Width returns the current terminal width in columns.
// It queries the OS via TIOCGWINSZ first, then falls back to the COLUMNS
// environment variable, and finally to the hard-coded default of 80.
func Width() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil && w > 0 {
		return w
	}
	// Try COLUMNS env var for environments that don't support TIOCGWINSZ
	// (e.g., CI systems, terminal emulators with limited ioctl support).
	if col := os.Getenv("COLUMNS"); col != "" {
		var n int
		if _, err := fmt.Sscanf(col, "%d", &n); err == nil && n > 0 {
			return n
		}
	}
	return defaultWidth
}

// IsTerminal reports whether stdout is connected to a terminal.
// When piped or redirected this returns false, allowing callers to
// strip ANSI codes and icon output automatically.
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
