.PHONY: all test build
export CGO_ENABLED:=1

all: test build

build:
	bash -x build.sh

test:
	go test -v -run=Test_DecodeTCPDump ./tcpdump > tcpdump.log
	go test -coverprofile cover.out -count=1 ./...
	go tool cover -html=cover.out