# Nada Makefile

.PHONY: build clean test install run-server run-analyze help

# Binary name
BINARY_NAME=nada
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) -v

# Build for Linux
build-linux:
	@echo "🐧 Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_UNIX) -v

# Build for Windows
build-windows:
	@echo "🪟 Building for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BINARY_WINDOWS) -v

# Build for all platforms
build-all: build-linux build-windows build

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Install the binary
install: build
	@echo "📲 Installing $(BINARY_NAME)..."
	go install

# Run the server
run-server:
	@echo "🌐 Starting server..."
	go run main.go server --port 3000

# Analyze current directory
run-analyze:
	@echo "🔍 Analyzing current directory..."
	go run main.go analyze .

# Analyze with output
run-analyze-json:
	@echo "🔍 Analyzing current directory with JSON output..."
	go run main.go analyze . --output report.json

# Format code
fmt:
	@echo "✨ Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Show help
help:
	@echo "Nada - Go Code Quality Analyzer"
	@echo ""
	@echo "Available commands:"
	@echo "  build         Build the application"
	@echo "  build-linux   Build for Linux"
	@echo "  build-windows Build for Windows"
	@echo "  build-all     Build for all platforms"
	@echo "  clean         Clean build artifacts"
	@echo "  test          Run tests"
	@echo "  deps          Install dependencies"
	@echo "  install       Install the binary"
	@echo "  run-server    Start web server on port 3000"
	@echo "  run-analyze   Analyze current directory"
	@echo "  run-analyze-json  Analyze with JSON output"
	@echo "  fmt           Format code"
	@echo "  lint          Run linter"
	@echo "  help          Show this help message"

# Default target
all: deps build