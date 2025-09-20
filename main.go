package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// --- Type Definitions ---
type Issue struct {
	Type        string      `json:"type"`
	Severity    string      `json:"severity"`
	File        string      `json:"file"`
	Line        int         `json:"line"`
	Column      int         `json:"column"`
	Message     string      `json:"message"`
	Rule        string      `json:"rule"`
	Description string      `json:"description"`
	Category    string      `json:"category,omitempty"`
	Priority    string      `json:"priority,omitempty"`
	Impact      IssueImpact `json:"impact,omitempty"`
}

type IssueImpact struct {
	EffortMinutes int `json:"effort_minutes"`
}

type QualityTrends struct {
	CyclomaticComplexityTrend string
	IssuesDensity             float64
	SecurityScore             float64
	MaintainabilityIndex      float64
	TechnicalDebtRatio        float64
}

type Report struct {
	ProjectPath   string         `json:"project_path"`
	Timestamp     time.Time      `json:"timestamp"`
	Grade         string         `json:"grade"`
	Score         float64        `json:"score"`
	Issues        []Issue        `json:"issues"`
	Metrics       Metrics        `json:"metrics"`
	FilesAnalyzed int            `json:"files_analyzed"`
	IssuesSummary map[string]int `json:"issues_summary"`
}

type Metrics struct {
	LinesOfCode          int     `json:"lines_of_code"`
	CyclomaticComplexity int     `json:"cyclomatic_complexity"`
	CodeDuplication      float64 `json:"code_duplication"`
	TestCoverage         float64 `json:"test_coverage"`
	TechnicalDebt        string  `json:"technical_debt"`
	Maintainability      string  `json:"maintainability"`
	Reliability          string  `json:"reliability"`
	Security             string  `json:"security"`
}

type CodeAnalyzer struct {
	fileSet    *token.FileSet
	issues     []Issue
	metrics    Metrics
	filesCount int
}

// --- Helper Methods ---
func (ca *CodeAnalyzer) getIssueCountByType(t string) int {
	count := 0
	for _, issue := range ca.issues {
		if issue.Type == t {
			count++
		}
	}
	return count
}

func (ca *CodeAnalyzer) getComplexityIssueCount() int {
	count := 0
	for _, issue := range ca.issues {
		if issue.Rule == "complex_function" {
			count++
		}
	}
	return count
}

func (ca *CodeAnalyzer) getIssueCountByRule(rule string) int {
	count := 0
	for _, issue := range ca.issues {
		if issue.Rule == rule {
			count++
		}
	}
	return count
}

// --- Helper Methods ---

// NewCodeAnalyzer creates a new code analyzer instance
func NewCodeAnalyzer() *CodeAnalyzer {
	return &CodeAnalyzer{
		fileSet: token.NewFileSet(),
		issues:  make([]Issue, 0),
	}
}

