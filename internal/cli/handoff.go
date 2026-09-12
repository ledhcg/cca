package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/ledhcg/cca/internal/app"
	"github.com/ledhcg/cca/internal/i18n"
	"github.com/ledhcg/cca/internal/profileenv"
	"github.com/ledhcg/cca/internal/session"
	"github.com/ledhcg/cca/internal/ui"
)

func cmdHandoff(a *app.App, argv []string) error {
	// Separate handoff flags and arguments before '--' from extra passthrough arguments
	var handoffArgs []string
	var extraArgs []string
	splitIdx := -1
	for i, arg := range argv {
		if arg == "--" {
			splitIdx = i
			break
		}
	}
	if splitIdx != -1 {
		handoffArgs = argv[:splitIdx]
		extraArgs = argv[splitIdx+1:]
	} else {
		handoffArgs = argv
	}

	boolFlags := []string{"--fork", "--yolo"}
	valueFlags := []string{"--from", "--id"}
	pa := parseArgs(handoffArgs, boolFlags, valueFlags)

	if len(pa.positional) == 0 {
		return errors.New(i18n.T(i18n.KeyCmdHandoffMissingTarget))
	}

	toProfile := pa.positional[0]
	// Any remaining positional args can also be passed to claude
	var passthroughPositional []string
	if len(pa.positional) > 1 {
		passthroughPositional = pa.positional[1:]
	}

	fromProfile := pa.values["--from"]
	if fromProfile == "" {
		fromProfile = profileenv.Active(a)
		if fromProfile == "" {
			fromProfile = profileenv.DefaultName
		}
	}

	if fromProfile == toProfile {
		return errors.New(i18n.T(i18n.KeyCmdHandoffSameProfile))
	}

	fromDir, err := profileenv.Require(a, fromProfile)
	if err != nil {
		return err
	}
	toDir, err := profileenv.Require(a, toProfile)
	if err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	var sid string
	var sessionTitle string

	if reqID := pa.values["--id"]; reqID != "" {
		resolved, err := session.ResolveSessionID(fromDir, cwd, reqID)
		if err != nil {
			return err
		}
		sid = resolved
		sessionTitle = sid
	} else {
		latest, err := session.LatestSession(fromDir, cwd)
		if err != nil {
			return err
		}
		if latest == nil {
			return errors.New(i18n.T(i18n.KeyCmdHandoffNoSessionFound, fromProfile))
		}
		sid = latest.ID
		sessionTitle = latest.Title
	}

	if !authStatus(a, toProfile).LoggedIn {
		ui.Warn(i18n.T(i18n.KeyCmdUseWarnNotLoggedIn, toProfile, toProfile))
	}

	shortID := sid
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	fmt.Printf("%s%s%s\n", ui.C.Dim, i18n.T(i18n.KeyCmdHandoffBanner, shortID, sessionTitle, fromProfile, toProfile), ui.C.Off)

	// Copy session from source to target
	if err := session.TransferSession(fromDir, toDir, cwd, sid, false); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyCmdHandoffCopyFailed), err)
	}

	// Prepare claude invocation
	claudeArgs := []string{"--resume", sid}
	if pa.bools["--fork"] {
		claudeArgs = append(claudeArgs, "--fork-session")
	}
	if pa.bools["--yolo"] {
		claudeArgs = append(claudeArgs, "--dangerously-skip-permissions")
		fmt.Println(i18n.T(i18n.KeyCmdUseYoloNote, ui.C.Yellow, ui.C.Off, ui.C.Dim, "--dangerously-skip-permissions", ui.C.Off))
	}

	claudeArgs = append(claudeArgs, passthroughPositional...)
	claudeArgs = append(claudeArgs, extraArgs...)

	return exitCodeErr{runInherited(claudeBin(a), claudeArgs, profileenv.EnvFor(a, toProfile))}
}
