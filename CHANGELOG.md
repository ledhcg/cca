# Changelog

All notable changes to this project are documented in this file.

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