// AnalyzeProject analyzes the entire project
func (ca *CodeAnalyzer) AnalyzeProject(projectPath string) (*Report, error) {
	fmt.Printf("🔍 Analyzing project: %s\n", projectPath)

	err := filepath.WalkDir(projectPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			if !strings.Contains(path, "vendor/") && !strings.Contains(path, ".git/") {
				ca.analyzeFile(path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	report := &Report{
		ProjectPath:   projectPath,
		Timestamp:     time.Now(),
		Issues:        ca.issues,
		Metrics:       ca.metrics,
		FilesAnalyzed: ca.filesCount,
		IssuesSummary: make(map[string]int),
	}

	report.Score = 100.0 // Placeholder for scoring logic
	report.Grade = "A"   // Placeholder for grading logic

	return report, nil
}

// generateRecommendations generates actionable recommendations
func (ca *CodeAnalyzer) generateRecommendations() []string {
	recommendations := []string{}

	vulnCount := ca.getIssueCountByType("vulnerability")
	if vulnCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🔒 URGENT: Address %d security vulnerabilities immediately to prevent potential breaches", vulnCount))
	}

	bugCount := ca.getIssueCountByType("bug")
	if bugCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🐛 HIGH: Fix %d bugs to improve application reliability", bugCount))
	}

	complexityIssues := ca.getComplexityIssueCount()
	if complexityIssues > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🔧 MEDIUM: Refactor %d complex functions to improve maintainability", complexityIssues))
	}

	if ca.metrics.TestCoverage < 70 {
		recommendations = append(recommendations,
			fmt.Sprintf("🧪 Increase test coverage from %.1f%% to at least 70%%", ca.metrics.TestCoverage))
	}

	debtMinutes := 0
	for _, issue := range ca.issues {
		debtMinutes += issue.Impact.EffortMinutes
	}

	if debtMinutes > 0 {
		hours := float64(debtMinutes) / 60
		recommendations = append(recommendations,
			fmt.Sprintf("⏰ Technical debt: Approximately %.1f hours of effort needed to address all issues", hours))
	}

	// Add specific recommendations based on common patterns
	hardcodedSecrets := ca.getIssueCountByRule("hardcoded_secret")
	if hardcodedSecrets > 0 {
		recommendations = append(recommendations,
			"🔑 Implement a secrets management solution (e.g., HashiCorp Vault, AWS Secrets Manager)")
	}

	sqlInjections := ca.getIssueCountByRule("sql_injection")
	if sqlInjections > 0 {
		recommendations = append(recommendations,
			"🛡️ Use parameterized queries and ORM frameworks to prevent SQL injection attacks")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "✅ Great job! Your code quality is excellent. Consider adding more tests to maintain this standard.")
	}

	return recommendations
}

// calculateQualityTrends calculates various quality trend indicators
func (ca *CodeAnalyzer) calculateQualityTrends() QualityTrends {
	trends := QualityTrends{}

	if ca.filesCount > 0 {
		avgComplexity := float64(ca.metrics.CyclomaticComplexity) / float64(ca.filesCount)

		// Complexity trend assessment
		switch {
		case avgComplexity > 15:
			trends.CyclomaticComplexityTrend = "Very High - Immediate refactoring needed"
		case avgComplexity > 10:
			trends.CyclomaticComplexityTrend = "High - Consider refactoring"
		case avgComplexity > 5:
			trends.CyclomaticComplexityTrend = "Moderate - Monitor closely"
		default:
			trends.CyclomaticComplexityTrend = "Good - Well structured code"
		}

		// Issues density (issues per 1000 lines of code)
		if ca.metrics.LinesOfCode > 0 {
			trends.IssuesDensity = float64(len(ca.issues)) / float64(ca.metrics.LinesOfCode) * 1000
		}

		// Security score (100 - security issues impact)
		securityIssues := ca.getIssueCountByType("vulnerability")
		trends.SecurityScore = max(0, 100-float64(securityIssues*10))

		// Maintainability index (based on complexity and documentation)
		documentationIssues := ca.getIssueCountByRule("missing_documentation")
		maintainabilityDeductions := float64(securityIssues*5 + documentationIssues*2)
		trends.MaintainabilityIndex = max(0, 100-maintainabilityDeductions)

		// Technical debt ratio (estimated hours / total LOC * 1000)
		totalDebtMinutes := 0
		for _, issue := range ca.issues {
			totalDebtMinutes += issue.Impact.EffortMinutes
		}
		if ca.metrics.LinesOfCode > 0 {
			trends.TechnicalDebtRatio = float64(totalDebtMinutes) / 60 / float64(ca.metrics.LinesOfCode) * 1000
		}
	}

	return trends
}

// Helper functions for risk analysis
func (ca *CodeAnalyzer) getRiskProbability(riskArea string) string {
	switch riskArea {
	case "vulnerability":
		count := ca.getIssueCountByType("vulnerability")
		if count > 5 {
			return "High"
		} else if count > 2 {
			return "Medium"
		}
		return "Low"
	case "complexity":
		avg := float64(ca.metrics.CyclomaticComplexity) / max(1, float64(ca.filesCount))
		if avg > 10 {
			return "High"
		} else if avg > 5 {
			return "Medium"
		}
		return "Low"
	default:
		return "Medium"
	}
}

