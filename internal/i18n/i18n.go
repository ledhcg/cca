package i18n

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// LangInfo holds metadata for a supported language.
type LangInfo struct {
	Code        string
	NativeName  string
	EnglishName string
}

var (
	mu          sync.RWMutex
	currentLang = "en"

	catalogs = map[string]map[Key]string{
		"en": enCatalog,
		"vi": viCatalog,
		"zh": zhCatalog,
		"ja": jaCatalog,
		"es": esCatalog,
	}

	supported = []LangInfo{
		{Code: "en", NativeName: "English", EnglishName: "English"},
		{Code: "vi", NativeName: "Tiếng Việt", EnglishName: "Vietnamese"},
		{Code: "zh", NativeName: "简体中文", EnglishName: "Simplified Chinese"},
		{Code: "ja", NativeName: "日本語", EnglishName: "Japanese"},
		{Code: "es", NativeName: "Español", EnglishName: "Spanish"},
	}
)

// SupportedLanguages returns the list of supported languages.
func SupportedLanguages() []LangInfo {
	return supported
}

// SupportedCodes returns a slice of supported language codes.
func SupportedCodes() []string {
	codes := make([]string, len(supported))
	for i, l := range supported {
		codes[i] = l.Code
	}
	return codes
}

// Normalize sanitizes a language string (e.g. "vi_VN.UTF-8" -> "vi", "VI" -> "vi").
func Normalize(raw string) string {
	s := strings.TrimSpace(strings.ToLower(raw))
	if idx := strings.IndexAny(s, "._-@"); idx != -1 {
		s = s[:idx]
	}
	switch s {
	case "vi", "vietnamese":
		return "vi"
	case "zh", "chinese", "cmn":
		return "zh"
	case "ja", "japanese":
		return "ja"
	case "es", "spanish":
		return "es"
	case "en", "english":
		return "en"
	default:
		return ""
	}
}

// Resolve determines the effective language following the priority hierarchy:
// 1. cliLang (if non-empty)
// 2. CCA_LANG environment variable
// 3. cfgLang from config.json (if non-empty)
// 4. LC_ALL, LC_MESSAGES, LANG system environment variables
// 5. Default fallback to "en"
func Resolve(cliLang, cfgLang string) string {
	if cliLang != "" {
		if l := Normalize(cliLang); l != "" {
			return l
		}
	}
	if env := os.Getenv("CCA_LANG"); env != "" {
		if l := Normalize(env); l != "" {
			return l
		}
	}
	if cfgLang != "" {
		if l := Normalize(cfgLang); l != "" {
			return l
		}
	}
	for _, envKey := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if val := os.Getenv(envKey); val != "" {
			if l := Normalize(val); l != "" {
				return l
			}
		}
	}
	return "en"
}

// Init sets the active language for translations.
func Init(lang string) {
	mu.Lock()
	defer mu.Unlock()
	normalized := Normalize(lang)
	if _, ok := catalogs[normalized]; ok {
		currentLang = normalized
	} else {
		currentLang = "en"
	}
}

// Current returns the code of the currently active language.
func Current() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLang
}

// Name returns the localized and English name for a given language code.
func Name(code string) string {
	for _, l := range supported {
		if l.Code == code {
			return fmt.Sprintf("%s (%s)", l.NativeName, l.Code)
		}
	}
	return code
}

// T looks up key in the active catalog, with automatic fallback to English,
// and returns the formatted string.
func T(key Key, args ...any) string {
	mu.RLock()
	lang := currentLang
	mu.RUnlock()

	cat, ok := catalogs[lang]
	if !ok {
		cat = catalogs["en"]
	}

	tmpl, found := cat[key]
	if !found {
		// Fall back to English
		tmpl, found = catalogs["en"][key]
		if !found {
			if len(args) == 0 {
				return key
			}
			return key + " " + fmt.Sprint(args...)
		}
	}

	if len(args) == 0 {
		return tmpl
	}
	return fmt.Sprintf(tmpl, args...)
}
