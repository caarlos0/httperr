SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

setup:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s
	go mod tidy
.PHONY: setup

build:
	go build
.PHONY: build

test:
	go test $(TEST_OPTIONS) -failfast -race -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m
.PHONY: test

cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

lint:
	golangci-lint run ./...
.PHONY: lint

ci: build test lint
.PHONY: ci

.DEFAULT_GOAL := ci
