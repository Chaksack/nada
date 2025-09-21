package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/chaksack/nada/internal/types"
)

// ComplexityRule checks for cyclomatic complexity and code structure issues
type ComplexityRule struct{}

// NewComplexityRule creates a new complexity rule
func NewComplexityRule() *ComplexityRule {
	return &ComplexityRule{}
}

// ID returns the rule identifier
func (r *ComplexityRule) ID() string {
	return "complexity"
}

// Name returns the rule name
func (r *ComplexityRule) Name() string {
	return "Complexity Analysis"
}

// Description returns the rule description
func (r *ComplexityRule) Description() string {
	return "Checks for high cyclomatic complexity, deep nesting, and large functions"
}

// Check analyzes the AST for complexity issues
func (r *ComplexityRule) Check(file string, node ast.Node, content string, fset *token.FileSet) []types.Issue {
	var issues []types.Issue

	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch x := n.(type) {
		case *ast.FuncDecl:
			issues = append(issues, r.checkFunctionComplexity(file, x, fset)...)
			issues = append(issues, r.checkFunctionSize(file, x, fset)...)
		case *ast.IfStmt:
			issues = append(issues, r.checkDeepNesting(file, x, content, fset)...)
		}

		return true
	})

	return issues
}

// checkFunctionComplexity calculates cyclomatic complexity for functions
func (r *ComplexityRule) checkFunctionComplexity(file string, fn *ast.FuncDecl, fset *token.FileSet) []types.Issue {
	var issues []types.Issue

	if fn.Body == nil || fn.Name == nil {
		return issues
	}

	complexity := r.calculateComplexity(fn)
	pos := fset.Position(fn.Pos())

	if complexity > 10 {
		severity := types.SeverityMedium
		if complexity > 15 {
			severity = types.SeverityHigh
		}

		issues = append(issues, types.Issue{
			Type:        types.TypeCodeSmell,
			Severity:    severity,
			File:        file,
			Line:        pos.Line,
			Column:      pos.Column,
			Rule:        "high_complexity",
			Message:     "High cyclomatic complexity",
			Description: fmt.Sprintf("Function '%s' has complexity %d (threshold: 10)", fn.Name.Name, complexity),
			Impact:      types.IssueImpact{EffortMinutes: complexity * 2},
		})
	}

	return issues
}

// calculateComplexity calculates the cyclomatic complexity of a function
func (r *ComplexityRule) calculateComplexity(fn *ast.FuncDecl) int {
	complexity := 1 // Base complexity

	ast.Inspect(fn, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt,
			*ast.TypeSwitchStmt, *ast.SelectStmt:
			complexity++
		case *ast.CaseClause:
			complexity++
		}
		return true
	})

	return complexity
}

// checkFunctionSize checks if functions are too large
func (r *ComplexityRule) checkFunctionSize(file string, fn *ast.FuncDecl, fset *token.FileSet) []types.Issue {
	var issues []types.Issue

	if fn.Body == nil || fn.Name == nil {
		return issues
	}

	start := fset.Position(fn.Body.Lbrace)
	end := fset.Position(fn.Body.Rbrace)
	lines := end.Line - start.Line

	if lines > 50 {
		severity := types.SeverityMedium
		if lines > 100 {
			severity = types.SeverityHigh
		}

		pos := fset.Position(fn.Pos())
		issues = append(issues, types.Issue{
			Type:        types.TypeCodeSmell,
			Severity:    severity,
			File:        file,
			Line:        pos.Line,
			Column:      pos.Column,
			Rule:        "large_function",
			Message:     "Function too large",
			Description: fmt.Sprintf("Function '%s' has %d lines (threshold: 50)", fn.Name.Name, lines),
			Impact:      types.IssueImpact{EffortMinutes: lines / 10},
		})
	}

	return issues
}

// checkDeepNesting checks for deeply nested code structures
func (r *ComplexityRule) checkDeepNesting(file string, stmt *ast.IfStmt, content string, fset *token.FileSet) []types.Issue {
	var issues []types.Issue

	pos := fset.Position(stmt.Pos())
	lines := strings.Split(content, "\n")

	if pos.Line-1 < len(lines) {
		line := lines[pos.Line-1]
		indentLevel := r.calculateIndentLevel(line)

		if indentLevel > 4 {
			issues = append(issues, types.Issue{
				Type:        types.TypeCodeSmell,
				Severity:    types.SeverityMedium,
				File:        file,
				Line:        pos.Line,
				Column:      pos.Column,
				Rule:        "deep_nesting",
				Message:     "Deep nesting detected",
				Description: fmt.Sprintf("Code is nested %d levels deep (threshold: 4)", indentLevel),
				Impact:      types.IssueImpact{EffortMinutes: 5},
			})
		}
	}

	return issues
}

// calculateIndentLevel calculates the indentation level of a line
func (r *ComplexityRule) calculateIndentLevel(line string) int {
	indentLevel := 0
	for _, char := range line {
		if char == '\t' {
			indentLevel++
		} else if char == ' ' {
			indentLevel++
			if indentLevel%4 == 0 {
				indentLevel = indentLevel / 4
			}
		} else {
			break
		}
	}
	return indentLevel
}
