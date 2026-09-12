package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ledhcg/cca/internal/app"
	"github.com/ledhcg/cca/internal/i18n"
	"github.com/ledhcg/cca/internal/profileenv"
	"github.com/ledhcg/cca/internal/session"
	"github.com/ledhcg/cca/internal/ui"
)

func cmdSession(a *app.App, argv []string) error {
	if len(argv) == 0 {
		return printSessionUsage()
	}

	sub := argv[0]
	rest := argv[1:]

	switch sub {
	case "ls", "list":
		return cmdSessionLs(a, rest)
	case "cp", "copy":
		return cmdSessionTransfer(a, rest, false)
	case "mv", "move":
		return cmdSessionTransfer(a, rest, true)
	case "rm", "delete":
		return cmdSessionRm(a, rest)
	case "help", "--help", "-h":
		return printSessionUsage()
	default:
		return errors.New(i18n.T(i18n.KeyCmdSessionUnknownSubcmd, sub))
	}
}

func printSessionUsage() error {
	fmt.Println(`cca session — manage Claude Code sessions across profiles

Usage:
  cca session ls   [--from <profile>] [--all]
  cca session cp   --from <profile> --to <profile> [--id <sessionId> | --all]
  cca session mv   --from <profile> --to <profile> [--id <sessionId> | --all]
  cca session rm   [--from <profile>] (--id <sessionId> | --all) [-y|--yes]

Aliases:
  copy -> cp, move -> mv, delete -> rm`)
	return nil
}

func cmdSessionLs(a *app.App, argv []string) error {
	pa := parseArgs(argv, []string{"--all"}, []string{"--from"})

	fromProfile := pa.values["--from"]
	if fromProfile == "" {
		fromProfile = profileenv.Active(a)
		if fromProfile == "" {
			fromProfile = profileenv.DefaultName
		}
	}

	fromDir, err := profileenv.Require(a, fromProfile)
	if err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	allProjects := pa.bools["--all"]

	sessions, err := session.ListSessions(fromDir, cwd, allProjects)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		msg := i18n.T(i18n.KeyCmdSessionNoSessionsFound, fromProfile)
		if allProjects {
			fmt.Println(ui.C.Dim + msg + ui.C.Off)
		} else {
			fmt.Printf("%s%s %s%s\n",
				ui.C.Dim,
				msg,
				i18n.T(i18n.KeyCmdSessionNoSessionsAllHint),
				ui.C.Off,
			)
		}
		return nil
	}

	var rows [][]string
	for _, s := range sessions {
		shortID := s.ID
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		msgs := strconv.Itoa(s.MessageCount)
		size := formatBytes(s.Size)
		mod := relativeTime(s.ModTime)

		if allProjects {
			rows = append(rows, []string{
				shortID,
				cleanOneLine(s.ProjectSlug, 28),
				s.Title,
				msgs,
				size,
				mod,
			})
		} else {
			rows = append(rows, []string{
				shortID,
				s.Title,
				msgs,
				size,
				mod,
			})
		}
	}

	if allProjects {
		printTable(rows, []string{
			i18n.T(i18n.KeyCmdSessionHeaderID),
			i18n.T(i18n.KeyCmdSessionHeaderProject),
			i18n.T(i18n.KeyCmdSessionHeaderTitle),
			i18n.T(i18n.KeyCmdSessionHeaderMessages),
			i18n.T(i18n.KeyCmdSessionHeaderSize),
			i18n.T(i18n.KeyCmdSessionHeaderModified),
		})
	} else {
		printTable(rows, []string{
			i18n.T(i18n.KeyCmdSessionHeaderID),
			i18n.T(i18n.KeyCmdSessionHeaderTitle),
			i18n.T(i18n.KeyCmdSessionHeaderMessages),
			i18n.T(i18n.KeyCmdSessionHeaderSize),
			i18n.T(i18n.KeyCmdSessionHeaderModified),
		})
	}

	return nil
}

