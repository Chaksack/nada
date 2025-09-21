BINARY_NAME=nada
MAIN_PACKAGE=./cmd/nada
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Commit=$(COMMIT)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Directories
BUILD_DIR=build
DIST_DIR=dist
COVERAGE_DIR=coverage

.PHONY: all build clean test test-unit test-integration test-coverage test-race test-benchmark fmt lint deps install uninstall run-server run-analyze help

# Default target
all: clean fmt lint test build

# Build the application
build: deps
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all: deps
	@echo "Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	
	# macOS
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	
	# Windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	
	@echo "Cross-platform builds complete in $(DIST_DIR)/"

# Install the application
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

# Uninstall the application
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME)"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run all tests
test: test-unit

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"
	$(GOCMD) tool cover -func=$(COVERAGE_DIR)/coverage.out

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -v -race ./...

# Run benchmark tests
test-benchmark:
	@echo "Running benchmark tests..."
	$(GOTEST) -v -bench=. -benchmem ./...

# Run specific test
test-verbose:
	@echo "Running tests with verbose output..."
	$(GOTEST) -v -count=1 ./...

# Test specific package
test-package:
	@echo "Usage: make test-package PKG=./internal/analyzer"
	@if [ -z "$(PKG)" ]; then \
		echo "Please specify PKG variable"; \
		exit 1; \
	fi
	$(GOTEST) -v $(PKG)

# Generate test fixtures
test-fixtures:
	@echo "Generating test fixtures..."
	@mkdir -p testdata
	@echo 'package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello, World!")\n}' > testdata/simple.go
	@echo 'package main\n\nfunc Complex() {\n\tif true {\n\t\tif true {\n\t\t\tif true {\n\t\t\t\tif true {\n\t\t\t\t\tif true {\n\t\t\t\t\t\tfmt.Println("deep")\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t}\n}' > testdata/complex.go
	@echo "Test fixtures created in testdata/"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	$(GOCMD) mod tidy

# Lint code
lint:
	@echo "Running linters..."
	$(GOLINT) run ./...

# Clean build artifacts and test files
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	@rm -f coverage.out coverage.html
	@find . -name "*.test" -delete
	@find . -name "*.prof" -delete

# Run server for development
run-server: build
	@echo "Starting development server..."
	./$(BUILD_DIR)/$(BINARY_NAME) server --port 3000

# Run analysis on current directory
run-analyze: build
	@echo "Analyzing current directory..."
	./$(BUILD_DIR)/$(BINARY_NAME) analyze .

# Run analysis with JSON output
run-analyze-json: build
	@echo "Analyzing current directory with JSON output..."
	./$(BUILD_DIR)/$(BINARY_NAME) analyze . --output nada-report.json

# Run analysis with verbose output
run-analyze-verbose: build
	@echo "Analyzing current directory with verbose output..."
	./$(BUILD_DIR)/$(BINARY_NAME) analyze . --verbose

# Setup development environment
dev-setup:
	@echo "Setting up development environment..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOMOD) download

# Create release archives
release: build-all
	@echo "Creating release archives..."
	@mkdir -p $(DIST_DIR)/archives
	
	# Linux
	tar -czf $(DIST_DIR)/archives/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(DIST_DIR) $(BINARY_NAME)-linux-amd64
	tar -czf $(DIST_DIR)/archives/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz -C $(DIST_DIR) $(BINARY_NAME)-linux-arm64
	
	# macOS
	tar -czf $(DIST_DIR)/archives/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz -C $(DIST_DIR) $(BINARY_NAME)-darwin-amd64
	tar -czf $(DIST_DIR)/archives/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(DIST_DIR) $(BINARY_NAME)-darwin-arm64
	
	# Windows
	zip -j $(DIST_DIR)/archives/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe
	
	@echo "Release archives created in $(DIST_DIR)/archives/"

# Generate documentation
docs:
	@echo "Generating documentation..."
	@mkdir -p docs
	$(GOCMD) doc -all . > docs/API.md

