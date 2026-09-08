package profileenv

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ledhcg/cca/internal/app"
)

func TestEnvSetGetUnset(t *testing.T) {
	env := []string{"FOO=1", "BAR=2"}

	env = EnvSet(env, "FOO", "99")
	if EnvGet(env, "FOO") != "99" {
		t.Errorf("EnvSet should overwrite an existing key; got %q", EnvGet(env, "FOO"))
	}
	if len(env) != 2 {
		t.Errorf("EnvSet on an existing key should not grow the slice; len=%d", len(env))
	}

	env = EnvSet(env, "BAZ", "3")
	if EnvGet(env, "BAZ") != "3" {
		t.Errorf("EnvSet should append a new key; got %q", EnvGet(env, "BAZ"))
	}
	if len(env) != 3 {
		t.Errorf("EnvSet on a new key should grow the slice; len=%d", len(env))
	}

	env = EnvUnset(env, "BAR")
	if EnvGet(env, "BAR") != "" {
		t.Errorf("EnvUnset should remove the key; still got %q", EnvGet(env, "BAR"))
	}
	if len(env) != 2 {
		t.Errorf("EnvUnset should shrink the slice; len=%d", len(env))
	}

	if EnvGet(env, "NOPE") != "" {
		t.Error("EnvGet on a missing key should return empty string")
	}
}

func TestValidate(t *testing.T) {
	valid := []string{"work", "work-2", "w.2_3", "9team"}
	for _, n := range valid {
		if err := Validate(n); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", n, err)
		}
	}

	invalid := []string{"", "-work", ".work", "wo rk", "wo/rk", "wo@rk", DefaultName, "config"}
	for _, n := range invalid {
		if err := Validate(n); err == nil {
			t.Errorf("Validate(%q) = nil, want an error", n)
		}
	}
}

func TestIsDefault(t *testing.T) {
	cases := map[string]bool{"": true, "default": true, "work": false, "Default": false}
	for name, want := range cases {
		if got := IsDefault(name); got != want {
			t.Errorf("IsDefault(%q) = %v, want %v", name, got, want)
		}
	}
}

func testApp(t *testing.T) *app.App {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CCA_ROOT", filepath.Join(home, ".claude-accounts"))
	a, err := app.New()
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestDir(t *testing.T) {
	a := testApp(t)

	if got := Dir(a, ""); got != a.Main {
		t.Errorf("Dir(\"\") = %q, want Main (%q)", got, a.Main)
	}
	if got := Dir(a, DefaultName); got != a.Main {
		t.Errorf("Dir(%q) = %q, want Main (%q)", DefaultName, got, a.Main)
	}
	want := filepath.Join(a.Root, "work")
	if got := Dir(a, "work"); got != want {
		t.Errorf("Dir(\"work\") = %q, want %q", got, want)
	}
}

func TestActive(t *testing.T) {
	a := testApp(t)

	if err := os.MkdirAll(Dir(a, "work"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.Unsetenv("CLAUDE_CONFIG_DIR"); err != nil {
		t.Fatal(err)
	}
	if got := Active(a); got != "" {
		t.Errorf("Active() with no CLAUDE_CONFIG_DIR = %q, want \"\" (default)", got)
	}

	t.Setenv("CLAUDE_CONFIG_DIR", Dir(a, "work"))
	if got := Active(a); got != "work" {
		t.Errorf("Active() = %q, want \"work\"", got)
	}
}