// analyzeFile analyzes a single Go file
func (ca *CodeAnalyzer) analyzeFile(filePath string) {
	ca.filesCount++

	content, err := os.ReadFile(filePath)
	if err != nil {
		ca.addIssue(Issue{
			Type:        "error",
			Severity:    "high",
			File:        filePath,
			Line:        1,
			Column:      1,
			Rule:        "read_error",
			Message:     "Cannot read file",
			Description: fmt.Sprintf("Error reading file: %v", err),
			Impact:      IssueImpact{EffortMinutes: 1},
		})
		return
	}

	// Parse the Go file
	node, err := parser.ParseFile(ca.fileSet, filePath, content, parser.ParseComments)
	if err != nil {
		ca.addIssue(Issue{
			Type:        "error",
			Severity:    "high",
			File:        filePath,
			Line:        1,
			Column:      1,
			Rule:        "parse_error",
			Message:     "Syntax error",
			Description: fmt.Sprintf("Parse error: %v", err),
			Impact:      IssueImpact{EffortMinutes: 1},
		})
		return
	}

	// Analyze AST
	ca.analyzeAST(filePath, node, string(content))

	// Analyze file content for additional issues
	ca.analyzeFileContent(filePath, string(content))

	// Update metrics
	ca.updateMetrics(string(content), node)
}

// analyzeAST analyzes the Abstract Syntax Tree
func (ca *CodeAnalyzer) analyzeAST(filePath string, node *ast.File, content string) {
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		pos := ca.fileSet.Position(n.Pos())

		switch x := n.(type) {
		case *ast.FuncDecl:
			ca.analyzeFunctionComplexity(filePath, x, pos)
			ca.checkFunctionNaming(filePath, x, pos)
			ca.checkFunctionSize(filePath, x, pos)

		case *ast.IfStmt:
			ca.checkDeepNesting(filePath, x, pos, content)

		case *ast.GenDecl:
			ca.checkVariableNaming(filePath, x, pos)

		case *ast.CallExpr:
			ca.checkDeprecatedFunctions(filePath, x, pos)
			ca.checkErrorHandling(filePath, x, pos)
		}

		return true
	})
}

// analyzeFileContent analyzes file content for patterns
func (ca *CodeAnalyzer) analyzeFileContent(filePath, content string) {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		lineNum := i + 1

		// Check for hardcoded passwords/secrets
		ca.checkHardcodedSecrets(filePath, lineNum, line)

		ca.checkTodoComments(filePath, lineNum, line)
		ca.checkTodoComments(filePath, lineNum, line)

		// Check for long lines
		ca.checkLineLength(filePath, lineNum, line)

		// Check for SQL injection patterns
		ca.checkSQLInjection(filePath, lineNum, line)

		// Check for unused imports/variables (simple pattern matching)
		ca.checkUnusedCode(filePath, lineNum, line)
	}

	// Check for missing documentation
	ca.checkDocumentation(filePath, content)
}

// analyzeFunctionComplexity calculates cyclomatic complexity
func (ca *CodeAnalyzer) analyzeFunctionComplexity(filePath string, fn *ast.FuncDecl, pos token.Position) {
	if fn.Body == nil {
		return
	}

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

	if complexity > 10 {
		severity := "medium"
		if complexity > 15 {
			severity = "high"
		}
		ca.addIssue(Issue{
			Type:        "code_smell",
			Severity:    severity,
			File:        filePath,
			Line:        pos.Line,
			Column:      pos.Column,
			Rule:        "high_complexity",
			Message:     "High cyclomatic complexity",
			Description: fmt.Sprintf("Function '%s' has complexity %d (threshold: 10)", fn.Name.Name, complexity),
			Impact:      IssueImpact{EffortMinutes: 10},
		})
	}

	ca.metrics.CyclomaticComplexity += complexity
}

