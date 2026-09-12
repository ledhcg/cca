// Package session manages discovering, parsing, transferring, and removing
// Claude Code sessions across different profile environments.
package session

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var slugRe = regexp.MustCompile(`[^a-zA-Z0-9]`)

// ProjectSlug derives Claude Code's project directory name from an absolute path.
// Non-alphanumeric characters are sanitized to '-'.
func ProjectSlug(dir string) string {
	if dir == "" {
		if wd, err := os.Getwd(); err == nil {
			dir = wd
		}
	} else if !filepath.IsAbs(dir) {
		if abs, err := filepath.Abs(dir); err == nil {
			dir = abs
		}
	}
	clean := filepath.Clean(dir)
	slug := slugRe.ReplaceAllString(clean, "-")
	if len(slug) > 200 {
		slug = slug[:200] + "-" + simpleHash36(clean)
	}
	return slug
}

func simpleHash36(s string) string {
	var hash uint32
	for i := 0; i < len(s); i++ {
		hash = hash*31 + uint32(s[i])
	}
	return strconv.FormatUint(uint64(hash), 36)
}

// SessionInfo holds metadata for a Claude Code session.
type SessionInfo struct {
	ID           string    `json:"id"`
	Path         string    `json:"path"`
	Title        string    `json:"title"`
	ModTime      time.Time `json:"modTime"`
	Size         int64     `json:"size"`
	MessageCount int       `json:"messageCount"`
	ProjectSlug  string    `json:"projectSlug,omitempty"`
}

// parseSessionFile reads a session .jsonl file, handling long lines safely.
func parseSessionFile(path string) (SessionInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return SessionInfo{}, err
	}

	sid := strings.TrimSuffix(filepath.Base(path), ".jsonl")
	info := SessionInfo{
		ID:      sid,
		Path:    path,
		ModTime: fi.ModTime(),
		Size:    fi.Size(),
	}

	f, err := os.Open(path)
	if err != nil {
		return info, err
	}
	defer func() { _ = f.Close() }()

	reader := bufio.NewReaderSize(f, 64*1024)

	var customTitle, aiTitle, lastPrompt, firstUserContent string
	msgCount := 0

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			line = bytes.TrimSpace(line)
			if len(line) > 0 {
				var entry struct {
					Type        string `json:"type"`
					CustomTitle string `json:"customTitle"`
					AiTitle     string `json:"aiTitle"`
					LastPrompt  string `json:"lastPrompt"`
					Message     struct {
						Role    string `json:"role"`
						Content any    `json:"content"`
					} `json:"message"`
				}

				if json.Unmarshal(line, &entry) == nil {
					switch entry.Type {
					case "custom-title":
						if entry.CustomTitle != "" {
							customTitle = entry.CustomTitle
						}
					case "ai-title":
						if entry.AiTitle != "" {
							aiTitle = entry.AiTitle
						}
					case "last-prompt":
						if entry.LastPrompt != "" {
							lastPrompt = entry.LastPrompt
						}
					case "user":
						msgCount++
						if firstUserContent == "" && entry.Message.Content != nil {
							firstUserContent = extractText(entry.Message.Content)
						}
					case "assistant", "message":
						msgCount++
					}
				}
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			break
		}
	}

	info.MessageCount = msgCount

	// Title precedence: customTitle > aiTitle > lastPrompt > firstUserContent > (untitled)
	switch {
	case customTitle != "":
		info.Title = customTitle
	case aiTitle != "":
		info.Title = aiTitle
	case lastPrompt != "":
		info.Title = lastPrompt
	case firstUserContent != "":
		info.Title = firstUserContent
	default:
		info.Title = "(untitled)"
	}

	info.Title = cleanOneLine(info.Title, 80)
	return info, nil
}

func extractText(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []any:
		var parts []string
		for _, item := range val {
			if m, ok := item.(map[string]any); ok {
				if text, ok := m["text"].(string); ok && text != "" {
					parts = append(parts, text)
				}
			}
		}
		return strings.Join(parts, " ")
	}
	return ""
}

func cleanOneLine(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.Join(strings.Fields(s), " ")
	if maxLen > 0 && len([]rune(s)) > maxLen {
		runes := []rune(s)
		return string(runes[:maxLen-1]) + "…"
	}
	return s
}

