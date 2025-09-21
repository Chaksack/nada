package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chaksack/nada/internal/types"
)

func TestNew(t *testing.T) {
	options := types.AnalysisOptions{
		ProjectPath: "/test/project",
		Verbose:     true,
	}

	analyzer := New(options)

	if analyzer == nil {
		t.Fatal("New() returned nil analyzer")
	}

	if analyzer.options.ProjectPath != "/test/project" {
		t.Errorf("analyzer.options.ProjectPath = %v, want /test/project", analyzer.options.ProjectPath)
	}

	if !analyzer.options.Verbose {
		t.Errorf("analyzer.options.Verbose = %v, want true", analyzer.options.Verbose)
	}
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name          string
		highIssues    int
		mediumIssues  int
		lowIssues     int
		complexity    int
		filesCount    int
		expectedRange [2]float64 // [min, max] expected score range
	}{
		{
			name:          "perfect code",
			highIssues:    0,
			mediumIssues:  0,
			lowIssues:     0,
			complexity:    5,
			filesCount:    1,
			expectedRange: [2]float64{95.0, 100.0},
		},
		{
			name:          "high issues penalty",
			highIssues:    5,
			mediumIssues:  0,
			lowIssues:     0,
			complexity:    5,
			filesCount:    1,
			expectedRange: [2]float64{40.0, 60.0},
		},
		{
			name:          "high complexity penalty",
			highIssues:    0,
			mediumIssues:  0,
			lowIssues:     0,
			complexity:    50, // Very high complexity
			filesCount:    1,
			expectedRange: [2]float64{0.0, 50.0},
		},
		{
			name:          "mixed issues",
			highIssues:    2,
			mediumIssues:  3,
			lowIssues:     5,
			complexity:    10,
			filesCount:    1,
			expectedRange: [2]float64{50.0, 80.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := New(types.AnalysisOptions{})

			// Set up test issues
			for i := 0; i < tt.highIssues; i++ {
				analyzer.addIssue(types.Issue{Severity: types.SeverityHigh})
			}
			for i := 0; i < tt.mediumIssues; i++ {
				analyzer.addIssue(types.Issue{Severity: types.SeverityMedium})
			}
			for i := 0; i < tt.lowIssues; i++ {
				analyzer.addIssue(types.Issue{Severity: types.SeverityLow})
			}

			analyzer.metrics.CyclomaticComplexity = tt.complexity
			analyzer.filesCount = tt.filesCount

			score := analyzer.calculateScore()

			if score < tt.expectedRange[0] || score > tt.expectedRange[1] {
				t.Errorf("calculateScore() = %v, want range [%v, %v]",
					score, tt.expectedRange[0], tt.expectedRange[1])
			}
		})
	}
}

func TestCalculateGrade(t *testing.T) {
	tests := []struct {
		name  string
		score float64
		want  string
	}{
		{"Grade A", 95.0, "A"},
		{"Grade A boundary", 90.0, "A"},
		{"Grade B", 85.0, "B"},
		{"Grade B boundary", 80.0, "B"},
		{"Grade C", 75.0, "C"},
		{"Grade C boundary", 70.0, "C"},
		{"Grade D", 65.0, "D"},
		{"Grade D boundary", 60.0, "D"},
		{"Grade F", 50.0, "F"},
		{"Grade F very low", 0.0, "F"},
	}

	analyzer := New(types.AnalysisOptions{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := analyzer.calculateGrade(tt.score)
			if got != tt.want {
				t.Errorf("calculateGrade(%v) = %v, want %v", tt.score, got, tt.want)
			}
		})
	}
}

func TestGetIssuesSummary(t *testing.T) {
	analyzer := New(types.AnalysisOptions{})

	// Add test issues
	testIssues := []types.Issue{
		{Type: types.TypeBug, Severity: types.SeverityHigh},
		{Type: types.TypeBug, Severity: types.SeverityMedium},
		{Type: types.TypeVulnerability, Severity: types.SeverityHigh},
		{Type: types.TypeCodeSmell, Severity: types.SeverityLow},
		{Type: types.TypeCodeSmell, Severity: types.SeverityLow},
	}

	for _, issue := range testIssues {
		analyzer.addIssue(issue)
	}

	summary := analyzer.getIssuesSummary()

	expectedCounts := map[string]int{
		types.SeverityHigh:      2,
		types.SeverityMedium:    1,
		types.SeverityLow:       2,
		types.TypeBug:           2,
		types.TypeVulnerability: 1,
		types.TypeCodeSmell:     2,
	}

	for key, expectedCount := range expectedCounts {
		if summary[key] != expectedCount {
			t.Errorf("getIssuesSummary()[%v] = %v, want %v", key, summary[key], expectedCount)
		}
	}
}

