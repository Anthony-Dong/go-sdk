# go env
export GO111MODULE := on
export GOPRIVATE :=
export GOFLAGS :=
export CGO_ENABLED := 1

.PHONY: ci
ci: check fmt build

.PHONY: init
init: ## init project
	@if [ ! -e $(shell go env GOPATH)/bin/golangci-lint ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1; fi

.PHONY: build
build:
	rm -rf bin
	go build -v -ldflags "-s -w" -o bin/gtool gtool/main.go
	bin/gtool --version

.PHONY: install
install:
	go build -v -ldflags "-s -w" -o $$(go env GOPATH)/bin/gtool  gtool/main.go
	$$(go env GOPATH)/bin/gtool --version

.PHONY: fmt
fmt:
	golangci-lint run  --fix -v

.PHONY: check
check:
	go run internal/cmd/clear.go

.PHONY: test
test: ## go tool cover -html=cover.out
	make -C gtool/tcpdump/test un_compress
	go test -coverprofile cover.out -count=1 ./...

.PHONY: help
help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)