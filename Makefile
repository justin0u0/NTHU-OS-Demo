all: build

.PHONY: build-all
build-all: build build-darwin-arm64 build-darwin-amd64 build-windows-amd64 build-linux-amd64

.PHONY: build
build:
	go build -o ./bin/demo .

.PHONY: clean
clean:
	rm -f ./bin/**

build-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 go build -o ./bin/demo_darwin_arm64 .

build-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/demo_darwin_amd64 .

build-windows-amd64:
	env GOOS=windows GOARCH=amd64 go build -o ./bin/demo_windows_amd64.exe .

build-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -o ./bin/demo_linux_amd64 .
