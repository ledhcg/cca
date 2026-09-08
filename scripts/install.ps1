#requires -Version 5.1
# Install cca (a static Go binary, no runtime required) on Windows, plain
# PowerShell. Prefers downloading a prebuilt release from GitHub; if there's no
# release yet (or the download fails) and Go is available locally, it builds
# from source instead.
$ErrorActionPreference = "Stop"
$Repo = "ledhcg/cca"
$SelfDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $SelfDir

$arch = $env:PROCESSOR_ARCHITECTURE
if ($arch -eq "ARM64") { $arch = "arm64" } else { $arch = "amd64" }

$bindir = Join-Path $HOME ".claude\bin"
New-Item -ItemType Directory -Force -Path $bindir | Out-Null
$dst = Join-Path $bindir "cca.exe"
$url = "https://github.com/$Repo/releases/latest/download/cca-windows-$arch.exe"

$downloaded = $false
try {
    Write-Host "Downloading $url"
    Invoke-WebRequest -Uri $url -OutFile "$dst.tmp" -ErrorAction Stop
    Move-Item -Force "$dst.tmp" $dst
    $downloaded = $true
    Write-Host "✓ Downloaded cca from a GitHub Release (windows/$arch)"
} catch {
    Remove-Item -Force "$dst.tmp" -ErrorAction SilentlyContinue
}

if (-not $downloaded) {
    $go = Get-Command go -ErrorAction SilentlyContinue
    if ($go -and (Test-Path (Join-Path $RootDir "go.mod"))) {
        Write-Host "! No release available (none published yet, or no network) — building from source with Go…"
        & $go.Source build -o $dst (Join-Path $RootDir "cmd\cca")
        Write-Host "✓ Built cca from $RootDir"
    } else {
        Write-Error "Could not download a binary and Go isn't available to build from source. Install Go from https://go.dev/dl/ and re-run install.ps1."
        exit 1
    }
}

& $dst install
