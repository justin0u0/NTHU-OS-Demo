VERSION = $(shell git tag --points-at HEAD)
COMMIT = $(shell git rev-parse HEAD)

LDFLAGS = -ldflags '-X "github.com/justin0u0/NTHU-OS-Demo/version.Version=$(VERSION)" -X "github.com/justin0u0/NTHU-OS-Demo/version.Commit=$(COMMIT)"'

all: build

.PHONY: build-all
build-all: build build-darwin-arm64 build-darwin-amd64 build-windows-amd64 build-linux-amd64

.PHONY: build
build:
	go build $(LDFLAGS) -o ./bin/demo .

.PHONY: clean
clean:
	rm -f ./bin/**

build-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ./bin/demo_darwin_arm64 .

build-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o ./bin/demo_darwin_amd64 .

build-windows-amd64:
	env GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./bin/demo_windows_amd64.exe .

build-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./bin/demo_linux_amd64 .