func cmdSessionTransfer(a *app.App, argv []string, move bool) error {
	pa := parseArgs(argv, []string{"--all"}, []string{"--from", "--to", "--id"})

	fromProfile := pa.values["--from"]
	toProfile := pa.values["--to"]

	if fromProfile == "" || toProfile == "" {
		return errors.New(i18n.T(i18n.KeyCmdSessionMissingFromTo))
	}
	if fromProfile == toProfile {
		return errors.New(i18n.T(i18n.KeyCmdSessionSameProfile))
	}

	fromDir, err := profileenv.Require(a, fromProfile)
	if err != nil {
		return err
	}
	toDir, err := profileenv.Require(a, toProfile)
	if err != nil {
		return err
	}

	hasAll := pa.bools["--all"]
	hasID := pa.values["--id"] != ""

	if (!hasAll && !hasID) || (hasAll && hasID) {
		return errors.New(i18n.T(i18n.KeyCmdSessionRequireIdOrAll))
	}

	cwd, _ := os.Getwd()

	if hasID {
		sid, err := session.ResolveSessionID(fromDir, cwd, pa.values["--id"])
		if err != nil {
			return err
		}
		if err := session.TransferSession(fromDir, toDir, cwd, sid, move); err != nil {
			return err
		}
		shortID := sid
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		if move {
			ui.Ok(i18n.T(i18n.KeyCmdSessionMoveSuccess, shortID, fromProfile, toProfile))
		} else {
			ui.Ok(i18n.T(i18n.KeyCmdSessionCopySuccess, shortID, fromProfile, toProfile))
		}
		return nil
	}

	// Bulk transfer (--all)
	count, err := session.TransferAllSessions(fromDir, toDir, cwd, move)
	if err != nil {
		return err
	}
	if count == 0 {
		ui.Warn(i18n.T(i18n.KeyCmdSessionNoSessionsFound, fromProfile))
		return nil
	}

	if move {
		ui.Ok(i18n.T(i18n.KeyCmdSessionMoveAllSuccess, count, fromProfile, toProfile))
	} else {
		ui.Ok(i18n.T(i18n.KeyCmdSessionCopyAllSuccess, count, fromProfile, toProfile))
	}
	return nil
}

func cmdSessionRm(a *app.App, argv []string) error {
	pa := parseArgs(argv, []string{"--all", "-y", "--yes"}, []string{"--from", "--id"})

	fromProfile := pa.values["--from"]
	if fromProfile == "" {
		fromProfile = profileenv.Active(a)
		if fromProfile == "" {
			fromProfile = profileenv.DefaultName
		}
	}

	fromDir, err := profileenv.Require(a, fromProfile)
	if err != nil {
		return err
	}

	hasAll := pa.bools["--all"]
	hasID := pa.values["--id"] != ""

	if (!hasAll && !hasID) || (hasAll && hasID) {
		return errors.New(i18n.T(i18n.KeyCmdSessionRequireIdOrAll))
	}

	cwd, _ := os.Getwd()

	if hasID {
		sid, err := session.ResolveSessionID(fromDir, cwd, pa.values["--id"])
		if err != nil {
			return err
		}

		shortID := sid
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}

		if !pa.bools["-y"] && !pa.bools["--yes"] {
			fmt.Print(i18n.T(i18n.KeyCmdSessionRmConfirmSingle, shortID, fromProfile))
			reader := bufio.NewReader(os.Stdin)
			ans, _ := reader.ReadString('\n')
			ans = strings.ToLower(strings.TrimSpace(ans))
			if ans != "y" && ans != "yes" && ans != "c" && ans != "co" && ans != "có" {
				fmt.Println(i18n.T(i18n.KeyCommonCancelled))
				return nil
			}
		}

		if err := session.RemoveSession(fromDir, cwd, sid); err != nil {
			return err
		}
		ui.Ok(i18n.T(i18n.KeyCmdSessionRmSuccess, shortID, fromProfile))
		return nil
	}

	// Bulk removal (--all)
	sessions, err := session.ListSessions(fromDir, cwd, false)
	if err != nil {
		return err
	}
	if len(sessions) == 0 {
		ui.Warn(i18n.T(i18n.KeyCmdSessionNoSessionsFound, fromProfile))
		return nil
	}

	if !pa.bools["-y"] && !pa.bools["--yes"] {
		fmt.Print(i18n.T(i18n.KeyCmdSessionRmConfirmAll, len(sessions), fromProfile))
		reader := bufio.NewReader(os.Stdin)
		ans, _ := reader.ReadString('\n')
		ans = strings.ToLower(strings.TrimSpace(ans))
		if ans != "y" && ans != "yes" && ans != "c" && ans != "co" && ans != "có" {
			fmt.Println(i18n.T(i18n.KeyCommonCancelled))
			return nil
		}
	}

	count := 0
	for _, s := range sessions {
		if err := session.RemoveSession(fromDir, cwd, s.ID); err == nil {
			count++
		}
	}

	ui.Ok(i18n.T(i18n.KeyCmdSessionRmAllSuccess, count, fromProfile))
	return nil
}
