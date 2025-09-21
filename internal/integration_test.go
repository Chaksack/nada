package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chaksack/nada/internal/analyzer"
	"github.com/chaksack/nada/internal/types"
)

// TestFullWorkflow tests the complete analysis workflow
func TestFullWorkflow(t *testing.T) {
	// Create a comprehensive test project
	tmpDir := createComplexTestProject(t)

	// Configure analysis options
	options := types.AnalysisOptions{
		ProjectPath:  tmpDir,
		IncludeTests: true,
		Verbose:      false,
	}

	// Run analysis
	codeAnalyzer := analyzer.New(options)
	report, err := codeAnalyzer.AnalyzeProject()
	if err != nil {
		t.Fatalf("AnalyzeProject() failed: %v", err)
	}

	// Verify report structure
	if report == nil {
		t.Fatal("AnalyzeProject() returned nil report")
	}

	// Check basic report fields
	if report.ProjectPath != tmpDir {
		t.Errorf("Report.ProjectPath = %v, want %v", report.ProjectPath, tmpDir)
	}

	if report.FilesAnalyzed == 0 {
		t.Error("Report.FilesAnalyzed should be > 0")
	}

	if report.Timestamp.IsZero() {
		t.Error("Report.Timestamp should be set")
	}

	// Verify issues were found (our test project has intentional issues)
	if len(report.Issues) == 0 {
		t.Error("Expected to find issues in test project")
	}

	// Check for specific issue types
	hasVulnerability := false
	hasBug := false
	hasCodeSmell := false

	for _, issue := range report.Issues {
		switch issue.Type {
		case types.TypeVulnerability:
			hasVulnerability = true
		case types.TypeBug:
			hasBug = true
		case types.TypeCodeSmell:
			hasCodeSmell = true
		}
	}

	if !hasVulnerability {
		t.Error("Expected to find vulnerability issues")
	}
	if !hasBug {
		t.Error("Expected to find bug issues")
	}
	if !hasCodeSmell {
		t.Error("Expected to find code smell issues")
	}

	// Verify score and grade are reasonable
	if report.Score < 0 || report.Score > 100 {
		t.Errorf("Report.Score = %v, want 0-100", report.Score)
	}

	validGrades := []string{"A", "B", "C", "D", "F"}
	gradeValid := false
	for _, grade := range validGrades {
		if report.Grade == grade {
			gradeValid = true
			break
		}
	}
	if !gradeValid {
		t.Errorf("Report.Grade = %v, want one of %v", report.Grade, validGrades)
	}

	// Verify issues summary
	if report.IssuesSummary == nil {
		t.Error("Report.IssuesSummary should not be nil")
	}

	// Verify recommendations
	if len(report.Recommendations) == 0 {
		t.Error("Report.Recommendations should not be empty")
	}

	// Verify trends
	if report.Trends.CyclomaticComplexityTrend == "" {
		t.Error("Report.Trends.CyclomaticComplexityTrend should be set")
	}
}

// TestAnalysisWithDifferentOptions tests various analysis configurations
func TestAnalysisWithDifferentOptions(t *testing.T) {
	tmpDir := createSimpleTestProject(t)

	tests := []struct {
		name    string
		options types.AnalysisOptions
	}{
		{
			name: "include tests",
			options: types.AnalysisOptions{
				ProjectPath:  tmpDir,
				IncludeTests: true,
				Verbose:      false,
			},
		},
		{
			name: "exclude tests",
			options: types.AnalysisOptions{
				ProjectPath:  tmpDir,
				IncludeTests: false,
				Verbose:      false,
			},
		},
		{
			name: "verbose mode",
			options: types.AnalysisOptions{
				ProjectPath: tmpDir,
				Verbose:     true,
			},
		},
		{
			name: "with exclusions",
			options: types.AnalysisOptions{
				ProjectPath:  tmpDir,
				ExcludeFiles: []string{"*_test.go", "vendor/*"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeAnalyzer := analyzer.New(tt.options)
			report, err := codeAnalyzer.AnalyzeProject()
			if err != nil {
				t.Fatalf("AnalyzeProject() failed: %v", err)
			}

			if report == nil {
				t.Fatal("AnalyzeProject() returned nil report")
			}

			// Basic validation
			if report.FilesAnalyzed < 0 {
				t.Error("FilesAnalyzed should be >= 0")
			}
		})
	}
}

// TestJSONReportSerialization tests JSON report generation
func TestJSONReportSerialization(t *testing.T) {
	tmpDir := createSimpleTestProject(t)

	options := types.AnalysisOptions{
		ProjectPath: tmpDir,
	}

	codeAnalyzer := analyzer.New(options)
	report, err := codeAnalyzer.AnalyzeProject()
	if err != nil {
		t.Fatalf("AnalyzeProject() failed: %v", err)
	}

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatalf("JSON marshaling failed: %v", err)
	}

	// Deserialize back
	var deserializedReport types.Report
	err = json.Unmarshal(jsonData, &deserializedReport)
	if err != nil {
		t.Fatalf("JSON unmarshaling failed: %v", err)
	}

	// Verify key fields
	if deserializedReport.ProjectPath != report.ProjectPath {
		t.Errorf("Deserialized ProjectPath mismatch")
	}
	if deserializedReport.FilesAnalyzed != report.FilesAnalyzed {
		t.Errorf("Deserialized FilesAnalyzed mismatch")
	}
	if len(deserializedReport.Issues) != len(report.Issues) {
		t.Errorf("Deserialized Issues count mismatch")
	}
}

