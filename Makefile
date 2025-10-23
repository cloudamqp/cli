# CloudAMQP CLI Makefile

# Variables
BINARY_NAME=cloudamqp
GO_BUILD_FLAGS=-v
GO_TEST_FLAGS=-v
GO_LDFLAGS=

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
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .

# Development workflow
.PHONY: dev
dev: fmt vet test build

openapi.yaml:
	curl -O https://docs.cloudamqp.com/openapi.yaml

openapi-instance.yaml:
	curl -O https://docs.cloudamqp.com/openapi-instance.yaml

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests"
	@echo "  integration-test - Run integration tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  deps          - Install dependencies"
	@echo "  install       - Install the binary"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  dev           - Run development workflow (fmt, vet, test, build)"
	@echo "  help          - Show this help message"
