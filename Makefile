# CloudAMQP CLI Makefile

# Variables
BINARY_NAME=cloudamqp
GO_BUILD_FLAGS=-v
GO_TEST_FLAGS=-v

# Version information (automatically extracted from git)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%d")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build with version information
GO_LDFLAGS=-X cloudamqp-cli/cmd.Version=$(VERSION) \
           -X cloudamqp-cli/cmd.BuildDate=$(BUILD_DATE) \
           -X cloudamqp-cli/cmd.GitCommit=$(GIT_COMMIT)

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build $(GO_BUILD_FLAGS) -ldflags "$(GO_LDFLAGS)" -o $(BINARY_NAME) .

# Run tests
.PHONY: test
test:
	go test $(GO_TEST_FLAGS) ./...

# Run integration tests
.PHONY: integration-test
integration-test:
	go test $(GO_TEST_FLAGS) -tags=integration .

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	go clean

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Vet code
.PHONY: vet
vet:
	go vet ./...

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Install the binary
.PHONY: install
install: build
	go install .

# Build for multiple platforms
.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(GO_LDFLAGS)" -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(GO_LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(GO_LDFLAGS)" -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(GO_LDFLAGS)" -o $(BINARY_NAME)-windows-amd64.exe .

# Development workflow
.PHONY: dev
dev: fmt vet test build

# Show version information that will be used
.PHONY: version-info
version-info:
	@echo "Version Info:"
	@echo "  VERSION:    $(VERSION)"
	@echo "  BUILD_DATE: $(BUILD_DATE)"
	@echo "  GIT_COMMIT: $(GIT_COMMIT)"

openapi.yaml:
	curl -O https://docs.cloudamqp.com/openapi.yaml

openapi-instance.yaml:
	curl -O https://docs.cloudamqp.com/openapi-instance.yaml

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary with version info"
	@echo "  test          - Run tests"
	@echo "  integration-test - Run integration tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  deps          - Install dependencies"
	@echo "  install       - Install the binary"
	@echo "  build-all     - Build for multiple platforms with version info"
	@echo "  dev           - Run development workflow (fmt, vet, test, build)"
	@echo "  version-info  - Show version information that will be used in build"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Version information is automatically extracted from git."
	@echo "Override with: make build VERSION=1.0.0 BUILD_DATE=2025-11-25"
