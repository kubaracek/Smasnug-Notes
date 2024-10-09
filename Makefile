launcher:
	go run github.com/akavel/rsrc -manifest weallonsamsung-launcher.exe.manifest -o ./bin/weallonsamsung.syso
	go build -o ./bin/weallonsamsung.exe ./cmd/launcher
setup:
	go build -o ./bin/setup.exe ./cmd/setup