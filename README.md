# cca — multiple Claude Code accounts on one machine

[![CI](https://github.com/ledhcg/cca/actions/workflows/ci.yml/badge.svg)](https://github.com/ledhcg/cca/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/ledhcg/cca)](https://github.com/ledhcg/cca/releases/latest)
[![Go](https://img.shields.io/badge/go-1.22%2B-00ADD8?logo=go&logoColor=white)](go.mod)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Run a personal and a work account side by side without `/login`/`/logout`
juggling. A single static Go binary — download and run, no Python/Node or any
other runtime required. Works on macOS, Linux, and Windows (PowerShell or Git
Bash).

```
$ cca ls
PROFILE    STATUS      ACCOUNT              PLAN
default ←  logged in   personal@gmail.com   max
work       logged in   you@company.com      team
```

## Features

- **Isolated profiles** — each account gets its own `CLAUDE_CONFIG_DIR`, so
  sessions, credentials, and history never collide.
- **Shared plugins/skills/agents** — install a plugin once, every profile
  sees it, via real symlinks with an automatic Windows fallback.
- **Multi-language support (i18n)** — out of the box in 5 languages (English,
  Vietnamese, Simplified Chinese, Japanese, Spanish) with `cca settings lang`.
- **In-place self-update** — check and update to the latest release seamlessly
  with `cca update` (plus 24h background checks).
- **No runtime required** — a single static Go binary; no Python, Node, or
  package manager needed to run `cca` itself.
- **Native credential storage per OS** — macOS Keychain, plain file on
  Linux/Windows — matching how Claude Code itself stores credentials.
- **One-command setup** — `cca install` wires up `PATH` and shell completion
  for bash/zsh/fish/PowerShell/cmd.exe.
- **Safety net** — `cca doctor` catches broken symlinks, orphaned
  credentials, and profiles that drifted out of sync.

## How it works

Claude Code reads the `CLAUDE_CONFIG_DIR` environment variable. Where it stores
the login credential differs by OS — verified by reading the `claude` binary
directly (build 2.1.211/2.1.234) and cross-checking against the official docs
at [code.claude.com/docs/en/authentication](https://code.claude.com/docs/en/authentication):

| OS | Credential storage |
|---|---|
| macOS | System Keychain, entry `Claude Code-credentials-<sha256(CLAUDE_CONFIG_DIR)[:8]>` |
| Linux | `<CLAUDE_CONFIG_DIR>/.credentials.json` (mode 0600) |
| Windows | `<CLAUDE_CONFIG_DIR>\.credentials.json` |

```js
// excerpt from the Claude Code binary — how macOS names its Keychain entry
let r = !env.CLAUDE_CONFIG_DIR,                                  // true ⇔ default directory
    n = In(),                                                    // directory in use
    o = r ? "" : `-${sha256(n).hex.slice(0,8)}`;
return `Claude Code-credentials${o}`
```

On macOS, each config directory gets its own Keychain entry keyed by a hash of
the path. On Linux/Windows it's simpler: the credential lives right inside the
profile directory, so removing the directory removes the session too — no
extra cleanup needed. `internal/credential` hides this difference behind a
single interface; the rest of `cca` doesn't need to know which OS it's
running on.

`cca` just provisions these profile directories under `~/.claude-accounts/`,
links in the shared files, and runs `claude` with the right environment variable.

Worth remembering: **the account is keyed by path**. Renaming a profile
directory breaks its session (on macOS, permanently — the hash changes; on
Linux/Windows it could technically survive a move, but `cca` intentionally
doesn't support that so behavior stays consistent across OSes). That's why
there's no `rename` command.

## Install

> **Note:** the one-liners below (`curl`/`irm` against `raw.githubusercontent.com`
> and GitHub Releases) only work once this repository is **public** — a
> private repo returns 404 to anyone without access. Until then, clone the
> repo and build from source (see "Building manually" below).

**macOS / Linux / Git Bash on Windows:**
```sh
curl -fsSL https://raw.githubusercontent.com/ledhcg/cca/main/scripts/install.sh | sh
source ~/.bashrc   # or ~/.zshrc / ~/.config/fish/config.fish depending on your shell
```

**Windows, plain PowerShell (no Git Bash needed):**
```powershell
irm https://raw.githubusercontent.com/ledhcg/cca/main/scripts/install.ps1 | iex
. $PROFILE
```

Or clone the repo and run `install.sh`/`install.ps1` from it directly (useful
if you'd rather build from source than download a prebuilt binary):
```sh
git clone https://github.com/ledhcg/cca
cd cca && ./scripts/install.sh
```

Both scripts detect the OS and CPU architecture, download the matching binary
from the latest [GitHub Release](https://github.com/ledhcg/cca/releases), drop
it into `~/.claude/bin/`, then call `cca install` to handle the rest:
- Detects the running shell (bash/zsh/fish/PowerShell) and appends a block to
  its rc file (PATH + completion), marked with `# >>> cca >>>` for easy removal.
- On Windows, also adds `~/.claude/bin` straight to the user's PATH via the
  registry (`HKCU\Environment`) — covers `cmd.exe`, which has no rc file to
  source, plus any shell opened after install. Idempotent: safe to run
  `cca install` again, it won't duplicate the entry. Takes effect in new
  terminals only; a window opened before running `cca install` won't see it.
- For `plugins`/`skills`/`agents`: tries a real symlink first; if Windows
  Developer Mode is off (the default), it falls back to a junction
  (directories) or a hard link (files) — no admin rights needed. Run
  `cca doctor` to see which one is in use.

No release available yet (or no network), but [Go](https://go.dev/dl/) 1.22+
is installed? Both scripts fall back to `go build` from source automatically.

### Building manually

```sh
git clone https://github.com/ledhcg/cca
cd cca
make build                # or: go build -o cca ./cmd/cca
./cca install
```

Cross-compiling for another OS (no need for that machine):
```sh
GOOS=windows GOARCH=amd64 go build -o cca-windows-amd64.exe ./cmd/cca
GOOS=darwin  GOARCH=arm64 go build -o cca-darwin-arm64 ./cmd/cca
```

## Commands

| Command | What it does |
|---|---|
| `cca new work --login` | create a profile and log in right away |
| `cca new yolo --yolo` | create a profile with `defaultMode: bypassPermissions` preset |
| `cca work [args…]` | run Claude under profile `work`; args are passed straight to `claude` |
| `cca default` | run the root account (`~/.claude`) |
| `cca ls` | which profile is logged into which account |
| `cca sync --all` | refresh shared files after installing a new plugin |
| `cca doctor` | broken links, orphaned credentials, drifted state |
| `cca info <name>` | directory, credential, account, disk usage |
| `cca sh <name>` | subshell with `CLAUDE_CONFIG_DIR` already set |
| `cca exec <name> -- <cmd>` | run an arbitrary command in a profile's environment |
| `cca rm <name>` | remove a profile and its credential (asks for confirmation) |
| `cca config` | open `~/.claude-accounts` in VS Code |
| `cca config --edit` | open `config.json` in `$EDITOR` |
| `cca config --print` | print `config.json` to the terminal |
| `cca install` | add cca to PATH + completion, update the installed binary in place |
| `cca version` | print the cca version |

`--yolo` is an alias for `--dangerously-skip-permissions`, and only has an
effect inside `cca` — running `claude` directly is untouched.

## Shared files

`~/.claude-accounts/config.json`:

```json
{
  "sharedLinks": ["plugins", "skills", "agents", "statusline-command.sh"],
  "sharedCopies": ["settings.json"],
  "syncStrategy": "merge",
  "profileLocalKeys": ["model", "effortLevel", "permissions"]
}
```

- `sharedLinks` — linked back to `~/.claude` (a symlink, or a junction/hard
  link when Windows can't grant symlink permission). Install a plugin once and
  every profile sees it.
- `sharedCopies` — each profile keeps its own editable copy.
- `syncStrategy` — `overwrite` (always take the original) · `keep-local` (keep
  local edits) · `merge` (take the original but preserve keys listed in
  `profileLocalKeys`).
- `profileLocalKeys` — the boundary between "shared config" and "this
  account's own choice". `model` lives here because a Max plan's
  `opus[1m]` isn't usable on a Pro plan.

`projects/`, `history.jsonl`, and `sessions/` are **never** shared — two
processes running in parallel would just overwrite each other.

## Known constraints

- **`claude daemon service install` only works with the default config
  directory** — the binary refuses outright: *"the launchd/systemd unit is a
  per-user singleton"*. Secondary profiles still work with Remote Control, just
  with the daemon starting on demand instead.
- **Bypass permissions set via `settings.json` is silently ignored in a
  session owned by VS Code** — the binary only logs a warning, nothing shows
  on screen. The command-line flag (`cca <name> --yolo`) is unaffected.
- **`--dangerously-skip-permissions` is blocked when running as root/sudo**,
  unless `IS_SANDBOX=1` or running inside bubblewrap.
- **The default profile keeps `.claude.json` at `~/.claude.json`** — outside
  the config directory — while secondary profiles keep it at
  `<dir>/.claude.json` inside. This asymmetry silently produces wrong numbers
  if read from the wrong place.
- **Windows with Developer Mode off (the default)**: real symlinks can't be
  created — `plugins`/`skills`/`agents` use a junction, `statusline-command.sh`
  uses a hard link. Neither needs admin rights. Enable Developer Mode
  (Settings → Privacy & security → For developers) for real symlinks instead.
- **Windows won't let you overwrite a running `.exe`** — when `cca install`
  updates itself (overwriting the binary currently executing it), it has to
  rename the old file to `cca.exe.old` first, then place the new one. Writing
  directly on top fails with "Access is denied". This is an implementation
  detail and doesn't affect normal use.

## Verified on

- macOS 15 (darwin-arm64), zsh — Claude Code `2.1.234` (native build).
- Windows 10, Git Bash, Go 1.26, Developer Mode off — Claude Code `2.1.211`.
  The `.credentials.json` mechanism was confirmed against the official docs
  (code.claude.com/docs/en/authentication) and by reading the binary directly.
  The safety of `os.RemoveAll`/`os.Symlink`/junctions on Windows was verified
  experimentally (a canary file placed in the real shared directory survives
  `cca rm` even when the profile contains a junction pointing at it).

`CLAUDE_CONFIG_DIR` and the credential storage location aren't fully documented
for every version — the mechanism above comes from reading the binary/docs and
verifying experimentally, so run `cca doctor` after every Claude Code
auto-update: it flags a profile that looks logged in but has no matching
credential.

## Development

```
cmd/cca/main.go            entry point: app.New() → cli.Run() → os.Exit(code)
internal/app                runtime paths (Home/Main/Root/ConfigPath), no globals
internal/cli                 dispatch, argument parser, every cmd* command
internal/config                shared config.json, seed()/sync
internal/credential              Keychain (macOS) / file (Linux, Windows) interface
internal/i18n                     multi-language catalogs (en, vi, zh, ja, es) and translation
internal/link                     symlink with junction/hard link/copy fallback
internal/profileenv                profile directories, environment variables
internal/shellrc                    shell detection + PATH/completion snippet
internal/ui                          colors, terminal output, Windows console setup, CJK width
internal/update                    version checking, GitHub releases API, in-place self-update
```

Every package other than `internal/cli` takes an explicit `*app.App` instead
of reading a global — see each package's doc comment for what it owns.
Command handlers return `error`; `internal/cli.Run` is the single place that
turns that into a process exit code.

No external module dependencies — standard library only. Run `make check`
(build + vet + gofmt + test + lint) before sending a PR — the same steps run
in `.github/workflows/ci.yml` on every push/PR to `main`. Releases are built
automatically via `.github/workflows/release.yml`: pushing a `vX.Y.Z` tag
builds and publishes binaries for all 6 OS/architecture combinations to a
GitHub Release.

## Contributing

Issues and pull requests are welcome. If you're changing anything in
`internal/link` or `internal/credential`, please test on the actual OS you're
targeting — the Windows-specific behavior in this project (junction
fallback, safe removal, self-update) was all verified experimentally, not
assumed from documentation alone.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.

## License

[MIT](LICENSE)
