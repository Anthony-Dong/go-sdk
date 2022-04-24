#!/usr/bin/env bash
set -ex

DIR=$(cd "$(dirname "${0}")" || exit 1; pwd)
rm -rf bin/{darwin,linux,windows}
mkdir -p bin/{darwin,linux,windows}


CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -ldflags "-s -w" -o bin/darwin/gtool main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags "-s -w" -o bin/linux/gtool main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -ldflags "-s -w" -o bin/windows/gtool main.go