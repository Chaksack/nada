package reporter

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/chaksack/nada/internal/types"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintConsoleReport(t *testing.T) {
	// Create test report
	report := &types.Report{
		ProjectPath:   "/test/project",
		Timestamp:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Grade:         "B",
		Score:         85.0,
		FilesAnalyzed: 10,
		Issues: []types.Issue{
			{
				Type:     types.TypeVulnerability,
				Severity: types.SeverityHigh,
				File:     "main.go",
				Line:     10,
				Message:  "SQL injection vulnerability",
				Rule:     "sql_injection",
			},
			{
				Type:     types.TypeBug,
				Severity: types.SeverityMedium,
				File:     "utils.go",
				Line:     25,
				Message:  "Missing error handling",
				Rule:     "error_handling",
			},
		},
		Metrics: types.Metrics{
			LinesOfCode:          1000,
			CyclomaticComplexity: 150,
			TestCoverage:         75.5,
		},
		IssuesSummary: map[string]int{
			types.SeverityHigh:      1,
			types.SeverityMedium:    1,
			types.SeverityLow:       0,
			types.TypeVulnerability: 1,
			types.TypeBug:           1,
			types.TypeCodeSmell:     0,
		},
		Trends: types.QualityTrends{
			CyclomaticComplexityTrend: "High - Consider refactoring",
			IssuesDensity:             2.0,
			SecurityScore:             90.0,
			MaintainabilityIndex:      85.0,
			TechnicalDebtRatio:        1.5,
		},
		Recommendations: []string{
			"🔒 URGENT: Address 1 security vulnerabilities",
			"🐛 HIGH: Fix 1 bugs to improve reliability",
		},
	}

	output := captureOutput(func() {
		PrintConsoleReport(report)
	})

	// Test for expected content
	expectedStrings := []string{
		"🎯 Nada Code Quality Report",
		"Project: /test/project",
		"Grade: 🔵 B (85.0/100)",
		"Files Analyzed: 10",
		"Lines of Code: 1000",
		"Test Coverage: 75.5%",
		"🔴 High: 1",
		"🟡 Medium: 1",
		"🟢 Low: 0",
		"🐛 Bugs: 1",
		"🔒 Vulnerabilities: 1",
		"💨 Code Smells: 0",
		"SQL injection vulnerability",

		"Quality Gates:",
		"Quality Trends:",
		"Recommendations:",
		"URGENT: Address 1 security vulnerabilities",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("PrintConsoleReport() output missing expected string: %q\nOutput: %s", expected, output)
		}
	}
}

func TestPrintConsoleReportNoIssues(t *testing.T) {
	report := &types.Report{
		ProjectPath:   "/perfect/project",
		Timestamp:     time.Now(),
		Grade:         "A",
		Score:         100.0,
		FilesAnalyzed: 5,
		Issues:        []types.Issue{},
		Metrics: types.Metrics{
			LinesOfCode:  500,
			TestCoverage: 95.0,
		},
		IssuesSummary: map[string]int{},
		Recommendations: []string{
			"✅ Great job! Your code quality is excellent.",
		},
	}

	output := captureOutput(func() {
		PrintConsoleReport(report)
	})

	expectedStrings := []string{
		"Grade: 🟢 A (100.0/100)",
		"✅ No issues found!",
		"Great job! Your code quality is excellent.",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("PrintConsoleReport() output missing expected string: %q", expected)
		}
	}
}

func TestGetGradeEmoji(t *testing.T) {
	tests := []struct {
		grade string
		want  string
	}{
		{"A", "🟢 A"},
		{"B", "🔵 B"},
		{"C", "🟡 C"},
		{"D", "🟠 D"},
		{"F", "🔴 F"},
		{"X", "X"}, // Unknown grade
	}

	for _, tt := range tests {
		t.Run(tt.grade, func(t *testing.T) {
			got := getGradeEmoji(tt.grade)
			if got != tt.want {
				t.Errorf("getGradeEmoji(%q) = %q, want %q", tt.grade, got, tt.want)
			}
		})
	}
}

