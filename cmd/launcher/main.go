package main

import (
	"fmt"
	"github.com/tawesoft/golib/v2/dialog"
	pkgcloak "smasnugnotes/pkg/cloak"
	pkglauncher "smasnugnotes/pkg/launcher"
	pkgui "smasnugnotes/pkg/ui"
)

const SAMSUNG_NOTES_INSTALL_ID = "9NBLGGH43VHV"
const SAMSUNG_NOTES_LAUNCH_ID = "wyx1vj98g3asy"

func main() {
	registerValues := pkgcloak.RegisterValues{
		"HARDWARE\\DESCRIPTION\\System\\BIOS": {
			"SystemProductName":  "NP960XFG-KC4UK",
			"SystemManufacturer": "Samsung",
		},
	}
	launcher := pkglauncher.NewLauncher()

	isInstalled, err := launcher.IsInstalledAppId(SAMSUNG_NOTES_INSTALL_ID)
	if err != nil {
		err := dialog.Error("Failed to query installed apps")
		if err != nil {
			panic(err)
		}
	}

	if !isInstalled {
		installed, err := pkgui.InstallDialog(func() error {
			return launcher.InstallAppId(SAMSUNG_NOTES_INSTALL_ID)
		})

		if err != nil {
			panic(err)
		}

		if !installed {
			return
		}
	}

	var cloak pkgcloak.Cloak
	cloak = pkgcloak.NewRegisterCloak(registerValues)

	err = cloak.CloakExecution(func() error {
		err := launcher.LaunchAppId(SAMSUNG_NOTES_LAUNCH_ID)
		if err != nil {
			err := dialog.Error(fmt.Sprintf("Failed to launch app: %+v", err))
			if err != nil {
				panic(err)
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return
}
