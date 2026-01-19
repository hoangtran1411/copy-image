# Makefile for copy-image project
# ================================

# Variables
BINARY_NAME=copyimage
BINARY_PATH=./cmd/copyimage
BUILD_DIR=./build
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOCLEAN=$(GOCMD) clean
GOLINT=golangci-lint

# Build flags
LDFLAGS=-ldflags "-s -w"

# Phony targets
.PHONY: all build clean test test-verbose coverage coverage-html lint fmt vet tidy run help

# Default target
all: tidy fmt vet test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME).exe $(BINARY_PATH)
	@echo "Build complete: $(BINARY_NAME).exe"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	set GOOS=windows&& set GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(BINARY_PATH)
	set GOOS=linux&& set GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(BINARY_PATH)
	set GOOS=darwin&& set GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(BINARY_PATH)
	@echo "Build complete for all platforms"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN) -testcache
	-del /f $(BINARY_NAME).exe 2>nul
	-rmdir /s /q $(BUILD_DIR) 2>nul
	-del /f $(COVERAGE_FILE) 2>nul
	-del /f $(COVERAGE_HTML) 2>nul
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) ./...
	@echo "Tests complete"

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	$(GOTEST) -v ./...
	@echo "Tests complete"

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOCLEAN) -testcache
	$(GOTEST) "-coverprofile=$(COVERAGE_FILE)" "-covermode=atomic" ./...
	$(GOCMD) tool cover "-func=$(COVERAGE_FILE)"
	@echo "Coverage report generated"

# Generate HTML coverage report
coverage-html: coverage
	@echo "Generating HTML coverage report..."
	$(GOCMD) tool cover "-html=$(COVERAGE_FILE)" -o $(COVERAGE_HTML)
	@echo "HTML report: $(COVERAGE_HTML)"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	$(GOLINT) run ./...
	@echo "Lint complete"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Format complete"

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet complete"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy
	@echo "Tidy complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	@echo "Dependencies downloaded"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME).exe

# Run with dry-run mode
run-dry:
	@echo "Running $(BINARY_NAME) in dry-run mode..."
	$(GOCMD) run $(BINARY_PATH) --dry-run --interactive=false

# Install to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install $(BINARY_PATH)
	@echo "Install complete"

# Show help
help:
	@echo ""
	@echo "Copy Image Tool - Makefile Commands"
	@echo "===================================="
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build targets:"
	@echo "  build       - Build the binary for current platform"
	@echo "  build-all   - Build for Windows, Linux, and macOS"
	@echo "  clean       - Remove build artifacts"
	@echo "  install     - Install to GOPATH/bin"
	@echo ""
	@echo "Test targets:"
	@echo "  test        - Run tests"
	@echo "  test-verbose- Run tests with verbose output"
	@echo "  coverage    - Run tests with coverage report (>70%%)"
	@echo "  coverage-html - Generate HTML coverage report"
	@echo ""
	@echo "Code quality:"
	@echo "  fmt         - Format code"
	@echo "  vet         - Run go vet"
	@echo "  lint        - Run golangci-lint"
	@echo ""
	@echo "Dependencies:"
	@echo "  tidy        - Tidy go.mod"
	@echo "  deps        - Download dependencies"
	@echo ""
	@echo "Run targets:"
	@echo "  run         - Build and run the application"
	@echo "  run-dry     - Run in dry-run mode"
	@echo ""
	@echo "  all         - Tidy, fmt, vet, test, build"
	@echo "  help        - Show this help message"
	@echo ""