// checkFunctionNaming checks function naming conventions
func (ca *CodeAnalyzer) checkFunctionNaming(filePath string, fn *ast.FuncDecl, pos token.Position) {
	if fn.Name == nil {
		return
	}

	name := fn.Name.Name

	// Check for proper naming conventions
	if len(name) < 2 {
		ca.addIssue(Issue{
			Type:        "code_smell",
			Severity:    "low",
			File:        filePath,
			Line:        pos.Line,
			Column:      pos.Column,
			Rule:        "short_function_name",
			Message:     "Function name too short",
			Description: fmt.Sprintf("Function name '%s' should be more descriptive", name),
			Impact:      IssueImpact{EffortMinutes: 2},
		})
	}

	// Check for acronyms in function names
	if strings.Contains(name, "URL") || strings.Contains(name, "HTTP") || strings.Contains(name, "API") {
		if !regexp.MustCompile(`^[A-Z][a-z]*([A-Z][a-z]*)*$`).MatchString(name) {
			ca.addIssue(Issue{
				Type:        "code_smell",
				Severity:    "low",
				File:        filePath,
				Line:        pos.Line,
				Column:      pos.Column,
				Rule:        "naming_convention",
				Message:     "Inconsistent naming convention",
				Description: fmt.Sprintf("Function '%s' should follow Go naming conventions", name),
				Impact:      IssueImpact{EffortMinutes: 2},
			})
		}
	}
}

// checkFunctionSize checks if function is too large
func (ca *CodeAnalyzer) checkFunctionSize(filePath string, fn *ast.FuncDecl, pos token.Position) {
	if fn.Body == nil {
		return
	}

	start := ca.fileSet.Position(fn.Body.Lbrace)
	end := ca.fileSet.Position(fn.Body.Rbrace)
	lines := end.Line - start.Line

	if lines > 50 {
		severity := "medium"
		if lines > 100 {
			severity = "high"
		}
		ca.addIssue(Issue{
			Type:        "code_smell",
			Severity:    severity,
			File:        filePath,
			Line:        pos.Line,
			Column:      pos.Column,
			Rule:        "large_function",
			Message:     "Function too large",
			Description: fmt.Sprintf("Function '%s' has %d lines (threshold: 50)", fn.Name.Name, lines),
			Impact:      IssueImpact{EffortMinutes: 8},
		})
	}
}

// checkDeepNesting checks for deeply nested code
func (ca *CodeAnalyzer) checkDeepNesting(filePath string, stmt *ast.IfStmt, pos token.Position, content string) {
	lines := strings.Split(content, "\n")
	if pos.Line-1 < len(lines) {
		line := lines[pos.Line-1]
		indentLevel := 0
		for _, char := range line {
			if char == '\t' {
				indentLevel++
			} else if char == ' ' {
				indentLevel++
				if indentLevel%4 == 0 { // Assuming 4 spaces = 1 tab
					indentLevel = indentLevel / 4
				}
			} else {
				break
			}
		}

		if indentLevel > 4 {
			ca.addIssue(Issue{
				Type:        "code_smell",
				Severity:    "medium",
				File:        filePath,
				Line:        pos.Line,
				Column:      pos.Column,
				Rule:        "deep_nesting",
				Message:     "Deep nesting detected",
				Description: fmt.Sprintf("Code is nested %d levels deep (threshold: 4)", indentLevel),
				Impact:      IssueImpact{EffortMinutes: 5},
			})
		}
	}
}

// checkVariableNaming checks variable naming conventions
func (ca *CodeAnalyzer) checkVariableNaming(filePath string, decl *ast.GenDecl, pos token.Position) {
	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range valueSpec.Names {
				if len(name.Name) == 1 && name.Name != "i" && name.Name != "j" && name.Name != "k" {
					ca.addIssue(Issue{
						Type:        "code_smell",
						Severity:    "low",
						File:        filePath,
						Line:        pos.Line,
						Column:      pos.Column,
						Rule:        "short_variable_name",
						Message:     "Variable name too short",
						Description: fmt.Sprintf("Variable '%s' should have a more descriptive name", name.Name),
						Impact:      IssueImpact{EffortMinutes: 1},
					})
				}
			}
		}
	}
}

// checkDeprecatedFunctions checks for deprecated function calls
func (ca *CodeAnalyzer) checkDeprecatedFunctions(filePath string, call *ast.CallExpr, pos token.Position) {
	if ident, ok := call.Fun.(*ast.Ident); ok {
		deprecatedFuncs := map[string]string{
			"ioutil.ReadFile":  "Use os.ReadFile instead",
			"ioutil.WriteFile": "Use os.WriteFile instead",
			"ioutil.ReadAll":   "Use io.ReadAll instead",
		}

		for deprecated, suggestion := range deprecatedFuncs {
			if ident.Name == deprecated {
				ca.addIssue(Issue{
					Type:        "code_smell",
					Severity:    "medium",
					File:        filePath,
					Line:        pos.Line,
					Column:      pos.Column,
					Rule:        "deprecated_function",
					Message:     "Deprecated function usage",
					Description: fmt.Sprintf("Function '%s' is deprecated. %s", deprecated, suggestion),
					Impact:      IssueImpact{EffortMinutes: 2},
				})
			}
		}
	}
}

