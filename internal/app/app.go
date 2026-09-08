// Package app holds the runtime environment cca operates in: where things
// live on disk. It replaces the package-level globals the pre-refactor code
// used (Home/Main/Root/ConfigPath), so every function that needs a path takes
// an explicit *App instead of reading a shared variable.
package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// App is the resolved set of directories cca operates on for the current
// process. Construct with New; it never changes after that.
type App struct {
	Home       string // the user's home directory
	Main       string // ~/.claude — the root/default Claude Code profile
	Root       string // ~/.claude-accounts (or $CCA_ROOT) — where extra profiles live
	ConfigPath string // Root/config.json
}

// New resolves App from the environment. CCA_ROOT overrides the default
// Root location — used by tests and by anyone who wants profiles to live
// somewhere other than ~/.claude-accounts.
func New() (*App, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}
	root := os.Getenv("CCA_ROOT")
	if root == "" {
		root = filepath.Join(home, ".claude-accounts")
	}
	return &App{
		Home:       home,
		Main:       filepath.Join(home, ".claude"),
		Root:       root,
		ConfigPath: filepath.Join(root, "config.json"),
	}, nil
}

// IsWindows reports whether cca is running on Windows.
func IsWindows() bool { return runtime.GOOS == "windows" }
