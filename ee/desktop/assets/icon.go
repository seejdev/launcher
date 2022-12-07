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

	//go:embed kolide-debug.png
	KolideDebugDesktopIcon []byte
)