// checkErrorHandling checks for proper error handling
func (ca *CodeAnalyzer) checkErrorHandling(filePath string, call *ast.CallExpr, pos token.Position) {
	// This is a simplified check - in practice, you'd want more sophisticated analysis
	if ident, ok := call.Fun.(*ast.Ident); ok {
		riskyFuncs := []string{"os.Open", "json.Marshal", "strconv.Atoi", "http.Get"}

		for _, risky := range riskyFuncs {
			if ident.Name == risky {
				// Check if error is being handled (this is a simplified check)
				ca.addIssue(Issue{
					Type:        "bug",
					Severity:    "medium",
					File:        filePath,
					Line:        pos.Line,
					Column:      pos.Column,
					Rule:        "missing_error_handling",
					Message:     "Potential missing error handling",
					Description: fmt.Sprintf("Function '%s' may return an error that should be handled", risky),
					Impact:      IssueImpact{EffortMinutes: 3},
				})
			}
		}
	}
}

// checkHardcodedSecrets checks for hardcoded secrets
func (ca *CodeAnalyzer) checkHardcodedSecrets(filePath string, lineNum int, line string) {
	secretPatterns := []struct {
		pattern string
		desc    string
	}{
		{`(?i)password\s*[:=]\s*["'][^"']+["']`, "Hardcoded password"},
		{`(?i)secret\s*[:=]\s*["'][^"']+["']`, "Hardcoded secret"},
		{`(?i)api[_-]?key\s*[:=]\s*["'][^"']+["']`, "Hardcoded API key"},
		{`(?i)token\s*[:=]\s*["'][^"']+["']`, "Hardcoded token"},
		{`(?i)aws[_-]?access[_-]?key\s*[:=]\s*["'][^"']+["']`, "AWS access key"},
	}

	for _, sp := range secretPatterns {
		if matched, _ := regexp.MatchString(sp.pattern, line); matched {
			ca.addIssue(Issue{
				Type:        "vulnerability",
				Severity:    "high",
				File:        filePath,
				Line:        lineNum,
				Column:      1,
				Rule:        "hardcoded_secret",
				Message:     sp.desc,
				Description: "Hardcoded secrets should be moved to environment variables or secure configuration",
				Impact:      IssueImpact{EffortMinutes: 10},
			})
		}
	}
}

func (ca *CodeAnalyzer) checkTodoComments(filePath string, lineNum int, line string) {
	patterns := []string{`(?i)//\s*todo`, `(?i)//\s*fixme`, `(?i)//\s*hack`}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, line); matched {
			ca.addIssue(Issue{
				Type:        "code_smell",
				Severity:    "low",
				File:        filePath,
				Line:        lineNum,
				Column:      1,
				Rule:        "todo_comment",
				Message:     "TODO/FIXME comment",
				Description: "Consider addressing this TODO/FIXME comment",
				Impact:      IssueImpact{EffortMinutes: 5},
			})
		}
	}
}

// checkLineLength checks for overly long lines
func (ca *CodeAnalyzer) checkLineLength(filePath string, lineNum int, line string) {
	if len(line) > 120 {
		ca.addIssue(Issue{
			Type:        "code_smell",
			Severity:    "low",
			File:        filePath,
			Line:        lineNum,
			Column:      1,
			Rule:        "long_line",
			Message:     "Line too long",
			Description: fmt.Sprintf("Line has %d characters (threshold: 120)", len(line)),
			Impact:      IssueImpact{EffortMinutes: 2},
		})
	}
}

