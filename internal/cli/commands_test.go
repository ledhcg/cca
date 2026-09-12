package cli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ledhcg/cca/internal/app"
	"github.com/ledhcg/cca/internal/config"
	"github.com/ledhcg/cca/internal/i18n"
	"github.com/ledhcg/cca/internal/update"
)

func TestExpandYolo(t *testing.T) {
	cases := []struct {
		name    string
		in      []string
		wantOut []string
		wantHit bool
	}{
		{
			name:    "replaces yolo",
			in:      []string{"--resume", "--yolo"},
			wantOut: []string{"--resume", "--dangerously-skip-permissions"},
			wantHit: true,
		},
		{
			name:    "no yolo present",
			in:      []string{"--resume", "-p", "hi"},
			wantOut: []string{"--resume", "-p", "hi"},
			wantHit: false,
		},
		{
			name:    "multiple yolo occurrences",
			in:      []string{"--yolo", "--yolo"},
			wantOut: []string{"--dangerously-skip-permissions", "--dangerously-skip-permissions"},
			wantHit: true,
		},
		{
			name:    "empty input",
			in:      []string{},
			wantOut: []string{},
			wantHit: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, hit := expandYolo(c.in)
			if !reflect.DeepEqual(out, c.wantOut) {
				t.Errorf("out = %#v, want %#v", out, c.wantOut)
			}
			if hit != c.wantHit {
				t.Errorf("hit = %v, want %v", hit, c.wantHit)
			}
		})
	}
}

func testApp(t *testing.T) *app.App {
	t.Helper()
	home := t.TempDir()
	return &app.App{
		Home:       home,
		Main:       filepath.Join(home, ".claude"),
		Root:       filepath.Join(home, ".claude-accounts"),
		ConfigPath: filepath.Join(home, ".claude-accounts", "config.json"),
	}
}

func TestCmdSettingsOverview(t *testing.T) {
	a := testApp(t)
	err := cmdSettings(a, parsedArgs{})
	if err != nil {
		t.Fatalf("cmdSettings overview failed: %v", err)
	}
}

func TestCmdSettingsSetLanguage(t *testing.T) {
	a := testApp(t)

	// Set to vi
	err := cmdSettings(a, parsedArgs{positional: []string{"lang", "vi"}})
	if err != nil {
		t.Fatalf("cmdSettings lang vi failed: %v", err)
	}
	cfg, err := config.Load(a)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Lang != "vi" {
		t.Errorf("expected cfg.Lang == 'vi', got %q", cfg.Lang)
	}
	if i18n.Current() != "vi" {
		t.Errorf("expected i18n.Current() == 'vi', got %q", i18n.Current())
	}

	// Reset to auto
	err = cmdSettings(a, parsedArgs{positional: []string{"lang", "auto"}})
	if err != nil {
		t.Fatalf("cmdSettings lang auto failed: %v", err)
	}
	cfg, err = config.Load(a)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Lang != "" {
		t.Errorf("expected cfg.Lang == '', got %q", cfg.Lang)
	}

	// Invalid language
	err = cmdSettings(a, parsedArgs{positional: []string{"lang", "invalid-lang"}})
	if err == nil {
		t.Errorf("expected error for invalid language, got nil")
	}
}

func TestCmdVersion(t *testing.T) {
	a := testApp(t)
	if err := cmdVersion(a); err != nil {
		t.Fatalf("cmdVersion failed: %v", err)
	}
}

func TestCmdUpdateCheckWithMock(t *testing.T) {
	a := testApp(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"tag_name":"v1.0.0","assets":[]}`)
	}))
	defer ts.Close()

	oldBase := update.GitHubAPIBaseURL
	update.GitHubAPIBaseURL = ts.URL
	defer func() { update.GitHubAPIBaseURL = oldBase }()

	err := cmdUpdate(a, parsedArgs{bools: map[string]bool{"--check": true}})
	if err != nil {
		t.Fatalf("cmdUpdate --check failed: %v", err)
	}
}
