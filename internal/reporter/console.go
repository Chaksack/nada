package reporter

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/chaksack/nada/internal/types"
)

// PrintConsoleReport prints a formatted report to the console
func PrintConsoleReport(report *types.Report) {
	printHeader(report)
	printMetrics(report)
	printIssuesSummary(report)
	printTopIssues(report)
	printRecommendations(report)
	printQualityGates(report)
	printTrends(report)
}

// printHeader prints the report header
func printHeader(report *types.Report) {
	fmt.Println("\n🎯 Nada Code Quality Report")
	fmt.Println("===========================")
	fmt.Printf("📁 Project: %s\n", report.ProjectPath)
	fmt.Printf("⏰ Analyzed: %s\n", report.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("📊 Grade: %s (%.1f/100)\n", getGradeEmoji(report.Grade), report.Score)
}

// printMetrics prints project metrics
func printMetrics(report *types.Report) {
	fmt.Println("\n📊 Project Metrics:")
	fmt.Printf("   📄 Files Analyzed: %d\n", report.FilesAnalyzed)
	fmt.Printf("   📏 Lines of Code: %d\n", report.Metrics.LinesOfCode)

	if report.Metrics.CyclomaticComplexity > 0 && report.FilesAnalyzed > 0 {
		avgComplexity := float64(report.Metrics.CyclomaticComplexity) / float64(report.FilesAnalyzed)
		fmt.Printf("   🔄 Avg Complexity: %.1f\n", avgComplexity)
	}

	fmt.Printf("   🧪 Test Coverage: %.1f%%\n", report.Metrics.TestCoverage)
}

// printIssuesSummary prints the issues summary
func printIssuesSummary(report *types.Report) {
	fmt.Println("\n📋 Issues Summary:")

	totalIssues := len(report.Issues)
	fmt.Printf("   📊 Total Issues: %d\n", totalIssues)

	// By severity
	fmt.Printf("   🔴 High: %d\n", report.IssuesSummary[types.SeverityHigh])
	fmt.Printf("   🟡 Medium: %d\n", report.IssuesSummary[types.SeverityMedium])
	fmt.Printf("   🟢 Low: %d\n", report.IssuesSummary[types.SeverityLow])

	// By type
	fmt.Printf("   🐛 Bugs: %d\n", report.IssuesSummary[types.TypeBug])
	fmt.Printf("   🔒 Vulnerabilities: %d\n", report.IssuesSummary[types.TypeVulnerability])
	fmt.Printf("   💨 Code Smells: %d\n", report.IssuesSummary[types.TypeCodeSmell])
}

// printTopIssues prints the most critical issues
func printTopIssues(report *types.Report) {
	if len(report.Issues) == 0 {
		fmt.Println("\n✅ No issues found!")
		return
	}

	fmt.Println("\n⚠️  Top Issues:")

	// Sort issues by priority: vulnerabilities first, then by severity
	sortedIssues := make([]types.Issue, len(report.Issues))
	copy(sortedIssues, report.Issues)

	sort.Slice(sortedIssues, func(i, j int) bool {
		// Vulnerabilities always come first
		if sortedIssues[i].Type == types.TypeVulnerability && sortedIssues[j].Type != types.TypeVulnerability {
			return true
		}
		if sortedIssues[j].Type == types.TypeVulnerability && sortedIssues[i].Type != types.TypeVulnerability {
			return false
		}

		// Then by severity
		severityOrder := map[string]int{
			types.SeverityHigh:   3,
			types.SeverityMedium: 2,
			types.SeverityLow:    1,
		}

		return severityOrder[sortedIssues[i].Severity] > severityOrder[sortedIssues[j].Severity]
	})

	// Show top 10 issues
	count := 0
	maxIssues := 10
	if len(sortedIssues) < maxIssues {
		maxIssues = len(sortedIssues)
	}

	for _, issue := range sortedIssues {
		if count >= maxIssues {
			break
		}

		// Show high priority issues
		if issue.Severity == types.SeverityHigh || issue.Type == types.TypeVulnerability {
			fmt.Printf("   %s:%d - %s [%s/%s]\n",
				filepath.Base(issue.File), issue.Line, issue.Message,
				getTypeEmoji(issue.Type), getSeverityEmoji(issue.Severity))
			count++
		}
	}

	if len(report.Issues) > maxIssues {
		fmt.Printf("   ... and %d more issues\n", len(report.Issues)-count)
	}
}

// printRecommendations prints actionable recommendations
func printRecommendations(report *types.Report) {
	if len(report.Recommendations) == 0 {
		return
	}

	fmt.Println("\n💡 Recommendations:")
	for _, rec := range report.Recommendations {
		fmt.Printf("   %s\n", rec)
	}
}

// printQualityGates prints quality gate results
func printQualityGates(report *types.Report) {
	fmt.Println("\n🚪 Quality Gates:")

	gates := []struct {
		name      string
		condition func(*types.Report) bool
		message   string
	}{
		{
			name:      "Grade A-C",
			condition: func(r *types.Report) bool { return r.Grade <= "C" },
			message:   "Maintain good code quality grade",
		},
		{
			name:      "High Issues < 5",
			condition: func(r *types.Report) bool { return r.IssuesSummary[types.SeverityHigh] < 5 },
			message:   "Keep high-severity issues under control",
		},
		{
			name:      "No Vulnerabilities",
			condition: func(r *types.Report) bool { return r.IssuesSummary[types.TypeVulnerability] == 0 },
			message:   "Ensure no security vulnerabilities",
		},
		{
			name:      "Score > 70",
			condition: func(r *types.Report) bool { return r.Score > 70 },
			message:   "Maintain acceptable quality score",
		},
	}

	for _, gate := range gates {
		status := "❌ FAILED"
		if gate.condition(report) {
			status = "✅ PASSED"
		}
		fmt.Printf("   %s: %s\n", gate.name, status)
	}
}

// printTrends prints quality trend analysis
func printTrends(report *types.Report) {
	if report.Trends == (types.QualityTrends{}) {
		return
	}

	fmt.Println("\n📈 Quality Trends:")

	if report.Trends.CyclomaticComplexityTrend != "" {
		fmt.Printf("   🔄 Complexity: %s\n", report.Trends.CyclomaticComplexityTrend)
	}

	if report.Trends.IssuesDensity > 0 {
		fmt.Printf("   📊 Issues Density: %.1f per 1000 LOC\n", report.Trends.IssuesDensity)
	}

	if report.Trends.SecurityScore > 0 {
		fmt.Printf("   🔒 Security Score: %.1f/100\n", report.Trends.SecurityScore)
	}

	if report.Trends.MaintainabilityIndex > 0 {
		fmt.Printf("   🛠️  Maintainability: %.1f/100\n", report.Trends.MaintainabilityIndex)
	}

	if report.Trends.TechnicalDebtRatio > 0 {
		fmt.Printf("   ⏰ Technical Debt: %.1f hours per 1000 LOC\n", report.Trends.TechnicalDebtRatio)
	}
}

// Helper functions for emojis and formatting

func getGradeEmoji(grade string) string {
	switch grade {
	case "A":
		return "🟢 A"
	case "B":
		return "🔵 B"
	case "C":
		return "🟡 C"
	case "D":
		return "🟠 D"
	case "F":
		return "🔴 F"
	default:
		return grade
	}
}

func getTypeEmoji(issueType string) string {
	switch issueType {
	case types.TypeBug:
		return "🐛"
	case types.TypeVulnerability:
		return "🔒"
	case types.TypeCodeSmell:
		return "💨"
	case types.TypeError:
		return "❌"
	default:
		return "❓"
	}
}

func getSeverityEmoji(severity string) string {
	switch severity {
	case types.SeverityHigh:
		return "🔴"
	case types.SeverityMedium:
		return "🟡"
	case types.SeverityLow:
		return "🟢"
	default:
		return "⚪"
	}
}
