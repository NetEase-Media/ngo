SHELL = /bin/bash

SCRIPT_DIR = $(shell pwd)/etc/script
PKG_LIST   = $(shell go list ./... | grep -v /vendor/ | grep -v /examples)

lint_code: ## Lint the files

dep: ## Get the dependencies
	go mod download

code_coverage: dep ## Generate global code coverage report
	sh ${SCRIPT_DIR}/coverage.sh

code_coverage_html: dep
	sh ${SCRIPT_DIR}/coverage.sh html;

race_detector: dep ## Run data race detector
	go test -gcflags=-l -race -short ${PKG_LIST}

unit_tests: dep ## Run unittests
	go test -gcflags=-l -v ${PKG_LIST}

.PHONY: lint_code dep code_coverage code_coverage_html race_detector unit_tests