GOLANG_VERSION?="1.16"
export PULL_BASE_SHA?=$(shell git show -s --format=%H main)
export PULL_PULL_SHA?=$(shell git show -s --format=%H)

BIN_DIR := bin
CONFIG_DIR := config

export CONSTANTS_CONFIG_FILE ?= $(CONFIG_DIR)/job-constants.yaml

.PHONY: default
default: build lint

.PHONY: build
build:  ## Build linter
	go build -o ./$(BIN_DIR)/prow-linter github.com/aws/eks-anywhere-prow-jobs/scripts/lint_prowjobs

.PHONY: lint
lint: ## Run linter
	./$(BIN_DIR)/prow-linter

.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[%\/0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