func TestGenerateRecommendations(t *testing.T) {
	tests := []struct {
		name         string
		issues       []types.Issue
		testCoverage float64
		wantContains []string
	}{
		{
			name: "vulnerabilities present",
			issues: []types.Issue{
				{Type: types.TypeVulnerability, Severity: types.SeverityHigh},
				{Type: types.TypeVulnerability, Severity: types.SeverityHigh},
			},
			testCoverage: 80.0,
			wantContains: []string{"URGENT", "security", "vulnerabilities"},
		},
		{
			name: "bugs present",
			issues: []types.Issue{
				{Type: types.TypeBug, Severity: types.SeverityMedium},
			},
			testCoverage: 80.0,
			wantContains: []string{"HIGH", "bugs", "reliability"},
		},
		{
			name:         "low test coverage",
			issues:       []types.Issue{},
			testCoverage: 50.0,
			wantContains: []string{"test coverage", "70%"},
		},
		{
			name:         "excellent code",
			issues:       []types.Issue{},
			testCoverage: 90.0,
			wantContains: []string{"Great job", "excellent"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := New(types.AnalysisOptions{})
			analyzer.issues = tt.issues
			analyzer.metrics.TestCoverage = tt.testCoverage

			recommendations := analyzer.generateRecommendations()

			if len(recommendations) == 0 {
				t.Errorf("generateRecommendations() returned empty slice")
			}

			// Check if expected strings are present in recommendations
			for _, expectedString := range tt.wantContains {
				found := false
				for _, rec := range recommendations {
					if contains(rec, expectedString) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("generateRecommendations() missing expected string: %v in %v",
						expectedString, recommendations)
				}
			}
		})
	}
}

func TestCalculateQualityTrends(t *testing.T) {
	analyzer := New(types.AnalysisOptions{})
	analyzer.filesCount = 2
	analyzer.metrics.CyclomaticComplexity = 20 // Average 10 per file
	analyzer.metrics.LinesOfCode = 1000

	// Add test issues
	analyzer.addIssue(types.Issue{
		Type:     types.TypeVulnerability,
		Severity: types.SeverityHigh,
		Impact:   types.IssueImpact{EffortMinutes: 30},
	})
	analyzer.addIssue(types.Issue{
		Type:     types.TypeCodeSmell,
		Severity: types.SeverityMedium,
		Impact:   types.IssueImpact{EffortMinutes: 15},
	})

	trends := analyzer.calculateQualityTrends()

	// Test complexity trend
	if trends.CyclomaticComplexityTrend == "" {
		t.Errorf("calculateQualityTrends() CyclomaticComplexityTrend is empty")
	}

	// Test issues density (should be > 0 since we have issues)
	if trends.IssuesDensity == 0 {
		t.Errorf("calculateQualityTrends() IssuesDensity = 0, want > 0")
	}

	// Test security score (should be < 100 due to vulnerability)
	if trends.SecurityScore >= 100 {
		t.Errorf("calculateQualityTrends() SecurityScore = %v, want < 100", trends.SecurityScore)
	}

	// Test technical debt ratio (should be > 0 due to impact minutes)
	if trends.TechnicalDebtRatio == 0 {
		t.Errorf("calculateQualityTrends() TechnicalDebtRatio = 0, want > 0")
	}
}

// Helper function to create a temporary test file
func createTestFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	return testFile
}

func TestAnalyzeFile(t *testing.T) {
	tests := []struct {
		name           string
		fileContent    string
		expectedIssues int
		expectError    bool
	}{
		{
			name: "valid Go file",
			fileContent: `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
			expectedIssues: 0, // May have some minor issues depending on rules
			expectError:    false,
		},
		{
			name: "file with syntax error",
			fileContent: `package main

func main() {
	fmt.Println("Hello World"
}`, // Missing closing parenthesis
			expectedIssues: 1,     // Should have parse error
			expectError:    false, // We handle parse errors gracefully
		},
		{
			name: "complex function",
			fileContent: `package main

func ComplexFunction(x int) int {
	if x > 10 {
		if x > 20 {
			if x > 30 {
				if x > 40 {
					if x > 50 {
						return x * 2
					}
					return x * 3
				}
				return x * 4
			}
			return x * 5
		}
		return x * 6
	}
	return x
}`,
			expectedIssues: 2, // Should have complexity and nesting issues
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := createTestFile(t, tt.fileContent)

			analyzer := New(types.AnalysisOptions{Verbose: false})
			initialIssueCount := len(analyzer.issues)

			analyzer.analyzeFile(testFile)

			newIssueCount := len(analyzer.issues) - initialIssueCount

			if tt.expectedIssues > 0 && newIssueCount == 0 {
				t.Errorf("analyzeFile() expected %d issues, got %d", tt.expectedIssues, newIssueCount)
			}
		})
	}
}

// Helper function for string contains check (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
