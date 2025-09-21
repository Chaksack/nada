package rules

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/chaksack/nada/internal/types"
)

type DocumentationRule struct{}

func NewDocumentationRule() *DocumentationRule {
	return &DocumentationRule{}
}

func (r *DocumentationRule) ID() string   { return "documentation" }
func (r *DocumentationRule) Name() string { return "Documentation" }
func (r *DocumentationRule) Description() string {
	return "Checks for missing function and type documentation"
}

func (r *DocumentationRule) Check(file string, node ast.Node, content string, fset *token.FileSet) []types.Issue {
	var issues []types.Issue
	lines := strings.Split(content, "\n")

	// Check for missing function documentation
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "func ") &&
			!strings.Contains(line, "func main") &&
			!strings.Contains(line, "func init") &&
			strings.Contains(line, "(") {

			// Check if previous line is a comment
			if i == 0 || !strings.HasPrefix(strings.TrimSpace(lines[i-1]), "//") {
				// Check if function is exported (starts with capital letter)
				funcName := r.extractFunctionName(line)
				if len(funcName) > 0 && funcName[0] >= 'A' && funcName[0] <= 'Z' {
					issues = append(issues, types.Issue{
						Type:        types.TypeCodeSmell,
						Severity:    types.SeverityLow,
						File:        file,
						Line:        i + 1,
						Column:      1,
						Rule:        "missing_documentation",
						Message:     "Missing function documentation",
						Description: "Exported functions should have documentation comments",
						Impact:      types.IssueImpact{EffortMinutes: 3},
					})
				}
			}
		}
	}

	return issues
}

func (r *DocumentationRule) extractFunctionName(line string) string {
	// Simple extraction of function name from function declaration
	parts := strings.Fields(line)
	for i, part := range parts {
		if part == "func" && i+1 < len(parts) {
			name := parts[i+1]
			if parenIndex := strings.Index(name, "("); parenIndex > 0 {
				return name[:parenIndex]
			}
			return name
		}
	}
	return ""
}
