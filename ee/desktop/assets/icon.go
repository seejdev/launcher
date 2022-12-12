//go:build darwin || linux
// +build darwin linux

package assets

import (
	_ "embed"
)

var (
	//go:embed kolide.png
	KolideDesktopIcon []byte

	//go:embed kolide_warn.png
	KolideDesktopIconWarn []byte
	//go:embed kolide_fail.png
	KolideDesktopIconFail []byte
	//go:embed kolide_idle.png
	KolideDesktopIconIdle []byte

	//go:embed circle-green.png
	KolideStatusGreen []byte
	//go:embed circle-yellow.png
	KolideStatusYellow []byte
	//go:embed circle-red.png
	KolideStatusRed []byte

	//go:embed kolide-debug.png
	KolideDebugDesktopIcon []byte
)
