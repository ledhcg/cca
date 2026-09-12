package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ledhcg/cca/internal/app"
	"github.com/ledhcg/cca/internal/config"
	"github.com/ledhcg/cca/internal/credential"
	"github.com/ledhcg/cca/internal/i18n"
	"github.com/ledhcg/cca/internal/link"
	"github.com/ledhcg/cca/internal/profileenv"
	"github.com/ledhcg/cca/internal/shellrc"
	"github.com/ledhcg/cca/internal/ui"
	"github.com/ledhcg/cca/internal/update"
)

// ── printing helpers ─────────────────────────────────────────────────────────
func printTable(rows [][]string, head []string) {
	widths := make([]int, len(head))
	for i, h := range head {
		widths[i] = ui.StringWidth(h)
	}
	for _, r := range rows {
		for i, c := range r {
			if l := ui.StringWidth(c); l > widths[i] {
				widths[i] = l
			}
		}
	}
	var headCells []string
	for i, h := range head {
		headCells = append(headCells, ui.PadRight(h, widths[i]))
	}
	fmt.Println(ui.C.Dim + strings.TrimRight(strings.Join(headCells, "  "), " ") + ui.C.Off)
	for _, r := range rows {
		var cells []string
		for i, c := range r {
			cells = append(cells, ui.PadRight(c, widths[i]))
		}
		fmt.Println(strings.TrimRight(strings.Join(cells, "  "), " "))
	}
}

func or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// ── commands ─────────────────────────────────────────────────────────────────
func cmdLs(a *app.App, args parsedArgs) error {
	names := profileenv.List(a)
	if args.bools["--names"] {
		fmt.Println(strings.Join(names, "\n"))
		return nil
	}
	cur := profileenv.Active(a)
	targets := append([]string{""}, names...)
	results := make([]AuthStatus, len(targets))
	var wg sync.WaitGroup
	for i, t := range targets {
		wg.Add(1)
		go func(i int, t string) {
			defer wg.Done()
			results[i] = authStatus(a, t)
		}(i, t)
	}
	wg.Wait()

	var rows [][]string
	for i, name := range targets {
		label := i18n.T(i18n.KeyCommonDefault)
		if name != "" {
			label = name
		}
		if name == cur {
			label += " ←"
		}
		st := results[i]
		status := i18n.T(i18n.KeyCmdLsStatusOut)
		if st.LoggedIn {
			status = i18n.T(i18n.KeyCmdLsStatusIn)
		}
		rows = append(rows, []string{
			label, status, or(st.Email, "—"), or(or(st.SubscriptionType, st.AuthMethod), i18n.T(i18n.KeyCommonNone)),
		})
	}
	printTable(rows, []string{
		i18n.T(i18n.KeyCmdLsHeaderProfile),
		i18n.T(i18n.KeyCmdLsHeaderStatus),
		i18n.T(i18n.KeyCmdLsHeaderAccount),
		i18n.T(i18n.KeyCmdLsHeaderPlan),
	})
	if len(names) == 0 {
		fmt.Printf("\n" + i18n.T(i18n.KeyCmdLsNoProfiles, ui.C.Dim, ui.C.Off) + "\n")
	}
	return nil
}

func cmdNew(a *app.App, args parsedArgs) error {
	if len(args.positional) == 0 {
		return errors.New(i18n.T(i18n.KeyCmdNewMissingName))
	}
	name := args.positional[0]
	if err := profileenv.Validate(name); err != nil {
		return err
	}
	cfg, err := config.Load(a)
	if err != nil {
		return err
	}
	d := profileenv.Dir(a, name)
	if _, err := os.Stat(d); err == nil {
		return errors.New(i18n.T(i18n.KeyCmdNewAlreadyExists, name, d))
	}
	if _, err := os.Stat(a.ConfigPath); err != nil {
		if err := config.Save(a, cfg); err != nil {
			return err
		}
	}
	log := config.Seed(a, d, cfg, "")
	ui.Ok(i18n.T(i18n.KeyCmdNewCreated, ui.C.Bold, name, ui.C.Off, d))
	for _, line := range log {
		fmt.Printf("  %s%s%s\n", ui.C.Dim, line, ui.C.Off)
	}
	backend := credential.GetBackend()
	fmt.Printf("  %s%s: %s%s\n", ui.C.Dim, backend.Label(), backend.DisplayID(d, profileenv.IsDefault(name)), ui.C.Off)
	if args.bools["--yolo"] {
		sp := filepath.Join(d, "settings.json")
		data := map[string]any{}
		if b, err := os.ReadFile(sp); err == nil {
			_ = json.Unmarshal(b, &data)
		}
		perms, _ := data["permissions"].(map[string]any)
		if perms == nil {
			perms = map[string]any{}
		}
		perms["defaultMode"] = "bypassPermissions"
		data["permissions"] = perms
		out, _ := json.MarshalIndent(data, "", "  ")
		_ = os.WriteFile(sp, append(out, '\n'), 0644)
		fmt.Println(i18n.T(i18n.KeyCmdNewYoloNote, ui.C.Yellow, ui.C.Off, ui.C.Dim, ui.C.Off))
	}
	if args.bools["--login"] {
		fmt.Println()
		doLogin(a, name)
	} else {
		fmt.Printf("\n" + i18n.T(i18n.KeyCmdNewLoginHint, ui.C.Cyan, name, ui.C.Off) + "\n")
	}
	return nil
}

