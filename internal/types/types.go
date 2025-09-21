package types

import "time"

// Issue represents a code quality issue
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

// IssueImpact represents the effort required to fix an issue
type IssueImpact struct {
	EffortMinutes int `json:"effort_minutes"`
}

// Severity levels
const (
	SeverityLow    = "low"
	SeverityMedium = "medium"
	SeverityHigh   = "high"
)

// Issue types
const (
	TypeBug           = "bug"
	TypeVulnerability = "vulnerability"
	TypeCodeSmell     = "code_smell"
	TypeError         = "error"
)

// QualityTrends represents quality trend analysis
type QualityTrends struct {
	CyclomaticComplexityTrend string  `json:"cyclomatic_complexity_trend"`
	IssuesDensity             float64 `json:"issues_density"`
	SecurityScore             float64 `json:"security_score"`
	MaintainabilityIndex      float64 `json:"maintainability_index"`
	TechnicalDebtRatio        float64 `json:"technical_debt_ratio"`
}

// Report represents the complete analysis report
type Report struct {
	ProjectPath     string         `json:"project_path"`
	Timestamp       time.Time      `json:"timestamp"`
	Grade           string         `json:"grade"`
	Score           float64        `json:"score"`
	Issues          []Issue        `json:"issues"`
	Metrics         Metrics        `json:"metrics"`
	FilesAnalyzed   int            `json:"files_analyzed"`
	IssuesSummary   map[string]int `json:"issues_summary"`
	Trends          QualityTrends  `json:"trends,omitempty"`
	Recommendations []string       `json:"recommendations,omitempty"`
}

// Metrics represents code quality metrics
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

// AnalysisOptions represents options for code analysis
type AnalysisOptions struct {
	ProjectPath  string
	OutputFile   string
	CoverageFile string
	DiffTarget   string
	ConfigFile   string
	ExcludeFiles []string
	IncludeTests bool
	Verbose      bool
}

// QualityGate represents a quality gate check
type QualityGate struct {
	Name      string `json:"name"`
	Condition string `json:"condition"`
	Threshold string `json:"threshold"`
	Passed    bool   `json:"passed"`
	Message   string `json:"message"`
}
