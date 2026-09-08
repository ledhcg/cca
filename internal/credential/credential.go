// Package credential answers "does this profile have a saved login, and
// where" across the two storage mechanisms Claude Code uses.
//
// Verified empirically (reading the claude.exe bundle + code.claude.com/docs/en/authentication):
//
//	macOS    → system Keychain, entry named "Claude Code-credentials-<sha256(config_dir)[:8]>".
//	Linux    → <config_dir>/.credentials.json (mode 0600).
//	Windows  → <config_dir>/.credentials.json.
//
// Linux and Windows share fileCredentialBackend since the credential lives
// right inside the profile directory — deleting the directory deletes it
// too, no separate cleanup needed like on macOS's Keychain.
package credential

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
)

// Backend abstracts credential storage for one profile. dir is the
// profile's directory (== its CLAUDE_CONFIG_DIR value); isDefault marks the
// root ~/.claude profile, which macOS keys differently (no path hash).
type Backend interface {
	Label() string
	DisplayID(dir string, isDefault bool) string
	Exists(dir string, isDefault bool) bool
	Forget(dir string, isDefault bool)
	// Orphans returns credential entries that don't belong to any known profile.
	// nil if the backend has no such concept.
	Orphans(knownIDs map[string]bool) []string
}

// GetBackend returns the credential backend for the current OS.
func GetBackend() Backend {
	if runtime.GOOS == "darwin" {
		return macKeychainBackend{}
	}
	return fileCredentialBackend{}
}

// ── Linux & Windows ───────────────────────────────────────────────────────────
type fileCredentialBackend struct{}

func (fileCredentialBackend) Label() string { return "Credential file" }

func (fileCredentialBackend) path(dir string) string {
	return filepath.Join(dir, ".credentials.json")
}

func (b fileCredentialBackend) DisplayID(dir string, isDefault bool) string {
	return b.path(dir)
}

func (b fileCredentialBackend) Exists(dir string, isDefault bool) bool {
	fi, err := os.Stat(b.path(dir))
	return err == nil && !fi.IsDir()
}

func (fileCredentialBackend) Forget(dir string, isDefault bool) {
	// Lives inside dir — os.RemoveAll(dir) already handles it.
}

func (fileCredentialBackend) Orphans(knownIDs map[string]bool) []string { return nil }

// ── macOS ─────────────────────────────────────────────────────────────────────
type macKeychainBackend struct{}

const keychainPrefix = "Claude Code-credentials-"

var keychainSuffixRe = regexp.MustCompile(`Claude Code-credentials-[0-9a-f]{8}`)

func (macKeychainBackend) Label() string { return "Keychain" }

func (macKeychainBackend) service(dir string, isDefault bool) string {
	if isDefault {
		return "Claude Code-credentials"
	}
	sum := sha256.Sum256([]byte(dir))
	return keychainPrefix + hex.EncodeToString(sum[:])[:8]
}

func (b macKeychainBackend) DisplayID(dir string, isDefault bool) string {
	return b.service(dir, isDefault)
}

func (b macKeychainBackend) Exists(dir string, isDefault bool) bool {
	svc := b.service(dir, isDefault)
	cmd := exec.Command("security", "find-generic-password", "-s", svc)
	return cmd.Run() == nil
}

func (b macKeychainBackend) Forget(dir string, isDefault bool) {
	svc := b.service(dir, isDefault)
	_ = exec.Command("security", "delete-generic-password", "-s", svc).Run()
}

func (macKeychainBackend) Orphans(knownIDs map[string]bool) []string {
	out, err := exec.Command("security", "dump-keychain").Output()
	if err != nil {
		return nil
	}
	found := map[string]bool{}
	for _, m := range keychainSuffixRe.FindAllString(string(out), -1) {
		found[m] = true
	}
	var orphans []string
	for s := range found {
		if !knownIDs[s] {
			orphans = append(orphans, s)
		}
	}
	return orphans
}
