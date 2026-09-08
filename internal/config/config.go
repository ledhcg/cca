// Package config owns cca's shared configuration (~/.claude-accounts/config.json)
// and syncing shared files (plugins, skills, settings.json, ...) into a
// profile directory.
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ledhcg/cca/internal/app"
	"github.com/ledhcg/cca/internal/link"
)

// Config is the shared configuration linked into every profile. See `cca sync --help`.
type Config struct {
	// Linked into ~/.claude — tooling, not account data.
	SharedLinks []string `json:"sharedLinks"`
	// Copied from ~/.claude — each profile keeps its own editable copy.
	SharedCopies []string `json:"sharedCopies"`
	// overwrite | keep-local | merge
	SyncStrategy string `json:"syncStrategy"`
	// With syncStrategy=merge: these keys belong to each profile, sync leaves them alone.
	ProfileLocalKeys []string `json:"profileLocalKeys"`
}

// Default is used until a config.json exists on disk.
var Default = Config{
	SharedLinks:      []string{"plugins", "skills", "agents", "statusline-command.sh"},
	SharedCopies:     []string{"settings.json"},
	SyncStrategy:     "merge",
	ProfileLocalKeys: []string{"model", "effortLevel", "permissions"},
}

// Load reads a.ConfigPath, falling back to Default if it doesn't exist yet.
func Load(a *app.App) (Config, error) {
	cfg := Default
	data, err := os.ReadFile(a.ConfigPath)
	if err != nil {
		return cfg, nil // no file yet → use defaults
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("%s is not valid JSON: %w", a.ConfigPath, err)
	}
	return cfg, nil
}

// Save writes cfg to a.ConfigPath, creating a.Root if needed.
func Save(a *app.App, cfg Config) error {
	if err := os.MkdirAll(a.Root, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.ConfigPath, append(data, '\n'), 0644)
}

// Merge merges ~/.claude's settings.json into a profile's copy.
//
// Rule: everything follows the original (~/.claude) — except keys listed in
// localKeys, which are each account's own choice.
func Merge(main, local map[string]any, localKeys []string) map[string]any {
	merged := make(map[string]any, len(main))
	for k, v := range main {
		merged[k] = v
	}
	for _, k := range localKeys {
		if v, ok := local[k]; ok {
			merged[k] = v
		}
	}
	return merged
}

// Seed syncs cfg's SharedLinks/SharedCopies from a.Main into dir (a profile
// directory, from profileenv.Dir), creating dir if needed. strategy overrides
// cfg.SyncStrategy when non-empty. Returns a human-readable log of what changed.
func Seed(a *app.App, dir string, cfg Config, strategy string) []string {
	_ = os.MkdirAll(dir, 0755)
	var log []string
	if strategy == "" {
		strategy = cfg.SyncStrategy
	}

	for _, item := range cfg.SharedLinks {
		src := filepath.Join(a.Main, item)
		dst := filepath.Join(dir, item)
		if _, err := os.Stat(src); err != nil {
			continue
		}
		if line := link.Make(dst, src); line != "" {
			log = append(log, line)
		}
	}

	for _, item := range cfg.SharedCopies {
		src := filepath.Join(a.Main, item)
		dst := filepath.Join(dir, item)
		sfi, err := os.Stat(src)
		if err != nil || sfi.IsDir() {
			continue
		}
		if _, err := os.Stat(dst); err != nil {
			_ = link.CopyFile(src, dst)
			log = append(log, "copied "+item)
			continue
		}
		srcBytes, _ := os.ReadFile(src)
		dstBytes, _ := os.ReadFile(dst)
		if bytes.Equal(srcBytes, dstBytes) {
			continue
		}
		switch {
		case strategy == "keep-local":
			log = append(log, "kept "+item+" as-is (locally modified)")
		case strategy == "overwrite":
			_ = link.CopyFile(src, dst)
			log = append(log, "overwrote "+item)
		case strategy == "merge" && strings.HasSuffix(item, ".json"):
			var mainData, localData map[string]any
			if err := json.Unmarshal(srcBytes, &mainData); err != nil {
				log = append(log, "skipped "+item+" (malformed JSON: "+err.Error()+")")
				continue
			}
			if err := json.Unmarshal(dstBytes, &localData); err != nil {
				log = append(log, "skipped "+item+" (malformed JSON: "+err.Error()+")")
				continue
			}
			merged := Merge(mainData, localData, cfg.ProfileLocalKeys)
			mergedBytes, _ := json.MarshalIndent(merged, "", "  ")
			mergedBytes = append(mergedBytes, '\n')
			var mergedCheck, currentCheck map[string]any
			_ = json.Unmarshal(mergedBytes, &mergedCheck)
			_ = json.Unmarshal(dstBytes, &currentCheck)
			if reflect.DeepEqual(mergedCheck, currentCheck) {
				continue
			}
			_ = os.WriteFile(dst, mergedBytes, 0644)
			var kept []string
			for _, k := range cfg.ProfileLocalKeys {
				if _, ok := localData[k]; ok {
					kept = append(kept, k)
				}
			}
			keptStr := "none"
			if len(kept) > 0 {
				keptStr = strings.Join(kept, ", ")
			}
			log = append(log, "merged "+item+" (kept local: "+keptStr+")")
		default:
			_ = link.CopyFile(src, dst)
			log = append(log, "overwrote "+item)
		}
	}
	return log
}
