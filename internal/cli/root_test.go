package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/chaksack/nada/internal/types"
)

func TestSetVersionInfo(t *testing.T) {
	testVersion := "1.0.0"
	testBuildTime := "2024-01-01T00:00:00Z"
	testCommit := "abc123"

	SetVersionInfo(testVersion, testBuildTime, testCommit)

	if version != testVersion {
		t.Errorf("SetVersionInfo() version = %v, want %v", version, testVersion)
	}
	if buildTime != testBuildTime {
		t.Errorf("SetVersionInfo() buildTime = %v, want %v", buildTime, testBuildTime)
	}
	if commit != testCommit {
		t.Errorf("SetVersionInfo() commit = %v, want %v", commit, testCommit)
	}
}

func TestRunVersion(t *testing.T) {
	// Set test version info
	SetVersionInfo("1.0.0-test", "2024-01-01", "abc123")

	// Capture stdout
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)

	// Mock the version command output
	runVersion(cmd, []string{})

	output := buf.String()

	// Check for expected content in output
	expectedStrings := []string{
		"Nada 1.0.0-test",
		"Build time: 2024-01-01",
		"Commit: abc123",
		"Andrew Chakdahah",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("runVersion() output missing expected string: %v\nGot: %v", expected, output)
		}
	}
}

