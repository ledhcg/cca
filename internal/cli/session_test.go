package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ledhcg/cca/internal/session"
)

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
	}

	for _, tc := range cases {
		if got := formatBytes(tc.in); got != tc.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRelativeTime(t *testing.T) {
	now := time.Now()
	cases := []struct {
		in   time.Time
		want string
	}{
		{time.Time{}, "—"},
		{now.Add(-20 * time.Second), "just now"},
		{now.Add(-5 * time.Minute), "5m ago"},
		{now.Add(-3 * time.Hour), "3h ago"},
		{now.Add(-48 * time.Hour), "2d ago"},
	}

	for _, tc := range cases {
		if got := relativeTime(tc.in); got != tc.want {
			t.Errorf("relativeTime(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestCmdSessionUsage(t *testing.T) {
	a := testApp(t)
	if err := cmdSession(a, nil); err != nil {
		t.Fatalf("cmdSession(nil) failed: %v", err)
	}
	if err := cmdSession(a, []string{"help"}); err != nil {
		t.Fatalf("cmdSession(help) failed: %v", err)
	}
}

func TestCmdSessionLs(t *testing.T) {
	a := testApp(t)

	// Ls on empty profile should succeed without error
	if err := cmdSession(a, []string{"ls"}); err != nil {
		t.Fatalf("cmdSession ls empty failed: %v", err)
	}
	if err := cmdSession(a, []string{"ls", "--all"}); err != nil {
		t.Fatalf("cmdSession ls --all empty failed: %v", err)
	}
}

func TestCmdSessionCpAndMvValidation(t *testing.T) {
	a := testApp(t)

	// Create profile "work" so target profile exists
	workDir := filepath.Join(a.Root, "work")
	_ = os.MkdirAll(workDir, 0755)

	// Missing --from and --to
	err := cmdSession(a, []string{"cp"})
	if err == nil || !strings.Contains(err.Error(), "--from and --to") {
		t.Errorf("expected missing --from and --to error, got %v", err)
	}

	// Same profile error
	err = cmdSession(a, []string{"cp", "--from", "default", "--to", "default", "--all"})
	if err == nil || !strings.Contains(err.Error(), "cannot be the same") {
		t.Errorf("expected same profile error, got %v", err)
	}

	// Neither --id nor --all
	err = cmdSession(a, []string{"cp", "--from", "default", "--to", "work"})
	if err == nil || !strings.Contains(err.Error(), "must specify either --id") {
		t.Errorf("expected require --id or --all error, got %v", err)
	}

	// Both --id and --all
	err = cmdSession(a, []string{"cp", "--from", "default", "--to", "work", "--id", "123", "--all"})
	if err == nil || !strings.Contains(err.Error(), "must specify either --id") {
		t.Errorf("expected mutual exclusion error, got %v", err)
	}
}

func TestCmdSessionTransferExecution(t *testing.T) {
	a := testApp(t)

	// Create profile "work"
	workDir := filepath.Join(a.Root, "work")
	_ = os.MkdirAll(workDir, 0755)

	// Create a mock session in default profile (a.Main)
	cwd, _ := os.Getwd()
	slug := session.ProjectSlug(cwd)
	defaultProjDir := filepath.Join(a.Main, "projects", slug)
	_ = os.MkdirAll(defaultProjDir, 0755)

	sid := "sess-test-uuid-456"
	transcript := filepath.Join(defaultProjDir, sid+".jsonl")
	_ = os.WriteFile(transcript, []byte(`{"type":"ai-title","aiTitle":"Test Copy Execution"}`), 0644)

	// Test cp --from default --to work --id sess-test-uuid-456
	err := cmdSession(a, []string{"cp", "--from", "default", "--to", "work", "--id", sid})
	if err != nil {
		t.Fatalf("cmdSession cp failed: %v", err)
	}

	// Verify target has the file
	targetFile := filepath.Join(workDir, "projects", slug, sid+".jsonl")
	if _, err := os.Stat(targetFile); err != nil {
		t.Errorf("target file not found after cp: %v", err)
	}

	// Test rm --from work --id sess-test-uuid-456 -y
	err = cmdSession(a, []string{"rm", "--from", "work", "--id", sid, "-y"})
	if err != nil {
		t.Fatalf("cmdSession rm failed: %v", err)
	}
	if _, err := os.Stat(targetFile); !os.IsNotExist(err) {
		t.Errorf("target file still exists after rm")
	}
}

func TestCmdHandoffValidation(t *testing.T) {
	a := testApp(t)

	// Missing target profile
	err := cmdHandoff(a, nil)
	if err == nil || !strings.Contains(err.Error(), "missing target profile") {
		t.Errorf("expected missing target profile error, got %v", err)
	}

	// Same profile error (from default to default)
	err = cmdHandoff(a, []string{"default", "--from", "default"})
	if err == nil || !strings.Contains(err.Error(), "cannot be the same") {
		t.Errorf("expected same profile error, got %v", err)
	}

	// Target profile not found
	err = cmdHandoff(a, []string{"nonexistent-profile"})
	if err == nil || !strings.Contains(err.Error(), "no profile") {
		t.Errorf("expected target profile not found error, got %v", err)
	}
}
