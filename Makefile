launcher:
	go run github.com/akavel/rsrc -manifest weallonsamsung-launcher.exe.manifest -o ./bin/weallonsamsung-launcher.syso
	go build -o ./bin/weallonsamsung-launcher.exe ./cmd/launcher