// ListSessions lists sessions in profileDir. If allProjects is false, it lists only sessions
// for projectDir. Results are sorted by ModTime descending.
func ListSessions(profileDir, projectDir string, allProjects bool) ([]SessionInfo, error) {
	projectsRoot := filepath.Join(profileDir, "projects")
	if fi, err := os.Stat(projectsRoot); err != nil || !fi.IsDir() {
		return nil, nil // No projects directory yet
	}

	var targetSlugs []string
	if allProjects {
		entries, err := os.ReadDir(projectsRoot)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
				targetSlugs = append(targetSlugs, e.Name())
			}
		}
	} else {
		slug := ProjectSlug(projectDir)
		targetSlugs = append(targetSlugs, slug)
	}

	var sessions []SessionInfo
	for _, slug := range targetSlugs {
		slugDir := filepath.Join(projectsRoot, slug)
		entries, err := os.ReadDir(slugDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
				p := filepath.Join(slugDir, e.Name())
				si, err := parseSessionFile(p)
				if err == nil {
					si.ProjectSlug = slug
					sessions = append(sessions, si)
				}
			}
		}
	}

	sort.Slice(sessions, func(i, j int) bool {
		if sessions[i].ModTime.Equal(sessions[j].ModTime) {
			return sessions[i].ID < sessions[j].ID
		}
		return sessions[i].ModTime.After(sessions[j].ModTime)
	})

	return sessions, nil
}

// LatestSession returns the most recently updated session for projectDir in profileDir, or nil.
func LatestSession(profileDir, projectDir string) (*SessionInfo, error) {
	sessions, err := ListSessions(profileDir, projectDir, false)
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}
	return &sessions[0], nil
}

// ResolveSessionID resolves a full session ID or unique prefix in projectDir within profileDir.
func ResolveSessionID(profileDir, projectDir, idOrPrefix string) (string, error) {
	idOrPrefix = strings.TrimSpace(idOrPrefix)
	if idOrPrefix == "" {
		return "", errors.New("empty session id")
	}

	slug := ProjectSlug(projectDir)
	directFile := filepath.Join(profileDir, "projects", slug, idOrPrefix+".jsonl")
	if _, err := os.Stat(directFile); err == nil {
		return idOrPrefix, nil
	}

	sessions, err := ListSessions(profileDir, projectDir, true)
	if err != nil {
		return "", err
	}

	var matched []string
	for _, s := range sessions {
		if s.ID == idOrPrefix {
			return s.ID, nil
		}
		if strings.HasPrefix(s.ID, idOrPrefix) {
			matched = append(matched, s.ID)
		}
	}

	if len(matched) == 1 {
		return matched[0], nil
	}
	if len(matched) > 1 {
		return "", fmt.Errorf("ambiguous session prefix '%s' (matches %d sessions)", idOrPrefix, len(matched))
	}
	return "", fmt.Errorf("no session found matching '%s'", idOrPrefix)
}

// TransferSession transfers a session (and its subagents, file-history, session-env, history.jsonl)
// between fromDir and toDir. If move is true, files are removed from fromDir.
func TransferSession(fromDir, toDir, projectDir, sessionID string, move bool) error {
	slug := ProjectSlug(projectDir)
	srcFile := filepath.Join(fromDir, "projects", slug, sessionID+".jsonl")
	fi, err := os.Stat(srcFile)
	if err != nil {
		// If not found in current slug, try to locate it across all projects in fromDir
		allSessions, lerr := ListSessions(fromDir, "", true)
		if lerr == nil {
			for _, s := range allSessions {
				if s.ID == sessionID {
					slug = s.ProjectSlug
					srcFile = s.Path
					fi, err = os.Stat(srcFile)
					break
				}
			}
		}
		if err != nil {
			return fmt.Errorf("session transcript not found: %s", srcFile)
		}
	}

	dstProjectDir := filepath.Join(toDir, "projects", slug)
	if err := os.MkdirAll(dstProjectDir, 0755); err != nil {
		return err
	}
	dstFile := filepath.Join(dstProjectDir, sessionID+".jsonl")

	// 1. Copy transcript preserving ModTime
	if err := copyFile(srcFile, dstFile, fi.ModTime()); err != nil {
		return err
	}

	// 2. Copy session subdirectory (subagents, tool results)
	srcSubdir := filepath.Join(fromDir, "projects", slug, sessionID)
	if sfi, err := os.Stat(srcSubdir); err == nil && sfi.IsDir() {
		dstSubdir := filepath.Join(toDir, "projects", slug, sessionID)
		if err := copyDir(srcSubdir, dstSubdir); err != nil {
			return err
		}
		if move {
			_ = os.RemoveAll(srcSubdir)
		}
	}

	// 3. Copy file-history (rollback snapshots)
	srcHistoryDir := filepath.Join(fromDir, "file-history", sessionID)
	if hfi, err := os.Stat(srcHistoryDir); err == nil && hfi.IsDir() {
		dstHistoryDir := filepath.Join(toDir, "file-history", sessionID)
		if err := copyDir(srcHistoryDir, dstHistoryDir); err != nil {
			return err
		}
		if move {
			_ = os.RemoveAll(srcHistoryDir)
		}
	}

	// 4. Copy session-env
	srcEnvDir := filepath.Join(fromDir, "session-env", sessionID)
	if efi, err := os.Stat(srcEnvDir); err == nil && efi.IsDir() {
		dstEnvDir := filepath.Join(toDir, "session-env", sessionID)
		if err := copyDir(srcEnvDir, dstEnvDir); err != nil {
			return err
		}
		if move {
			_ = os.RemoveAll(srcEnvDir)
		}
	}

	// 5. Transfer matching history.jsonl lines
	transferHistoryLines(fromDir, toDir, sessionID, move)

	if move {
		_ = os.Remove(srcFile)
	}

	return nil
}

