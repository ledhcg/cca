# Changelog

All notable changes to this project are documented in this file.

## [1.1.0] — 2026-09-13

### Added
- **Multi-language support (i18n)**: 5 popular languages out of the box —
  English (`en`, default), Vietnamese (`vi`), Simplified Chinese (`zh`),
  Japanese (`ja`), and Spanish (`es`), plus system auto-detection (`auto`).
  All command outputs, guides, table headers, and error messages are fully
  localized with zero external dependencies.
- **Dedicated `cca settings` command**: Interactive dashboard and language
  picker for configuring user preferences and display language (`cca settings lang`).
  Cleanly separated from `cca config` (which opens `~/.claude-accounts` in VS Code).
- **In-place self-update (`cca update`)**: Download and replace the current binary
  directly from GitHub Releases (`ledhcg/cca`) across macOS, Linux, and Windows
  with safe in-place replacement mechanics (using the `.old` executable rename pattern on Windows).
  Supports `cca update --check` (or `-c`) to inspect for updates without downloading.
- **Non-blocking background update checks**: Periodically checks for new releases
  every 24 hours in a background goroutine without slowing down CLI execution.
  Notifies users via a subtle, non-intrusive alert on `stderr` when an update is available.
- **Enhanced `cca version`**: Displays current version and highlights available updates
  from the local cache.
- **Rune- and CJK-aware terminal formatting**: `ui.StringWidth` and `ui.PadRight`
  ensure visual alignment across Vietnamese accented text and full-width CJK (Chinese, Japanese)
  characters in tables and menus.
- **Session transfer & handoff (`cca handoff` & `cca session`)**:
  - `cca handoff <to-profile> [--from <profile>] [--id <sessionId>] [--fork] [--yolo]`:
    Seamlessly hand off an active session from one account to another (e.g. when quota runs out)
    and resume work immediately without losing conversation history.
  - `cca session ls [--from <profile>] [--all]`:
    Inspect recent sessions, titles, message counts, sizes, and modification dates.
  - `cca session cp --from <profile> --to <profile> [--id <sessionId> | --all]`:
    Copy sessions across profile environments.
  - `cca session mv --from <profile> --to <profile> [--id <sessionId> | --all]`:
    Move sessions between profiles, pruning them from the source profile.
  - `cca session rm [--from <profile>] (--id <sessionId> | --all) [-y]`:
    Remove session transcripts, subagents, rollback snapshots, and prompt history.

## [1.0.0] — 2026-09-08

Initial public release.

### Added
- Multi-account profile management for Claude Code: `new`, `ls`, `sh`,
  `exec`, `info`, `rm`, `default`, `config`, `doctor`, `sync`, `install`,
  `version`, `guide`/`help`.
- Per-OS credential handling: macOS Keychain (`internal/credential`), and a
  plain `.credentials.json` file on Linux/Windows — verified by reading the
  `claude` binary directly and cross-checked against the official docs.
- Shared-files system (`~/.claude-accounts/config.json`): `sharedLinks`,
  `sharedCopies`, and a `merge` sync strategy that preserves
  `profileLocalKeys` (e.g. `model`) per profile.
- `--yolo` flag (alias for `--dangerously-skip-permissions`, scoped to `cca`
  only) and `cca new <name> --login` for a create-and-log-in-immediately flow.
- `cca install`: adds `cca` to `PATH` plus shell completion for
  bash/zsh/fish/PowerShell, and — on Windows — writes `PATH` directly to
  `HKCU\Environment` so `cmd.exe` and freshly opened shells pick it up too.
  Idempotent; safe to re-run.
- Cross-platform symlink fallback (`internal/link`): real symlink first,
  falling back to a junction (directories) or hard link (files) on Windows
  when Developer Mode is off — no admin rights required.
- Rewrite from a shell script into a single static Go binary
  (`cmd/cca` + `internal/*`), so `cca` runs with no Python/Node or other
  runtime dependency on macOS, Linux, and Windows.
- Install scripts (`scripts/install.sh`, `scripts/install.ps1`) that detect
  OS/arch, pull the matching binary from the latest GitHub Release, and fall
  back to `go build` from source when no release is available or there's no
  network.
- CI (`build`, `vet`, `gofmt`, `test -race -cover`, `golangci-lint`) and an
  automated release pipeline that builds and publishes binaries for all 6
  OS/architecture combinations on every `vX.Y.Z` tag.

### Fixed
- `cca exec <name> -- <cmd>` never actually running the given command.

