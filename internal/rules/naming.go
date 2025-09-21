package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/chaksack/nada/internal/types"
)

// NamingRule checks for proper naming conventions
type NamingRule struct{}

// NewNamingRule creates a new naming rule
func NewNamingRule() *NamingRule {
	return &NamingRule{}
}

// ID returns the rule identifier
func (r *NamingRule) ID() string {
	return "naming"
}

// Name returns the rule name
func (r *NamingRule) Name() string {
	return "Naming Conventions"
}

// Description returns the rule description
func (r *NamingRule) Description() string {
	return "Checks for proper Go naming conventions and descriptive names"
}

// Check analyzes naming conventions in the code
func (r *NamingRule) Check(file string, node ast.Node, content string, fset *token.FileSet) []types.Issue {
	var issues []types.Issue

	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		pos := fset.Position(n.Pos())

		switch x := n.(type) {
		case *ast.FuncDecl:
			issues = append(issues, r.checkFunctionNaming(file, x, pos)...)
		case *ast.GenDecl:
			issues = append(issues, r.checkVariableNaming(file, x, pos)...)
		case *ast.TypeSpec:
			issues = append(issues, r.checkTypeNaming(file, x, pos)...)
		}

		return true
	})

	return issues
}

// checkFunctionNaming validates function naming conventions
func (r *NamingRule) checkFunctionNaming(file string, fn *ast.FuncDecl, pos token.Position) []types.Issue {
	var issues []types.Issue

	if fn.Name == nil {
		return issues
	}

	name := fn.Name.Name

	// Skip special functions
	if name == "main" || name == "init" || strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") {
		return issues
	}

	// Check for too short names
	if len(name) < 2 {
		issues = append(issues, types.Issue{
			Type:        types.TypeCodeSmell,
			Severity:    types.SeverityLow,
			File:        file,
			Line:        pos.Line,
			Column:      pos.Column,
			Rule:        "short_function_name",
			Message:     "Function name too short",
			Description: fmt.Sprintf("Function name '%s' should be more descriptive", name),
			Impact:      types.IssueImpact{EffortMinutes: 2},
		})
	}

	// Check for proper camelCase
	if !r.isValidCamelCase(name) {
		issues = append(issues, types.Issue{
			Type:        types.TypeCodeSmell,
			Severity:    types.SeverityLow,
			File:        file,
			Line:        pos.Line,
			Column:      pos.Column,
			Rule:        "naming_convention",
			Message:     "Invalid naming convention",
			Description: fmt.Sprintf("Function '%s' should follow camelCase convention", name),
			Impact:      types.IssueImpact{EffortMinutes: 2},
		})
	}

	// Check for common abbreviations that should be uppercase
	if r.hasImproperAbbreviations(name) {
		issues = append(issues, types.Issue{
			Type:        types.TypeCodeSmell,
			Severity:    types.SeverityLow,
			File:        file,
			Line:        pos.Line,
			Column:      pos.Column,
			Rule:        "abbreviation_convention",
			Message:     "Improper abbreviation capitalization",
			Description: fmt.Sprintf("Function '%s' should capitalize common abbreviations (HTTP, URL, API, etc.)", name),
			Impact:      types.IssueImpact{EffortMinutes: 2},
		})
	}

	return issues
}

// checkVariableNaming validates variable naming conventions
func (r *NamingRule) checkVariableNaming(file string, decl *ast.GenDecl, pos token.Position) []types.Issue {
	var issues []types.Issue

	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range valueSpec.Names {
				// Skip common short variable names
				if r.isAcceptableShortName(name.Name) {
					continue
				}

				if len(name.Name) == 1 {
					issues = append(issues, types.Issue{
						Type:        types.TypeCodeSmell,
						Severity:    types.SeverityLow,
						File:        file,
						Line:        pos.Line,
						Column:      pos.Column,
						Rule:        "short_variable_name",
						Message:     "Variable name too short",
						Description: fmt.Sprintf("Variable '%s' should have a more descriptive name", name.Name),
						Impact:      types.IssueImpact{EffortMinutes: 1},
					})
				}
			}
		}
	}

	return issues
}

// checkTypeNaming validates type naming conventions
func (r *NamingRule) checkTypeNaming(file string, typeSpec *ast.TypeSpec, pos token.Position) []types.Issue {
	var issues []types.Issue

	if typeSpec.Name == nil {
		return issues
	}

	name := typeSpec.Name.Name

	// Check if type name starts with uppercase (should be exported) or lowercase (private)
	if len(name) > 0 {
		firstChar := name[0]
		if firstChar >= 'A' && firstChar <= 'Z' {
			// Exported type - check for proper naming
			if !r.isValidPascalCase(name) {
				issues = append(issues, types.Issue{
					Type:        types.TypeCodeSmell,
					Severity:    types.SeverityLow,
					File:        file,
					Line:        pos.Line,
					Column:      pos.Column,
					Rule:        "type_naming_convention",
					Message:     "Invalid type naming convention",
					Description: fmt.Sprintf("Exported type '%s' should follow PascalCase convention", name),
					Impact:      types.IssueImpact{EffortMinutes: 2},
				})
			}
		}
	}

	return issues
}

// isValidCamelCase checks if a name follows camelCase convention
func (r *NamingRule) isValidCamelCase(name string) bool {
	// Simple camelCase validation
	if len(name) == 0 {
		return false
	}

	// First character should be lowercase for unexported functions
	if name[0] >= 'A' && name[0] <= 'Z' {
		return true // Exported function, should be PascalCase
	}

	// Check for underscores (not camelCase)
	if strings.Contains(name, "_") {
		return false
	}

	return true
}

// isValidPascalCase checks if a name follows PascalCase convention
func (r *NamingRule) isValidPascalCase(name string) bool {
	if len(name) == 0 {
		return false
	}

	// First character should be uppercase
	if name[0] < 'A' || name[0] > 'Z' {
		return false
	}

	// Check for underscores (not PascalCase)
	if strings.Contains(name, "_") {
		return false
	}

	return true
}

// hasImproperAbbreviations checks for common abbreviations that should be uppercase
func (r *NamingRule) hasImproperAbbreviations(name string) bool {
	improperPatterns := []string{
		`(?i)\bhttp\b`, `(?i)\burl\b`, `(?i)\bapi\b`, `(?i)\bjson\b`,
		`(?i)\bxml\b`, `(?i)\bhtml\b`, `(?i)\bid\b`, `(?i)\bsql\b`,
	}

	for _, pattern := range improperPatterns {
		if matched, _ := regexp.MatchString(pattern, name); matched {
			// Check if it's already properly capitalized
			properPattern := strings.ToUpper(strings.Trim(pattern, `(?i)\b`))
			if !strings.Contains(name, properPattern) {
				return true
			}
		}
	}

	return false
}

// isAcceptableShortName checks if a short variable name is acceptable
func (r *NamingRule) isAcceptableShortName(name string) bool {
	acceptableShortNames := []string{
		"i", "j", "k", "n", "m", // Loop counters
		"x", "y", "z", // Coordinates
		"w", "h", // Width, height
		"r", "w", // Reader, writer (in context)
		"s", // String (in context)
		"b", // Byte/boolean (in context)
		"t", // Time/testing (in context)
		"c", // Channel/context (in context)
		"f", // File/function (in context)
	}

	for _, acceptable := range acceptableShortNames {
		if name == acceptable {
			return true
		}
	}

	return false
}
