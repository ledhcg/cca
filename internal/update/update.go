package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ledhcg/cca/internal/app"
	"github.com/ledhcg/cca/internal/config"
)

// DefaultCooldown is 24 hours between automated background checks.
const DefaultCooldown = 24 * time.Hour

// GitHubRepo is the upstream repository identifier.
const GitHubRepo = "ledhcg/cca"

// GitHubAPIBaseURL is overridable in tests with a mock server.
var GitHubAPIBaseURL = "https://api.github.com"

// ReleaseInfo contains metadata about a published GitHub release.
type ReleaseInfo struct {
	Version     string    `json:"version"`
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
	Changelog   string    `json:"changelog"`
	AssetName   string    `json:"asset_name"`
	AssetURL    string    `json:"asset_url"`
	AssetSize   int64     `json:"asset_size"`
}

type githubRelease struct {
	TagName     string        `json:"tag_name"`
	HTMLURL     string        `json:"html_url"`
	PublishedAt time.Time     `json:"published_at"`
	Body        string        `json:"body"`
	Assets      []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// CurrentAssetName returns the expected release binary filename for the current OS/Arch.
func CurrentAssetName() string {
	name := fmt.Sprintf("cca-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

// CompareVersions compares two semver versions (e.g. "v1.0.0" vs "v1.1.0").
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
// "dev" is treated as older than any published release (-1).
func CompareVersions(v1, v2 string) int {
	v1 = strings.TrimSpace(v1)
	v2 = strings.TrimSpace(v2)

	if v1 == v2 {
		return 0
	}
	if v1 == "dev" {
		return -1
	}
	if v2 == "dev" {
		return 1
	}

	p1 := parseSemver(v1)
	p2 := parseSemver(v2)

	for i := 0; i < 3; i++ {
		if p1.nums[i] < p2.nums[i] {
			return -1
		}
		if p1.nums[i] > p2.nums[i] {
			return 1
		}
	}

	// Compare prerelease: a version without prerelease has higher precedence than with prerelease
	// e.g. 1.0.0 > 1.0.0-rc1
	if p1.prerelease == "" && p2.prerelease != "" {
		return 1
	}
	if p1.prerelease != "" && p2.prerelease == "" {
		return -1
	}
	if p1.prerelease < p2.prerelease {
		return -1
	}
	if p1.prerelease > p2.prerelease {
		return 1
	}

	return 0
}

type parsedVersion struct {
	nums       [3]int
	prerelease string
}

func parseSemver(s string) parsedVersion {
	s = strings.TrimPrefix(s, "v")
	s = strings.TrimPrefix(s, "V")

	// Strip build metadata after '+'
	if idx := strings.IndexByte(s, '+'); idx != -1 {
		s = s[:idx]
	}

	var prerelease string
	if idx := strings.IndexByte(s, '-'); idx != -1 {
		prerelease = s[idx+1:]
		s = s[:idx]
	}

	parts := strings.Split(s, ".")
	var nums [3]int
	for i := 0; i < len(parts) && i < 3; i++ {
		n, _ := strconv.Atoi(parts[i])
		nums[i] = n
	}

	return parsedVersion{nums: nums, prerelease: prerelease}
}

// CheckLatest checks the GitHub Releases API for the latest release.
func CheckLatest(currentVersion string) (*ReleaseInfo, bool, error) {
	return CheckLatestWithTimeout(currentVersion, 10*time.Second)
}

// CheckLatestWithTimeout checks the GitHub Releases API with a specific timeout.
func CheckLatestWithTimeout(currentVersion string, timeout time.Duration) (*ReleaseInfo, bool, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", GitHubAPIBaseURL, GitHubRepo)
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", fmt.Sprintf("cca/%s (%s; %s)", currentVersion, runtime.GOOS, runtime.GOARCH))

	resp, err := client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, false, fmt.Errorf("no releases found for %s", GitHubRepo)
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, false, fmt.Errorf("GitHub API rate limit exceeded")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var ghRel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&ghRel); err != nil {
		return nil, false, fmt.Errorf("invalid release JSON: %w", err)
	}

	expectedAsset := CurrentAssetName()
	var matchedAsset *githubAsset
	for _, a := range ghRel.Assets {
		if a.Name == expectedAsset {
			matchedAsset = &a
			break
		}
	}

	rel := &ReleaseInfo{
		Version:     ghRel.TagName,
		TagName:     ghRel.TagName,
		PublishedAt: ghRel.PublishedAt,
		HTMLURL:     ghRel.HTMLURL,
		Changelog:   ghRel.Body,
		AssetName:   expectedAsset,
	}

	if matchedAsset != nil {
		rel.AssetURL = matchedAsset.BrowserDownloadURL
		rel.AssetSize = matchedAsset.Size
	}

	hasUpdate := CompareVersions(currentVersion, rel.Version) < 0
	return rel, hasUpdate, nil
}

// ShouldCheck returns whether an automated update check should be performed based on cooldown.
func ShouldCheck(lastCheck int64, now time.Time) bool {
	if lastCheck == 0 {
		return true
	}
	last := time.Unix(lastCheck, 0)
	if now.Before(last) { // Clock skew
		return true
	}
	return now.Sub(last) >= DefaultCooldown
}

// TriggerBackgroundCheck triggers a non-blocking asynchronous update check if cooldown has passed.
func TriggerBackgroundCheck(a *app.App, currentVersion string) {
	if currentVersion == "dev" || os.Getenv("CCA_NO_UPDATE_CHECK") != "" || os.Getenv("CI") != "" {
		return
	}

	cfg, err := config.Load(a)
	if err != nil {
		return
	}

	if !ShouldCheck(cfg.LastUpdateCheck, time.Now()) {
		return
	}

	go func() {
		rel, hasUpdate, err := CheckLatestWithTimeout(currentVersion, 5*time.Second)
		now := time.Now().Unix()

		cfg, loadErr := config.Load(a)
		if loadErr != nil {
			return
		}

		cfg.LastUpdateCheck = now
		if err == nil {
			if hasUpdate {
				cfg.LatestVersion = rel.Version
			} else {
				cfg.LatestVersion = ""
			}
		}
		_ = config.Save(a, cfg)
	}()
}

// Apply downloads the new release asset and replaces the running binary.
func Apply(a *app.App, release *ReleaseInfo) error {
	if release == nil || release.AssetURL == "" {
		return fmt.Errorf("no download URL available for release %s (asset %s)", release.TagName, CurrentAssetName())
	}

	// Determine primary binary destination
	destPath, err := os.Executable()
	if err != nil {
		destPath = canonicalPath(a)
	} else if resolved, err := filepath.EvalSymlinks(destPath); err == nil {
		destPath = resolved
	}

	canonPath := canonicalPath(a)

	// Download new binary to a temporary file in the same directory as destPath
	tmpDir := filepath.Dir(destPath)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		tmpDir = os.TempDir()
	}

	tmpFile, err := os.CreateTemp(tmpDir, ".cca-update-*")
	if err != nil {
		return fmt.Errorf("could not create temporary update file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Get(release.AssetURL)
	if err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_ = tmpFile.Close()
		return fmt.Errorf("download failed with HTTP %s", resp.Status)
	}

	written, err := io.Copy(tmpFile, resp.Body)
	_ = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to write update: %w", err)
	}
	if release.AssetSize > 0 && written != release.AssetSize {
		return fmt.Errorf("corrupted download: got %d bytes, expected %d", written, release.AssetSize)
	}

	if err := os.Chmod(tmpPath, 0755); err != nil {
		return fmt.Errorf("could not set executable permissions: %w", err)
	}

	// In-place replacement
	if err := replaceExecutable(tmpPath, destPath); err != nil {
		return err
	}

	// If the current executable is outside ~/.claude/bin, but ~/.claude/bin/cca exists, update it too
	if destPath != canonPath {
		if _, err := os.Stat(canonPath); err == nil {
			_ = replaceExecutable(destPath, canonPath)
		}
	}

	return nil
}

func canonicalPath(a *app.App) string {
	name := "cca"
	if runtime.GOOS == "windows" {
		name = "cca.exe"
	}
	return filepath.Join(a.Main, "bin", name)
}

func replaceExecutable(src, dst string) error {
	_ = os.MkdirAll(filepath.Dir(dst), 0755)

	if runtime.GOOS == "windows" {
		// Windows: rename running file aside to .old, then move new file into place
		if _, err := os.Stat(dst); err == nil {
			old := dst + ".old"
			_ = os.Remove(old)
			if err := os.Rename(dst, old); err != nil {
				return fmt.Errorf("could not move existing binary aside: %w", err)
			}
		}
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0755)
	}

	// Unix: atomic rename
	return os.Rename(src, dst)
}
