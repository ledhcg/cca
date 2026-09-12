// Package cli is cca's command layer: argument parsing, dispatch, and every
// cmdX command handler. Handlers return an error instead of calling
// os.Exit directly — Run is the single place that turns an error (or an
// exitCodeErr, for commands that pass through a child process's exit code)
// into a process exit code.
package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledhcg/cca/internal/app"
	"github.com/ledhcg/cca/internal/config"
	"github.com/ledhcg/cca/internal/i18n"
	"github.com/ledhcg/cca/internal/profileenv"
	"github.com/ledhcg/cca/internal/ui"
	"github.com/ledhcg/cca/internal/update"
)

// Version is overridden at build time via
// -ldflags "-X github.com/ledhcg/cca/internal/cli.Version=vX.Y.Z"
// (see .github/workflows/release.yml). Local builds report "dev".
var Version = "dev"

// exitCodeErr lets a command hand Run an exact exit code (e.g. passed
// through from a child `claude` process) without Run printing anything —
// the command has already said everything worth saying.
type exitCodeErr struct{ code int }

func (e exitCodeErr) Error() string { return "" }

// ── command-line arguments ────────────────────────────────────────────────────
// Deliberately not using the "flag" package: it stops parsing at the first
// non-flag token (e.g. "work" in `cca new work --login`), which breaks the
// flags-after-positional syntax every cca command relies on. parseArgs below
// scans the whole argv regardless of order.
type parsedArgs struct {
	positional []string
	bools      map[string]bool
	values     map[string]string
}

func parseArgs(args []string, boolNames, valueNames []string) parsedArgs {
	boolSet := map[string]bool{}
	for _, n := range boolNames {
		boolSet[n] = true
	}
	valueSet := map[string]bool{}
	for _, n := range valueNames {
		valueSet[n] = true
	}
	pa := parsedArgs{bools: map[string]bool{}, values: map[string]string{}}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case boolSet[a]:
			pa.bools[a] = true
		case valueSet[a]:
			if i+1 < len(args) {
				pa.values[a] = args[i+1]
				i++
			}
		default:
			pa.positional = append(pa.positional, a)
		}
	}
	return pa
}

func requirePositional(args parsedArgs) (string, error) {
	if len(args.positional) == 0 {
		return "", fmt.Errorf("missing profile name")
	}
	return args.positional[0], nil
}

func extractLang(argv []string) ([]string, string) {
	var filtered []string
	var lang string
	isPassthrough := false
	if len(argv) > 0 && (argv[0] == "use" || argv[0] == "exec") {
		isPassthrough = true
	}

	for i := 0; i < len(argv); i++ {
		arg := argv[i]
		if strings.HasPrefix(arg, "--lang=") {
			if !isPassthrough || i == 0 {
				lang = strings.TrimPrefix(arg, "--lang=")
				continue
			}
		} else if arg == "--lang" && i+1 < len(argv) {
			if !isPassthrough || i == 0 {
				lang = argv[i+1]
				i++
				continue
			}
		}
		filtered = append(filtered, arg)
	}
	return filtered, lang
}

var subcommands = map[string]bool{
	"ls": true, "new": true, "use": true, "login": true, "logout": true,
	"info": true, "sh": true, "exec": true, "rm": true, "sync": true,
	"doctor": true, "settings": true, "update": true, "config": true, "guide": true, "help": true, "install": true,
	"version": true,
}

// Run parses argv (os.Args[1:]) and dispatches to the matching command. It
// returns the process exit code — the only caller, cmd/cca/main.go, does
// nothing but pass it to os.Exit.
func Run(a *app.App, argv []string) int {
	update.TriggerBackgroundCheck(a, Version)

	code := dispatch(a, argv)
	cmd := "ls"
	if len(argv) > 0 {
		cmd = argv[0]
	}
	if code == 0 {
		notifyUpdate(a, cmd)
	}
	return code
}

func notifyUpdate(a *app.App, cmd string) {
	if cmd == "update" || cmd == "version" || Version == "dev" {
		return
	}
	if os.Getenv("CCA_NO_UPDATE_CHECK") != "" || os.Getenv("CI") != "" {
		return
	}
	if !ui.IsTerminal(os.Stderr) {
		return
	}
	cfg, err := config.Load(a)
	if err != nil || cfg.LatestVersion == "" {
		return
	}
	if update.CompareVersions(Version, cfg.LatestVersion) < 0 {
		fmt.Fprintf(os.Stderr, "\n%s! %s%s\n  %s\n\n",
			ui.C.Yellow,
			i18n.T(i18n.KeyCmdUpdateAvailableNotice, Version, cfg.LatestVersion),
			ui.C.Off,
			i18n.T(i18n.KeyCmdUpdateRunHint),
		)
	}
}

