package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chaksack/nada/internal/rules"
	"github.com/chaksack/nada/internal/types"
)

// CodeAnalyzer performs static code analysis on Go projects
type CodeAnalyzer struct {
	fileSet    *token.FileSet
	issues     []types.Issue
	metrics    types.Metrics
	filesCount int
	options    types.AnalysisOptions
	ruleEngine *rules.Engine
}

// New creates a new CodeAnalyzer instance
func New(options types.AnalysisOptions) *CodeAnalyzer {
	return &CodeAnalyzer{
		fileSet:    token.NewFileSet(),
		issues:     make([]types.Issue, 0),
		options:    options,
		ruleEngine: rules.NewEngine(),
	}
}

// AnalyzeProject analyzes the entire project and returns a report
func (ca *CodeAnalyzer) AnalyzeProject() (*types.Report, error) {
	if ca.options.Verbose {
		fmt.Printf("🔍 Analyzing project: %s\n", ca.options.ProjectPath)
	}

	// Walk through project directory
	err := filepath.WalkDir(ca.options.ProjectPath, ca.walkFunc)
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	// Calculate final metrics and scores
	score := ca.calculateScore()
	grade := ca.calculateGrade(score)
	trends := ca.calculateQualityTrends()
	recommendations := ca.generateRecommendations()

	report := &types.Report{
		ProjectPath:     ca.options.ProjectPath,
		Timestamp:       time.Now(),
		Issues:          ca.issues,
		Metrics:         ca.metrics,
		FilesAnalyzed:   ca.filesCount,
		IssuesSummary:   ca.getIssuesSummary(),
		Score:           score,
		Grade:           grade,
		Trends:          trends,
		Recommendations: recommendations,
	}

	return report, nil
}

// walkFunc is called for each file during directory traversal
func (ca *CodeAnalyzer) walkFunc(path string, d os.DirEntry, err error) error {
	if err != nil {
		return err
	}

	// Skip directories and non-Go files
	if d.IsDir() || !strings.HasSuffix(path, ".go") {
		return nil
	}

	// Skip vendor and .git directories
	if strings.Contains(path, "vendor/") || strings.Contains(path, ".git/") {
		return nil
	}

	// Skip test files if not included
	if !ca.options.IncludeTests && strings.HasSuffix(path, "_test.go") {
		return nil
	}

	// Skip excluded files
	for _, exclude := range ca.options.ExcludeFiles {
		if matched, _ := filepath.Match(exclude, filepath.Base(path)); matched {
			return nil
		}
	}

	ca.analyzeFile(path)
	return nil
}

// analyzeFile analyzes a single Go file
func (ca *CodeAnalyzer) analyzeFile(filePath string) {
	ca.filesCount++

	if ca.options.Verbose {
		fmt.Printf("📄 Analyzing: %s\n", filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		ca.addIssue(types.Issue{
			Type:        types.TypeError,
			Severity:    types.SeverityHigh,
			File:        filePath,
			Line:        1,
			Column:      1,
			Rule:        "read_error",
			Message:     "Cannot read file",
			Description: fmt.Sprintf("Error reading file: %v", err),
			Impact:      types.IssueImpact{EffortMinutes: 1},
		})
		return
	}

	// Parse Go file
	node, err := parser.ParseFile(ca.fileSet, filePath, content, parser.ParseComments)
	if err != nil {
		ca.addIssue(types.Issue{
			Type:        types.TypeError,
			Severity:    types.SeverityHigh,
			File:        filePath,
			Line:        1,
			Column:      1,
			Rule:        "parse_error",
			Message:     "Syntax error",
			Description: fmt.Sprintf("Parse error: %v", err),
			Impact:      types.IssueImpact{EffortMinutes: 1},
		})
		return
	}

	// Apply rule engine
	fileIssues := ca.ruleEngine.AnalyzeFile(filePath, node, string(content), ca.fileSet)
	ca.issues = append(ca.issues, fileIssues...)

	// Update metrics
	ca.updateMetrics(string(content), node)
}

// updateMetrics updates the code metrics based on file analysis
func (ca *CodeAnalyzer) updateMetrics(content string, node *ast.File) {
	lines := strings.Split(content, "\n")
	ca.metrics.LinesOfCode += len(lines)

	// Calculate cyclomatic complexity
	complexity := ca.calculateFileComplexity(node)
	ca.metrics.CyclomaticComplexity += complexity

	// Update test coverage estimation
	if strings.Contains(content, "func Test") {
		ca.metrics.TestCoverage = float64(ca.filesCount) / float64(max(ca.filesCount, 1)) * 100
	}
}

// calculateFileComplexity calculates the cyclomatic complexity of a file
func (ca *CodeAnalyzer) calculateFileComplexity(node *ast.File) int {
	complexity := 0

	ast.Inspect(node, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt,
			*ast.TypeSwitchStmt, *ast.SelectStmt:
			complexity++
		case *ast.CaseClause:
			complexity++
		case *ast.FuncDecl:
			complexity++ // Base complexity for each function
		}
		return true
	})

	return complexity
}

// addIssue adds a new issue to the analyzer
func (ca *CodeAnalyzer) addIssue(issue types.Issue) {
	ca.issues = append(ca.issues, issue)
}

// getIssuesSummary returns a summary of issues by severity and type
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
	summary := ca.getIssuesSummary()

	// Deduct points based on issues
	score := baseScore - float64(summary[types.SeverityHigh]*10) -
		float64(summary[types.SeverityMedium]*5) -
		float64(summary[types.SeverityLow]*1)

	// Factor in complexity
	if ca.filesCount > 0 {
		avgComplexity := float64(ca.metrics.CyclomaticComplexity) / float64(ca.filesCount)
		if avgComplexity > 10 {
			score -= (avgComplexity - 10) * 2
		}
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

// calculateQualityTrends calculates quality trend indicators
func (ca *CodeAnalyzer) calculateQualityTrends() types.QualityTrends {
	trends := types.QualityTrends{}
	summary := ca.getIssuesSummary()

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

		// Calculate various trend metrics
		if ca.metrics.LinesOfCode > 0 {
			trends.IssuesDensity = float64(len(ca.issues)) / float64(ca.metrics.LinesOfCode) * 1000
		}

		securityIssues := summary[types.TypeVulnerability]
		trends.SecurityScore = max(0, 100-float64(securityIssues*10))

		maintainabilityDeductions := float64(securityIssues*5 + summary["missing_documentation"]*2)
		trends.MaintainabilityIndex = max(0, 100-maintainabilityDeductions)

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

// generateRecommendations generates actionable recommendations
func (ca *CodeAnalyzer) generateRecommendations() []string {
	var recommendations []string
	summary := ca.getIssuesSummary()

	if vulnCount := summary[types.TypeVulnerability]; vulnCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🔒 URGENT: Address %d security vulnerabilities immediately", vulnCount))
	}

	if bugCount := summary[types.TypeBug]; bugCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🐛 HIGH: Fix %d bugs to improve reliability", bugCount))
	}

	if ca.metrics.TestCoverage < 70 {
		recommendations = append(recommendations,
			fmt.Sprintf("🧪 Increase test coverage from %.1f%% to at least 70%%", ca.metrics.TestCoverage))
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations,
			"✅ Great job! Your code quality is excellent.")
	}

	return recommendations
}

// Helper function for max
func max[T ~int | ~float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}
