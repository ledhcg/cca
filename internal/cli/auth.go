package cli

import (
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ledhcg/cca/internal/app"
	"github.com/ledhcg/cca/internal/profileenv"
)

func claudeBin(a *app.App) string {
	if p, err := exec.LookPath("claude"); err == nil {
		return p
	}
	name := "claude"
	if app.IsWindows() {
		name = "claude.exe"
	}
	return filepath.Join(a.Home, ".local", "bin", name)
}

// AuthStatus is the parsed output of `claude auth status`.
type AuthStatus struct {
	LoggedIn         bool   `json:"loggedIn"`
	Email            string `json:"email"`
	SubscriptionType string `json:"subscriptionType"`
	AuthMethod       string `json:"authMethod"`
}

func authStatus(a *app.App, name string) AuthStatus {
	var st AuthStatus
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, claudeBin(a), "auth", "status")
	cmd.Env = profileenv.EnvFor(a, name)
	out, err := cmd.Output()
	if err != nil {
		return st // claude missing / errored / timed out ⇒ treat as logged out
	}
	_ = json.Unmarshal(out, &st) // malformed JSON ⇒ keep zero-value (logged out)
	return st
}
