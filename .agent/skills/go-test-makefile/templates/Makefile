# Makefile Template for Go Projects
# Change these variables as needed
BINARY_NAME=MyApp
BUILD_DIR=build
MAIN_FILE=main.go

.PHONY: all build clean test coverage lint deps help

all: build

# Help command to display available targets
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run unit tests"
	@echo "  make coverage - Generate coverage report"
	@echo "  make lint     - Run linter"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make deps     - Download dependencies"

# Build the application
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Run the application
run:
	go run $(MAIN_FILE)

# Cross-platform builds
build-all: build-windows build-linux build-mac

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_FILE)

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_FILE)

build-mac:
	@echo "Building for Mac..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-mac $(MAIN_FILE)

# Testing
test:
	@echo "Running tests..."
	go test ./... -v

# Code Coverage
coverage:
	@echo "Generating coverage report..."
	go test ./... -coverprofile=$(BUILD_DIR)/coverage.out
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report saved to $(BUILD_DIR)/coverage.html"

# Linting (requires golangci-lint)
lint:
	@echo "Linting..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Dependencies
deps:
	@echo "Downloading dependencies..."
	go mod tidy
	go mod download
