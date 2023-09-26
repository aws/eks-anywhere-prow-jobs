ifneq ($(filter true,$(CI) $(CODEBUILD_CI)),)
GOLANG_VERSION?="1.20"
GO_VERSION ?= $(shell source scripts/go.sh && get_go_path $(GOLANG_VERSION))
GO=$(GO_VERSION)/go
else
GO=$(shell which go)
endif
export PULL_BASE_SHA?=$(shell git show -s --format=%H upstream/main)
export PULL_PULL_SHA?=$(shell git show -s --format=%H)

BIN_DIR := bin
CONFIG_DIR := config

export CONSTANTS_CONFIG_FILE ?= $(CONFIG_DIR)/job-constants.yaml

.PHONY: default
default: build lint

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO) fmt scripts/lint_prowjobs/main.go

.PHONY: vet
vet: ## Run go vet against code.
	$(GO) vet scripts/lint_prowjobs/main.go

.PHONY: build
build: fmt vet ## Build linter
	$(GO) build -o ./$(BIN_DIR)/prow-linter github.com/aws/eks-anywhere-prow-jobs/scripts/lint_prowjobs

.PHONY: lint
lint: build ## Run linter
	./$(BIN_DIR)/prow-linter

.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[%\/0-9A-Za-z_-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
