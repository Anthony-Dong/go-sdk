#!/usr/bin/env bash
set -ex

DIR=$(cd "$(dirname "${0}")" || exit 1; pwd)
ARCH=$(go env GOHOSTARCH)

rm -rf bin/darwin bin/linux
mkdir -p bin/darwin bin/linux

case $(uname) in
  "Darwin")
  CGO_ENABLED=1 go build -v -ldflags "-s -w" -o bin/darwin/gtool-${ARCH} main.go
  cp -f bin/darwin/gtool-${ARCH} bin/gtool
  ;;
  "Linux")
  CGO_ENABLED=1 go build -v -ldflags "-s -w" -o bin/linux/gtool-${ARCH} main.go
  cp -f bin/linux/gtool-${ARCH} bin/gtool
  ;;
  "*")
  echo "Not Support!"
esac

# 不支持windows
# GOOS=windows GOARCH=amd64 go build -v -ldflags "-s -w" -o bin/windows/gtool main.go