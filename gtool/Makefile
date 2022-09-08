.PHONY: all test build
export CGO_ENABLED:=1

all: test build

install:
	go build -v -ldflags "-s -w" -o bin/gtool main.go

build:
	bash -x build.sh

test:
	go test -v -run=Test_DecodeTCPDump ./tcpdump > tcpdump.log
	echo "follow form console log...." >> tcpdump.log
	go test -v -run=TestNewConsulSource ./tcpdump >> tcpdump.log
	go test -coverprofile cover.out -count=1 ./...
	#go tool cover -html=cover.out