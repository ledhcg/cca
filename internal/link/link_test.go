package link

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMakeLinkCreatesAndIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	line := Make(dst, src)
	if line == "" {
		t.Fatal("makeLink on a fresh dst should return a non-empty log line")
	}
	if !IsLink(dst) {
		t.Fatal("dst should be a link after makeLink")
	}
	got, err := os.ReadFile(dst)
	if err != nil || string(got) != "hello" {
		t.Fatalf("dst content = %q, err=%v, want %q", got, err, "hello")
	}

	// Calling again with an already-correct link should be a no-op.
	if line := Make(dst, src); line != "" {
		t.Errorf("makeLink on an already-correct link should return \"\", got %q", line)
	}
}

func TestMakeLinkSkipsRealDirectory(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst-is-a-real-dir")
	if err := os.WriteFile(src, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(dst, 0755); err != nil {
		t.Fatal(err)
	}

	line := Make(dst, src)
	if line == "" {
		t.Fatal("expected a 'skipped' log line, got empty string")
	}
	fi, err := os.Lstat(dst)
	if err != nil || !fi.IsDir() {
		t.Error("makeLink must not touch a real directory at dst")
	}
}

func TestMakeLinkSkipsMissingSource(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "does-not-exist")
	dst := filepath.Join(dir, "dst.txt")

	line := Make(dst, src)
	if line == "" {
		t.Fatal("expected a 'skipped' log line when source is missing, got empty string")
	}
	if _, err := os.Lstat(dst); err == nil {
		t.Error("dst should not have been created when source is missing")
	}
}

func TestAlreadyLinked(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	other := filepath.Join(dir, "other.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(other, []byte("b"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Symlink(src, dst); err != nil {
		t.Fatal(err)
	}
	if !AlreadyLinked(dst, src) {
		t.Error("alreadyLinked should be true when dst already points at src")
	}
	if AlreadyLinked(dst, other) {
		t.Error("alreadyLinked should be false when dst points elsewhere")
	}
}

func TestCopyFileAndCopyDir(t *testing.T) {
	dir := t.TempDir()
	srcDir := filepath.Join(dir, "src")
	dstDir := filepath.Join(dir, "dst")
	if err := os.MkdirAll(filepath.Join(srcDir, "nested"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("A"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "nested", "b.txt"), []byte("B"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := CopyDir(srcDir, dstDir); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dstDir, "nested", "b.txt"))
	if err != nil || string(got) != "B" {
		t.Fatalf("copyDir did not copy nested file correctly: %q, err=%v", got, err)
	}
}
