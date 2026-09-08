// Package ui is cca's presentation layer: colored terminal output. It never
// calls os.Exit — that decision belongs solely to cmd/cca/main.go, driven by
// what internal/cli.Run returns. Business-logic packages return errors;
// only Ok/Warn (non-fatal, informational) print directly.
package ui

import (
	"fmt"
	"os"
)

// Colors holds ANSI escape codes, empty when stdout isn't a TTY or NO_COLOR is set.
type Colors struct {
	Dim, Bold, Red, Green, Yellow, Cyan, Off string
}

// C holds the active color set for this process, set once by Init.
var C Colors

// Init detects whether stdout is a terminal and NO_COLOR isn't set, and
// populates C accordingly. Call once at startup.
func Init() {
	if IsTerminal(os.Stdout) && os.Getenv("NO_COLOR") == "" {
		C = Colors{
			Dim: "\033[2m", Bold: "\033[1m", Red: "\033[31m",
			Green: "\033[32m", Yellow: "\033[33m", Cyan: "\033[36m", Off: "\033[0m",
		}
	}
}

// IsTerminal reports whether f is attached to a character device (a TTY).
func IsTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// PrintError prints err to stderr in the same "✗ msg" style cca has always
// used for fatal errors. Called exactly once, from cmd/cca/main.go.
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "%s✗%s %s\n", C.Red, C.Off, err)
}

// Ok prints a non-fatal success message ("✓ ...").
func Ok(format string, a ...any) {
	fmt.Printf("%s✓%s %s\n", C.Green, C.Off, fmt.Sprintf(format, a...))
}

// Warn prints a non-fatal warning message ("! ...").
func Warn(format string, a ...any) {
	fmt.Printf("%s!%s %s\n", C.Yellow, C.Off, fmt.Sprintf(format, a...))
}
