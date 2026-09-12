package i18n

import (
	"strings"
	"sync"
	"testing"
)

func TestCatalogParity(t *testing.T) {
	catalogsToTest := []struct {
		code string
		cat  map[Key]string
	}{
		{"vi", viCatalog},
		{"zh", zhCatalog},
		{"ja", jaCatalog},
		{"es", esCatalog},
	}

	for k := range enCatalog {
		for _, tc := range catalogsToTest {
			val, ok := tc.cat[k]
			if !ok || val == "" {
				t.Errorf("[%s] missing translation for key: %s", tc.code, k)
			}
		}
	}
}

func TestResolveHierarchy(t *testing.T) {
	// Clear env vars for testing
	t.Setenv("CCA_LANG", "")
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "")

	// 1. Default fallback
	if got := Resolve("", ""); got != "en" {
		t.Errorf("expected 'en', got '%s'", got)
	}

	// 2. System LANG fallback
	t.Setenv("LANG", "vi_VN.UTF-8")
	if got := Resolve("", ""); got != "vi" {
		t.Errorf("expected 'vi' from LANG, got '%s'", got)
	}

	// 3. Config overrides LANG
	if got := Resolve("", "ja"); got != "ja" {
		t.Errorf("expected 'ja' from config, got '%s'", got)
	}

	// 4. CCA_LANG overrides config
	t.Setenv("CCA_LANG", "zh")
	if got := Resolve("", "ja"); got != "zh" {
		t.Errorf("expected 'zh' from CCA_LANG, got '%s'", got)
	}

	// 5. CLI flag overrides CCA_LANG
	if got := Resolve("es", "ja"); got != "es" {
		t.Errorf("expected 'es' from CLI flag, got '%s'", got)
	}
}

func TestFallback(t *testing.T) {
	Init("vi")
	defer Init("en")

	// Existing key in vi
	if got := T(KeyCommonDefault); got != "mặc định" {
		t.Errorf("expected 'mặc định', got '%s'", got)
	}

	// Nonexistent key returns key itself
	if got := T("nonexistent.key"); got != "nonexistent.key" {
		t.Errorf("expected 'nonexistent.key', got '%s'", got)
	}
}

func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			Init("vi")
			_ = T(KeyCommonOk)
		}()
		go func() {
			defer wg.Done()
			Init("en")
			_ = T(KeyCommonOk)
		}()
	}
	wg.Wait()
}

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"en":         "en",
		"EN":         "en",
		"en_US.UTF8": "en",
		"vi":         "vi",
		"vi_VN":      "vi",
		"vietnamese": "vi",
		"zh":         "zh",
		"zh-CN":      "zh",
		"ja":         "ja",
		"ja_JP":      "ja",
		"es":         "es",
		"es_ES":      "es",
		"unknown":    "",
	}

	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPrintUsageGuideKeys(t *testing.T) {
	for _, code := range []string{"en", "vi", "zh", "ja", "es"} {
		Init(code)
		usage := T(KeyGuideUsage)
		if !strings.Contains(usage, "cca") {
			t.Errorf("[%s] KeyGuideUsage does not contain 'cca'", code)
		}
		guide := T(KeyGuideFull, "", "", "", "")
		if !strings.Contains(guide, "cca") {
			t.Errorf("[%s] KeyGuideFull does not contain 'cca'", code)
		}
	}
}
