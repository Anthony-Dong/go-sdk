export CGO_ENABLED:=1

.PHONY: all
all: test build

.PHONY: build
build:
	bash -x build.sh

.PHONY: test
test:
	make -C tcpdump/test un_compress
	go test -v -run=Test_DecodeTCPDump ./tcpdump > tcpdump.log
	echo "follow form console log...." >> tcpdump.log
	go test -v -run=TestNewConsulSource ./tcpdump >> tcpdump.log
	go test -coverprofile cover.out -count=1 ./...

.PHONY : fmt
fmt: ## fmt
	golangci-lint run  --fix -v