// checkSQLInjection checks for potential SQL injection vulnerabilities
func (ca *CodeAnalyzer) checkSQLInjection(filePath string, lineNum int, line string) {
	sqlPatterns := []string{
		`(?i)query\s*[:=]\s*["'].*%s.*["']`,
		`(?i)fmt\.Sprintf\s*\(\s*["'].*SELECT.*%s.*["']`,
		`(?i)fmt\.Sprintf\s*\(\s*["'].*INSERT.*%s.*["']`,
		`(?i)fmt\.Sprintf\s*\(\s*["'].*UPDATE.*%s.*["']`,
	}

	for _, pattern := range sqlPatterns {
		if matched, _ := regexp.MatchString(pattern, line); matched {
			ca.addIssue(Issue{
				Type:        "vulnerability",
				Severity:    "high",
				File:        filePath,
				Line:        lineNum,
				Column:      1,
				Rule:        "sql_injection",
				Message:     "Potential SQL injection",
				Description: "Use parameterized queries to prevent SQL injection attacks",
				Impact:      IssueImpact{EffortMinutes: 15},
			})
		}
	}
}

// checkUnusedCode checks for potentially unused code
func (ca *CodeAnalyzer) checkUnusedCode(filePath string, lineNum int, line string) {
	// Simplified check for unused imports
	if strings.Contains(line, "import") && strings.Contains(line, "_") {
		ca.addIssue(Issue{
			Type:        "code_smell",
			Severity:    "low",
			File:        filePath,
			Line:        lineNum,
			Column:      1,
			Rule:        "blank_import",
			Message:     "Blank import",
			Description: "Consider if this blank import is necessary",
			Impact:      IssueImpact{EffortMinutes: 1},
		})
	}
}

// checkDocumentation checks for missing documentation
func (ca *CodeAnalyzer) checkDocumentation(filePath, content string) {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "func ") &&
			!strings.Contains(line, "func main") &&
			strings.Contains(line, "(") {
			// Check if previous line is a comment
			if i == 0 || !strings.HasPrefix(strings.TrimSpace(lines[i-1]), "//") {
				ca.addIssue(Issue{
					Type:        "code_smell",
					Severity:    "low",
					File:        filePath,
					Line:        i + 1,
					Column:      1,
					Rule:        "missing_documentation",
					Message:     "Missing function documentation",
					Description: "Public functions should have documentation comments",
					Impact:      IssueImpact{EffortMinutes: 3},
				})
			}
		}
	}
}

// updateMetrics updates the code metrics
func (ca *CodeAnalyzer) updateMetrics(content string, node *ast.File) {
	lines := strings.Split(content, "\n")
	ca.metrics.LinesOfCode += len(lines)

	// Calculate test coverage (simplified - would need actual test execution)
	testFiles := 0
	if strings.Contains(content, "func Test") {
		testFiles++
	}
	if ca.filesCount > 0 {
		ca.metrics.TestCoverage = float64(testFiles) / float64(ca.filesCount) * 100
	}
}

// addIssue adds a new issue to the analyzer
// addIssue adds a new issue to the analyzer
// addIssue adds a new issue to the analyzer (refactored to take Issue struct)
func (ca *CodeAnalyzer) addIssue(issue Issue) {
	ca.issues = append(ca.issues, issue)
}

// getIssuesSummary returns a summary of issues by severity
func (ca *CodeAnalyzer) getIssuesSummary() map[string]int {
	summary := make(map[string]int)

	for _, issue := range ca.issues {
		summary[issue.Severity]++
		summary[issue.Type]++
	}

	return summary
}