func doLogin(a *app.App, name string) {
	fmt.Printf("%s"+i18n.T(i18n.KeyCmdLoginOpeningBrowser, name)+"%s\n", ui.C.Dim, ui.C.Off)
	cmd := exec.Command(claudeBin(a), "auth", "login")
	cmd.Env = profileenv.EnvFor(a, name)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	_ = cmd.Run() // success/failure is read back via authStatus below
	st := authStatus(a, name)
	if st.LoggedIn {
		ui.Ok("'%s' → %s (%s)", name, st.Email, or(st.SubscriptionType, st.AuthMethod))
	} else {
		ui.Warn(i18n.T(i18n.KeyCmdLoginStillOut, name))
	}
}

func cmdLogin(a *app.App, args parsedArgs) error {
	name, err := requirePositional(args)
	if err != nil {
		return err
	}
	if _, err := profileenv.Require(a, name); err != nil {
		return err
	}
	doLogin(a, name)
	return nil
}

func cmdLogout(a *app.App, args parsedArgs) error {
	name, err := requirePositional(args)
	if err != nil {
		return err
	}
	if _, err := profileenv.Require(a, name); err != nil {
		return err
	}
	cmd := exec.Command(claudeBin(a), "auth", "logout")
	cmd.Env = profileenv.EnvFor(a, name)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	_ = cmd.Run()
	ui.Ok(i18n.T(i18n.KeyCmdLogoutSuccess, name))
	return nil
}

const (
	yoloFlag   = "--yolo"
	bypassFlag = "--dangerously-skip-permissions"
)

// expandYolo: `--yolo` is an alias for `--dangerously-skip-permissions`.
func expandYolo(argv []string) ([]string, bool) {
	hit := false
	out := make([]string, len(argv))
	for i, a := range argv {
		if a == yoloFlag {
			out[i] = bypassFlag
			hit = true
		} else {
			out[i] = a
		}
	}
	return out, hit
}

func runInherited(bin string, argv []string, env []string) int {
	cmd := exec.Command(bin, argv...)
	cmd.Env = env
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return ee.ExitCode()
		}
		return 1
	}
	return 0
}

func cmdUse(a *app.App, name string, rest []string) error {
	if _, err := profileenv.Require(a, name); err != nil {
		return err
	}
	if !authStatus(a, name).LoggedIn {
		ui.Warn(i18n.T(i18n.KeyCmdUseWarnNotLoggedIn, name, name))
	}
	rest, yolo := expandYolo(rest)
	if yolo {
		fmt.Println(i18n.T(i18n.KeyCmdUseYoloNote, ui.C.Yellow, ui.C.Off, ui.C.Dim, bypassFlag, ui.C.Off))
	}
	return exitCodeErr{runInherited(claudeBin(a), rest, profileenv.EnvFor(a, name))}
}

func cmdSh(a *app.App, args parsedArgs) error {
	name, err := requirePositional(args)
	if err != nil {
		return err
	}
	if _, err := profileenv.Require(a, name); err != nil {
		return err
	}
	env := profileenv.EnvFor(a, name)
	shell := shellrc.DefaultInteractiveShell()
	fmt.Printf("%s"+i18n.T(i18n.KeyCmdShSubshellBanner, profileenv.EnvGet(env, "CLAUDE_CONFIG_DIR"))+"%s\n", ui.C.Dim, ui.C.Off)
	return exitCodeErr{runInherited(shell[0], shell[1:], env)}
}

