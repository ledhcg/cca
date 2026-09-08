#!/bin/sh
# Install cca (a static Go binary, no runtime required) on macOS / Linux / Git
# Bash (Windows). Prefers downloading a prebuilt release from GitHub; if there's
# no release yet (or the download fails) and Go is available locally, it builds
# from source instead.
set -eu
REPO="ledhcg/cca"
SELFDIR=$(cd "$(dirname "$0")" && pwd)
ROOTDIR=$(cd "$SELFDIR/.." && pwd)

os=""
case "$(uname -s)" in
  Darwin) os=darwin ;;
  Linux) os=linux ;;
  MINGW*|MSYS*|CYGWIN*) os=windows ;;
  *) echo "✗ Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

arch=""
case "$(uname -m)" in
  x86_64|amd64) arch=amd64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) echo "✗ Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

ext=""
[ "$os" = "windows" ] && ext=".exe"

bindir="$HOME/.claude/bin"
mkdir -p "$bindir"
dst="$bindir/cca$ext"
url="https://github.com/$REPO/releases/latest/download/cca-$os-$arch$ext"

download() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$dst.tmp" 2>/dev/null
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$url" -O "$dst.tmp" 2>/dev/null
  else
    return 1
  fi
}

if download; then
  mv "$dst.tmp" "$dst"
  chmod +x "$dst"
  echo "✓ Downloaded cca from a GitHub Release ($os/$arch)"
elif command -v go >/dev/null 2>&1 && [ -f "$ROOTDIR/go.mod" ]; then
  echo "! No release available (none published yet, or no network) — building from source with Go…"
  rm -f "$dst.tmp"
  go build -o "$dst" "$ROOTDIR/cmd/cca"
  chmod +x "$dst"
  echo "✓ Built cca from $ROOTDIR"
else
  rm -f "$dst.tmp"
  echo "✗ Could not download a binary and Go isn't available to build from source." >&2
  echo "  Install Go from https://go.dev/dl/ and re-run install.sh." >&2
  exit 1
fi

"$dst" install
