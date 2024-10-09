package main

import (
	pkgcloak "weallonsamsung/pkg/cloak"
	pkglauncher "weallonsamsung/pkg/launcher"
)

func main() {
	registerValues := pkgcloak.RegisterValues{
		"HARDWARE\\DESCRIPTION\\System\\BIOS": {
			"SystemProductName":  "NP960XFG-KC4UK",
			"SystemManufacturer": "Samsung",
		},
	}
	launchAppId := "wyx1vj98g3asy"

	var launcher pkglauncher.Launcher
	launcher = pkglauncher.NewSamsungNotes()

	var cloak pkgcloak.Cloak
	cloak = pkgcloak.NewRegisterCloak(registerValues)

	err := cloak.CloakExecution(func() error {
		err := launcher.LaunchAppId(launchAppId)
		if err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return
}