func TestGetTypeEmoji(t *testing.T) {
	tests := []struct {
		issueType string
		want      string
	}{
		{types.TypeBug, "🐛"},
		{types.TypeVulnerability, "🔒"},
		{types.TypeCodeSmell, "💨"},
		{types.TypeError, "❌"},
		{"unknown", "❓"},
	}

	for _, tt := range tests {
		t.Run(tt.issueType, func(t *testing.T) {
			got := getTypeEmoji(tt.issueType)
			if got != tt.want {
				t.Errorf("getTypeEmoji(%q) = %q, want %q", tt.issueType, got, tt.want)
			}
		})
	}
}

func TestGetSeverityEmoji(t *testing.T) {
	tests := []struct {
		severity string
		want     string
	}{
		{types.SeverityHigh, "🔴"},
		{types.SeverityMedium, "🟡"},
		{types.SeverityLow, "🟢"},
		{"unknown", "⚪"},
	}

	for _, tt := range tests {
		t.Run(tt.severity, func(t *testing.T) {
			got := getSeverityEmoji(tt.severity)
			if got != tt.want {
				t.Errorf("getSeverityEmoji(%q) = %q, want %q", tt.severity, got, tt.want)
			}
		})
	}
}

func TestPrintHeader(t *testing.T) {
	report := &types.Report{
		ProjectPath: "/test/project",
		Timestamp:   time.Date(2024, 1, 1, 12, 30, 45, 0, time.UTC),
		Grade:       "A",
		Score:       95.5,
	}

	output := captureOutput(func() {
		printHeader(report)
	})

	expectedStrings := []string{
		"🎯 Nada Code Quality Report",
		"===========================",
		"📁 Project: /test/project",
		"⏰ Analyzed: 2024-01-01 12:30:45",
		"📊 Grade: 🟢 A (95.5/100)",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printHeader() output missing expected string: %q", expected)
		}
	}
}

func TestPrintMetrics(t *testing.T) {
	report := &types.Report{
		FilesAnalyzed: 15,
		Metrics: types.Metrics{
			LinesOfCode:          2500,
			CyclomaticComplexity: 75,
			TestCoverage:         82.3,
		},
	}

	output := captureOutput(func() {
		printMetrics(report)
	})

	expectedStrings := []string{
		"📊 Project Metrics:",
		"📄 Files Analyzed: 15",
		"📏 Lines of Code: 2500",
		"🔄 Avg Complexity: 5.0", // 75/15
		"🧪 Test Coverage: 82.3%",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printMetrics() output missing expected string: %q", expected)
		}
	}
}

func TestPrintIssuesSummary(t *testing.T) {
	report := &types.Report{
		Issues: []types.Issue{
			{Type: types.TypeBug, Severity: types.SeverityHigh},
			{Type: types.TypeVulnerability, Severity: types.SeverityHigh},
			{Type: types.TypeCodeSmell, Severity: types.SeverityMedium},
		},
		IssuesSummary: map[string]int{
			types.SeverityHigh:      2,
			types.SeverityMedium:    1,
			types.SeverityLow:       0,
			types.TypeBug:           1,
			types.TypeVulnerability: 1,
			types.TypeCodeSmell:     1,
		},
	}

	output := captureOutput(func() {
		printIssuesSummary(report)
	})

	expectedStrings := []string{
		"📋 Issues Summary:",
		"📊 Total Issues: 3",
		"🔴 High: 2",
		"🟡 Medium: 1",
		"🟢 Low: 0",
		"🐛 Bugs: 1",
		"🔒 Vulnerabilities: 1",
		"💨 Code Smells: 1",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printIssuesSummary() output missing expected string: %q", expected)
		}
	}
}

