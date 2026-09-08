//go:build windows

package ui

import "syscall"

// SetupConsole forces Windows's default legacy code page (e.g. cp1252) to
// UTF-8 so non-ASCII output renders correctly. Uses syscall directly (no
// golang.org/x/sys dependency needed) — just two kernel32 calls, best-effort
// so errors are ignored.
func SetupConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleOutputCP := kernel32.NewProc("SetConsoleOutputCP")
	setConsoleCP := kernel32.NewProc("SetConsoleCP")
	const cpUTF8 = 65001
	setConsoleOutputCP.Call(uintptr(cpUTF8))
	setConsoleCP.Call(uintptr(cpUTF8))
}
