package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/chaksack/nada/internal/analyzer"
	"github.com/chaksack/nada/internal/reporter"
	"github.com/chaksack/nada/internal/types"
)

var (
	version   = "dev"
	buildTime = "unknown"
	commit    = "unknown"
)

// SetVersionInfo sets version information from build flags
func SetVersionInfo(v, bt, c string) {
	version = v
	buildTime = bt
	commit = c
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "nada",
	Short: "A comprehensive Go code quality analyzer",
	Long: `Nada is a static code analysis tool for Go projects that detects bugs,
vulnerabilities, code smells, and provides quality metrics similar to SonarQube.

Examples:
  nada analyze .                    # Analyze current directory
  nada analyze /path/to/project     # Analyze specific project
  nada analyze . --output report.json  # Export to JSON
  nada server --port 3000           # Start web API server`,
}

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze a Go codebase for quality issues",
	Args:  cobra.ExactArgs(1),
	RunE:  runAnalyze,
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run:   runVersion,
}

// runAnalyze executes the code analysis
func runAnalyze(cmd *cobra.Command, args []string) error {
	projectPath := args[0]

	// Validate project path
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", projectPath)
	}

	// Get command flags
	outputFile, _ := cmd.Flags().GetString("output")
	coverageFile, _ := cmd.Flags().GetString("coverage")
	diffTarget, _ := cmd.Flags().GetString("diff")
	configFile, _ := cmd.Flags().GetString("config")
	verbose, _ := cmd.Flags().GetBool("verbose")
	includeTests, _ := cmd.Flags().GetBool("include-tests")
	excludeFiles, _ := cmd.Flags().GetStringSlice("exclude")

	// Create analysis options
	options := types.AnalysisOptions{
		ProjectPath:  projectPath,
		OutputFile:   outputFile,
		CoverageFile: coverageFile,
		DiffTarget:   diffTarget,
		ConfigFile:   configFile,
		ExcludeFiles: excludeFiles,
		IncludeTests: includeTests,
		Verbose:      verbose,
	}

	// Create analyzer and run analysis
	codeAnalyzer := analyzer.New(options)
	report, err := codeAnalyzer.AnalyzeProject()
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Print report to console
	reporter.PrintConsoleReport(report)

	// Save report if requested
	if outputFile != "" {
		if err := saveReport(report, outputFile); err != nil {
			return fmt.Errorf("failed to save report: %w", err)
		}
		fmt.Printf("💾 Report saved to: %s\n", outputFile)
	}

	// Exit with error code based on quality gates
	if shouldFailBuild(report) {
		os.Exit(1)
	}

	return nil
}

// runVersion shows version information
func runVersion(cmd *cobra.Command, args []string) {
	fmt.Fprintln(cmd.OutOrStdout(), "Nada", version)
	fmt.Fprintln(cmd.OutOrStdout(), "Build time:", buildTime)
	fmt.Fprintln(cmd.OutOrStdout(), "Commit:", commit)
	fmt.Fprintln(cmd.OutOrStdout(), "Author: Andrew Chakdahah (@chaksack)")
}

// saveReport saves the analysis report to a JSON file
func saveReport(report *types.Report, filename string) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal report to JSON
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// shouldFailBuild determines if the build should fail based on quality gates
func shouldFailBuild(report *types.Report) bool {
	// Fail if there are high severity issues
	if report.IssuesSummary[types.SeverityHigh] > 0 {
		return true
	}

	// Fail if there are vulnerabilities
	if report.IssuesSummary[types.TypeVulnerability] > 0 {
		return true
	}

	// Fail if grade is F
	if report.Grade == "F" {
		return true
	}

	return false
}

// Execute runs the CLI
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add commands
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(versionCmd)

	// Analyze command flags
	analyzeCmd.Flags().StringP("output", "o", "", "Output file for JSON report")
	analyzeCmd.Flags().StringP("coverage", "c", "", "Coverage profile file")
	analyzeCmd.Flags().StringP("diff", "d", "", "Analyze only changes (staged, unstaged, HEAD, branch)")
	analyzeCmd.Flags().StringP("config", "", "", "Configuration file path")
	analyzeCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	analyzeCmd.Flags().BoolP("include-tests", "t", false, "Include test files in analysis")
	analyzeCmd.Flags().StringSliceP("exclude", "e", []string{}, "Exclude file patterns")

	// Global flags
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress output")
}
