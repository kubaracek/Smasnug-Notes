package main

import (
	pkgcloak "weallonsamsung/pkg/cloak"
	pkglauncher "weallonsamsung/pkg/launcher"
)

func main() {
	var launcher pkglauncher.Launcher
	launcher = pkglauncher.NewSamsungNotes()

	var cloak pkgcloak.Cloak
	cloak = pkgcloak.NewSamsung()

	err := cloak.CloakExecuting(func() error {
		err := launcher.LaunchApp()
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
