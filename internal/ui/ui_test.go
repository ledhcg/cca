package ui

import "testing"

func TestStringWidth(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"PROFILE", 7},
		{"Đã đăng nhập", 12}, // 12 display columns (Latin accented), 18 bytes in UTF-8
		{"HỒ SƠ", 5},
		{"已登录", 6},     // 3 CJK characters = 6 display columns
		{"プロファイル", 12}, // 6 Japanese Katakana characters = 12 display columns
		{"Cài đặt", 7},
		{"", 0},
	}

	for _, c := range cases {
		if got := StringWidth(c.in); got != c.want {
			t.Errorf("StringWidth(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestPadRight(t *testing.T) {
	// Pad "Đã đăng nhập" (width 12) to width 18 -> should add 6 spaces
	padded := PadRight("Đã đăng nhập", 18)
	if got := StringWidth(padded); got != 18 {
		t.Errorf("StringWidth(PadRight(%q, 18)) = %d, want 18", "Đã đăng nhập", got)
	}

	// Pad "已登录" (width 6) to width 10 -> should add 4 spaces
	paddedZh := PadRight("已登录", 10)
	if got := StringWidth(paddedZh); got != 10 {
		t.Errorf("StringWidth(PadRight(%q, 10)) = %d, want 10", "已登录", got)
	}
}
