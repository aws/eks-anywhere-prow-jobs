GOLANG_VERSION?="1.16"
export PULL_BASE_SHA?=$(shell git show -s --format=%H main)
export PULL_PULL_SHA?=$(shell git show -s --format=%H)

BIN_DIR := bin

.PHONY: default
default: build lint

.PHONY: build
build:  ## Build linter
	go build -o ./$(BIN_DIR)/prow-linter github.com/aws/eks-anywhere-prow-jobs/scripts/lint_prowjobs

.PHONY: lint
lint: ## Run linter
	./$(BIN_DIR)/prow-linter