// TransferAllSessions transfers all sessions of projectDir from fromDir to toDir.
func TransferAllSessions(fromDir, toDir, projectDir string, move bool) (int, error) {
	sessions, err := ListSessions(fromDir, projectDir, false)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, s := range sessions {
		if err := TransferSession(fromDir, toDir, projectDir, s.ID, move); err == nil {
			count++
		}
	}
	return count, nil
}

// RemoveSession deletes all data for sessionID in profileDir.
func RemoveSession(profileDir, projectDir, sessionID string) error {
	slug := ProjectSlug(projectDir)
	srcFile := filepath.Join(profileDir, "projects", slug, sessionID+".jsonl")

	if _, err := os.Stat(srcFile); err != nil {
		// Fallback: search across all project slugs
		allSessions, lerr := ListSessions(profileDir, "", true)
		if lerr == nil {
			for _, s := range allSessions {
				if s.ID == sessionID {
					slug = s.ProjectSlug
					srcFile = s.Path
					break
				}
			}
		}
	}

	_ = os.Remove(srcFile)
	_ = os.RemoveAll(filepath.Join(profileDir, "projects", slug, sessionID))
	_ = os.RemoveAll(filepath.Join(profileDir, "file-history", sessionID))
	_ = os.RemoveAll(filepath.Join(profileDir, "session-env", sessionID))
	removeHistoryLines(profileDir, sessionID)
	return nil
}

func copyFile(src, dst string, mtime time.Time) error {
	_ = os.MkdirAll(filepath.Dir(dst), 0755)
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	_ = out.Close()

	if !mtime.IsZero() {
		_ = os.Chtimes(dst, mtime, mtime)
	}
	return nil
}

func copyDir(src, dst string) error {
	_ = os.MkdirAll(dst, 0755)
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		fi, err := d.Info()
		var mtime time.Time
		if err == nil {
			mtime = fi.ModTime()
		}
		return copyFile(path, target, mtime)
	})
}

func transferHistoryLines(fromDir, toDir, sessionID string, move bool) {
	srcHistory := filepath.Join(fromDir, "history.jsonl")
	data, err := os.ReadFile(srcHistory)
	if err != nil {
		return
	}

	lines := bytes.Split(data, []byte("\n"))
	var matchedLines [][]byte
	var remainingLines [][]byte

	matchToken := fmt.Sprintf(`"sessionId":"%s"`, sessionID)

	for _, line := range lines {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		if bytes.Contains(line, []byte(matchToken)) {
			matchedLines = append(matchedLines, line)
		} else {
			remainingLines = append(remainingLines, line)
		}
	}

	if len(matchedLines) > 0 {
		dstHistory := filepath.Join(toDir, "history.jsonl")
		_ = os.MkdirAll(toDir, 0755)
		f, err := os.OpenFile(dstHistory, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			for _, ml := range matchedLines {
				_, _ = f.Write(append(ml, '\n'))
			}
			_ = f.Close()
		}
	}

	if move && len(matchedLines) > 0 {
		var buf bytes.Buffer
		for _, rl := range remainingLines {
			buf.Write(append(rl, '\n'))
		}
		tmp := srcHistory + ".tmp"
		if err := os.WriteFile(tmp, buf.Bytes(), 0644); err == nil {
			_ = os.Rename(tmp, srcHistory)
		}
	}
}

func removeHistoryLines(profileDir, sessionID string) {
	srcHistory := filepath.Join(profileDir, "history.jsonl")
	data, err := os.ReadFile(srcHistory)
	if err != nil {
		return
	}

	lines := bytes.Split(data, []byte("\n"))
	var remainingLines [][]byte
	matchToken := fmt.Sprintf(`"sessionId":"%s"`, sessionID)

	for _, line := range lines {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		if !bytes.Contains(line, []byte(matchToken)) {
			remainingLines = append(remainingLines, line)
		}
	}

	var buf bytes.Buffer
	for _, rl := range remainingLines {
		buf.Write(append(rl, '\n'))
	}
	tmp := srcHistory + ".tmp"
	if err := os.WriteFile(tmp, buf.Bytes(), 0644); err == nil {
		_ = os.Rename(tmp, srcHistory)
	}
}
