# SPDX-License-Identifier: Apache-2.0
# SPDX-FileCopyrightText: 2025 Scott Friedman and Project Contributors

.PHONY: help build install test lint clean run fmt vet tidy

# Variables
BINARY_NAME=globus-connect-server
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/globus-connect-server
	@echo "Build complete: ./$(BINARY_NAME)"

## install: Install the binary to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) ./cmd/globus-connect-server
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Coverage report: coverage.out"

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  brew install golangci-lint  # macOS"; \
		echo "  or visit https://golangci-lint.run/usage/install/"; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## tidy: Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	go mod tidy
	@echo "go.mod tidied"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

## run: Build and run the binary
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

## dev: Run in development mode (with verbose output)
dev: build
	@echo "Running $(BINARY_NAME) in development mode..."
	./$(BINARY_NAME) --verbose

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	@echo "Dependencies downloaded"

## verify: Run all checks (fmt, vet, lint, test)
verify: fmt vet lint test
	@echo "All checks passed!"

## build-all: Build for all platforms
build-all:
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/globus-connect-server
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/globus-connect-server
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/globus-connect-server
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/globus-connect-server
	@echo "Build complete: ./dist/"

## version: Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Built: $(DATE)"
