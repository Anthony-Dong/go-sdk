#!/usr/bin/env bash
set -ex

DIR=$(cd "$(dirname "${0}")" || exit 1; pwd)

rm -rf bin/darwin bin/linux
mkdir -p bin/darwin bin/linux

case $(uname) in
  "Darwin")
  CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -v -ldflags "-s -w" -o bin/darwin/gtool-amd64 main.go
  CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -v -ldflags "-s -w" -o bin/darwin/gtool-arm64 main.go
  cp bin/darwin/gtool-amd64 bin/gtool
  ;;
  "Linux")
  CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -ldflags "-s -w" -o bin/linux/gtool-amd64 main.go
  # CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -v -ldflags "-s -w" -o bin/linux/gtool-arm64 main.go # cgo arm64 编译有问题
  cp bin/linux/gtool-amd64 bin/gtool
  ;;
  "*")
  echo "Not Support!"
esac

# 不支持windows
# GOOS=windows GOARCH=amd64 go build -v -ldflags "-s -w" -o bin/windows/gtool main.go