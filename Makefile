default:
	go install github.com/akavel/rsrc
.PHONY: default

launcher:
	go run github.com/akavel/rsrc -manifest ./cmd/launcher/weallonsamsung.exe.manifest -o ./cmd/launcher/resource.syso
	go build -o ./bin/weallonsamsung.exe ./cmd/launcher
.PHONY: launcher

setup:
	go run github.com/akavel/rsrc -manifest ./cmd/setup/setup.exe.manifest -o ./cmd/setup/resource.syso
	go build -o ./bin/setup.exe ./cmd/setup
.PHONY: setup