# Verify the module
verify:
	@echo "Verifying module..."
	$(GOMOD) verify

# Security check
security:
	@echo "Running security checks..."
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; $(GOGET) -u github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; }
	gosec ./...

# Check for updates
check-updates:
	@echo "Checking for dependency updates..."
	$(GOCMD) list -u -m all

# Run static analysis on the project itself (dogfooding)
self-analyze: build
	@echo "Running Nada on itself..."
	./$(BUILD_DIR)/$(BINARY_NAME) analyze . --output self-analysis.json --verbose

# Profile CPU usage during analysis
profile-cpu: build
	@echo "Profiling CPU usage..."
	$(GOTEST) -cpuprofile=cpu.prof -bench=BenchmarkFullAnalysis ./...
	$(GOCMD) tool pprof cpu.prof

# Profile memory usage during analysis
profile-memory: build
	@echo "Profiling memory usage..."
	$(GOTEST) -memprofile=mem.prof -bench=BenchmarkFullAnalysis ./...
	$(GOCMD) tool pprof mem.prof

# Generate mocks (if needed)
generate-mocks:
	@echo "Generating mocks..."
	@command -v mockgen >/dev/null 2>&1 || { echo "Installing mockgen..."; $(GOGET) -u github.com/golang/mock/mockgen@latest; }
	# Add mockgen commands here if needed

# Run fuzzing tests
fuzz:
	@echo "Running fuzz tests..."
	$(GOTEST) -fuzz=. -fuzztime=30s ./...

# Continuous testing (watch mode)
test-watch:
	@echo "Running tests in watch mode..."
	@command -v entr >/dev/null 2>&1 || { echo "Please install 'entr' for watch mode"; exit 1; }
	find . -name "*.go" | entr -r make test-unit

# Quality gate check
quality-gate: test-coverage lint security
	@echo "Running quality gate checks..."
	@echo "✅ All quality gates passed!"

# Pre-commit hook setup
setup-hooks:
	@echo "Setting up git hooks..."
	@mkdir -p .git/hooks
	@echo '#!/bin/sh\nmake fmt lint test-unit' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Pre-commit hooks installed"

# Help
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  build-all          - Build for multiple platforms"
	@echo "  install            - Install the application"
	@echo "  uninstall          - Uninstall the application"
	@echo "  deps               - Download dependencies"
	@echo ""
	@echo "Testing:"
	@echo "  test               - Run all tests"
	@echo "  test-unit          - Run unit tests"
	@echo "  test-integration   - Run integration tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  test-race          - Run tests with race detection"
	@echo "  test-benchmark     - Run benchmark tests"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  test-package       - Test specific package (specify PKG=...)"
	@echo "  test-fixtures      - Generate test fixtures"
	@echo "  test-watch         - Run tests in watch mode"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo "  security           - Run security checks"
	@echo "  quality-gate       - Run all quality checks"
	@echo "  self-analyze       - Run Nada on itself"
	@echo ""
	@echo "Development:"
	@echo "  run-server         - Run development server"
	@echo "  run-analyze        - Analyze current directory"
	@echo "  run-analyze-json   - Analyze with JSON output"
	@echo "  run-analyze-verbose - Analyze with verbose output"
	@echo "  dev-setup          - Setup development environment"
	@echo "  setup-hooks        - Setup git pre-commit hooks"
	@echo ""
	@echo "Profiling:"
	@echo "  profile-cpu        - Profile CPU usage"
	@echo "  profile-memory     - Profile memory usage"
	@echo "  fuzz               - Run fuzz tests"
	@echo ""
	@echo "Release:"
	@echo "  release            - Create release archives"
	@echo "  docs               - Generate documentation"
	@echo "  verify             - Verify module"
	@echo "  check-updates      - Check for dependency updates"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean              - Clean build artifacts"
	@echo "  help               - Show this help"