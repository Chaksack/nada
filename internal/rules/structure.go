package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/chaksack/nada/internal/types"
)

// StructureRule checks for code structure issues
type StructureRule struct{}

func NewStructureRule() *StructureRule {
	return &StructureRule{}
}

func (r *StructureRule) ID() string   { return "structure" }
func (r *StructureRule) Name() string { return "Code Structure" }
func (r *StructureRule) Description() string {
	return "Checks for structural issues like TODO comments, long lines, and unused imports"
}

func (r *StructureRule) Check(file string, node ast.Node, content string, fset *token.FileSet) []types.Issue {
	var issues []types.Issue

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lineNum := i + 1
		issues = append(issues, r.checkTodoComments(file, lineNum, line)...)
		issues = append(issues, r.checkLineLength(file, lineNum, line)...)
		issues = append(issues, r.checkUnusedCode(file, lineNum, line)...)
	}

	fmt.Printf("Returning issues: %+v\n", issues)
	return issues
}

func (r *StructureRule) checkTodoComments(file string, lineNum int, line string) []types.Issue {
	var issues []types.Issue
	patterns := []string{`(?i)//\s*todo`, `(?i)//\s*fixme`, `(?i)//\s*hack`}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, line); matched {
			issues = append(issues, types.Issue{
				Type:        types.TypeCodeSmell,
				Severity:    types.SeverityLow,
				File:        file,
				Line:        lineNum,
				Column:      1,
				Rule:        "todo_comment",
				Message:     "TODO/FIXME comment",
				Description: "Consider addressing this TODO/FIXME comment",
				Impact:      types.IssueImpact{EffortMinutes: 5},
			})
		}
	}
	return issues
}

func (r *StructureRule) checkLineLength(file string, lineNum int, line string) []types.Issue {
	var issues []types.Issue
	if len(line) > 120 {
		fmt.Printf("Found long line: %d\n", len(line))
		issues = append(issues, types.Issue{
			Type:        types.TypeCodeSmell,
			Severity:    types.SeverityLow,
			File:        file,
			Line:        lineNum,
			Column:      1,
			Rule:        "long_line",
			Message:     "Line too long",
			Description: fmt.Sprintf("Line has %d characters (threshold: 120)", len(line)),
			Impact:      types.IssueImpact{EffortMinutes: 2},
		})
	}
	return issues
}

func (r *StructureRule) checkUnusedCode(file string, lineNum int, line string) []types.Issue {
	var issues []types.Issue
	if strings.Contains(line, "import") && strings.Contains(line, "_") {
		issues = append(issues, types.Issue{
			Type:        types.TypeCodeSmell,
			Severity:    types.SeverityLow,
			File:        file,
			Line:        lineNum,
			Column:      1,
			Rule:        "blank_import",
			Message:     "Blank import",
			Description: "Consider if this blank import is necessary",
			Impact:      types.IssueImpact{EffortMinutes: 1},
		})
	}
	return issues
}
