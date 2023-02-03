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
.PHONY : all
all: build ## Let's go!

init: ## init project
	@if [ ! -e $(shell go env GOPATH)/bin/golangci-lint ]; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1; fi

.PHONY: ci
ci: check test build ## ci

.PHONY : build
build: ## gtool build
	make -C gtool build
	rm -rf bin; cp -r gtool/bin bin

.PHONY : fmt
fmt: ## fmt
	golangci-lint run  --fix -v
	make -C gtool fmt

.PHONY : check
check: ## check custom rule
	go run internal/cmd/clear.go

.PHONY : test
test: ## go test
	go test -coverprofile cover.out -count=1 ./...
	make -C gtool test

.PHONY : release
release: ## release new version
	@for elem in `find . -name '*.md'`; do sed -i 's/1.0.4/1.0.5/g' $$elem ; done

.PHONY : help
help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)