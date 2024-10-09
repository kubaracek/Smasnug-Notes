default:
	go install github.com/akavel/rsrc
.PHONY: default

launcher:
	go run github.com/akavel/rsrc -manifest ./cmd/launcher/weallonsamsung.exe.manifest -o ./cmd/launcher/resource.syso
	go build -o ./bin/weallonsamsung.exe ./cmd/launcher
.PHONY: launcher

zip:
	make -j 2 launcher setup
	zip