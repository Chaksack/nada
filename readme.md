# 🔍 Nada - Go Code Quality Analyzer

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/chaksack/nada)](https://goreportcard.com/report/github.com/chaksack/nada)

A comprehensive static code analysis tool for Go codebases, providing bug detection, security vulnerability scanning, code complexity measurement, and quality grading similar to SonarQube.

## ✨ Features

- **🔍 Static Code Analysis** - AST-based deep inspection of Go code
- **🐛 Bug Detection** - Identifies potential runtime issues and logic errors
- **🔒 Security Scanning** - Detects SQL injection, hardcoded secrets, and other vulnerabilities
- **📊 Quality Metrics** - Cyclomatic complexity, maintainability scoring, and overall grading (A-F)
- **🛠️ Multiple Interfaces** - CLI tool, REST API, and JSON export
- **⚙️ Configurable Rules** - Custom rule definitions via YAML configuration
- **🔄 Git Integration** - Analyze only changed files and lines
- **🧪 Coverage Integration** - Include test coverage metrics in analysis

## 🚀 Installation

### Install via Go

```bash
go install github.com/chaksack/nada@latest
```

### Build from Source

```bash
git clone https://github.com/chaksack/nada.git
cd nada
go build -o nada ./cmd/nada
```

### Using Make

```bash
git clone https://github.com/chaksack/nada.git
cd nada
make build
```

## 📋 Usage

### Basic Analysis

```bash
# Analyze current directory
nada analyze .

# Analyze specific project
nada analyze /path/to/project

# Export to JSON
nada analyze . --output report.json

# Analyze with coverage profile
nada analyze . --coverage coverage.out
```

### Git Integration

```bash
# Analyze only staged changes
nada analyze . --diff staged

# Analyze changes since main branch
nada analyze . --diff origin/main
```

### Web API Server

```bash
# Start server on default port (3000)
nada server

# Start server on custom port
nada server --port 8080
```

#### API Endpoints

```bash
# Health check
curl http://localhost:3000/health

# Analyze project
curl -X POST http://localhost:3000/analyze \
  -H "Content-Type: application/json" \
  -d '{"project_path": "/path/to/project"}'
```

## ⚙️ Configuration

Create a `nada.yaml` file in your project root to customize analysis rules:

```yaml
rules:
  - id: no-todo-comments
    description: "Disallow TODO comments in production code"
    pattern: "TODO|FIXME"
    severity: medium
    type: code_smell

  - id: max-function-length
    description: "Functions should not exceed 50 lines"
    max_lines: 50
    severity: low
    type: code_smell

  - id: no-hardcoded-secrets
    description: "Detect hardcoded passwords and API keys"
    pattern: "(password|api[_-]?key|secret)\\s*[:=]\\s*['\"][^'\"]{8,}"
    severity: high
    type: vulnerability

quality_gates:
  min_grade: C
  max_high_issues: 0
  max_medium_issues: 10
```

## 📊 Quality Grading

| Grade | Score | Description |
|-------|-------|-------------|
| A | 90-100 | Excellent code quality |
| B | 80-89 | Good code quality |
| C | 70-79 | Average code quality |
| D | 60-69 | Below average quality |
| F | <60 | Poor code quality |

## 🚦 CI/CD Integration

### GitHub Actions

```yaml
name: Code Quality Check
on: [push, pull_request]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install Nada
        run: go install github.com/chaksack/nada@latest
      
      - name: Run Tests with Coverage
        run: go test -coverprofile=coverage.out ./...
      
      - name: Run Quality Analysis
        run: nada analyze . --coverage coverage.out --output report.json
      
      - name: Upload Report
        uses: actions/upload-artifact@v4
        with:
          name: nada-report
          path: report.json
```

### GitLab CI

```yaml
stages:
  - test
  - analyze

test:
  stage: test
  image: golang:1.21
  script:
    - go test -coverprofile=coverage.out ./...
  artifacts:
    paths:
      - coverage.out

code_quality:
  stage: analyze
  image: golang:1.21
  dependencies:
    - test
  script:
    - go install github.com/chaksack/nada@latest
    - nada analyze . --coverage coverage.out --output report.json
  artifacts:
    reports:
      codequality: report.json
```

## 📈 Example Output

```
🎯 Nada Code Quality Report
===========================
📁 Project: example-service
⏰ Analyzed: 2024-01-15 14:30:45
📊 Grade: B (84.2/100)

📊 Metrics:
   📄 Files: 24
   📏 Lines: 3,247
   🔄 Avg Complexity: 3.8
   🧪 Coverage: 78.5%

🚨 Issues (25 total):
   🔴 High: 1
   🟡 Medium: 6
   🟢 Low: 18

🔍 Top Issues:
   auth.go:45    SQL injection vulnerability    [high]
   handler.go:123 Missing error handling        [medium]
   utils.go:78   Function complexity: 12        [medium]

🚪 Quality Gates: ✅ PASSED
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Make (optional)

### Development Commands

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Run with race detection
go test -race ./...

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Build for current platform
go build -o nada ./cmd/nada

# Cross-compile for multiple platforms
make build-all
```

### Project Structure

```
nada/
├── cmd/nada/           # CLI entry point
├── internal/           # Private application code
│   ├── analyzer/       # Core analysis engine
│   ├── rules/          # Rule definitions
│   ├── server/         # Web API server
│   └── config/         # Configuration handling
├── pkg/                # Public API
├── testdata/           # Test fixtures
├── Makefile
├── go.mod
└── README.md
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## 📊 Supported Analyses

### 🐛 Bug Detection
- Missing error handling
- Potential null pointer dereferences
- Type assertion without checks
- Resource leaks (unclosed files, connections)

### 🔒 Security Vulnerabilities
- SQL injection patterns
- Hardcoded secrets and API keys
- Weak cryptographic practices
- Command injection risks
- Path traversal vulnerabilities

### 💨 Code Smells
- Functions exceeding complexity thresholds
- Long parameter lists
- Deep nesting levels
- Poor naming conventions
- Missing documentation
- Code duplication

### 📏 Metrics
- Lines of code (physical and logical)
- Cyclomatic complexity
- Cognitive complexity
- Test coverage integration
- Technical debt estimation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [SonarQube](https://www.sonarqube.org/)'s code quality methodology
- Built with [Cobra](https://github.com/spf13/cobra) for CLI interface
- Web API powered by [Fiber](https://github.com/gofiber/fiber)
- Leverages Go's powerful AST parsing capabilities

## 📞 Support

- 📖 [Documentation](https://github.com/chaksack/nada/wiki)
- 🐛 [Issue Tracker](https://github.com/chaksack/nada/issues)
- 💬 [Discussions](https://github.com/chaksack/nada/discussions)

---

**Made with ❤️ for the Go community by [Andrew Chakdahah](https://github.com/chaksack)**