func cmdExec(a *app.App, name string, rest []string) error {
	if _, err := profileenv.Require(a, name); err != nil {
		return err
	}
	if len(rest) == 0 {
		return errors.New(i18n.T(i18n.KeyCmdExecMissingCmd))
	}
	return exitCodeErr{runInherited(rest[0], rest[1:], profileenv.EnvFor(a, name))}
}

func cmdRm(a *app.App, args parsedArgs) error {
	name, err := requirePositional(args)
	if err != nil {
		return err
	}
	if name == profileenv.DefaultName {
		return errors.New(i18n.T(i18n.KeyCmdRmCannotRmDefault))
	}
	d, err := profileenv.Require(a, name)
	if err != nil {
		return err
	}
	backend := credential.GetBackend()
	hasCred := backend.Exists(d, false)
	st := authStatus(a, name)
	if !args.bools["-y"] && !args.bools["--yes"] {
		fmt.Println(i18n.T(i18n.KeyCmdRmConfirmHeading, ui.C.Bold, name, ui.C.Off))
		fmt.Printf("  %-10s %s\n", i18n.T(i18n.KeyCmdRmDirectoryLabel), d)
		credLine := backend.DisplayID(d, false)
		if hasCred {
			credLine += " " + i18n.T(i18n.KeyCmdRmHasSession)
		} else {
			credLine += " (" + i18n.T(i18n.KeyCommonNone) + ")"
		}
		fmt.Printf("  %-10s %s\n", strings.ToLower(backend.Label()), credLine)
		if st.Email != "" {
			fmt.Printf("  %-10s %s\n", i18n.T(i18n.KeyCmdRmAccountLabel), st.Email)
		}
		fmt.Printf("  %s%s%s\n", ui.C.Dim, i18n.T(i18n.KeyCmdRmSharedHint), ui.C.Off)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(i18n.T(i18n.KeyCmdRmConfirmPrompt))
		answer, _ := reader.ReadString('\n')
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer != "y" && answer != "yes" && answer != "c" && answer != "co" && answer != "có" {
			fmt.Println(i18n.T(i18n.KeyCommonCancelled))
			return exitCodeErr{0}
		}
	}
	if !args.bools["--keep-keychain"] && hasCred {
		backend.Forget(d, false)
	}
	if err := link.RemoveProfileDir(d); err != nil {
		return fmt.Errorf("could not remove %s: %w", d, err)
	}
	ui.Ok(i18n.T(i18n.KeyCmdRmRemoved, name))
	return nil
}

func cmdSync(a *app.App, args parsedArgs) error {
	all := args.bools["--all"]
	var name string
	if len(args.positional) > 0 {
		name = args.positional[0]
	}
	if name == profileenv.DefaultName {
		return errors.New(i18n.T(i18n.KeyCmdSyncCannotSyncDefault))
	}
	cfg, err := config.Load(a)
	if err != nil {
		return err
	}
	var names []string
	if all {
		names = profileenv.List(a)
	} else {
		if name == "" {
			return errors.New(i18n.T(i18n.KeyCmdSyncMissingName))
		}
		if _, err := profileenv.Require(a, name); err != nil {
			return err
		}
		names = []string{name}
	}
	if len(names) == 0 {
		return errors.New(i18n.T(i18n.KeyCmdSyncNoProfiles))
	}
	strategy := args.values["--strategy"]
	for _, n := range names {
		log := config.Seed(a, profileenv.Dir(a, n), cfg, strategy)
		suffix := ""
		if len(log) == 0 {
			suffix = fmt.Sprintf(" %s%s%s", ui.C.Dim, i18n.T(i18n.KeyCmdSyncUpToDate), ui.C.Off)
		}
		fmt.Printf("%s%s%s%s\n", ui.C.Bold, n, ui.C.Off, suffix)
		for _, line := range log {
			fmt.Printf("  %s\n", line)
		}
	}
	return nil
}

