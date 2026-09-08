package config

import "testing"

func TestMergeSettings(t *testing.T) {
	main := map[string]any{
		"model":       "opus",
		"statusLine":  "custom.sh",
		"permissions": map[string]any{"defaultMode": "default"},
	}
	local := map[string]any{
		"model":       "sonnet", // profileLocalKey → local wins
		"statusLine":  "old.sh", // not a profileLocalKey → main wins
		"extraLocal":  "kept",   // not present in main, and not a profileLocalKey → dropped
		"effortLevel": "high",   // profileLocalKey, absent from main → main has no such key so it's added
	}
	localKeys := []string{"model", "effortLevel", "permissions"}

	got := Merge(main, local, localKeys)

	if got["model"] != "sonnet" {
		t.Errorf("model = %v, want local value 'sonnet'", got["model"])
	}
	if got["statusLine"] != "custom.sh" {
		t.Errorf("statusLine = %v, want main value 'custom.sh' (not a profileLocalKey)", got["statusLine"])
	}
	if _, ok := got["extraLocal"]; ok {
		t.Errorf("extraLocal should not survive the merge: it's local-only and not a profileLocalKey")
	}
	if got["effortLevel"] != "high" {
		t.Errorf("effortLevel = %v, want local value 'high' (profileLocalKey absent from main)", got["effortLevel"])
	}
	if perms, ok := got["permissions"].(map[string]any); !ok || perms["defaultMode"] != "default" {
		t.Errorf("permissions = %v, want main's value (local has no 'permissions' key to override with)", got["permissions"])
	}
}

func TestMergeSettingsNoLocalOverrides(t *testing.T) {
	main := map[string]any{"a": 1.0, "b": 2.0}
	local := map[string]any{"a": 99.0, "b": 99.0}

	got := Merge(main, local, nil)

	if got["a"] != 1.0 || got["b"] != 2.0 {
		t.Errorf("with no profileLocalKeys, everything should follow main; got %#v", got)
	}
}
