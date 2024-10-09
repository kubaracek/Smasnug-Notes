package main

import (
	pkglauncher "weallonsamsung/pkg/launcher"
)

func main() {
	installAppId := "9nblggh43vhv"

	var launcher pkglauncher.Launcher
	launcher = pkglauncher.NewSamsungNotes()

	err := launcher.InstallAppId(installAppId)
	if err != nil {
		panic(err)
	}

	return
}
