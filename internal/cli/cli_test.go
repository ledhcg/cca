package cli

import (
	"os"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	boolNames := []string{"--login", "--yolo"}
	valueNames := []string{"--strategy"}

	cases := []struct {
		name       string
		args       []string
		wantPos    []string
		wantBools  map[string]bool
		wantValues map[string]string
	}{
		{
			name:       "flag before positional",
			args:       []string{"--login", "work"},
			wantPos:    []string{"work"},
			wantBools:  map[string]bool{"--login": true},
			wantValues: map[string]string{},
		},
		{
			// This ordering is exactly why parseArgs exists instead of the
			// "flag" package, which stops at the first non-flag token.
			name:       "flag after positional",
			args:       []string{"work", "--login"},
			wantPos:    []string{"work"},
			wantBools:  map[string]bool{"--login": true},
			wantValues: map[string]string{},
		},
		{
			name:       "value flag consumes next token",
			args:       []string{"work", "--strategy", "overwrite"},
			wantPos:    []string{"work"},
			wantBools:  map[string]bool{},
			wantValues: map[string]string{"--strategy": "overwrite"},
		},
		{
			name:       "value flag at end of argv with nothing to consume",
			args:       []string{"work", "--strategy"},
			wantPos:    []string{"work"},
			wantBools:  map[string]bool{},
			wantValues: map[string]string{},
		},
		{
			name:       "mixed bools, value, and multiple positionals",
			args:       []string{"work", "--yolo", "--strategy", "merge", "extra"},
			wantPos:    []string{"work", "extra"},
			wantBools:  map[string]bool{"--yolo": true},
			wantValues: map[string]string{"--strategy": "merge"},
		},
		{
			name:       "no args",
			args:       nil,
			wantPos:    nil,
			wantBools:  map[string]bool{},
			wantValues: map[string]string{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := parseArgs(c.args, boolNames, valueNames)
			if !reflect.DeepEqual(got.positional, c.wantPos) {
				t.Errorf("positional = %#v, want %#v", got.positional, c.wantPos)
			}
			if !reflect.DeepEqual(got.bools, c.wantBools) {
				t.Errorf("bools = %#v, want %#v", got.bools, c.wantBools)
			}
			if !reflect.DeepEqual(got.values, c.wantValues) {
				t.Errorf("values = %#v, want %#v", got.values, c.wantValues)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	dir := t.TempDir()
	if !isDir(dir) {
		t.Errorf("isDir(%q) = false, want true", dir)
	}

	file := dir + "/f.txt"
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if isDir(file) {
		t.Errorf("isDir(%q) = true, want false (it's a file)", file)
	}

	if isDir(dir + "/does-not-exist") {
		t.Error("isDir on a nonexistent path = true, want false")
	}
}

func TestExtractLang(t *testing.T) {
	cases := []struct {
		name     string
		in       []string
		wantArgv []string
		wantLang string
	}{
		{
			name:     "flag before command",
			in:       []string{"--lang", "vi", "ls"},
			wantArgv: []string{"ls"},
			wantLang: "vi",
		},
		{
			name:     "flag with equals before command",
			in:       []string{"--lang=ja", "doctor"},
			wantArgv: []string{"doctor"},
			wantLang: "ja",
		},
		{
			name:     "flag after command",
			in:       []string{"ls", "--lang", "es"},
			wantArgv: []string{"ls"},
			wantLang: "es",
		},
		{
			name:     "passthrough use command preserves trailing flag",
			in:       []string{"use", "work", "--lang", "fr"},
			wantArgv: []string{"use", "work", "--lang", "fr"},
			wantLang: "",
		},
		{
			name:     "settings lang positional is untouched",
			in:       []string{"settings", "lang", "vi"},
			wantArgv: []string{"settings", "lang", "vi"},
			wantLang: "",
		},
		{
			name:     "no lang flag",
			in:       []string{"ls", "--names"},
			wantArgv: []string{"ls", "--names"},
			wantLang: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotArgv, gotLang := extractLang(tc.in)
			if !reflect.DeepEqual(gotArgv, tc.wantArgv) {
				t.Errorf("gotArgv = %#v, want %#v", gotArgv, tc.wantArgv)
			}
			if gotLang != tc.wantLang {
				t.Errorf("gotLang = %q, want %q", gotLang, tc.wantLang)
			}
		})
	}
}
