// Package link handles syncing shared files (plugins/skills/agents/settings)
// from ~/.claude into a profile directory: a real symlink where possible,
// degrading gracefully to a junction, hard link, or plain copy on platforms
// that can't grant one.
package link

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// IsLink reports whether path is a symlink OR a junction. On Windows,
// Lstat().Mode() sets os.ModeSymlink for both reparse-point kinds — verified
// experimentally, no separate handling needed like some other languages require.
func IsLink(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink != 0
}

// AlreadyLinked reports whether dst is already a link (or, on a Windows
// hard-link fallback, the same file) pointing at src.
func AlreadyLinked(dst, src string) bool {
	if IsLink(dst) {
		target, err := os.Readlink(dst)
		if err != nil {
			return false
		}
		ta, err1 := filepath.Abs(target)
		sa, err2 := filepath.Abs(src)
		return err1 == nil && err2 == nil && ta == sa
	}
	// A Windows hard-link fallback leaves no reparse point for Readlink to see —
	// compare by file identity (same inode/file-id) instead of path.
	dfi, err1 := os.Stat(dst)
	sfi, err2 := os.Stat(src)
	if err1 != nil || err2 != nil || dfi.IsDir() {
		return false
	}
	return os.SameFile(dfi, sfi)
}

// Make syncs dst into a link pointing at src. Tries a real symlink first;
// on Windows, if that's denied (Developer Mode off / not admin), it degrades in
// order: junction (directories, no special privilege needed) → hard link
// (files) → copy.
//
// Returns a log line describing the action taken, or "" if dst was already correct.
func Make(dst, src string) string {
	if AlreadyLinked(dst, src) {
		return ""
	}

	if IsLink(dst) {
		_ = os.Remove(dst) // safe for both symlinks and junctions — verified experimentally
	} else if fi, err := os.Stat(dst); err == nil {
		if fi.IsDir() {
			return "skipped " + filepath.Base(dst) + " (it's a real directory, not a link)"
		}
		_ = os.Remove(dst)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return "skipped " + filepath.Base(dst) + " (source does not exist)"
	}
	srcIsDir := srcInfo.IsDir()

	if err := os.Symlink(src, dst); err == nil {
		return "linked " + filepath.Base(dst) + " → " + src
	}

	if runtime.GOOS == "windows" && srcIsDir {
		cmd := exec.Command("cmd", "/c", "mklink", "/J", dst, src)
		if err := cmd.Run(); err == nil {
			return "junction " + filepath.Base(dst) + " → " + src + "  (symlink needs Developer Mode — used a junction instead)"
		}
	}

	if !srcIsDir {
		if err := os.Link(src, dst); err == nil {
			return "hard link " + filepath.Base(dst) + " → " + src
		}
	}

	var copyErr error
	if srcIsDir {
		copyErr = CopyDir(src, dst)
	} else {
		copyErr = CopyFile(src, dst)
	}
	if copyErr != nil {
		return "failed to sync " + filepath.Base(dst) + " (no link, and copy also failed: " + copyErr.Error() + ")"
	}
	hint := "check symlink permissions"
	if runtime.GOOS == "windows" {
		hint = "enable Developer Mode: Settings → Privacy & security → For developers"
	}
	return "copied " + filepath.Base(dst) + " (could not create a link — " + hint + "; edits here will NOT sync back)"
}

// RemoveProfileDir deletes a profile directory. os.RemoveAll has been verified
// safe with Windows junctions — it unlinks the reparse point without recursing
// into the real target — so no "unlink links first" safety dance is needed here
// unlike some other languages/runtimes.
func RemoveProfileDir(d string) error {
	return os.RemoveAll(d)
}

// CopyFile copies src to dst, preserving src's file mode.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }() // read-only handle, nothing to flush on close

	fi, err := in.Stat()
	if err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fi.Mode())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close() // surfaces a flush error a deferred Close would silently drop
}

// CopyDir recursively copies src into dst.
func CopyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := CopyDir(s, d); err != nil {
				return err
			}
		} else {
			if err := CopyFile(s, d); err != nil {
				return err
			}
		}
	}
	return nil
}
