##
# Build vars
#
BUILD_DIR ?= ./build
SHELL=/bin/bash

##
# Dependencies
#
GOLANGCI_VER = v1.33.0
GOTESTSUM_VER = 0.4.2
OS = $(shell uname | tr A-Z a-z)

.PHONY: run
run: ## Run ./cmd/uampnotif
	go run ./cmd/uampnotif $(ARGS)

build:
	@mkdir -p ${BUILD_DIR}

clean: ## Clean up build directory
	rm -rf ${BUILD_DIR}

.PHONY: build-linux
build-linux: build ## Build targeting for linux
	env GOOS=linux GOARCH=arm64 go build -o build/uampnotif-linux ./cmd/uampnotif

.PHONY: build-mac
build-mac: build ## Build targeting for Mac
	go build -o ${BUILD_DIR}/uampnotif ./cmd/uampnotif
	ls -la ${BUILD_DIR}

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VER}
	@ln -sf gotestsum-${GOTESTSUM_VER} bin/gotestsum
bin/gotestsum-${GOTESTSUM_VER}:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VER}/gotestsum_${GOTESTSUM_VER}_${OS}_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VER} && chmod +x ./bin/gotestsum-${GOTESTSUM_VER}

bin/golangci-lint:
	mkdir -p ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin ${GOLANGCI_VER}

verify: bin/golangci-lint ## Run static code analyzers
	bin/golangci-lint run

.PHONY: test
test: FMT ?= "standard-verbose"
test: bin/gotestsum ## Run whole test suite
	bin/gotestsum -f $(FMT) $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
