package types

import (
	"testing"
	"time"
)

func TestIssue(t *testing.T) {
	tests := []struct {
		name  string
		issue Issue
		want  string
	}{
		{
			name: "valid issue",
			issue: Issue{
				Type:        TypeBug,
				Severity:    SeverityHigh,
				File:        "test.go",
				Line:        10,
				Column:      5,
				Message:     "Test issue",
				Rule:        "test_rule",
				Description: "Test description",
				Impact:      IssueImpact{EffortMinutes: 15},
			},
			want: TypeBug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.issue.Type != tt.want {
				t.Errorf("Issue.Type = %v, want %v", tt.issue.Type, tt.want)
			}
		})
	}
}

func TestReport(t *testing.T) {
	tests := []struct {
		name   string
		report Report
	}{
		{
			name: "valid report",
			report: Report{
				ProjectPath:   "/test/project",
				Timestamp:     time.Now(),
				Grade:         "A",
				Score:         95.5,
				Issues:        []Issue{},
				FilesAnalyzed: 10,
				IssuesSummary: map[string]int{
					SeverityHigh:   0,
					SeverityMedium: 2,
					SeverityLow:    5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.report.Grade != "A" {
				t.Errorf("Report.Grade = %v, want A", tt.report.Grade)
			}
			if tt.report.Score != 95.5 {
				t.Errorf("Report.Score = %v, want 95.5", tt.report.Score)
			}
			if tt.report.FilesAnalyzed != 10 {
				t.Errorf("Report.FilesAnalyzed = %v, want 10", tt.report.FilesAnalyzed)
			}
		})
	}
}

func TestAnalysisOptions(t *testing.T) {
	tests := []struct {
		name    string
		options AnalysisOptions
		wantErr bool
	}{
		{
			name: "valid options",
			options: AnalysisOptions{
				ProjectPath:  "/test/project",
				OutputFile:   "report.json",
				IncludeTests: true,
				Verbose:      false,
			},
			wantErr: false,
		},
		{
			name: "minimal options",
			options: AnalysisOptions{
				ProjectPath: "/test/project",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.options.ProjectPath == "" && !tt.wantErr {
				t.Errorf("AnalysisOptions should have ProjectPath")
			}
		})
	}
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"SeverityLow", SeverityLow, "low"},
		{"SeverityMedium", SeverityMedium, "medium"},
		{"SeverityHigh", SeverityHigh, "high"},
		{"TypeBug", TypeBug, "bug"},
		{"TypeVulnerability", TypeVulnerability, "vulnerability"},
		{"TypeCodeSmell", TypeCodeSmell, "code_smell"},
		{"TypeError", TypeError, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Constant %s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}
