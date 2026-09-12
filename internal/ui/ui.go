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

// RuneWidth returns the display column width of a single rune in a monospace terminal.
// Latin, Vietnamese accents (NFC), and standard ASCII occupy 1 column.
// CJK ideographs, Katakana, Hiragana, and Fullwidth symbols occupy 2 columns.
func RuneWidth(r rune) int {
	if r == 0 {
		return 0
	}
	// Common zero-width combining characters
	if r < 32 || (r >= 0x7F && r < 0xA0) || (r >= 0x300 && r <= 0x36F) {
		return 0
	}
	// CJK, Kana, Hangul, Fullwidth punctuation
	if (r >= 0x1100 && r <= 0x115F) ||
		(r >= 0x2E80 && r <= 0xA4CF && r != 0x303F) ||
		(r >= 0xAC00 && r <= 0xD7A3) ||
		(r >= 0xF900 && r <= 0xFAFF) ||
		(r >= 0xFE10 && r <= 0xFE19) ||
		(r >= 0xFE30 && r <= 0xFE6F) ||
		(r >= 0xFF00 && r <= 0xFF60) ||
		(r >= 0xFFE0 && r <= 0xFFE6) ||
		(r >= 0x20000 && r <= 0x3FFFF) {
		return 2
	}
	return 1
}

// StringWidth returns the visual monospace column width of string s.
func StringWidth(s string) int {
	w := 0
	inEscape := false
	for _, r := range s {
		if r == 0x1b { // ANSI escape sequence start \033
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' || r == 'K' || r == 'H' || r == 'J' {
				inEscape = false
			}
			continue
		}
		w += RuneWidth(r)
	}
	return w
}

// PadRight pads s on the right with spaces until its display width equals targetWidth.
func PadRight(s string, targetWidth int) string {
	w := StringWidth(s)
	if w >= targetWidth {
		return s
	}
	pad := targetWidth - w
	spaces := make([]byte, pad)
	for i := range spaces {
		spaces[i] = ' '
	}
	return s + string(spaces)
}