// TestPerformanceWithLargeProject tests performance on a larger codebase
func TestPerformanceWithLargeProject(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	tmpDir := createLargeTestProject(t, 50) // 50 files

	options := types.AnalysisOptions{
		ProjectPath: tmpDir,
		Verbose:     false,
	}

	start := time.Now()
	codeAnalyzer := analyzer.New(options)
	report, err := codeAnalyzer.AnalyzeProject()
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("AnalyzeProject() failed: %v", err)
	}

	if report.FilesAnalyzed != 50 {
		t.Errorf("Expected 50 files analyzed, got %d", report.FilesAnalyzed)
	}

	// Performance check (should complete within reasonable time)
	if duration > 30*time.Second {
		t.Errorf("Analysis took too long: %v", duration)
	}

	t.Logf("Analyzed %d files in %v", report.FilesAnalyzed, duration)
}

// TestErrorHandling tests error conditions
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		projectPath string
		expectError bool
	}{
		{
			name:        "non-existent directory",
			projectPath: "/non/existent/path",
			expectError: true,
		},
		{
			name:        "empty directory",
			projectPath: t.TempDir(),
			expectError: false, // Should not error, just return empty report
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := types.AnalysisOptions{
				ProjectPath: tt.projectPath,
			}

			codeAnalyzer := analyzer.New(options)
			report, err := codeAnalyzer.AnalyzeProject()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError && report == nil {
				t.Error("Expected report but got nil")
			}
		})
	}
}

// Helper functions for creating test projects

func createSimpleTestProject(t *testing.T) string {
	files := map[string]string{
		"main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
		"utils_test.go": `package main

import "testing"

func TestUtils(t *testing.T) {
	// Test implementation
}`,
	}

	return createTestProject(t, files)
}

func createComplexTestProject(t *testing.T) string {
	files := map[string]string{
		"main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
		"bug.go": `package main

import "os"

func main() {
	f, _ := os.Open("file.txt")
	_ = f
}`,
		"security.go": `package main

import "fmt"

func BadSecurity() {
	// This should trigger security rules
	password := "hardcoded_password_123"
	apiKey := "sk-1234567890abcdef"
	
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userInput)
	
	fmt.Println(password, apiKey, query)
}`,
		"complexity.go": `package main

func VeryComplexFunction(x, y, z int) int {
	if x > 0 {
		if y > 0 {
			if z > 0 {
				switch x + y + z {
				case 1:
					if x == 1 {
						return 1
					}
					return 2
				case 2:
					if y == 2 {
						return 3
					}
					return 4
				case 3:
					if z == 3 {
						return 5
					}
					return 6
				default:
					if x > y {
						if y > z {
							return 7
						}
						return 8
					}
					return 9
				}
			}
		}
	}
	return 0
}`,
		"naming.go": `package main

func f() {} // Short function name

func BadNamingExample() {
	var x int    // Short variable name
	var y string // Short variable name
	_ = x
	_ = y
}`,
		"structure.go": `package main

import "fmt"

func TodoExample() {
	// TODO: implement this function properly
	// FIXME: this is a temporary hack
	fmt.Println("This line is way too long and exceeds the recommended line length limit which should trigger a long line warning")
}`,
		"documentation.go": `package main

// This function has documentation
func GoodFunction() {
	// Implementation
}

func BadFunction() {
	// This function lacks documentation
}`,
		"pkg/helper.go": `package pkg

import _ "fmt" // Blank import

func Helper() {
	// Helper function
}`,
	}

	return createTestProject(t, files)
}

func createLargeTestProject(t *testing.T, fileCount int) string {
	files := make(map[string]string)

	// Add main.go
	files["main.go"] = `package main

import "fmt"

func main() {
	fmt.Println("Large test project")
}`

	// Generate multiple files with various issues
	for i := 0; i < fileCount-1; i++ {
		filename := fmt.Sprintf("file%d.go", i)
		content := fmt.Sprintf(`package main

import "fmt"

// TODO: optimize file %d
func Function%d() {
	password := "secret_%d"
	if true {
		if true {
			if true {
				if true {
					fmt.Println("deeply nested", password)
				}
			}
		}
	}
}

func ComplexFunction%d(x int) int {
	switch x {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return 4
	case 5:
		return 5
	default:
		if x > 10 {
			return x * 2
		}
		return x
	}
}`, i, i, i, i)

		files[filename] = content
	}

	return createTestProject(t, files)
}

func createTestProject(t *testing.T, files map[string]string) string {
	tmpDir := t.TempDir()

	for filename, content := range files {
		filePath := filepath.Join(tmpDir, filename)

		// Create directory if needed
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	return tmpDir
}

// Benchmark tests
func BenchmarkFullAnalysis(b *testing.B) {
	tmpDir := createSimpleTestProject(&testing.T{}) // Create test project once

	options := types.AnalysisOptions{
		ProjectPath: tmpDir,
		Verbose:     false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codeAnalyzer := analyzer.New(options)
		_, err := codeAnalyzer.AnalyzeProject()
		if err != nil {
			b.Fatalf("AnalyzeProject() failed: %v", err)
		}
	}
}

func BenchmarkComplexAnalysis(b *testing.B) {
	tmpDir := createComplexTestProject(&testing.T{})

	options := types.AnalysisOptions{
		ProjectPath: tmpDir,
		Verbose:     false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codeAnalyzer := analyzer.New(options)
		_, err := codeAnalyzer.AnalyzeProject()
		if err != nil {
			b.Fatalf("AnalyzeProject() failed: %v", err)
		}
	}
}
