##
# Build vars
#
BUILD_DIR ?= build
SHELL=/bin/bash

##
# Dependencies
#
GOTESTSUM_VERSION = 0.4.2
OS = $(shell uname | tr A-Z a-z)

.PHONY: run
run:
	go run cmd/uampnotif/main.go

build:
	@mkdir -p build

.PHONY: build-linux
build-linux: build
	env GOOS=linux GOARCH=arm64 go build -o build/uampnotif-linux cmd/uampnotif/main.go

.PHONY: build-mac
build-mac: build
	env GOOS=darwin GOARCH=386 go build -o build/uampnotif cmd/uampnotif/main.go

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum
bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_${OS}_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}

.PHONY: test
test: FMT ?= "standard-verbose"
test: bin/gotestsum
	bin/gotestsum -f $(FMT) $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
