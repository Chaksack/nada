# Security Policy

## Supported Versions

We actively support the following versions of Nada with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability in Nada, please report it responsibly.

### How to Report

**Please DO NOT report security vulnerabilities through public GitHub issues.**

Instead, please report security vulnerabilities by emailing **security@chaksack.dev** with the following information:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact of the vulnerability
- Any suggested fixes or mitigation strategies
- Your contact information (if you'd like to be credited)

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours
- **Initial Assessment**: We will provide an initial assessment within 5 business days
- **Regular Updates**: We will send you regular updates on our progress
- **Resolution**: We aim to resolve security issues within 90 days
- **Disclosure**: We will coordinate with you on public disclosure timing

### Security Response Process

1. **Report Received**: Vulnerability report is received and assigned to a primary handler
2. **Confirmation**: The vulnerability is confirmed and its impact is assessed
3. **Fix Development**: A fix is developed and tested
4. **Release Preparation**: A security release is prepared
5. **Public Disclosure**: The vulnerability is disclosed publicly after the fix is available

### Security Best Practices for Users

When using Nada in your environment:

- Always use the latest stable version
- Run Nada with minimal required privileges
- Be cautious when analyzing untrusted code
- Review security advisories regularly
- Use official installation methods (go install, official releases)

### Security Features

Nada includes several security features:

- **Safe AST Parsing**: Uses Go's built-in AST parser with safety checks
- **Sandboxed Analysis**: Code analysis runs in a controlled environment
- **Input Validation**: All user inputs are validated and sanitized
- **No Code Execution**: Nada only analyzes code; it never executes analyzed code

### Scope

This security policy applies to:

- The main Nada repository (`github.com/chaksack/nada`)
- Official releases and binaries
- Docker images
- Documentation and examples

### Hall of Fame

We recognize security researchers who help keep Nada secure:

<!-- Future security researchers will be listed here -->

*No security issues have been reported yet.*

---

Thank you for helping keep Nada and the Go community safe!

**Contact**: security@chaksack.dev  
**PGP Key**: Available upon request