func TestSaveReport(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Test report
	report := &types.Report{
		ProjectPath:   "/test/project",
		Grade:         "A",
		Score:         95.0,
		FilesAnalyzed: 10,
		Issues:        []types.Issue{},
		IssuesSummary: map[string]int{
			"high":   0,
			"medium": 2,
			"low":    5,
		},
	}

	// Test file path
	reportFile := filepath.Join(tmpDir, "test-report.json")

	// Save report
	err := saveReport(report, reportFile)
	if err != nil {
		t.Fatalf("saveReport() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(reportFile); os.IsNotExist(err) {
		t.Errorf("saveReport() did not create file: %v", reportFile)
	}

	// Read and verify file content
	data, err := os.ReadFile(reportFile)
	if err != nil {
		t.Fatalf("Failed to read saved report: %v", err)
	}

	// Parse JSON to verify it's valid
	var savedReport types.Report
	err = json.Unmarshal(data, &savedReport)
	if err != nil {
		t.Fatalf("Saved report is not valid JSON: %v", err)
	}

	// Verify key fields
	if savedReport.Grade != report.Grade {
		t.Errorf("Saved report Grade = %v, want %v", savedReport.Grade, report.Grade)
	}
	if savedReport.Score != report.Score {
		t.Errorf("Saved report Score = %v, want %v", savedReport.Score, report.Score)
	}
}

func TestSaveReportWithNestedDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	report := &types.Report{
		ProjectPath: "/test/project",
		Grade:       "B",
		Score:       85.0,
	}

	// Test nested directory creation
	reportFile := filepath.Join(tmpDir, "reports", "nested", "test-report.json")

	err := saveReport(report, reportFile)
	if err != nil {
		t.Fatalf("saveReport() with nested directory failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(reportFile); os.IsNotExist(err) {
		t.Errorf("saveReport() did not create nested file: %v", reportFile)
	}
}

func TestShouldFailBuild(t *testing.T) {
	tests := []struct {
		name   string
		report *types.Report
		want   bool
	}{
		{
			name: "should fail - high severity issues",
			report: &types.Report{
				Grade: "B",
				IssuesSummary: map[string]int{
					types.SeverityHigh:   2,
					types.SeverityMedium: 1,
				},
			},
			want: true,
		},
		{
			name: "should fail - vulnerabilities",
			report: &types.Report{
				Grade: "B",
				IssuesSummary: map[string]int{
					types.TypeVulnerability: 1,
				},
			},
			want: true,
		},
		{
			name: "should fail - grade F",
			report: &types.Report{
				Grade: "F",
				IssuesSummary: map[string]int{
					types.SeverityMedium: 5,
				},
			},
			want: true,
		},
		{
			name: "should pass - good quality",
			report: &types.Report{
				Grade: "B",
				IssuesSummary: map[string]int{
					types.SeverityMedium: 2,
					types.SeverityLow:    5,
				},
			},
			want: false,
		},
		{
			name: "should pass - perfect quality",
			report: &types.Report{
				Grade:         "A",
				IssuesSummary: map[string]int{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldFailBuild(tt.report)
			if got != tt.want {
				t.Errorf("shouldFailBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunAnalyze(t *testing.T) {
	// Create a temporary Go project for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "main.go")

	testCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`

	err := os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name         string
		args         []string
		flags        map[string]string
		expectError  bool
		expectedFile string
	}{
		{
			name:        "analyze current directory",
			args:        []string{tmpDir},
			expectError: false,
		},
		{
			name: "analyze with output file",
			args: []string{tmpDir},
			flags: map[string]string{
				"output": filepath.Join(tmpDir, "report.json"),
			},
			expectError:  false,
			expectedFile: filepath.Join(tmpDir, "report.json"),
		},
		{
			name:        "analyze non-existent directory",
			args:        []string{"/non/existent/path"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command for each test
			cmd := &cobra.Command{
				Use:  "analyze",
				RunE: runAnalyze,
			}

			// Add flags
			cmd.Flags().String("output", "", "Output file")
			cmd.Flags().String("coverage", "", "Coverage file")
			cmd.Flags().String("diff", "", "Diff target")
			cmd.Flags().String("config", "", "Config file")
			cmd.Flags().Bool("verbose", false, "Verbose output")
			cmd.Flags().Bool("include-tests", false, "Include tests")
			cmd.Flags().StringSlice("exclude", []string{}, "Exclude patterns")

			// Set flags
			for flag, value := range tt.flags {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					t.Fatalf("Failed to set flag %s: %v", flag, err)
				}
			}

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// Run command
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("runAnalyze() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("runAnalyze() unexpected error: %v", err)
			}

			// Check if expected output file was created
			if tt.expectedFile != "" {
				if _, err := os.Stat(tt.expectedFile); os.IsNotExist(err) {
					t.Errorf("runAnalyze() expected output file not created: %v", tt.expectedFile)
				}
			}
		})
	}
}

func TestCommandInitialization(t *testing.T) {
	// Test that commands are properly initialized
	if analyzeCmd == nil {
		t.Error("analyzeCmd is nil")
	}

	if versionCmd == nil {
		t.Error("versionCmd is nil")
	}

	if rootCmd == nil {
		t.Error("rootCmd is nil")
	}

	// Test analyze command flags
	outputFlag := analyzeCmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("analyze command missing output flag")
	}

	verboseFlag := analyzeCmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("analyze command missing verbose flag")
	}

	includeTestsFlag := analyzeCmd.Flags().Lookup("include-tests")
	if includeTestsFlag == nil {
		t.Error("analyze command missing include-tests flag")
	}
}



// Helper function to create test files
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

func TestRunAnalyzeWithComplexProject(t *testing.T) {
	// Create a more complex test project
	files := map[string]string{
		"main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
		"utils/helper.go": `package utils

// TODO: implement this
func Helper() {
	password := "secret123"  // This should trigger security rule
	fmt.Println(password)
}`,
		"pkg/complex.go": `package pkg

func VeryComplexFunction(x int) int {
	if x > 0 {
		if x > 10 {
			if x > 20 {
				if x > 30 {
					if x > 40 {
						return x * 2
					}
				}
			}
		}
	}
	return x
}`,
	}

	tmpDir := createTestProject(t, files)
	outputFile := filepath.Join(tmpDir, "report.json")

	cmd := &cobra.Command{
		Use:  "analyze",
		RunE: runAnalyze,
	}

	// Add flags
	cmd.Flags().String("output", "", "Output file")
	cmd.Flags().String("coverage", "", "Coverage file")
	cmd.Flags().String("diff", "", "Diff target")
	cmd.Flags().String("config", "", "Config file")
	cmd.Flags().Bool("verbose", false, "Verbose output")
	cmd.Flags().Bool("include-tests", false, "Include tests")
	cmd.Flags().StringSlice("exclude", []string{}, "Exclude patterns")

	// Set output flag
	err := cmd.Flags().Set("output", outputFile)
	if err != nil {
		t.Fatalf("Failed to set output flag: %v", err)
	}

	// Run analysis
	cmd.SetArgs([]string{tmpDir})
	err = cmd.Execute()
	if err != nil {
		t.Fatalf("runAnalyze() failed: %v", err)
	}

	// Verify report was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected report file not created: %v", outputFile)
	}

	// Read and verify report content
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read report file: %v", err)
	}

	var report types.Report
	err = json.Unmarshal(data, &report)
	if err != nil {
		t.Fatalf("Report is not valid JSON: %v", err)
	}

	// Verify report contains expected data
	if report.FilesAnalyzed == 0 {
		t.Error("Report should show analyzed files")
	}

	if len(report.Issues) == 0 {
		t.Error("Report should contain issues for the test code")
	}

	if report.ProjectPath != tmpDir {
		t.Errorf("Report ProjectPath = %v, want %v", report.ProjectPath, tmpDir)
	}
}
