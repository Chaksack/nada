package rules

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/chaksack/nada/internal/types"
)

// SecurityRule checks for security vulnerabilities
type SecurityRule struct{}

// NewSecurityRule creates a new security rule
func NewSecurityRule() *SecurityRule {
	return &SecurityRule{}
}

// ID returns the rule identifier
func (r *SecurityRule) ID() string {
	return "security"
}

// Name returns the rule name
func (r *SecurityRule) Name() string {
	return "Security Analysis"
}

// Description returns the rule description
func (r *SecurityRule) Description() string {
	return "Detects security vulnerabilities including hardcoded secrets, SQL injection, and unsafe practices"
}

// Check analyzes the code for security issues
func (r *SecurityRule) Check(file string, node ast.Node, content string, fset *token.FileSet) []types.Issue {
	var issues []types.Issue

	// Check file content line by line
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lineNum := i + 1
		issues = append(issues, r.checkHardcodedSecrets(file, lineNum, line)...)
		issues = append(issues, r.checkSQLInjection(file, lineNum, line)...)
	}

	return issues
}

// checkHardcodedSecrets detects hardcoded secrets and credentials
func (r *SecurityRule) checkHardcodedSecrets(file string, lineNum int, line string) []types.Issue {
	var issues []types.Issue

	secretPatterns := []struct {
		pattern     string
		description string
		severity    string
	}{
		{`(?i)password\s*[:=]\s*["'][^"']{3,}["']`, "Hardcoded password", types.SeverityHigh},
		{`(?i)secret\s*[:=]\s*["'][^"']{8,}["']`, "Hardcoded secret", types.SeverityHigh},
		{`(?i)api[_-]?key\s*[:=]\s*["'][^"']{8,}["']`, "Hardcoded API key", types.SeverityHigh},
		{`(?i)token\s*[:=]\s*["'][^"']{16,}["']`, "Hardcoded token", types.SeverityHigh},
		{`(?i)aws[_-]?access[_-]?key\s*[:=]\s*["'][^"']+["']`, "AWS access key", types.SeverityHigh},
		{`(?i)private[_-]?key\s*[:=]\s*["'][^"']+["']`, "Private key", types.SeverityHigh},
	}

	for _, sp := range secretPatterns {
		if matched, err := regexp.MatchString(sp.pattern, line); err == nil && matched {
			// Additional check to avoid false positives
			if !r.isFalsePositive(line) {
				issues = append(issues, types.Issue{
					Type:        types.TypeVulnerability,
					Severity:    sp.severity,
					File:        file,
					Line:        lineNum,
					Column:      1,
					Rule:        "hardcoded_secret",
					Message:     sp.description,
					Description: "Hardcoded secrets should be moved to environment variables or secure configuration",
					Impact:      types.IssueImpact{EffortMinutes: 10},
				})
			}
		}
	}

	return issues
}

// checkSQLInjection detects potential SQL injection vulnerabilities
func (r *SecurityRule) checkSQLInjection(file string, lineNum int, line string) []types.Issue {
	var issues []types.Issue

	sqlPatterns := []string{
		`(?i)query\s*[:=]\s*["'].*%[sv].*["']`,
		`(?i)fmt\.Sprintf\s*\(\s*["'].*SELECT.*%[sv].*["']`,
		`(?i)fmt\.Sprintf\s*\(\s*["'].*INSERT.*%[sv].*["']`,
		`(?i)fmt\.Sprintf\s*\(\s*["'].*UPDATE.*%[sv].*["']`,
		`(?i)fmt\.Sprintf\s*\(\s*["'].*DELETE.*%[sv].*["']`,
		`(?i)["'].*SELECT.*\+.*["']`,
		`(?i)["'].*INSERT.*\+.*["']`,
	}

	for _, pattern := range sqlPatterns {
		if matched, err := regexp.MatchString(pattern, line); err == nil && matched {
			issues = append(issues, types.Issue{
				Type:        types.TypeVulnerability,
				Severity:    types.SeverityHigh,
				File:        file,
				Line:        lineNum,
				Column:      1,
				Rule:        "sql_injection",
				Message:     "Potential SQL injection",
				Description: "Use parameterized queries to prevent SQL injection attacks",
				Impact:      types.IssueImpact{EffortMinutes: 15},
			})
		}
	}

	return issues
}

// isFalsePositive checks if a potential secret detection is a false positive
func (r *SecurityRule) isFalsePositive(line string) bool {
	// Common false positives
	falsePositives := []string{
		"password=\"\"",
		"password=''",
		"secret=\"\"",
		"secret=''",
		"token=\"\"",
		"token=''",
		"password=\"placeholder\"",
		"password=\"example\"",
		"password=\"test\"",
		"password=\"dummy\"",
		"secret=\"placeholder\"",
		"api_key=\"your_key_here\"",
		"token=\"your_token_here\"",
	}

	lowerLine := strings.ToLower(line)
	for _, fp := range falsePositives {
		if strings.Contains(lowerLine, fp) {
			return true
		}
	}

	// Check for comments
	if strings.Contains(line, "//") || strings.Contains(line, "/*") {
		return true
	}

	return false
}