func cmdInfo(a *app.App, args parsedArgs) error {
	name, err := requirePositional(args)
	if err != nil {
		return err
	}
	d, err := profileenv.Require(a, name)
	if err != nil {
		return err
	}
	def := profileenv.IsDefault(name)
	cd := profileenv.ConfigDirStr(a, name)
	st := authStatus(a, name)
	backend := credential.GetBackend()

	var dotclaude string
	if def {
		dotclaude = filepath.Join(a.Home, ".claude.json")
	} else {
		dotclaude = filepath.Join(d, ".claude.json")
	}
	projects := 0
	if b, err := os.ReadFile(dotclaude); err == nil {
		var data struct {
			Projects map[string]any `json:"projects"`
		}
		if json.Unmarshal(b, &data) == nil {
			projects = len(data.Projects)
		}
	}

	label := name
	if def {
		label = profileenv.DefaultName
	}
	fmt.Printf("%s%s%s\n", ui.C.Bold, label, ui.C.Off)
	if def {
		fmt.Printf("  %s %s  %s%s%s\n", ui.PadRight(i18n.T(i18n.KeyCmdInfoConfigDir), 20), d, ui.C.Dim, i18n.T(i18n.KeyCmdInfoConfigDirNotSet), ui.C.Off)
	} else {
		fmt.Printf("  %s %s\n", ui.PadRight("CLAUDE_CONFIG_DIR", 20), cd)
	}
	haveStr := i18n.T(i18n.KeyCommonNone)
	if backend.Exists(d, def) {
		haveStr = i18n.T(i18n.KeyCommonPresent)
	}
	fmt.Printf("  %s %s  %s\n", ui.PadRight(backend.Label(), 20), backend.DisplayID(d, def), haveStr)
	fmt.Printf("  %s %s  (%s)\n", ui.PadRight(i18n.T(i18n.KeyCmdInfoLoggedIn), 20), or(st.Email, "—"), or(or(st.SubscriptionType, st.AuthMethod), "—"))
	fmt.Printf("  %s %s\n", ui.PadRight(i18n.T(i18n.KeyCmdInfoStatus), 20), i18n.T(i18n.KeyCmdInfoProjectsOpened, filepath.Base(dotclaude), projects))
	if !def {
		var size int64
		_ = filepath.WalkDir(d, func(path string, entry os.DirEntry, err error) error {
			if err != nil || entry.IsDir() {
				return nil
			}
			if link.IsLink(path) {
				return nil
			}
			if fi, err := entry.Info(); err == nil {
				size += fi.Size()
			}
			return nil
		})
		fmt.Printf("  %s %.0f KB\n", ui.PadRight(i18n.T(i18n.KeyCmdInfoPrivateSize), 20), float64(size)/1024)
		entries, _ := os.ReadDir(d)
		var links []string
		for _, e := range entries {
			p := filepath.Join(d, e.Name())
			if link.IsLink(p) {
				links = append(links, e.Name())
			}
		}
		if len(links) > 0 {
			sort.Strings(links)
			fmt.Printf("  %s %s\n", ui.PadRight(i18n.T(i18n.KeyCmdInfoShared), 20), strings.Join(links, ", "))
		}
	}
	return nil
}

