default:
	go install github.com/akavel/rsrc
	sudo apt-get install gcc-mingw-w64
.PHONY: default

launcher:
	go run github.com/akavel/rsrc -manifest ./cmd/launcher/weallonsamsung.exe.manifest -o ./cmd/launcher/resource.syso
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o ./bin/weallonsamsung.exe ./cmd/launcher
.PHONY: launcher

setup:
	go run github.com/akavel/rsrc -manifest ./cmd/setup/setup.exe.manifest -o ./cmd/setup/resource.syso
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o ./bin/setup.exe ./cmd/setup
.PHONY: setup

version:
	@echo "Updating all .manifest files to version $(VERSION)"
	@find . -name '*.manifest' -exec sed -i'' -e 's/{{version}}/$(VERSION)/g' {} +
.PHONY: version

zip:
	make -j 2 launcher setup
	zip