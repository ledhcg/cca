// Package profileenv computes per-profile paths and the environment a child
// `claude` process runs under, given the app's resolved directories.
package profileenv

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ledhcg/cca/internal/app"
)

// DefaultName is the name used for the root ~/.claude profile.
const DefaultName = "default"

var nameRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,31}$`)

var reserved = map[string]bool{DefaultName: true, "config.json": true, "config": true, "settings": true}

// IsDefault reports whether name refers to the root ~/.claude profile.
func IsDefault(name string) bool {
	return name == "" || name == DefaultName
}

// Dir returns the directory for profile name.
func Dir(a *app.App, name string) string {
	if IsDefault(name) {
		return a.Main
	}
	return filepath.Join(a.Root, name)
}

// ConfigDirStr returns the value CLAUDE_CONFIG_DIR is set to for name — the
// same as Dir, kept as a separate name so call sites can say which meaning
// they intend.
func ConfigDirStr(a *app.App, name string) string {
	return Dir(a, name)
}

// List returns the names of every extra profile under a.Root, sorted.
func List(a *app.App) []string {
	entries, err := os.ReadDir(a.Root)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names
}

// Require returns name's directory, or an error if the profile doesn't exist.
func Require(a *app.App, name string) (string, error) {
	if IsDefault(name) {
		return a.Main, nil
	}
	d := Dir(a, name)
	if fi, err := os.Stat(d); err != nil || !fi.IsDir() {
		known := strings.Join(List(a), ", ")
		if known == "" {
			known = "(no profiles yet)"
		}
		return "", fmt.Errorf("no profile '%s'. Available: %s\n  Create it with: cca new %s", name, known, name)
	}
	return d, nil
}

// Validate checks that name is a legal, non-reserved profile name.
func Validate(name string) error {
	if reserved[name] {
		return fmt.Errorf("'%s' is a reserved name — the default profile is already ~/.claude", name)
	}
	if !nameRe.MatchString(name) {
		return fmt.Errorf("profile names may only contain letters, digits, and . _ - and must start with a letter or digit")
	}
	return nil
}

// EnvFor returns the "KEY=VALUE" list for a child process running under profile name.
func EnvFor(a *app.App, name string) []string {
	env := os.Environ()
	if IsDefault(name) {
		env = EnvUnset(env, "CLAUDE_CONFIG_DIR") // default account: let Claude use ~/.claude on its own
		env = EnvSet(env, "CCA_PROFILE", DefaultName)
	} else {
		env = EnvSet(env, "CLAUDE_CONFIG_DIR", ConfigDirStr(a, name))
		env = EnvSet(env, "CCA_PROFILE", name)
	}
	return env
}

// EnvSet returns env with key set to val, overwriting any existing entry.
func EnvSet(env []string, key, val string) []string {
	prefix := key + "="
	for i, e := range env {
		if strings.HasPrefix(e, prefix) {
			env[i] = prefix + val
			return env
		}
	}
	return append(env, prefix+val)
}

// EnvUnset returns env with key removed, if present.
func EnvUnset(env []string, key string) []string {
	prefix := key + "="
	out := env[:0]
	for _, e := range env {
		if !strings.HasPrefix(e, prefix) {
			out = append(out, e)
		}
	}
	return out
}

// EnvGet returns key's value in env, or "" if unset.
func EnvGet(env []string, key string) string {
	prefix := key + "="
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			return e[len(prefix):]
		}
	}
	return ""
}

// Active guesses which profile is active based on the current process's
// CLAUDE_CONFIG_DIR (e.g. when already inside a `cca sh <name>` subshell).
func Active(a *app.App) string {
	cur := os.Getenv("CLAUDE_CONFIG_DIR")
	if cur == "" {
		return "" // unset ⇒ default (correct: running `claude` bare right now hits default)
	}
	curAbs, err := filepath.Abs(cur)
	if err != nil {
		return ""
	}
	for _, n := range List(a) {
		if pAbs, err := filepath.Abs(Dir(a, n)); err == nil && pAbs == curAbs {
			return n
		}
	}
	return ""
}
