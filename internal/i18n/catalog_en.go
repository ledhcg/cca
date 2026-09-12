package i18n

var enCatalog = map[Key]string{
	// Common / Generic
	KeyCommonDefault:   "default",
	KeyCommonNone:      "none",
	KeyCommonCancelled: "Cancelled",
	KeyCommonOk:        "ok",
	KeyCommonPresent:   "present",
	KeyCommonActive:    "Active",

	// Commands - ls
	KeyCmdLsHeaderProfile: "PROFILE",
	KeyCmdLsHeaderStatus:  "STATUS",
	KeyCmdLsHeaderAccount: "ACCOUNT",
	KeyCmdLsHeaderPlan:    "PLAN",
	KeyCmdLsStatusIn:      "logged in",
	KeyCmdLsStatusOut:     "logged out",
	KeyCmdLsNoProfiles:    "No extra profiles yet. Create one: %scca new work --login%s",

	// Commands - info
	KeyCmdInfoConfigDir:       "Config dir",
	KeyCmdInfoConfigDirNotSet: "(CLAUDE_CONFIG_DIR not set)",
	KeyCmdInfoLoggedIn:        "Logged in",
	KeyCmdInfoStatus:          "Status",
	KeyCmdInfoProjectsOpened:  "%s · %d project(s) opened",
	KeyCmdInfoPrivateSize:     "Private size",
	KeyCmdInfoShared:          "Shared",

	// Commands - doctor
	KeyCmdDoctorProfilesHeading: "Profiles",
	KeyCmdDoctorBrokenLinks:     "broken link(s): %s",
	KeyCmdDoctorMissingLinks:    "missing link(s): %s (run cca sync --all)",
	KeyCmdDoctorLoggedInNoCred:  "logged in but no %s found: %s",
	KeyCmdDoctorHasCredNoLogin:  "has a %s but couldn't log in",
	KeyCmdDoctorOrphansHeading:  "%s not owned by any profile",
	KeyCmdDoctorOrphanHint:      "(profile removed by hand?)",

	// Commands - new
	KeyCmdNewMissingName:   "missing profile name — example: cca new work",
	KeyCmdNewAlreadyExists: "profile '%s' already exists at %s",
	KeyCmdNewCreated:       "Created profile %s%s%s → %s",
	KeyCmdNewLoginHint:     "Log in:  %scca login %s%s",
	KeyCmdNewYoloNote:      "  %s⚡ defaultMode = bypassPermissions%s %s(this profile only)%s",

	// Commands - login / logout
	KeyCmdLoginOpeningBrowser: "Opening the browser to log in profile '%s'…",
	KeyCmdLoginStillOut:       "'%s' is still logged out",
	KeyCmdLogoutSuccess:       "Logged out profile '%s'",

	// Commands - use / sh / exec
	KeyCmdUseWarnNotLoggedIn: "Profile '%s' is not logged in — run: cca login %s",
	KeyCmdUseYoloNote:        "%s⚡ bypass permissions%s %s(--yolo → %s)%s",
	KeyCmdShSubshellBanner:   "Subshell with CLAUDE_CONFIG_DIR=%s — type exit to leave",
	KeyCmdExecMissingCmd:     "missing command to run — example: cca exec work -- claude auth status",

	// Commands - rm
	KeyCmdRmCannotRmDefault: "cannot remove the default profile — it's just ~/.claude",
	KeyCmdRmConfirmHeading:  "About to remove profile %s%s%s:",
	KeyCmdRmDirectoryLabel:  "directory",
	KeyCmdRmAccountLabel:    "account",
	KeyCmdRmHasSession:      "(still holds a logged-in session)",
	KeyCmdRmSharedHint:      "Shared links are only unlinked — ~/.claude itself is untouched.",
	KeyCmdRmConfirmPrompt:   "Confirm? [y/N] ",
	KeyCmdRmCancelled:       "Cancelled",
	KeyCmdRmRemoved:         "Removed profile '%s'",

	// Commands - sync
	KeyCmdSyncCannotSyncDefault: "cannot sync the default profile — it's just ~/.claude",
	KeyCmdSyncMissingName:       "missing profile name — example: cca sync work   (or cca sync --all)",
	KeyCmdSyncNoProfiles:        "no profiles to sync yet",
	KeyCmdSyncUpToDate:          "(already up to date)",

	// Commands - settings
	KeyCmdSettingsTitle:          "⚙  cca Settings",
	KeyCmdSettingsLangLabel:      "Display language",
	KeyCmdSettingsSyncLabel:      "Sync strategy",
	KeyCmdSettingsSourceConfig:   "source: config.json",
	KeyCmdSettingsSourceEnv:      "source: environment",
	KeyCmdSettingsSourceDefault:  "source: system default",
	KeyCmdSettingsOptionsHeading: "Available options:",
	KeyCmdSettingsOptLangMenu:    "cca settings lang          Select language via interactive menu",
	KeyCmdSettingsOptLangCode:    "cca settings lang <code>   Quickly set language (en, vi, zh, ja, es, auto)",
	KeyCmdSettingsPromptTitle:    "🌐  Select Display Language",
	KeyCmdSettingsPromptChoice:   "Select [1-%d] or 'q' to cancel: ",
	KeyCmdSettingsInvalidChoice:  "Invalid choice. Please choose between 1 and %d.",
	KeyCmdSettingsLangUpdated:    "Display language updated: %s",
	KeyCmdSettingsLangAutoSet:    "Language reset to auto (current system: %s)",
	KeyCmdSettingsLangErrInvalid: "Unsupported language '%s'. Supported: %s, or 'auto'",

	// Commands - update
	KeyCmdUpdateChecking:        "Checking for updates…",
	KeyCmdUpdateAlreadyLatest:   "cca is already up to date (%s)",
	KeyCmdUpdateAvailableNotice: "A new version of cca is available: %s → %s",
	KeyCmdUpdateRunHint:         "Run 'cca update' to upgrade.",
	KeyCmdUpdateDownloading:     "Downloading cca %s…",
	KeyCmdUpdateSuccess:         "Successfully updated cca to %s",
	KeyCmdUpdateCheckFailed:     "failed to check for updates: %v",
	KeyCmdUpdateApplyFailed:     "failed to apply update: %v",
	KeyCmdUpdateNoAsset:         "no prebuilt binary available for %s/%s in release %s",

	// Commands - handoff
	KeyCmdHandoffMissingTarget:  "missing target profile name — example: cca handoff work",
	KeyCmdHandoffSameProfile:    "source and target profile cannot be the same",
	KeyCmdHandoffNoSessionFound: "no sessions found in profile '%s' for this project",
	KeyCmdHandoffBanner:         "Handoff session %s (%s) from '%s' to '%s'…",
	KeyCmdHandoffCopyFailed:     "failed to hand off session: %w",

	// Commands - session
	KeyCmdSessionUnknownSubcmd:     "unknown session command '%s' — try: ls, cp, mv, rm",
	KeyCmdSessionMissingFromTo:     "--from and --to are both required — example: cca session cp --from work --to personal --all",
	KeyCmdSessionSameProfile:       "--from and --to cannot be the same profile",
	KeyCmdSessionRequireIdOrAll:    "must specify either --id <sessionId> or --all",
	KeyCmdSessionNoSessionsFound:   "No sessions found in profile '%s' for this project",
	KeyCmdSessionNoSessionsAllHint: "(use --all to view sessions across all projects)",
	KeyCmdSessionCopySuccess:       "Copied session %s from '%s' to '%s'",
	KeyCmdSessionCopyAllSuccess:    "Copied %d session(s) from '%s' to '%s'",
	KeyCmdSessionMoveSuccess:       "Moved session %s from '%s' to '%s'",
	KeyCmdSessionMoveAllSuccess:    "Moved %d session(s) from '%s' to '%s'",
	KeyCmdSessionRmSuccess:         "Removed session %s from '%s'",
	KeyCmdSessionRmAllSuccess:      "Removed %d session(s) from '%s'",
	KeyCmdSessionRmConfirmSingle:   "About to remove session %s from profile '%s'. Confirm? [y/N] ",
	KeyCmdSessionRmConfirmAll:      "About to remove %d session(s) for this project from profile '%s'. Confirm? [y/N] ",
	KeyCmdSessionHeaderID:          "SESSION ID",
	KeyCmdSessionHeaderProject:     "PROJECT",
	KeyCmdSessionHeaderTitle:       "TITLE",
	KeyCmdSessionHeaderMessages:    "MSGS",
	KeyCmdSessionHeaderSize:        "SIZE",
	KeyCmdSessionHeaderModified:    "MODIFIED",

	// Guide & Usage
	KeyGuideUsage: `cca — multiple Claude Code accounts on one machine

Usage: cca <command> [args...]

Commands:
  ls                    list profiles and which account each is logged into
  new <name>              create a new profile (--login, --yolo)
  use <name> [args…]      run Claude Code under a profile
  handoff <to> [args…]     continue current session in another profile (--fork)
  session <subcommand>     manage sessions: ls, cp, mv, rm (--from, --to, --all, --id)
  login <name>             log in a profile
  logout <name>            log out a profile
  info <name>               show details for a profile
  sh <name>                  open a subshell in a profile's environment
  exec <name> -- <cmd>        run an arbitrary command in a profile's environment
  rm <name>                    remove a profile (-y, --keep-keychain)
  sync [<name>|--all]           refresh shared files (--strategy)
  doctor                         check links, credentials, orphaned profiles
  settings [lang]                 manage settings & display language
  update [--check]                check for updates and self-update
  config                          open ~/.claude-accounts in VS Code (--edit, --print)
  guide                            full usage guide
  install                          add cca to PATH + completion
  version                          print the cca version

Examples:
  cca new work --login       create profile 'work' and log in right away
  cca work                   run Claude Code under profile 'work'
  cca work --resume          any trailing args are passed straight to claude
  cca handoff work           handoff active session to 'work' and continue
  cca session ls             list recent sessions in current directory
  cca session cp --from default --to work --all   copy all sessions
  cca ls                     see which profile is logged into which account
  cca sync --all             refresh shared files after installing a new plugin
  cca settings lang          change display language
  cca work --yolo            alias for --dangerously-skip-permissions
  cca sh work                open a subshell with CLAUDE_CONFIG_DIR already set
  cca exec work -- git log   run any command inside a profile's environment

See also: cca guide`,

	KeyGuideFull: `
%[1]scca — multiple Claude Code accounts on one machine%[2]s

%[1]sQUICK START%[2]s
  %[3]scca new work --login%[2]s     create profile 'work', open the browser to log in
  %[3]scca ls%[2]s                   see which profile is logged into which account
  %[3]scca work%[2]s                 run Claude Code under profile 'work'

  The root profile (~/.claude) is always available as %[3]sdefault%[2]s — no need to create it.

%[1]sHOW IT WORKS%[2]s
  Each profile is its own config directory under ~/.claude-accounts/. The login
  session lives in the Keychain (macOS) or in a .credentials.json file inside
  the profile directory itself (Linux/Windows) — each profile keeps an
  independent session, so running two accounts side by side in two terminals
  just works.

%[1]sRUNNING CLAUDE%[2]s
  %[3]scca work%[2]s                     normal interactive session
  %[3]scca work --resume%[2]s            any trailing args are passed straight to claude
  %[3]scca work -p "question"%[2]s       one-shot print mode
  %[3]scca work --yolo%[2]s              alias for --dangerously-skip-permissions
  %[3]scca default --yolo%[2]s           the root account, skipping permission checks
  %[4]s--yolo only has an effect inside cca; running ` + "`claude`" + ` directly is untouched.%[2]s

%[1]sMANAGING PROFILES%[2]s
  %[3]scca new <name> [--login] [--yolo]%[2]s   create; --yolo presets bypassPermissions
  %[3]scca login <name>%[2]s / %[3]scca logout <name>%[2]s     log a profile in / out
  %[3]scca info <name>%[2]s                     directory, credential, account, size
  %[3]scca rm <name>%[2]s                       remove a profile (asks for confirmation)
  %[3]scca settings [lang]%[2]s                 configure display language and options

%[1]sSHARED FILES%[2]s
  plugins, skills, agents are linked back to ~/.claude — install a plugin once
  and every profile sees it. On Windows without Developer Mode enabled, cca
  falls back to junctions/hard links/copies when it can't create a real
  symlink (see %[3]scca doctor%[2]s).

  %[3]scca sync --all%[2]s      refresh shared files (run after installing a new plugin)
  %[3]scca config%[2]s          open ~/.claude-accounts in VS Code to edit config.json

%[1]sTROUBLESHOOTING%[2]s
  %[3]scca doctor%[2]s   broken links, a profile that looks logged in but has no credential.

%[1]sSEE ALSO%[2]s
  %[3]scca <command> --help%[2]s   per-command details
  %[3]scca version%[2]s            print the cca version
`,

	// Profile validation errors
	KeyProfileErrReserved: "'%s' is a reserved name — the default profile is already ~/.claude",
	KeyProfileErrInvalid:  "profile names may only contain letters, digits, and . _ - and must start with a letter or digit",
	KeyProfileErrNotFound: "no profile '%s'. Available: %s",
	KeyProfileErrCreateIt: "  Create it with: cca new %s",

	// General CLI errors
	KeyCliErrMissingProfileName: "missing profile name — example: cca use work",
	KeyCliErrNotCommandOrProf:   "'%s' is neither a command nor an existing profile.",
	KeyCliErrExistingProfiles:   "  Existing profiles: %s",
	KeyCliErrNoExtraProfiles:    "  No extra profiles yet (only 'default' = ~/.claude).",
	KeyCliErrCreateThisProfile:  "  Create this profile:  cca new %s --login",
	KeyCliErrSeeCommandsGuide:   "  See commands:       cca --help   ·   Full guide: cca guide",
}
