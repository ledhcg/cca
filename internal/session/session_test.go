package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestProjectSlug(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"/Users/dinhcuong/Workspace/cca", "-Users-dinhcuong-Workspace-cca"},
		{"/Users/dinhcuong/.claude", "-Users-dinhcuong--claude"},
		{"/Users/dinhcuong/.codex", "-Users-dinhcuong--codex"},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := ProjectSlug(tc.in)
			if got != tc.want {
				t.Errorf("ProjectSlug(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestParseSessionFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test-session-123.jsonl")

	// Create JSONL with multiple lines, including a very long line (>100KB)
	longText := strings.Repeat("A", 120*1024)
	content := fmt.Sprintf(`{"type":"mode","mode":"normal"}
{"type":"last-prompt","lastPrompt":"Last prompt text"}
{"type":"ai-title","aiTitle":"AI Generated Title"}
{"type":"user","message":{"role":"user","content":"Initial prompt"}}
{"type":"assistant","message":{"role":"assistant","content":"%s"}}
{"type":"user","message":{"role":"user","content":"Second prompt"}}
`, longText)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	info, err := parseSessionFile(filePath)
	if err != nil {
		t.Fatalf("parseSessionFile failed: %v", err)
	}

	if info.ID != "test-session-123" {
		t.Errorf("ID = %q, want 'test-session-123'", info.ID)
	}
	if info.Title != "AI Generated Title" {
		t.Errorf("Title = %q, want 'AI Generated Title'", info.Title)
	}
	if info.MessageCount != 3 { // 2 user + 1 assistant
		t.Errorf("MessageCount = %d, want 3", info.MessageCount)
	}
}

func TestListAndLatestSessions(t *testing.T) {
	profileDir := t.TempDir()
	projectDir := "/Users/test/my-project"
	slug := ProjectSlug(projectDir)

	slugDir := filepath.Join(profileDir, "projects", slug)
	_ = os.MkdirAll(slugDir, 0755)

	// Create session 1 (older)
	s1 := filepath.Join(slugDir, "sess-1.jsonl")
	_ = os.WriteFile(s1, []byte(`{"type":"ai-title","aiTitle":"Session 1"}`), 0644)
	t1 := time.Now().Add(-1 * time.Hour)
	_ = os.Chtimes(s1, t1, t1)

	// Create session 2 (newer)
	s2 := filepath.Join(slugDir, "sess-2.jsonl")
	_ = os.WriteFile(s2, []byte(`{"type":"ai-title","aiTitle":"Session 2"}`), 0644)
	t2 := time.Now()
	_ = os.Chtimes(s2, t2, t2)

	// Test ListSessions
	sessions, err := ListSessions(profileDir, projectDir, false)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
	if sessions[0].ID != "sess-2" {
		t.Errorf("expected newest session first ('sess-2'), got %q", sessions[0].ID)
	}

	// Test LatestSession
	latest, err := LatestSession(profileDir, projectDir)
	if err != nil {
		t.Fatalf("LatestSession failed: %v", err)
	}
	if latest == nil || latest.ID != "sess-2" {
		t.Errorf("expected latest to be 'sess-2', got %v", latest)
	}

	// Test ResolveSessionID (exact and prefix)
	resolved, err := ResolveSessionID(profileDir, projectDir, "sess-1")
	if err != nil || resolved != "sess-1" {
		t.Errorf("ResolveSessionID exact failed: %v, got %q", err, resolved)
	}
	resolvedPrefix, err := ResolveSessionID(profileDir, projectDir, "sess-2")
	if err != nil || resolvedPrefix != "sess-2" {
		t.Errorf("ResolveSessionID prefix failed: %v, got %q", err, resolvedPrefix)
	}
}

