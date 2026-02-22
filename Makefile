# Hello World Plugin — build, test, lint.
# Standalone repo; requires github.com/opentalon/opentalon (pkg/plugin).

.PHONY: build test lint

BINARY_NAME ?= hello-world-plugin

build:
	go build -o $(BINARY_NAME) .
	@echo "Built: $(BINARY_NAME)"

test:
	go test -race -count=1 ./...

lint:
	golangci-lint run
