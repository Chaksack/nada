package rules

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/chaksack/nada/internal/types"
)

type ErrorHandlingRule struct{}

func NewErrorHandlingRule() *ErrorHandlingRule {
	return &ErrorHandlingRule{}
}

func (r *ErrorHandlingRule) ID() string   { return "error_handling" }
func (r *ErrorHandlingRule) Name() string { return "Error Handling" }
func (r *ErrorHandlingRule) Description() string {
	return "Checks for proper error handling and deprecated function usage"
}

func (r *ErrorHandlingRule) Check(file string, node ast.Node, content string, fset *token.FileSet) []types.Issue {
	var issues []types.Issue

	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		if call, ok := n.(*ast.CallExpr); ok {
			pos := fset.Position(call.Pos())
			issues = append(issues, r.checkDeprecatedFunctions(file, call, pos)...)
			issues = append(issues, r.checkErrorHandling(file, call, pos)...)
		}

		return true
	})

	return issues
}

func (r *ErrorHandlingRule) checkDeprecatedFunctions(file string, call *ast.CallExpr, pos token.Position) []types.Issue {
	var issues []types.Issue

	if ident, ok := call.Fun.(*ast.Ident); ok {
		deprecatedFuncs := map[string]string{
			"ioutil.ReadFile":  "Use os.ReadFile instead",
			"ioutil.WriteFile": "Use os.WriteFile instead",
			"ioutil.ReadAll":   "Use io.ReadAll instead",
		}

		for deprecated, suggestion := range deprecatedFuncs {
			if ident.Name == deprecated {
				issues = append(issues, types.Issue{
					Type:        types.TypeCodeSmell,
					Severity:    types.SeverityMedium,
					File:        file,
					Line:        pos.Line,
					Column:      pos.Column,
					Rule:        "deprecated_function",
					Message:     "Deprecated function usage",
					Description: fmt.Sprintf("Function '%s' is deprecated. %s", deprecated, suggestion),
					Impact:      types.IssueImpact{EffortMinutes: 2},
				})
			}
		}
	}

	return issues
}

func (r *ErrorHandlingRule) checkErrorHandling(file string, call *ast.CallExpr, pos token.Position) []types.Issue {
	var issues []types.Issue

	// This is a simplified check for demonstration
	if ident, ok := call.Fun.(*ast.Ident); ok {
		riskyFuncs := []string{"os.Open", "json.Marshal", "strconv.Atoi", "http.Get"}

		for _, risky := range riskyFuncs {
			if ident.Name == risky {
				issues = append(issues, types.Issue{
					Type:        types.TypeBug,
					Severity:    types.SeverityMedium,
					File:        file,
					Line:        pos.Line,
					Column:      pos.Column,
					Rule:        "missing_error_handling",
					Message:     "Potential missing error handling",
					Description: fmt.Sprintf("Function '%s' may return an error that should be handled", risky),
					Impact:      types.IssueImpact{EffortMinutes: 3},
				})
			}
		}
	}

	return issues
}
