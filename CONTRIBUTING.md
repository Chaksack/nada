# Contributing to Nada

Thank you for your interest in contributing to Nada! We welcome contributions from the community and are excited to see what you'll bring to the project.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/chaksack/nada.git
   cd nada
   ```

3. Add the upstream repository as a remote:
   ```bash
   git remote add upstream https://github.com/chaksack/nada.git
   ```

4. Install development dependencies:
   ```bash
   make dev-setup
   ```

5. Verify your setup by running tests:
   ```bash
   make test
   ```

## 🔄 Development Workflow

### Creating a Feature Branch

```bash
# Fetch the latest changes from upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create a new feature branch
git checkout -b feature/your-feature-name
```

### Making Changes

1. Make your changes in the feature branch
2. Follow the coding standards (see below)
3. Write or update tests for your changes
4. Update documentation if necessary

### Testing Your Changes

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Run security checks
make security

# Format your code
make fmt
```

### Committing Your Changes

We follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages:

```bash
# Examples of good commit messages
git commit -m "feat: add support for custom rule configurations"
git commit -m "fix: resolve null pointer dereference in analyzer"
git commit -m "docs: update installation instructions"
git commit -m "test: add unit tests for complexity calculator"
```

**Commit Types:**
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

### Submitting Your Pull Request

1. Push your feature branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Create a pull request from your fork to the main repository
3. Fill out the pull request template completely
4. Ensure all CI checks pass
5. Address any feedback from reviewers

## 📋 Coding Standards

### Go Code Guidelines

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Write comprehensive comments for exported functions
- Keep functions focused and small (ideally < 50 lines)
- Use interfaces to define behavior
- Handle errors explicitly
- Write table-driven tests where appropriate

### Code Quality Checklist

- [ ] Code follows Go best practices
- [ ] All new code has corresponding tests
- [ ] Tests pass locally (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] No linting errors (`make lint`)
- [ ] Documentation is updated if necessary
- [ ] Security checks pass (`make security`)

## 🧪 Testing Guidelines

### Writing Tests

- Write unit tests for all new functionality
- Use table-driven tests for testing multiple scenarios
- Mock external dependencies
- Aim for high test coverage (>80%)
- Include both positive and negative test cases

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if tt.wantErr && err == nil {
                t.Errorf("expected error but got none")
            }
            if !tt.wantErr && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            if got != tt.expected {
                t.Errorf("got %v, want %v", got, tt.expected)
            }
        })
    }
}
```

## 📖 Documentation

### Code Documentation

- Document all exported functions, types, and constants
- Use Go doc conventions
- Include examples in documentation where helpful
- Keep comments up-to-date with code changes

### README Updates

When adding new features:
- Update the feature list
- Add usage examples
- Update configuration options if applicable
- Include any new dependencies

## 🐛 Reporting Issues

### Bug Reports

When reporting bugs, please include:

- Go version (`go version`)
- Operating system and version
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Minimal code example that reproduces the issue
- Error messages or logs

### Feature Requests

When requesting new features:
- Describe the problem you're trying to solve
- Explain why this feature would be valuable
- Provide examples of how it would be used
- Consider potential implementation approaches

## 🔄 Release Process

Releases are managed by maintainers following semantic versioning:

- **Patch** (1.0.1): Bug fixes and minor improvements
- **Minor** (1.1.0): New features that don't break existing functionality
- **Major** (2.0.0): Breaking changes

## 🏷️ Issue Labels

We use the following labels to categorize issues:

- `bug`: Something isn't working correctly
- `enhancement`: New feature or improvement
- `documentation`: Documentation needs improvement
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention needed
- `question`: Further information requested
- `priority/high`: High priority issue
- `priority/medium`: Medium priority issue
- `priority/low`: Low priority issue

## 👥 Community Guidelines

- Be respectful and inclusive
- Provide constructive feedback
- Help others learn and grow
- Follow the [Go Code of Conduct](https://golang.org/conduct)

## 🎯 Areas for Contribution

We especially welcome contributions in these areas:

### Core Features
- New static analysis rules
- Performance improvements
- Security vulnerability detection
- Code complexity metrics

### Tools & Integrations
- IDE plugins and extensions
- Additional CI/CD integrations
- Docker improvements
- Package manager integrations

### Documentation
- API documentation
- Tutorials and guides
- Code examples
- Translations

### Testing
- Unit test coverage
- Integration tests
- Performance benchmarks
- Test data and fixtures

## 🚀 Advanced Development

### Adding New Analysis Rules

1. Create a new rule in `internal/rules/`
2. Implement the `Rule` interface:
   ```go
   type Rule interface {
       ID() string
       Description() string
       Check(node ast.Node) []Issue
   }
   ```
3. Register the rule in the analyzer
4. Add comprehensive tests
5. Update documentation

### Extending the Web API

1. Add new endpoints in `internal/server/`
2. Follow RESTful conventions
3. Include proper error handling
4. Add API documentation
5. Test with various HTTP clients

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Go AST Package](https://pkg.go.dev/go/ast)
- [Static Analysis in Go](https://go.dev/blog/analysis)
- [Effective Go](https://golang.org/doc/effective_go.html)

## 🙏 Recognition

Contributors will be recognized in:
- The project's AUTHORS file
- Release notes for significant contributions
- GitHub's contributor insights

Thank you for helping make Nada better! 🎉

## 📞 Getting Help

If you need help or have questions:

1. Check existing [issues](https://github.com/chaksack/nada/issues)
2. Start a [discussion](https://github.com/chaksack/nada/discussions)
3. Reach out to [@chaksack](https://github.com/chaksack)

---

*This contributing guide is inspired by best practices from the Go community and other successful open-source projects.*