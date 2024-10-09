package main

import (
	"fmt"
	"github.com/tawesoft/golib/v2/dialog"
	"syscall"
	pkgcloak "weallonsamsung/pkg/cloak"
	pkglauncher "weallonsamsung/pkg/launcher"
)

const SAMSUNG_NOTES_INSTALL_ID = "9NBLGGH43VHV"
const SAMSUNG_NOTES_LAUNCH_ID = "wyx1vj98g3asy"

func main() {
	hideConsole()
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

	shouldInstall := false

	if !isInstalled {
		answer, err := dialog.Ask("Samsung Notes is not installed. Do you wish to install it?\n\n" +
			"Samsung Notes will start Automatically after the installation.\n" +
			"It might take a minute or two\n\n" +
			"If you choose 'No' you can install Samsung Notes yourself directly from Microsoft Store.\n")
		if err != nil {
			panic(err)
		}

		shouldInstall = answer
	}

	if shouldInstall {
		err := launcher.InstallAppId(SAMSUNG_NOTES_INSTALL_ID)
		if err != nil {
			err := dialog.Error(fmt.Sprintf("Unable to install the app: %+v", err))
			if err != nil {
				panic(err)
			}
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

func hideConsole() {
	getConsoleWindow := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	if getConsoleWindow.Find() != nil {
		return
	}

	showWindow := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	if showWindow.Find() != nil {
		return
	}

	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return
	}

	showWindow.Call(hwnd, 0)
}
