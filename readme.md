# 🔍 Nada - Go Code Quality Analyzer

A comprehensive code quality analysis tool similar to SonarQube, built specifically for Go codebases. Nada provides static code analysis, detects bugs and vulnerabilities, measures code complexity, and assigns quality grades.

## ✨ Features

### 🔍 Static Code Analysis
- **AST-based Analysis** - Deep inspection of Go syntax trees
- **Cyclomatic Complexity** - Function complexity measurement
- **Code Smell Detection** - Anti-patterns and maintainability issues
- **Bug Detection** - Potential runtime issues and logic errors
- **Security Vulnerabilities** - SQL injection, hardcoded secrets, etc.
- **Naming Conventions** - Go-specific naming standard checks

### 📊 Quality Metrics
- Lines of code counting
- Test coverage estimation
- Technical debt assessment
- Maintainability scoring
- Overall quality grading (A-F)

### 🛠️ Multiple Interfaces
- **CLI Tool** - Command-line interface with Cobra
- **Web API** - REST API with Go Fiber
- **JSON Export** - Detailed reports in JSON format

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd nada

# Install dependencies
go mod download

# Build the tool
make build

# Or install directly
make install
```

### Basic Usage

```bash
# Analyze a Go project
./nada analyze /path/to/your/go/project

# Analyze current directory
./nada analyze .

# Export results to JSON
./nada analyze . --output report.json

# Start web server
./nada server --port 3000
```

## 📋 Commands

### CLI Commands

| Command | Description | Options |
|---------|-------------|---------|
| `analyze [path]` | Analyze Go codebase | `--output, -o` - JSON output file |
| `server` | Start web API server | `--port, -p` - Server port (default: 3000) |

### Web API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Server information |
| `/analyze` | POST | Analyze project via API |
| `/health` | GET | Health check |

#### API Usage Example

```bash
# Analyze via API
curl -X POST http://localhost:3000/analyze \
  -H "Content-Type: application/json" \
  -d '{"project_path": "/path/to/project"}'
```

## 🔍 Analysis Categories

### 🐛 Bugs
- Missing error handling
- Potential null pointer dereferences
- Type mismatches
- Logic errors

### 🔒 Vulnerabilities
- SQL injection patterns
- Hardcoded secrets/passwords
- Insecure cryptographic practices
- Path traversal vulnerabilities

### 💨 Code Smells
- Long functions (>50 lines)
- High cyclomatic complexity (>10)
- Deep nesting (>4 levels)
- Poor naming conventions
- Missing documentation
- TODO/FIXME comments

## 📊 Quality Grading

Nada assigns grades based on:

- **A (90-100)** - Excellent code quality
- **B (80-89)** - Good code quality
- **C (70-79)** - Average code quality
- **D (60-69)** - Below average quality
- **F (<60)** - Poor code quality

### Scoring Factors
- Issue severity and count
- Cyclomatic complexity
- Test coverage
- Code duplication
- Security vulnerabilities

## 🛠️ Development

### Prerequisites
- Go 1.21 or higher
- Make (optional, for convenience commands)

### Building from Source

```bash
# Install dependencies
make deps

# Run tests
make test

# Format code
make fmt

# Build for current platform
make build

# Build for all platforms
make build-all
```

### Development Commands

```bash
# Run server in development
make run-server

# Analyze current directory
make run-analyze

# Generate JSON report
make run-analyze-json

# Clean build artifacts
make clean
```

## 📈 Example Output

```
🎯 Code Quality Analysis Report
================================
📁 Project: ./example-project
⏰ Analyzed at: 2025-01-15 10:30:45
📊 Overall Grade: B (82.5/100)
📄 Files Analyzed: 15
📏 Lines of Code: 2,847
🔄 Avg Cyclomatic Complexity: 4.2
🧪 Test Coverage: 67.3%

📋 Issues Summary:
   🔴 High: 2
   🟡 Medium: 8
   🟢 Low: 15
   🐛 Bugs: 3
   🔒 Vulnerabilities: 1
   💨 Code Smells: 21

⚠️  Top Issues:
   main.go:45 - Potential SQL injection [vulnerability/high]
   handler.go:123 - Missing error handling [bug/medium]
   utils.go:78 - Function too complex (complexity: 12) [code_smell/medium]

🚪 Quality Gates:
   Grade A-C: ✅ PASSED
   High Issues < 5: ✅ PASSED
   No Vulnerabilities: ❌ FAILED
```

## 🔧 Configuration


Nada supports custom rule configuration files (`nada.yaml` or `.nada.yaml`) in your project root. You can define your own static analysis rules using YAML.

### Example: nada.yaml

```yaml
rules:
   - id: no-todo-comments
      description: "Disallow TODO comments in code."
      pattern: "TODO"
      severity: "medium"
      type: "code_smell"

   - id: max-function-length
      description: "Warn if a function exceeds 40 lines."
      max_lines: 40
      severity: "low"
      type: "code_smell"

   - id: no-hardcoded-password
      description: "Disallow hardcoded password strings."
      pattern: "password"
      severity: "high"
      type: "vulnerability"
