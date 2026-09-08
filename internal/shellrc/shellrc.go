// Package shellrc detects the user's shell and generates the PATH +
// completion snippet `cca install` appends to its rc file, plus the
// Windows-specific user-PATH registry update.
package shellrc

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ledhcg/cca/internal/app"
)

const (
	MarkStart = "# >>> cca >>>"
	MarkEnd   = "# <<< cca <<<"
	ccaCmds   = "ls new use sh exec login logout rm sync info doctor config guide help install version"
)

// Snippet returns the PATH + completion block to insert into an rc file for
// the given shell kind ("bash", "zsh", "fish", "powershell").
func Snippet(kind string) string {
	switch kind {
	case "fish":
		return MarkStart + "\nset -gx PATH $HOME/.claude/bin $PATH\n" + MarkEnd
	case "powershell":
		return MarkStart + "\n" +
			`if ($env:Path -notlike "*.claude\bin*") {` + "\n" +
			`  $env:Path = "$HOME\.claude\bin;" + $env:Path` + "\n" +
			"}\n" + MarkEnd
	case "zsh":
		body := "_cca() {\n" +
			"  local -a cmds\n" +
			"  cmds=(" + ccaCmds + ")\n" +
			"  if (( CURRENT == 2 )); then compadd -- $cmds\n" +
			`  else compadd -- ${(f)"$(cca ls --names 2>/dev/null)"}; fi` + "\n" +
			"}\n" +
			"compdef _cca cca 2>/dev/null"
		return MarkStart + "\n" + `export PATH="$HOME/.claude/bin:$PATH"` + "\n" + body + "\n" + MarkEnd
	default: // bash (default POSIX shell, also covers Git Bash on Windows)
		body := "_cca() {\n" +
			"  local cur=${COMP_WORDS[COMP_CWORD]}\n" +
			`  if [ "$COMP_CWORD" -eq 1 ]; then` + "\n" +
			`    COMPREPLY=($(compgen -W "` + ccaCmds + `" -- "$cur"))` + "\n" +
			"  else\n" +
			`    COMPREPLY=($(compgen -W "$(cca ls --names 2>/dev/null)" -- "$cur"))` + "\n" +
			"  fi\n" +
			"}\n" +
			"complete -F _cca cca 2>/dev/null"
		return MarkStart + "\n" + `export PATH="$HOME/.claude/bin:$PATH"` + "\n" + body + "\n" + MarkEnd
	}
}

// DetectRC guesses the right rc file to append PATH + completion to.
func DetectRC(a *app.App) (path, kind string) {
	if app.IsWindows() {
		shellEnv := os.Getenv("SHELL")
		if strings.Contains(shellEnv, "bash") || os.Getenv("MSYSTEM") != "" {
			return filepath.Join(a.Home, ".bashrc"), "bash"
		}
		pwsh, err := exec.LookPath("pwsh")
		if err != nil {
			pwsh, err = exec.LookPath("powershell")
		}
		if err == nil {
			out, runErr := exec.Command(pwsh, "-NoProfile", "-Command", "$PROFILE").Output()
			profile := strings.TrimSpace(string(out))
			if runErr == nil && profile != "" {
				return profile, "powershell"
			}
		}
		return a.Home, "unknown" // plain cmd.exe — no rc file concept
	}
	shell := filepath.Base(os.Getenv("SHELL"))
	switch shell {
	case "zsh":
		return filepath.Join(a.Home, ".zshrc"), "zsh"
	case "fish":
		return filepath.Join(a.Home, ".config", "fish", "config.fish"), "fish"
	case "bash":
		return filepath.Join(a.Home, ".bashrc"), "bash"
	default:
		return filepath.Join(a.Home, ".bashrc"), "bash"
	}
}

// AddToWindowsUserPath persists dir into the current user's PATH environment
// variable via the registry (HKCU\Environment). Unlike a shell rc file, this
// covers cmd.exe — which has no rc-file concept — and any shell opened after
// this call, without waiting for the user to source anything. It complements
// (doesn't replace) the rc snippet, which still owns completion.
//
// Shells out to powershell for [Environment]::SetEnvironmentVariable rather
// than the "setx" builtin, which silently truncates values over 1024 chars —
// a real risk given how long PATH gets in practice.
func AddToWindowsUserPath(dir string) (added bool, err error) {
	ps, err := exec.LookPath("powershell")
	if err != nil {
		ps, err = exec.LookPath("pwsh")
		if err != nil {
			return false, err
		}
	}
	const script = `
$dir = [Environment]::GetEnvironmentVariable('Path','User')
$target = $env:CCA_ADD_DIR
$parts = @()
if ($dir) { $parts = $dir -split ';' | Where-Object { $_ -ne '' } }
if ($parts -contains $target) {
  Write-Output 'unchanged'
} else {
  $new = if ($dir) { $dir.TrimEnd(';') + ';' + $target } else { $target }
  [Environment]::SetEnvironmentVariable('Path', $new, 'User')
  Write-Output 'added'
}
`
	cmd := exec.Command(ps, "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Env = append(os.Environ(), "CCA_ADD_DIR="+dir)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "added", nil
}

// DefaultInteractiveShell returns the command to launch a subshell for `cca sh`.
func DefaultInteractiveShell() []string {
	if app.IsWindows() {
		if shellEnv := os.Getenv("SHELL"); shellEnv != "" {
			if _, err := os.Stat(shellEnv); err == nil {
				return []string{shellEnv}
			}
		}
		if bash, err := exec.LookPath("bash"); err == nil {
			return []string{bash}
		}
		if pwsh, err := exec.LookPath("pwsh"); err == nil {
			return []string{pwsh}
		}
		if pwsh, err := exec.LookPath("powershell"); err == nil {
			return []string{pwsh}
		}
		return []string{"powershell"}
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	return []string{shell}
}
