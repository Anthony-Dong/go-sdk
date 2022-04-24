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

# go test
GO_TEST_PKG_NAME := $(shell go list ./...)

# go env
export GO111MODULE := on
export GOPROXY := https://goproxy.cn,direct
export GOPRIVATE :=
export GOFLAGS :=

# PHONY
.PHONY : all init build fmt build_cmd deploy test help

all: fmt build_tool ## Let's go!

init: ## init project and init env
	go mod tidy
	@if [ ! -e $(shell go env GOPATH)/bin/golangci-lint ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1; fi

build_tool: ## build this project
	@cd gtool; go build -v -ldflags "-s -w"  -o bin/gtool main.go; cd -; mkdir -p bin ;mv gtool/bin/gtool bin/gtool

fmt: ## fmt
	@for elem in `find . -name '*.go' | grep -v 'internal/pkg'`;do goimports -w $$elem; gofmt -w $$elem; done

deploy: fmt build_tool clear_tool

clear_tool:
	go run internal/cmd/clear.go

test: ## go test
	go test -count=1 ./...

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)