// calculateScore calculates the overall code quality score
func (ca *CodeAnalyzer) calculateScore() float64 {
	if ca.filesCount == 0 {
		return 0
	}

	baseScore := 100.0
	highIssues := ca.getIssuesSummary()["high"]
	mediumIssues := ca.getIssuesSummary()["medium"]
	lowIssues := ca.getIssuesSummary()["low"]

	// Deduct points based on issues
	score := baseScore - float64(highIssues*10) - float64(mediumIssues*5) - float64(lowIssues*1)

	// Factor in complexity
	avgComplexity := float64(ca.metrics.CyclomaticComplexity) / float64(ca.filesCount)
	if avgComplexity > 10 {
		score -= (avgComplexity - 10) * 2
	}

	// Ensure score is within bounds
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// calculateGrade calculates letter grade based on score
func (ca *CodeAnalyzer) calculateGrade(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

// CLI Commands
var rootCmd = &cobra.Command{
	Use:   "nada",
	Short: "A code quality analyzer similar to SonarQube",
	Long: `Nada is a comprehensive code quality analysis tool that scans Go codebases
for bugs, vulnerabilities, code smells, and provides quality metrics and grades.`,
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze a Go codebase",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectPath := args[0]

		analyzer := NewCodeAnalyzer()
		report, err := analyzer.AnalyzeProject(projectPath)
		if err != nil {
			fmt.Printf("❌ Error analyzing project: %v\n", err)
			os.Exit(1)
		}

		printReport(report)

		// Save report if requested
		if outputFile, _ := cmd.Flags().GetString("output"); outputFile != "" {
			saveReport(report, outputFile)
		}
	},
}

// printReport prints the analysis report to console
func printReport(report *Report) {
	fmt.Println("\n🎯 Code Quality Analysis Report")
	fmt.Println("================================")
	fmt.Printf("📁 Project: %s\n", report.ProjectPath)
	fmt.Printf("⏰ Analyzed at: %s\n", report.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("📊 Overall Grade: %s (%.1f/100)\n", report.Grade, report.Score)
	fmt.Printf("📄 Files Analyzed: %d\n", report.FilesAnalyzed)
	fmt.Printf("📏 Lines of Code: %d\n", report.Metrics.LinesOfCode)

	if report.Metrics.CyclomaticComplexity > 0 {
		avgComplexity := float64(report.Metrics.CyclomaticComplexity) / float64(report.FilesAnalyzed)
		fmt.Printf("🔄 Avg Cyclomatic Complexity: %.1f\n", avgComplexity)
	}

	fmt.Printf("🧪 Test Coverage: %.1f%%\n", report.Metrics.TestCoverage)

	// Issues summary
	fmt.Println("\n📋 Issues Summary:")
	fmt.Printf("   🔴 High: %d\n", report.IssuesSummary["high"])
	fmt.Printf("   🟡 Medium: %d\n", report.IssuesSummary["medium"])
	fmt.Printf("   🟢 Low: %d\n", report.IssuesSummary["low"])

	fmt.Printf("   🐛 Bugs: %d\n", report.IssuesSummary["bug"])
	fmt.Printf("   🔒 Vulnerabilities: %d\n", report.IssuesSummary["vulnerability"])
	fmt.Printf("   💨 Code Smells: %d\n", report.IssuesSummary["code_smell"])

	// Show top issues
	if len(report.Issues) > 0 {
		fmt.Println("\n⚠️  Top Issues:")
		count := 0
		for _, issue := range report.Issues {
			if count >= 10 {
				break
			}
			if issue.Severity == "high" || issue.Type == "vulnerability" {
				fmt.Printf("   %s:%d - %s [%s/%s]\n",
					filepath.Base(issue.File), issue.Line, issue.Message,
					issue.Type, issue.Severity)
				count++
			}
		}

		if len(report.Issues) > 10 {
			fmt.Printf("   ... and %d more issues\n", len(report.Issues)-10)
		}
	}

	// Quality gates
	fmt.Println("\n🚪 Quality Gates:")
	fmt.Printf("   Grade A-C: %s\n", getStatus(report.Grade <= "C"))
	fmt.Printf("   High Issues < 5: %s\n", getStatus(report.IssuesSummary["high"] < 5))
	fmt.Printf("   No Vulnerabilities: %s\n", getStatus(report.IssuesSummary["vulnerability"] == 0))
}

func getStatus(passed bool) string {
	if passed {
		return "✅ PASSED"
	}
	return "❌ FAILED"
}

// saveReport saves the report to a file
func saveReport(report *Report, filename string) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("❌ Error marshaling report: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("❌ Error saving report: %v\n", err)
		return
	}

	fmt.Printf("💾 Report saved to: %s\n", filename)
}

func init() {
	// Analyze command flags
	analyzeCmd.Flags().StringP("output", "o", "", "Output file for JSON report")
	rootCmd.AddCommand(analyzeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}
