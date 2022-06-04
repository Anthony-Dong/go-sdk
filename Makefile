# #######################################################
# Function :Makefile for go                             #
# Platform :All Linux Based Platform                    #
# Version  :1.0                                         #
# Date     :2020-12-17                                  #
# Author   :fanhaodong516@gmail.com                     #
# Usage    :make help									#
# #######################################################

# dir
PROJECT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# go env
export GO111MODULE := on
export GOPROXY := https://goproxy.cn,direct
export GOPRIVATE :=
export GOFLAGS :=

# PHONY
.PHONY : all init build build_tool fmt build_cmd check deploy test help release

all: fmt build_tool ## Let's go!

init: ## init project and init env
	go mod tidy
	@if [ ! -e $(shell go env GOPATH)/bin/golangci-lint ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1; fi

build_tool: ## build gtool
	@cd gtool; go build -v -ldflags "-s -w"  -o bin/gtool main.go; cd -; mkdir -p bin ;mv gtool/bin/gtool bin/gtool

build: ## cross compiling
	@cd gtool; sh -x build.sh; cd -

fmt: ## fmt
	@for elem in `find . -name '*.go' | grep -v 'internal/pkg'`;do goimports -w $$elem; gofmt -w $$elem; done

deploy: fmt test check build_tool ## deploy this project

check: ## check custom rule
	go run internal/cmd/clear.go

test: ## go test
	go test -coverprofile cover.out -count=1 ./...
	go tool cover -html=cover.out

release: ## release new version
	@for elem in `find . -name '*.md'`;do sed -i 's/v1.0.0/v1.0.1/g' $$elem ; done

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)