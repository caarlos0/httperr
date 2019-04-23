SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

export GO111MODULE := on

setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod download
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
	./bin/golangci-lint run --tests=false --enable-all --disable=lll ./...
.PHONY: lint

ci: build test lint
.PHONY: ci

card:
	wget -O card.png -c "https://og.caarlos0.dev/**httperr**:%20better%20handle%20errors%20on%20HTTP%20handlers.png?theme=light&md=1&fontSize=100px&images=https://github.com/caarlos0.png"
.PHONY: card

.DEFAULT_GOAL := ci