func dispatch(a *app.App, argv []string) int {
	ui.Init()
	ui.SetupConsole()

	var cliLang string
	argv, cliLang = extractLang(argv)

	cfg, _ := config.Load(a)
	i18n.Init(i18n.Resolve(cliLang, cfg.Lang))

	if len(argv) == 0 {
		return runCmd(func() error { return cmdLs(a, parsedArgs{bools: map[string]bool{}}) })
	}

	// `cca work [args…]` = `cca use work [args…]`
	if !subcommands[argv[0]] && !strings.HasPrefix(argv[0], "-") {
		name := argv[0]
		if profileenv.IsDefault(name) || isDir(profileenv.Dir(a, name)) {
			argv = append([]string{"use"}, argv...)
		} else {
			known := profileenv.List(a)
			msg := i18n.T(i18n.KeyCliErrNotCommandOrProf, name)
			if len(known) > 0 {
				msg += "\n" + i18n.T(i18n.KeyCliErrExistingProfiles, strings.Join(known, ", "))
			} else {
				msg += "\n" + i18n.T(i18n.KeyCliErrNoExtraProfiles)
			}
			if profileenv.Validate(name) == nil {
				msg += "\n" + i18n.T(i18n.KeyCliErrCreateThisProfile, name)
			}
			msg += "\n" + i18n.T(i18n.KeyCliErrSeeCommandsGuide)
			return runCmd(func() error { return errors.New(msg) })
		}
	}

	cmd, rest := argv[0], argv[1:]
	switch cmd {
	case "ls":
		return runCmd(func() error { return cmdLs(a, parseArgs(rest, []string{"--names"}, nil)) })
	case "new":
		return runCmd(func() error { return cmdNew(a, parseArgs(rest, []string{"--login", "--yolo"}, nil)) })
	case "use":
		if len(rest) == 0 {
			return runCmd(func() error { return errors.New(i18n.T(i18n.KeyCliErrMissingProfileName)) })
		}
		return runCmd(func() error { return cmdUse(a, rest[0], rest[1:]) })
	case "login":
		return runCmd(func() error { return cmdLogin(a, parseArgs(rest, nil, nil)) })
	case "logout":
		return runCmd(func() error { return cmdLogout(a, parseArgs(rest, nil, nil)) })
	case "info":
		return runCmd(func() error { return cmdInfo(a, parseArgs(rest, nil, nil)) })
	case "sh":
		return runCmd(func() error { return cmdSh(a, parseArgs(rest, nil, nil)) })
	case "exec":
		if len(rest) == 0 {
			return runCmd(func() error {
				return errors.New("missing profile name — example: cca exec work -- claude auth status")
			})
		}
		name, cmdArgs := rest[0], rest[1:]
		if len(cmdArgs) > 0 && cmdArgs[0] == "--" {
			cmdArgs = cmdArgs[1:]
		}
		return runCmd(func() error { return cmdExec(a, name, cmdArgs) })
	case "rm":
		return runCmd(func() error { return cmdRm(a, parseArgs(rest, []string{"-y", "--yes", "--keep-keychain"}, nil)) })
	case "sync":
		return runCmd(func() error { return cmdSync(a, parseArgs(rest, []string{"--all"}, []string{"--strategy"})) })
	case "doctor":
		return runCmd(func() error { return cmdDoctor(a, parseArgs(rest, nil, nil)) })
	case "settings":
		return runCmd(func() error { return cmdSettings(a, parseArgs(rest, nil, nil)) })
	case "update":
		return runCmd(func() error { return cmdUpdate(a, parseArgs(rest, []string{"--check", "-c"}, nil)) })
	case "config":
		return runCmd(func() error { return cmdConfig(a, parseArgs(rest, []string{"--edit", "--print"}, nil)) })
	case "guide":
		printGuide()
		return 0
	case "help", "--help", "-h":
		printUsage()
		return 0
	case "version", "--version", "-v":
		return runCmd(func() error { return cmdVersion(a) })
	case "install":
		return runCmd(func() error { return cmdInstall(a, parseArgs(rest, nil, []string{"--rc", "--rc-kind"})) })
	default:
		printUsage()
		return 0
	}
}

// runCmd executes a command handler and turns its returned error into a
// process exit code, printing it first unless it's an exitCodeErr (which
// means the command already said everything worth saying).
func runCmd(fn func() error) int {
	err := fn()
	if err == nil {
		return 0
	}
	var ec exitCodeErr
	if errors.As(err, &ec) {
		return ec.code
	}
	ui.PrintError(err)
	return 1
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// installSelf copies the currently-running binary into ~/.claude/bin/cca(.exe).
// No symlink/wrapper trick needed like the old Python version: Go produces a
// single standalone binary, so "installing" just means placing that file
// somewhere PATH can find it.
func installSelf(a *app.App) {
	self, err := os.Executable()
	if err != nil {
		ui.Warn("could not determine the current binary's path: %v", err)
		return
	}
	if resolved, err := filepath.EvalSymlinks(self); err == nil {
		self = resolved
	}
	bindir := filepath.Join(a.Main, "bin")
	_ = os.MkdirAll(bindir, 0755)
	name := "cca"
	if app.IsWindows() {
		name = "cca.exe"
	}
	dst := filepath.Join(bindir, name)
	if err := copyExecutable(self, dst); err != nil {
		ui.Warn("could not copy the binary to %s: %v", dst, err)
		return
	}
	fmt.Printf("  %s%s%s\n", ui.C.Dim, dst, ui.C.Off)
}

// copyExecutable writes src to dst, including the case where dst is the binary
// currently executing (self-update). Verified experimentally: Windows locks a
// running process's own image file — renaming ANOTHER file ON TOP of it fails
// with "Access is denied", even though renaming that running file to a
// DIFFERENT name is allowed (renaming doesn't touch its contents). So the old
// file must be moved out of the way first, then the new one written in place.
func copyExecutable(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if _, err := os.Stat(dst); err == nil {
		old := dst + ".old"
		_ = os.Remove(old) // clean up leftovers from a previous update, if any
		if err := os.Rename(dst, old); err != nil {
			return err
		}
	}
	return os.WriteFile(dst, data, 0755)
}