func cmdDoctor(a *app.App, args parsedArgs) error {
	cfg, err := config.Load(a)
	if err != nil {
		return err
	}
	names := profileenv.List(a)
	problems := 0
	fmt.Printf("%s%s%s\n", ui.C.Bold, i18n.T(i18n.KeyCmdDoctorProfilesHeading), ui.C.Off)
	backend := credential.GetBackend()
	knownIDs := map[string]bool{}
	for _, n := range names {
		d := profileenv.Dir(a, n)
		def := profileenv.IsDefault(n)
		credID := backend.DisplayID(d, def)
		knownIDs[credID] = true
		st := authStatus(a, n)
		var broken []string
		entries, _ := os.ReadDir(d)
		for _, e := range entries {
			p := filepath.Join(d, e.Name())
			if link.IsLink(p) {
				if _, err := os.Stat(p); err != nil {
					broken = append(broken, e.Name())
				}
			}
		}
		var missing []string
		for _, item := range cfg.SharedLinks {
			if _, err := os.Stat(filepath.Join(a.Main, item)); err == nil {
				if _, err := os.Stat(filepath.Join(d, item)); err != nil {
					missing = append(missing, item)
				}
			}
		}
		hasCred := backend.Exists(d, def)
		var bits []string
		if len(broken) > 0 {
			bits = append(bits, fmt.Sprintf("%s%s%s", ui.C.Red, i18n.T(i18n.KeyCmdDoctorBrokenLinks, strings.Join(broken, ", ")), ui.C.Off))
			problems++
		}
		if len(missing) > 0 {
			bits = append(bits, fmt.Sprintf("%s%s%s", ui.C.Yellow, i18n.T(i18n.KeyCmdDoctorMissingLinks, strings.Join(missing, ", ")), ui.C.Off))
			problems++
		}
		if st.LoggedIn && !hasCred {
			bits = append(bits, fmt.Sprintf("%s%s%s", ui.C.Yellow, i18n.T(i18n.KeyCmdDoctorLoggedInNoCred, strings.ToLower(backend.Label()), credID), ui.C.Off))
			problems++
		}
		if !st.LoggedIn && hasCred {
			bits = append(bits, fmt.Sprintf("%s%s%s", ui.C.Yellow, i18n.T(i18n.KeyCmdDoctorHasCredNoLogin, strings.ToLower(backend.Label())), ui.C.Off))
			problems++
		}
		if len(bits) == 0 {
			fmt.Printf("  %s: %s%s%s\n", n, ui.C.Green, i18n.T(i18n.KeyCommonOk), ui.C.Off)
		} else {
			fmt.Printf("  %s: %s\n", n, strings.Join(bits, "; "))
		}
	}

	orphans := backend.Orphans(knownIDs)
	if len(orphans) > 0 {
		fmt.Printf("\n%s"+i18n.T(i18n.KeyCmdDoctorOrphansHeading, backend.Label())+"%s\n", ui.C.Bold, ui.C.Off)
		sort.Strings(orphans)
		for _, s := range orphans {
			fmt.Printf("  %s%s%s  %s%s%s\n", ui.C.Yellow, s, ui.C.Off, ui.C.Dim, i18n.T(i18n.KeyCmdDoctorOrphanHint), ui.C.Off)
		}
	}

	if problems > 0 {
		return exitCodeErr{1}
	}
	return nil
}

func cmdSettings(a *app.App, args parsedArgs) error {
	cfg, err := config.Load(a)
	if err != nil {
		return err
	}

	if len(args.positional) == 0 {
		// Dashboard overview
		fmt.Printf("%s%s%s\n\n", ui.C.Bold, i18n.T(i18n.KeyCmdSettingsTitle), ui.C.Off)

		sourceStr := i18n.T(i18n.KeyCmdSettingsSourceDefault)
		if cfg.Lang != "" {
			sourceStr = i18n.T(i18n.KeyCmdSettingsSourceConfig)
		} else if os.Getenv("CCA_LANG") != "" {
			sourceStr = i18n.T(i18n.KeyCmdSettingsSourceEnv)
		}

		curLangName := i18n.Name(i18n.Current())
		fmt.Printf("  %-25s %s  %s(%s)%s\n", i18n.T(i18n.KeyCmdSettingsLangLabel)+":", curLangName, ui.C.Dim, sourceStr, ui.C.Off)
		fmt.Printf("  %-25s %s\n\n", i18n.T(i18n.KeyCmdSettingsSyncLabel)+":", cfg.SyncStrategy)

		fmt.Printf("%s%s%s\n", ui.C.Dim, i18n.T(i18n.KeyCmdSettingsOptionsHeading), ui.C.Off)
		fmt.Printf("  %s\n", i18n.T(i18n.KeyCmdSettingsOptLangMenu))
		fmt.Printf("  %s\n", i18n.T(i18n.KeyCmdSettingsOptLangCode))
		return nil
	}

	sub := args.positional[0]
	if sub == "lang" || sub == "language" {
		if len(args.positional) > 1 {
			target := args.positional[1]
			if target == "auto" || target == "system" {
				cfg.Lang = ""
				if err := config.Save(a, cfg); err != nil {
					return err
				}
				i18n.Init(i18n.Resolve("", ""))
				ui.Ok(i18n.T(i18n.KeyCmdSettingsLangAutoSet, i18n.Name(i18n.Current())))
				return nil
			}
			norm := i18n.Normalize(target)
			if norm == "" {
				return errors.New(i18n.T(i18n.KeyCmdSettingsLangErrInvalid, target, strings.Join(i18n.SupportedCodes(), ", ")))
			}
			cfg.Lang = norm
			if err := config.Save(a, cfg); err != nil {
				return err
			}
			i18n.Init(norm)
			ui.Ok(i18n.T(i18n.KeyCmdSettingsLangUpdated, i18n.Name(norm)))
			return nil
		}

		// Interactive menu
		langs := i18n.SupportedLanguages()
		fmt.Printf("\n%s%s%s\n\n", ui.C.Bold, i18n.T(i18n.KeyCmdSettingsPromptTitle), ui.C.Off)

		cur := i18n.Current()
		for i, l := range langs {
			active := ""
			if l.Code == cur && cfg.Lang != "" {
				active = fmt.Sprintf("  %s✓ [%s]%s", ui.C.Green, i18n.T(i18n.KeyCommonActive), ui.C.Off)
			}
			label := fmt.Sprintf("[%d] %-18s (%s)", i+1, l.NativeName, l.Code)
			fmt.Printf("  %s%s\n", label, active)
		}
		autoActive := ""
		if cfg.Lang == "" {
			autoActive = fmt.Sprintf("  %s✓ [%s]%s", ui.C.Green, i18n.T(i18n.KeyCommonActive), ui.C.Off)
		}
		fmt.Printf("  [%d] %-18s (auto: %s)%s\n\n", len(langs)+1, "Auto (System)", cur, autoActive)

		fmt.Print("  " + i18n.T(i18n.KeyCmdSettingsPromptChoice, len(langs)+1))
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" || input == "q" || input == "Q" {
			fmt.Println("  " + i18n.T(i18n.KeyCommonCancelled))
			return nil
		}

		var choice int
		_, err := fmt.Sscanf(input, "%d", &choice)
		if err != nil || choice < 1 || choice > len(langs)+1 {
			return errors.New(i18n.T(i18n.KeyCmdSettingsInvalidChoice, len(langs)+1))
		}

		if choice == len(langs)+1 {
			cfg.Lang = ""
			if err := config.Save(a, cfg); err != nil {
				return err
			}
			i18n.Init(i18n.Resolve("", ""))
			ui.Ok(i18n.T(i18n.KeyCmdSettingsLangAutoSet, i18n.Name(i18n.Current())))
			return nil
		}

		chosen := langs[choice-1].Code
		cfg.Lang = chosen
		if err := config.Save(a, cfg); err != nil {
			return err
		}
		i18n.Init(chosen)
		ui.Ok(i18n.T(i18n.KeyCmdSettingsLangUpdated, i18n.Name(chosen)))
		return nil
	}

	return fmt.Errorf("unknown settings option '%s' — try: cca settings lang", sub)
}