func TestPrintTopIssues(t *testing.T) {
	report := &types.Report{
		Issues: []types.Issue{
			{
				Type:     types.TypeVulnerability,
				Severity: types.SeverityHigh,
				File:     "/path/to/file.go",
				Line:     42,
				Message:  "SQL injection found",
			},
			{
				Type:     types.TypeBug,
				Severity: types.SeverityHigh,
				File:     "/path/to/another.go",
				Line:     10,
				Message:  "Null pointer dereference",
			},
			{
				Type:     types.TypeCodeSmell,
				Severity: types.SeverityMedium,
				File:     "/path/to/smelly.go",
				Line:     5,
				Message:  "Function too complex",
			},
		},
	}

	output := captureOutput(func() {
		printTopIssues(report)
	})

	expectedStrings := []string{
		"⚠️  Top Issues:",
		"file.go:42 - SQL injection found [🔒/🔴]",
		"another.go:10 - Null pointer dereference [🐛/🔴]",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printTopIssues() output missing expected string: %q", expected)
		}
	}
}

func TestPrintTopIssuesNoIssues(t *testing.T) {
	report := &types.Report{
		Issues: []types.Issue{},
	}

	output := captureOutput(func() {
		printTopIssues(report)
	})

	if !strings.Contains(output, "✅ No issues found!") {
		t.Errorf("printTopIssues() should show 'No issues found' message")
	}
}

func TestPrintRecommendations(t *testing.T) {
	report := &types.Report{
		Recommendations: []string{
			"🔒 URGENT: Address security vulnerabilities",
			"🧪 Increase test coverage to 80%",
			"🔧 Refactor complex functions",
		},
	}

	output := captureOutput(func() {
		printRecommendations(report)
	})

	expectedStrings := []string{
		"💡 Recommendations:",
		"🔒 URGENT: Address security vulnerabilities",
		"🧪 Increase test coverage to 80%",
		"🔧 Refactor complex functions",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printRecommendations() output missing expected string: %q", expected)
		}
	}
}

func TestPrintRecommendationsEmpty(t *testing.T) {
	report := &types.Report{
		Recommendations: []string{},
	}

	output := captureOutput(func() {
		printRecommendations(report)
	})

	// Should not print anything if no recommendations
	if strings.Contains(output, "💡 Recommendations:") {
		t.Errorf("printRecommendations() should not print header when no recommendations")
	}
}

func TestPrintQualityGates(t *testing.T) {
	tests := []struct {
		name   string
		report *types.Report
		expect []string
	}{
		{
			name: "some gates fail",
			report: &types.Report{
				Grade: "F",
				Score: 45.0,
				IssuesSummary: map[string]int{
					types.SeverityHigh:      10,
					types.TypeVulnerability: 2,
				},
			},
			expect: []string{
				"❌ FAILED",
				"❌ FAILED",
				"❌ FAILED",
				"❌ FAILED",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				printQualityGates(tt.report)
			})

			if !strings.Contains(output, "🚪 Quality Gates:") {
				t.Errorf("printQualityGates() output missing header")
			}

			for _, expected := range tt.expect {
				if !strings.Contains(output, expected) {
					t.Errorf("printQualityGates() output missing expected status: %q", expected)
				}
			}
		})
	}
}

func TestPrintTrends(t *testing.T) {
	report := &types.Report{
		Trends: types.QualityTrends{
			CyclomaticComplexityTrend: "High - Consider refactoring",
			IssuesDensity:             2.5,
			SecurityScore:             85.0,
			MaintainabilityIndex:      78.5,
			TechnicalDebtRatio:        1.2,
		},
	}

	output := captureOutput(func() {
		printTrends(report)
	})

	expectedStrings := []string{
		"📈 Quality Trends:",
		"🔄 Complexity: High - Consider refactoring",
		"📊 Issues Density: 2.5 per 1000 LOC",
		"🔒 Security Score: 85.0/100",
		"🛠️  Maintainability: 78.5/100",
		"⏰ Technical Debt: 1.2 hours per 1000 LOC",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printTrends() output missing expected string: %q", expected)
		}
	}
}

func TestPrintTrendsEmpty(t *testing.T) {
	report := &types.Report{
		Trends: types.QualityTrends{}, // Empty trends
	}

	output := captureOutput(func() {
		printTrends(report)
	})

	// Should not print anything if trends are empty
	if strings.Contains(output, "📈 Quality Trends:") {
		t.Errorf("printTrends() should not print header when trends are empty")
	}
}
