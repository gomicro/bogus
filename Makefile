SHELL = bash
APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')
DATE := $(shell date -u +%Y-%m-%d%Z%H:%M:%S)
CI_BUILD_NUMBER ?= dev
BUILD_NUMBER ?= $(CI_BUILD_NUMBER)
BUILD_VERSION = $(VERSION)-$(BUILD_NUMBER)
CI_COMMIT ?= $(shell git rev-parse HEAD)
GIT_COMMIT_HASH ?= $(CI_COMMIT)
NO_COLOR := \033[0m
INFO_COLOR := \033[0;36m
VERSION = 0.0.1
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test -v $(shell $(GOCMD) list ./... | grep -v /vendor/)
GOFMT = go fmt
CGO_ENABLED ?= 0
GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')


.PHONY: all
all: test

.PHONY: clean
clean: ## Clean out all generated items
	-@$(GOCLEAN)

.PHONY: coverage
coverage: ## Generates the total code coverage of the project
	@$(eval COVERAGE_DIR=$(shell mktemp -d))
	@mkdir -p $(COVERAGE_DIR)/tmp
	@for j in $$(go list ./... | grep -v '/vendor/' | grep -v '/ext/'); do go test -covermode=count -coverprofile=$(COVERAGE_DIR)/$$(basename $$j).out $$j > /dev/null 2>&1; done
	@echo 'mode: count' > $(COVERAGE_DIR)/tmp/full.out
	@tail -q -n +2 $(COVERAGE_DIR)/*.out >> $(COVERAGE_DIR)/tmp/full.out
	@$(GOCMD) tool cover -func=$(COVERAGE_DIR)/tmp/full.out | tail -n 1 | sed -e 's/^.*statements)[[:space:]]*//' -e 's/%//'

.PHONY: fmt
fmt: ## Run go fmt
	$(GOFMT)

.PHONY: help
help: ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: test
test: unit_test ## Run all available tests

.PHONY: unit_test
unit_test: ## Run all available unit tests
	$(GOTEST)