func TestTransferSessionCopyAndMove(t *testing.T) {
	fromProfile := t.TempDir()
	toProfile := t.TempDir()
	projectDir := "/Users/test/app"
	slug := ProjectSlug(projectDir)
	sid := "session-abc-123"

	// Setup fromProfile files
	fromProjectDir := filepath.Join(fromProfile, "projects", slug)
	_ = os.MkdirAll(fromProjectDir, 0755)
	fromTranscript := filepath.Join(fromProjectDir, sid+".jsonl")
	_ = os.WriteFile(fromTranscript, []byte(`{"type":"ai-title","aiTitle":"Transfer Test"}`), 0644)

	// Subdir (subagents)
	fromSubdir := filepath.Join(fromProjectDir, sid, "subagents")
	_ = os.MkdirAll(fromSubdir, 0755)
	_ = os.WriteFile(filepath.Join(fromSubdir, "agent-1.jsonl"), []byte(`{"agent":1}`), 0644)

	// File-history
	fromFileHistory := filepath.Join(fromProfile, "file-history", sid)
	_ = os.MkdirAll(fromFileHistory, 0755)
	_ = os.WriteFile(filepath.Join(fromFileHistory, "hash@v1"), []byte("backup"), 0644)

	// Session-env
	fromEnv := filepath.Join(fromProfile, "session-env", sid)
	_ = os.MkdirAll(fromEnv, 0755)
	_ = os.WriteFile(filepath.Join(fromEnv, "env.sh"), []byte("export FOO=1"), 0644)

	// History.jsonl
	fromHistory := filepath.Join(fromProfile, "history.jsonl")
	historyContent := fmt.Sprintf(`{"display":"other prompt","sessionId":"other-sess"}
{"display":"my prompt","sessionId":"%s"}
`, sid)
	_ = os.WriteFile(fromHistory, []byte(historyContent), 0644)

	// 1. Test COPY (move=false)
	err := TransferSession(fromProfile, toProfile, projectDir, sid, false)
	if err != nil {
		t.Fatalf("TransferSession copy failed: %v", err)
	}

	// Verify target has all files
	toTranscript := filepath.Join(toProfile, "projects", slug, sid+".jsonl")
	if _, err := os.Stat(toTranscript); err != nil {
		t.Errorf("target transcript missing: %v", err)
	}
	toSubagent := filepath.Join(toProfile, "projects", slug, sid, "subagents", "agent-1.jsonl")
	if _, err := os.Stat(toSubagent); err != nil {
		t.Errorf("target subagent missing: %v", err)
	}
	toFileHistory := filepath.Join(toProfile, "file-history", sid, "hash@v1")
	if _, err := os.Stat(toFileHistory); err != nil {
		t.Errorf("target file history missing: %v", err)
	}
	toEnv := filepath.Join(toProfile, "session-env", sid, "env.sh")
	if _, err := os.Stat(toEnv); err != nil {
		t.Errorf("target session env missing: %v", err)
	}
	toHistoryData, _ := os.ReadFile(filepath.Join(toProfile, "history.jsonl"))
	if !strings.Contains(string(toHistoryData), sid) {
		t.Errorf("target history.jsonl missing session lines")
	}

	// Verify source STILL has files after copy
	if _, err := os.Stat(fromTranscript); err != nil {
		t.Errorf("source transcript was removed during copy")
	}

	// 2. Test MOVE (move=true) to a third profile
	thirdProfile := t.TempDir()
	err = TransferSession(fromProfile, thirdProfile, projectDir, sid, true)
	if err != nil {
		t.Fatalf("TransferSession move failed: %v", err)
	}

	// Verify source transcript is removed
	if _, err := os.Stat(fromTranscript); !os.IsNotExist(err) {
		t.Errorf("source transcript should have been removed after move")
	}

	// Verify source history lines were filtered
	fromHistoryData, _ := os.ReadFile(fromHistory)
	if strings.Contains(string(fromHistoryData), sid) {
		t.Errorf("source history.jsonl still contains session lines after move")
	}
}

func TestRemoveSession(t *testing.T) {
	profileDir := t.TempDir()
	projectDir := "/Users/test/rm-project"
	slug := ProjectSlug(projectDir)
	sid := "sess-to-remove"

	slugDir := filepath.Join(profileDir, "projects", slug)
	_ = os.MkdirAll(slugDir, 0755)
	transcript := filepath.Join(slugDir, sid+".jsonl")
	_ = os.WriteFile(transcript, []byte(`{"type":"ai-title","aiTitle":"Delete me"}`), 0644)

	history := filepath.Join(profileDir, "history.jsonl")
	_ = os.WriteFile(history, []byte(fmt.Sprintf(`{"sessionId":"%s"}\n{"sessionId":"keep-me"}`, sid)), 0644)

	err := RemoveSession(profileDir, projectDir, sid)
	if err != nil {
		t.Fatalf("RemoveSession failed: %v", err)
	}

	if _, err := os.Stat(transcript); !os.IsNotExist(err) {
		t.Errorf("transcript was not deleted")
	}
	histData, _ := os.ReadFile(history)
	if strings.Contains(string(histData), sid) {
		t.Errorf("history was not cleaned")
	}
}
