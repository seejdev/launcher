//go:build linux
// +build linux

package menu

func isDarkMode() bool {
	return false
}

func RegisterThemeChangeListener(f func()) {
	// no-op
}
