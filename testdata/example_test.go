package testdata

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/chaksack/nada/internal/rules"
	"github.com/chaksack/nada/internal/types"
)

// ExampleTestData provides sample Go code for testing rules
var ExampleTestData = map[string]string{
	"simple.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,

	"complex.go": `package main

import "fmt"

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
					for i := 0; i < 10; i++ {
						if i%2 == 0 {
							if i%4 == 0 {
								fmt.Println(i)
							}
						}
					}
					return 0
				}
			}
		}
	}
	return -1
}`,

	"security_issues.go": `package main

import (
	"fmt"
	"database/sql"
)

func SecurityIssues() {
	// Hardcoded secrets
	password := "admin123"
	apiKey := "sk-1234567890abcdef"
	token := "bearer_token_here"
	awsKey := "AKIAIOSFODNN7EXAMPLE"
	
	// SQL injection vulnerability
	userInput := "1; DROP TABLE users;"
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userInput)
	
	db, _ := sql.Open("mysql", "connection_string")
	rows, _ := db.Query(query)
	defer rows.Close()
	
	fmt.Println(password, apiKey, token, awsKey)
}`,

	"naming_issues.go": `package main

import "fmt"

// Short function names
func f() {}
func g(x int) int { return x }
func h() string { return "" }

// Bad naming conventions
func httpGetURL() {}  // Should be HTTPGetURL
func jsonData() {}    // Should be JSONData
func apiHandler() {}  // Should be APIHandler

func BadExample() {
	// Short variable names
	var x int
	var y string
	var z float64
	
	// Non-descriptive names
	var a, b, c int
	var temp string
	var data interface{}
	
	fmt.Println(x, y, z, a, b, c, temp, data)
}`,

	"structure_issues.go": `package main

import (
	"fmt"
	_ "unused_import_example"
)

func StructureIssues() {
	// TODO: implement proper error handling
	// FIXME: this is a temporary workaround
	// HACK: quick fix for production
	
	fmt.Println("This is an extremely long line that exceeds the typical 120-character limit and should trigger a line length warning in our static analysis tool")
	
	// Deeply nested code
	if true {
		if true {
			if true {
				if true {
					if true {
						if true {
							fmt.Println("Too deeply nested!")
						}
					}
				}
			}
		}
	}
}`,

	"documentation_issues.go": `package main

import "fmt"

// GoodFunction has proper documentation
func GoodFunction() {
	fmt.Println("This function is documented")
}

func BadFunction() {
	fmt.Println("This exported function lacks documentation")
}

func anotherBadFunction() {
	fmt.Println("This is also missing docs")
}

type BadType struct {
	Field1 string
	Field2 int
}

// GoodType has documentation
type GoodType struct {
	Name string
	Age  int
}`,

	"error_handling_issues.go": `package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"  // Deprecated package
	"os"
	"strconv"
)

func ErrorHandlingIssues() {
	// Missing error handling
	file, _ := os.Open("nonexistent.txt")
	defer file.Close()
	
	data, _ := ioutil.ReadAll(file)  // Deprecated function
	
	var result map[string]interface{}
	json.Unmarshal(data, &result)  // Missing error handling
	
	number, _ := strconv.Atoi("not_a_number")  // Missing error handling
	
	fmt.Println(result, number)
}`,

	"performance_issues.go": `package main

import "fmt"

func PerformanceIssues() {
	// Inefficient string concatenation
	result := ""
	for i := 0; i < 1000; i++ {
		result += fmt.Sprintf("item-%d ", i)
	}
	
	// Unnecessary slice allocations
	var items []string
	for i := 0; i < 100; i++ {
		items = append(items, fmt.Sprintf("item-%d", i))
	}
	
	fmt.Println(result, len(items))
}`,
}

// TestExampleAnalysis demonstrates how to use the test data
func TestExampleAnalysis(t *testing.T) {
	engine := rules.NewEngine()

	for filename, code := range ExampleTestData {
		t.Run(filename, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, filename, code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", filename, err)
			}

			issues := engine.AnalyzeFile(filename, node, code, fset)

			// Log issues found for manual inspection
			if len(issues) > 0 {
				t.Logf("Found %d issues in %s:", len(issues), filename)
				for _, issue := range issues {
					t.Logf("  - %s:%d [%s/%s] %s",
						issue.File, issue.Line, issue.Type, issue.Severity, issue.Message)
				}
			} else {
				t.Logf("No issues found in %s", filename)
			}
		})
	}
}

