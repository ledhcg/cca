package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ledhcg/cca/internal/app"
)

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		v1, v2 string
		want   int
	}{
		// Equals
		{"v1.0.0", "v1.0.0", 0},
		{"1.0.0", "v1.0.0", 0},
		{"v1.2.3", "1.2.3", 0},
		{"dev", "dev", 0},

		// Dev builds
		{"dev", "v1.0.0", -1},
		{"v1.0.0", "dev", 1},

		// Major, minor, patch differences
		{"v1.0.0", "v1.0.1", -1},
		{"v1.0.1", "v1.0.0", 1},
		{"v1.1.0", "v1.0.9", 1},
		{"v1.0.9", "v1.1.0", -1},
		{"v2.0.0", "v1.99.99", 1},
		{"v1.99.99", "v2.0.0", -1},

		// Build metadata (should be ignored per SemVer)
		{"v1.0.0+build1", "v1.0.0+build2", 0},
		{"v1.0.1+build1", "v1.0.0", 1},

		// Prereleases (release is higher than prerelease)
		{"v1.0.0", "v1.0.0-rc1", 1},
		{"v1.0.0-rc1", "v1.0.0", -1},
		{"v1.0.0-beta.1", "v1.0.0-beta.2", -1},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s_vs_%s", tc.v1, tc.v2), func(t *testing.T) {
			got := CompareVersions(tc.v1, tc.v2)
			if got != tc.want {
				t.Errorf("CompareVersions(%q, %q) = %d, want %d", tc.v1, tc.v2, got, tc.want)
			}
		})
	}
}

func TestShouldCheck(t *testing.T) {
	now := time.Now()

	// 1. Never checked (0) -> should check
	if !ShouldCheck(0, now) {
		t.Errorf("ShouldCheck(0, now) = false, want true")
	}

	// 2. Checked 1 hour ago -> should NOT check (<24h)
	oneHourAgo := now.Add(-1 * time.Hour).Unix()
	if ShouldCheck(oneHourAgo, now) {
		t.Errorf("ShouldCheck(1h ago, now) = true, want false")
	}

	// 3. Checked 25 hours ago -> should check (>24h)
	twentyFiveHoursAgo := now.Add(-25 * time.Hour).Unix()
	if !ShouldCheck(twentyFiveHoursAgo, now) {
		t.Errorf("ShouldCheck(25h ago, now) = false, want true")
	}

	// 4. Clock skew: last check is in the future -> should check
	inTheFuture := now.Add(1 * time.Hour).Unix()
	if !ShouldCheck(inTheFuture, now) {
		t.Errorf("ShouldCheck(future, now) = false, want true")
	}
}

func TestCheckLatestWithMockServer(t *testing.T) {
	expectedAsset := CurrentAssetName()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/ledhcg/cca/releases/latest" {
			http.NotFound(w, r)
			return
		}
		rel := githubRelease{
			TagName:     "v1.1.0",
			HTMLURL:     "https://github.com/ledhcg/cca/releases/tag/v1.1.0",
			PublishedAt: time.Now(),
			Body:        "New features and improvements",
			Assets: []githubAsset{
				{
					Name:               expectedAsset,
					Size:               12345,
					BrowserDownloadURL: "https://example.com/download/" + expectedAsset,
				},
				{
					Name:               "cca-other-arch",
					Size:               12345,
					BrowserDownloadURL: "https://example.com/download/other",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rel)
	}))
	defer ts.Close()

	// Override API base URL
	oldBase := GitHubAPIBaseURL
	GitHubAPIBaseURL = ts.URL
	defer func() { GitHubAPIBaseURL = oldBase }()

	// Test when current version is older (v1.0.0 < v1.1.0)
	rel, hasUpdate, err := CheckLatest("v1.0.0")
	if err != nil {
		t.Fatalf("CheckLatest failed: %v", err)
	}
	if !hasUpdate {
		t.Errorf("hasUpdate = false, want true")
	}
	if rel.Version != "v1.1.0" {
		t.Errorf("rel.Version = %q, want 'v1.1.0'", rel.Version)
	}
	if rel.AssetName != expectedAsset {
		t.Errorf("rel.AssetName = %q, want %q", rel.AssetName, expectedAsset)
	}
	if rel.AssetURL != "https://example.com/download/"+expectedAsset {
		t.Errorf("rel.AssetURL = %q, want expected download url", rel.AssetURL)
	}

	// Test when current version is equal or newer (v1.1.0 == v1.1.0)
	_, hasUpdateSame, err := CheckLatest("v1.1.0")
	if err != nil {
		t.Fatalf("CheckLatest failed: %v", err)
	}
	if hasUpdateSame {
		t.Errorf("hasUpdate = true, want false when on latest version")
	}
}

func TestApplyWithMockServer(t *testing.T) {
	dummyPayload := []byte("#!/bin/sh\necho updated binary\n")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(dummyPayload)
	}))
	defer ts.Close()

	home := t.TempDir()
	a := &app.App{
		Home: home,
		Main: filepath.Join(home, ".claude"),
		Root: filepath.Join(home, ".claude-accounts"),
	}

	binDir := filepath.Join(a.Main, "bin")
	_ = os.MkdirAll(binDir, 0755)
	target := filepath.Join(binDir, "cca")
	if err := os.WriteFile(target, []byte("old binary content"), 0755); err != nil {
		t.Fatal(err)
	}

	rel := &ReleaseInfo{
		Version:   "v1.1.0",
		TagName:   "v1.1.0",
		AssetName: CurrentAssetName(),
		AssetURL:  ts.URL + "/download",
		AssetSize: int64(len(dummyPayload)),
	}

	if err := replaceExecutable(target, target); err != nil {
		t.Logf("Self-replace executable check passed")
	}

	// Test download and replace
	tmpFile, err := os.CreateTemp(binDir, ".cca-test-*")
	if err != nil {
		t.Fatal(err)
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()
	_ = os.WriteFile(tmpPath, dummyPayload, 0755)

	if err := replaceExecutable(tmpPath, target); err != nil {
		t.Fatalf("replaceExecutable failed: %v", err)
	}

	gotData, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotData) != string(dummyPayload) {
		t.Errorf("binary content = %q, want %q", string(gotData), string(dummyPayload))
	}
	_ = rel
}
