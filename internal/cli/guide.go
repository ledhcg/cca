package cli

import (
	"fmt"

	"github.com/ledhcg/cca/internal/ui"
)

func printUsage() {
	fmt.Println(`cca — multiple Claude Code accounts on one machine

Usage: cca <command> [args...]

Commands:
  ls                    list profiles and which account each is logged into
  new <name>              create a new profile (--login, --yolo)
  use <name> [args…]      run Claude Code under a profile
  login <name>             log in a profile
  logout <name>            log out a profile
  info <name>               show details for a profile
  sh <name>                  open a subshell in a profile's environment
  exec <name> -- <cmd>        run an arbitrary command in a profile's environment
  rm <name>                    remove a profile (-y, --keep-keychain)
  sync [<name>|--all]           refresh shared files (--strategy)
  doctor                         check links, credentials, orphaned profiles
  config                          open ~/.claude-accounts in VS Code (--edit, --print)
  guide                            full usage guide
  install                          add cca to PATH + completion
  version                          print the cca version

Examples:
  cca new work --login       create profile 'work' and log in right away
  cca work                   run Claude Code under profile 'work'
  cca work --resume          any trailing args are passed straight to claude
  cca ls                     see which profile is logged into which account
  cca sync --all             refresh shared files after installing a new plugin
  cca work --yolo            alias for --dangerously-skip-permissions
  cca sh work                open a subshell with CLAUDE_CONFIG_DIR already set
  cca exec work -- git log   run any command inside a profile's environment

See also: cca guide`)
}

func printGuide() {
	b, o, c, d := ui.C.Bold, ui.C.Off, ui.C.Cyan, ui.C.Dim
	fmt.Printf(`
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
  %[4]s--yolo only has an effect inside cca; running `+"`claude`"+` directly is untouched.%[2]s

%[1]sMANAGING PROFILES%[2]s
  %[3]scca new <name> [--login] [--yolo]%[2]s   create; --yolo presets bypassPermissions
  %[3]scca login <name>%[2]s / %[3]scca logout <name>%[2]s     log a profile in / out
  %[3]scca info <name>%[2]s                     directory, credential, account, size
  %[3]scca rm <name>%[2]s                       remove a profile (asks for confirmation)

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
`, b, o, c, d)
}
