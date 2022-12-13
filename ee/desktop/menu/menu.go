package menu

import (
	"fyne.io/systray"
	"github.com/kolide/launcher/ee/desktop/assets"
)

var dynMenuItem *systray.MenuItem

func Init(hostname string) {
	onReady := func() {
		systray.SetTemplateIcon(assets.KolideDesktopIcon, assets.KolideDesktopIcon)
		systray.SetTooltip("Kolide")

		dynMenuItem = systray.AddMenuItem("Kolide Agent is running", "")

		msg := "Kolide Agent is blocking access"
		systray.SetTemplateIcon(assets.KolideDesktopIconFail, assets.KolideDesktopIconFail)
		systray.SetTooltip(msg)
		dynMenuItem.SetTitle(msg)
		dynMenuItem.SetIcon(assets.KolideStatusRed)

		systray.AddSeparator()

		dynMenuItem.Disable()
		// dynMenuItem.SetIcon(assets.KolideStatusGreen)

		// versionItem := systray.AddMenuItem(fmt.Sprintf("Version %s", version.Version().Version), "")
		// versionItem.Disable()

		failingChecksMenuItem := systray.AddMenuItem("Failing checks", "")
		failingChecksMenuItem.SetIcon(assets.KolideAlert)
		failingChecksMenuItem.AddSubMenuItem("Require Google Chrome is up-to-date", "")
		failingChecksMenuItem.AddSubMenuItem("Ensure Kolide Agent Has Full Disk Access Entitlement", "")
		failingChecksMenuItem.AddSubMenuItem("Ensure Important Updates Are Installed", "")
		failingChecksMenuItem.AddSubMenuItem("Help Guide...", "")

		systray.AddSeparator()
		privacyCenter := systray.AddMenuItem("My Privacy Center", "")
		privacyCenter.SetIcon(assets.KolidePrivacyCenter)

		systray.AddSeparator()
		systray.AddMenuItem("About Kolide agent", "")

		// if prod environment, return
		if hostname == "k2device-preprod.kolide.com" || hostname == "k2device.kolide.com" {
			return
		}
		/*
			// in non prod environment
			systray.SetTemplateIcon(assets.KolideDebugDesktopIcon, assets.KolideDebugDesktopIcon)
			systray.AddSeparator()
			systray.AddMenuItem("--- DEBUG ---", "").Disable()
			systray.AddMenuItem(fmt.Sprintf("Hostname: %s", hostname), "").Disable()
		*/
	}

	systray.Run(onReady, func() {})
}

func SetStatus(st string) {
	switch st {
	case "good":
		msg := "Kolide agent is running"
		systray.SetTemplateIcon(assets.KolideDesktopIcon, assets.KolideDesktopIcon)
		systray.SetTooltip(msg)
		dynMenuItem.SetTitle(msg)
		dynMenuItem.SetIcon(assets.KolideStatusGreen)
	case "warn":
		msg := "Kolide agent has detected problems"
		systray.SetTemplateIcon(assets.KolideDesktopIconWarn, assets.KolideDesktopIconWarn)
		systray.SetTooltip(msg)
		dynMenuItem.SetTitle(msg)
		dynMenuItem.SetIcon(assets.KolideStatusYellow)
	case "fail":
		msg := "Kolide agent is blocking access"
		systray.SetTemplateIcon(assets.KolideDesktopIconFail, assets.KolideDesktopIconFail)
		systray.SetTooltip(msg)
		dynMenuItem.SetTitle(msg)
		dynMenuItem.SetIcon(assets.KolideStatusRed)
	case "idle":
		msg := "Kolide agent is running"
		systray.SetTemplateIcon(assets.KolideDesktopIconIdle, assets.KolideDesktopIconIdle)
		systray.SetTooltip(msg)
		dynMenuItem.SetTitle(msg)
		dynMenuItem.SetIcon(assets.KolideStatusGreen)
	}
}

func Shutdown() {
	systray.Quit()
}