```

Supported fields:
- `id`: Unique rule identifier
- `description`: Description of the rule
- `pattern`: Regex pattern to match in code (optional)
- `max_lines`: Maximum allowed lines for a function (optional)
- `severity`: Issue severity (`low`, `medium`, `high`)
- `type`: Issue type (`code_smell`, `vulnerability`, etc)

Nada will automatically load and apply these rules during analysis if the config file is present.

## 🚦 CI/CD Integration

You can run Nada as part of your CI/CD pipeline to enforce code quality automatically. Below are examples for GitHub Actions and GitLab CI:

### GitHub Actions Example

Create a workflow file at `.github/workflows/nada.yml`:

```yaml
name: Nada Code Quality
on: [push, pull_request]
jobs:
   analyze:
      runs-on: ubuntu-latest
      steps:
         - uses: actions/checkout@v3
         - name: Set up Go
            uses: actions/setup-go@v5
            with:
               go-version: '1.21'
         - name: Install Nada
            run: |
               make build
         - name: Run Nada Analysis
            run: |
               ./nada analyze . --output report.json
         - name: Upload Report Artifact
            uses: actions/upload-artifact@v4
            with:
               name: nada-report
               path: report.json
```

### GitLab CI Example

Add to your `.gitlab-ci.yml`:

```yaml
stages:
   - analyze

nada_analysis:
   stage: analyze
   image: golang:1.21
   script:
      - go mod download
      - make build
      - ./nada analyze . --output report.json
   artifacts:
      paths:
         - report.json
```

**Tip:** You can fail the pipeline if the report contains high-severity issues by parsing `report.json` in a custom script step.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make test` and `make fmt`
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Inspired by SonarQube's code quality methodology
- Built with [Cobra CLI](https://github.com/spf13/cobra)
- Web server powered by [Go Fiber](https://github.com/gofiber/fiber)
- Go's powerful AST parsing capabilities

## 🔮 Future Enhancements

- [ ] Web dashboard interface

## 🖥️ Web Dashboard Interface (Planned)

Nada will offer a web dashboard for interactive code quality exploration and reporting. This dashboard will be served by the built-in web API and provide:

- Project overview and quality grade
- Interactive issue explorer (filter by type, severity, file)
- Visualizations: grade history, complexity, coverage, and more
- File-by-file and function-by-function drilldown
- Downloadable and shareable reports
- Live updates as code changes (with server running)

### Example Usage (future)

```bash
# Start the dashboard server
./nada server --port 3000
# Open your browser to http://localhost:3000/dashboard
```

### Planned Endpoints

- `/dashboard` — Main dashboard UI
- `/api/report` — Get latest analysis report (JSON)
- `/api/issues` — List/filter issues
- `/api/metrics` — Project metrics and trends

**Feedback and feature requests are welcome!**
- [ ] Support for additional languages
- [ ] Performance benchmarking
- [ ] IDE plugins

### 🟦 Git Integration for Diff Analysis

Nada can analyze only the code changes in your Git repository, making it ideal for code reviews, pull requests, and CI pipelines. This feature restricts static analysis to changed files and lines, so you can focus on new and modified code.

**Usage:**

#### CLI

```bash
# Analyze only staged changes (pre-commit)
./nada analyze . --diff=staged

# Analyze only unstaged changes
./nada analyze . --diff=unstaged

# Analyze changes since last commit
./nada analyze . --diff=HEAD

# Analyze changes compared to a remote branch (e.g., for PRs)
./nada analyze . --diff=origin/main
```

#### Web API

```json
{
   "project_path": ".",
   "diff": "HEAD"
}
```

**How it works:**
- Nada detects if the target path is a Git repo.
- Runs `git diff` (with the appropriate ref) to get changed files and line ranges.
- Restricts static analysis to only those files/lines.
- Reports issues only for changed code.

**Output:**
- Same as normal, but only for changed code.
- Optionally, issues can be annotated as "new in diff" for PR review.

**Example Output:**

```
main.go:45 - Potential SQL injection [vulnerability/high] [diff]
utils.go:78 - Function too complex (complexity: 12) [code_smell/medium] [diff]
```

**Supported diff targets:**
- `staged`, `unstaged`, `HEAD`, any branch or commit (e.g., `origin/main`)

This enables focused code review and CI workflows, ensuring only new and changed code is flagged.

## 🧪 Code Coverage Integration

Nada supports Go code coverage profile integration. You can supply a coverage profile (from `go test -coverprofile=coverage.out`) to include detailed coverage metrics in your analysis reports and dashboard.

### CLI Usage

```bash
# Run Go tests with coverage
 go test -coverprofile=coverage.out ./...

# Analyze project with coverage
 ./nada analyze . --coverage=coverage.out
```

### Web API Usage

Send the `coverage` field in your POST request:

```json
{
  "project_path": ".",
  "coverage": "coverage.out"
}
```

### Dashboard UI

The dashboard Overview and Metrics views will display:
- Test coverage percent
- Coverage statements and covered statements

### Output Example

```
🧪 Test Coverage: 67.3%
Statements: 1200 | Covered: 808
```

**Tip:** Use coverage integration in CI to ensure new code is tested!

## 🚦 Performance Benchmarking

Nada now measures and reports analysis performance metrics:

- **Analysis Duration (ms):** Total time taken to analyze the project.
- **Files Per Second:** Throughput of the analyzer (files analyzed per second).

These metrics are included in the CLI, API, and dashboard output for every analysis run.

### Example Output

```
Analysis Duration: 512.4 ms
Files Per Second: 29.3
```

You can use these metrics to:
- Track analysis speed as your project grows
- Benchmark Nada on different hardware or CI environments
- Detect regressions in analyzer performance

Performance metrics are available in the `metrics` section of the JSON report and API response.