// TestSpecificRuleOnExamples tests specific rules against example code
func TestSecurityRuleOnExamples(t *testing.T) {
	rule := rules.NewSecurityRule()
	fset := token.NewFileSet()

	code := ExampleTestData["security_issues.go"]
	node, err := parser.ParseFile(fset, "security_issues.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse security test code: %v", err)
	}

	issues := rule.Check("security_issues.go", node, code, fset)

	if len(issues) == 0 {
		t.Error("SecurityRule should find issues in security_issues.go")
	}

	// Count different types of security issues
	hardcodedSecrets := 0
	sqlInjections := 0

	for _, issue := range issues {
		switch issue.Rule {
		case "hardcoded_secret":
			hardcodedSecrets++
		case "sql_injection":
			sqlInjections++
		}
	}

	if hardcodedSecrets == 0 {
		t.Error("SecurityRule should find hardcoded secrets")
	}

	if sqlInjections == 0 {
		t.Error("SecurityRule should find SQL injection issues")
	}

	t.Logf("Found %d hardcoded secrets and %d SQL injection issues",
		hardcodedSecrets, sqlInjections)
}

func TestComplexityRuleOnExamples(t *testing.T) {
	rule := rules.NewComplexityRule()
	fset := token.NewFileSet()

	code := ExampleTestData["complex.go"]
	node, err := parser.ParseFile(fset, "complex.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse complexity test code: %v", err)
	}

	issues := rule.Check("complex.go", node, code, fset)

	if len(issues) == 0 {
		t.Error("ComplexityRule should find issues in complex.go")
	}

	// Should find high complexity issues
	hasComplexityIssue := false
	for _, issue := range issues {
		if issue.Rule == "high_complexity" {
			hasComplexityIssue = true
			break
		}
	}

	if !hasComplexityIssue {
		t.Error("ComplexityRule should find high complexity issues")
	}
}

// Helper function to create test report
func CreateTestReport() *types.Report {
	return &types.Report{
		ProjectPath:   "/test/project",
		Grade:         "B",
		Score:         82.5,
		FilesAnalyzed: 10,
		Issues: []types.Issue{
			{
				Type:     types.TypeVulnerability,
				Severity: types.SeverityHigh,
				File:     "security.go",
				Line:     15,
				Message:  "Hardcoded password detected",
				Rule:     "hardcoded_secret",
				Impact:   types.IssueImpact{EffortMinutes: 10},
			},
			{
				Type:     types.TypeBug,
				Severity: types.SeverityMedium,
				File:     "main.go",
				Line:     25,
				Message:  "Missing error handling",
				Rule:     "error_handling",
				Impact:   types.IssueImpact{EffortMinutes: 5},
			},
			{
				Type:     types.TypeCodeSmell,
				Severity: types.SeverityLow,
				File:     "utils.go",
				Line:     8,
				Message:  "TODO comment found",
				Rule:     "todo_comment",
				Impact:   types.IssueImpact{EffortMinutes: 2},
			},
		},
		Metrics: types.Metrics{
			LinesOfCode:          1500,
			CyclomaticComplexity: 85,
			TestCoverage:         75.5,
		},
		IssuesSummary: map[string]int{
			types.SeverityHigh:      1,
			types.SeverityMedium:    1,
			types.SeverityLow:       1,
			types.TypeVulnerability: 1,
			types.TypeBug:           1,
			types.TypeCodeSmell:     1,
		},
		Trends: types.QualityTrends{
			CyclomaticComplexityTrend: "Moderate - Monitor closely",
			IssuesDensity:             2.0,
			SecurityScore:             90.0,
			MaintainabilityIndex:      82.0,
			TechnicalDebtRatio:        1.1,
		},
		Recommendations: []string{
			"🔒 URGENT: Address 1 security vulnerabilities",
			"🐛 HIGH: Fix 1 bugs to improve reliability",
			"💡 Consider increasing test coverage to 80%+",
		},
	}
}

// Benchmarks using example data
func BenchmarkSecurityRuleAnalysis(b *testing.B) {
	rule := rules.NewSecurityRule()
	fset := token.NewFileSet()

	code := ExampleTestData["security_issues.go"]
	node, err := parser.ParseFile(fset, "security_issues.go", code, parser.ParseComments)
	if err != nil {
		b.Fatalf("Failed to parse test code: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rule.Check("security_issues.go", node, code, fset)
	}
}

func BenchmarkComplexityRuleAnalysis(b *testing.B) {
	rule := rules.NewComplexityRule()
	fset := token.NewFileSet()

	code := ExampleTestData["complex.go"]
	node, err := parser.ParseFile(fset, "complex.go", code, parser.ParseComments)
	if err != nil {
		b.Fatalf("Failed to parse test code: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rule.Check("complex.go", node, code, fset)
	}
}
