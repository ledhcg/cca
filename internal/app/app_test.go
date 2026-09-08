package app

import (
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CCA_ROOT", "")

	a, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if a.Home != home {
		t.Errorf("Home = %q, want %q", a.Home, home)
	}
	if want := filepath.Join(home, ".claude"); a.Main != want {
		t.Errorf("Main = %q, want %q", a.Main, want)
	}
	if want := filepath.Join(home, ".claude-accounts"); a.Root != want {
		t.Errorf("Root = %q, want %q (default when CCA_ROOT is unset)", a.Root, want)
	}
	if want := filepath.Join(a.Root, "config.json"); a.ConfigPath != want {
		t.Errorf("ConfigPath = %q, want %q", a.ConfigPath, want)
	}
}

func TestNewRespectsCCARoot(t *testing.T) {
	home := t.TempDir()
	custom := filepath.Join(home, "custom-root")
	t.Setenv("HOME", home)
	t.Setenv("CCA_ROOT", custom)

	a, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if a.Root != custom {
		t.Errorf("Root = %q, want CCA_ROOT override %q", a.Root, custom)
	}
}