func cmdConfig(a *app.App, args parsedArgs) error {
	cfg, err := config.Load(a)
	if err != nil {
		return err
	}
	if args.bools["--print"] {
		if _, err := os.Stat(a.ConfigPath); err != nil {
			fmt.Printf("%s(no %s yet — using defaults)%s\n", ui.C.Dim, a.ConfigPath, ui.C.Off)
		}
		out, _ := json.MarshalIndent(cfg, "", "  ")
		fmt.Println(string(out))
		return nil
	}
	if _, err := os.Stat(a.ConfigPath); err != nil {
		if err := config.Save(a, cfg); err != nil {
			return err
		}
	}
	if args.bools["--edit"] {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano"
			if app.IsWindows() {
				editor = "notepad"
			}
		}
		cmd := exec.Command(editor, a.ConfigPath)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		_ = cmd.Run()
		return nil
	}
	codeBin, err := exec.LookPath("code")
	if err != nil {
		return fmt.Errorf("could not find 'code' in PATH.\n"+
			"  Install: open VS Code → Ctrl+Shift+P → \"Shell Command: Install 'code' command in PATH\"\n"+
			"  Or view it manually: %s\n"+
			"  Print the JSON to the terminal instead of opening VS Code: cca config --print", a.Root)
	}
	_ = exec.Command(codeBin, a.Root).Run()
	ui.Ok("Opened %s in VS Code", a.Root)
	return nil
}

func cmdVersion(a *app.App) error {
	fmt.Println("cca " + Version)
	cfg, err := config.Load(a)
	if err == nil && cfg.LatestVersion != "" && update.CompareVersions(Version, cfg.LatestVersion) < 0 {
		fmt.Printf("\n%s! %s%s\n  %s\n",
			ui.C.Yellow,
			i18n.T(i18n.KeyCmdUpdateAvailableNotice, Version, cfg.LatestVersion),
			ui.C.Off,
			i18n.T(i18n.KeyCmdUpdateRunHint),
		)
	}
	return nil
}

func cmdUpdate(a *app.App, args parsedArgs) error {
	fmt.Println(ui.C.Dim + i18n.T(i18n.KeyCmdUpdateChecking) + ui.C.Off)

	rel, hasUpdate, err := update.CheckLatest(Version)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyCmdUpdateCheckFailed), err)
	}

	cfg, _ := config.Load(a)
	cfg.LastUpdateCheck = time.Now().Unix()
	if hasUpdate {
		cfg.LatestVersion = rel.Version
	} else {
		cfg.LatestVersion = ""
	}
	_ = config.Save(a, cfg)

	if !hasUpdate {
		ui.Ok(i18n.T(i18n.KeyCmdUpdateAlreadyLatest, Version))
		return nil
	}

	if args.bools["--check"] || args.bools["-c"] {
		ui.Warn(i18n.T(i18n.KeyCmdUpdateAvailableNotice, Version, rel.Version))
		fmt.Printf("  %s\n", i18n.T(i18n.KeyCmdUpdateRunHint))
		if rel.HTMLURL != "" {
			fmt.Printf("  %s%s%s\n", ui.C.Dim, rel.HTMLURL, ui.C.Off)
		}
		return nil
	}

	if Version == "dev" {
		ui.Warn("Running a development build. Proceeding to update to latest release: %s", rel.Version)
	}

	fmt.Println(i18n.T(i18n.KeyCmdUpdateDownloading, rel.Version))
	if err := update.Apply(a, rel); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyCmdUpdateApplyFailed), err)
	}

	ui.Ok(i18n.T(i18n.KeyCmdUpdateSuccess, rel.Version))
	return nil
}

func cmdInstall(a *app.App, args parsedArgs) error {
	rcPath := args.values["--rc"]
	kind := args.values["--rc-kind"]

	if rcPath == "" {
		rcPath, kind = shellrc.DetectRC(a)
	} else if kind == "" {
		kind = "bash"
	}

	bindir := filepath.Join(a.Main, "bin")
	if app.IsWindows() {
		switch added, perr := shellrc.AddToWindowsUserPath(bindir); {
		case perr != nil:
			ui.Warn("could not update the Windows PATH automatically: %v", perr)
			fmt.Printf("  Add this to PATH manually: %s%s%s\n", ui.C.Cyan, bindir, ui.C.Off)
			fmt.Println("  Windows: Settings → System → About → Advanced system settings → Environment Variables → Path")
		case added:
			ui.Ok("Added %s to your Windows user PATH", bindir)
			fmt.Printf("  %sTakes effect in any terminal opened from now on — this one included.%s\n", ui.C.Dim, ui.C.Off)
		}
	}

	if kind == "unknown" {
		if !app.IsWindows() {
			fmt.Printf("%sCould not auto-detect a shell with an rc file.%s\n", ui.C.Yellow, ui.C.Off)
			fmt.Printf("  Add this to PATH manually: %s%s%s\n", ui.C.Cyan, bindir, ui.C.Off)
		}
		installSelf(a)
		return nil
	}

	snippet := shellrc.Snippet(kind)
	_ = os.MkdirAll(filepath.Dir(rcPath), 0755)
	text := ""
	if b, err := os.ReadFile(rcPath); err == nil {
		text = string(b)
	}
	if strings.Contains(text, shellrc.MarkStart) {
		start := strings.Index(text, shellrc.MarkStart)
		end := strings.Index(text, shellrc.MarkEnd) + len(shellrc.MarkEnd)
		_ = os.WriteFile(rcPath, []byte(text[:start]+snippet+text[end:]), 0644)
		ui.Ok("Updated the cca block in %s", rcPath)
	} else {
		_ = os.WriteFile(rcPath, []byte(strings.TrimRight(text, "\n")+"\n\n"+snippet+"\n"), 0644)
		ui.Ok("Added cca to %s", rcPath)
	}
	fmt.Printf("  %sRemove it: delete the block between '%s' and '%s'%s\n", ui.C.Dim, shellrc.MarkStart, shellrc.MarkEnd, ui.C.Off)
	reloadHint := rcPath
	if kind == "powershell" {
		reloadHint = "$PROFILE"
	}
	fmt.Printf("  Reload: %s. %s%s\n", ui.C.Cyan, reloadHint, ui.C.Off)

	installSelf(a)
	return nil
}
