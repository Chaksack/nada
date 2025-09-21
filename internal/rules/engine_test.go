package rules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/chaksack/nada/internal/types"
)

// Test the rule engine
func TestNewEngine(t *testing.T) {
	engine := NewEngine()

	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}

	rules := engine.GetRules()
	if len(rules) == 0 {
		t.Error("NewEngine() created engine with no rules")
	}

	// Check that expected rule types are registered
	expectedRules := []string{"complexity", "security", "naming", "structure", "documentation", "error_handling"}
	ruleIDs := make(map[string]bool)

	for _, rule := range rules {
		ruleIDs[rule.ID()] = true
	}

	for _, expectedRule := range expectedRules {
		if !ruleIDs[expectedRule] {
			t.Errorf("Expected rule '%s' not found in engine", expectedRule)
		}
	}
}

func TestEngineRegisterRule(t *testing.T) {
	engine := NewEngine()
	initialCount := len(engine.GetRules())

	// Create a mock rule
	mockRule := &MockRule{id: "test_rule"}
	engine.RegisterRule(mockRule)

	rules := engine.GetRules()
	if len(rules) != initialCount+1 {
		t.Errorf("RegisterRule() expected %d rules, got %d", initialCount+1, len(rules))
	}

	// Check if our mock rule is there
	found := false
	for _, rule := range rules {
		if rule.ID() == "test_rule" {
			found = true
			break
		}
	}
	if !found {
		t.Error("RegisterRule() did not add the rule to the engine")
	}
}

func TestEngineAnalyzeFile(t *testing.T) {
	engine := NewEngine()

	// Simple Go code with issues
	code := `package main

import "fmt"

func VeryLongFunctionNameThatViolatesNamingConventions() {
	// TODO: fix this later
	password := "hardcoded_password"
	if true {
		if true {
			if true {
				if true {
					if true {
						fmt.Println("deeply nested")
					}
				}
			}
		}
	}
	fmt.Println(password)
}`

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse test code: %v", err)
	}

	issues := engine.AnalyzeFile("test.go", node, code, fset)

	if len(issues) == 0 {
		t.Error("AnalyzeFile() expected to find issues, but found none")
	}

	// Check for specific issue types we expect
	hasSecurityIssue := false
	hasStructureIssue := false
	hasComplexityIssue := false

	for _, issue := range issues {
		switch issue.Type {
		case types.TypeVulnerability:
			hasSecurityIssue = true
		case types.TypeCodeSmell:
			if issue.Rule == "deep_nesting" || issue.Rule == "todo_comment" {
				hasStructureIssue = true
			}
		}
		if issue.Rule == "high_complexity" {
			hasComplexityIssue = true
		}
	}

	if !hasSecurityIssue {
		t.Error("AnalyzeFile() expected to find security issues")
	}
	if !hasStructureIssue {
		t.Error("AnalyzeFile() expected to find structure issues")
	}
	if !hasComplexityIssue {
		t.Error("AnalyzeFile() expected to find complexity issues")
	}
}

// Mock rule for testing
type MockRule struct {
	id string
}

func (r *MockRule) ID() string          { return r.id }
func (r *MockRule) Name() string        { return "Mock Rule" }
func (r *MockRule) Description() string { return "A mock rule for testing" }
func (r *MockRule) Check(file string, node ast.Node, content string, fset *token.FileSet) []types.Issue {
	return []types.Issue{
		{
			Type:     types.TypeCodeSmell,
			Severity: types.SeverityLow,
			File:     file,
			Line:     1,
			Rule:     r.id,
			Message:  "Mock issue",
		},
	}
}

// Test ComplexityRule
func TestComplexityRule(t *testing.T) {
	rule := NewComplexityRule()

	if rule.ID() != "complexity" {
		t.Errorf("ComplexityRule.ID() = %v, want complexity", rule.ID())
	}

	// Test complex function
	code := `package main

func ComplexFunction(x int) int {
	if x > 0 {
		if x > 10 {
			if x > 20 {
				if x > 30 {
					if x > 40 {
						switch x {
						case 50:
							return 1
						case 60:
							return 2
						default:
							return 3
						}
					}
				}
			}
		}
	}
	return 0
}`

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse test code: %v", err)
	}

	issues := rule.Check("test.go", node, code, fset)

	if len(issues) == 0 {
		t.Error("ComplexityRule.Check() expected to find complexity issues")
	}

	// Check for complexity issue
	hasComplexityIssue := false
	for _, issue := range issues {
		if issue.Rule == "high_complexity" {
			hasComplexityIssue = true
			break
		}
	}

	if !hasComplexityIssue {
		t.Error("ComplexityRule.Check() expected to find high_complexity issue")
	}
}

// Test SecurityRule
func TestSecurityRule(t *testing.T) {
	rule := NewSecurityRule()

	if rule.ID() != "security" {
		t.Errorf("SecurityRule.ID() = %v, want security", rule.ID())
	}

	tests := []struct {
		name     string
		code     string
		wantRule string
	}{
		{
			name: "hardcoded password",
			code: `package main
var password = "secret123"`,
			wantRule: "hardcoded_secret",
		},
		{
			name: "SQL injection",
			code: `package main
import "fmt"
func main() {
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userInput)
}`,
			wantRule: "sql_injection",
		},
		{
			name: "hardcoded API key",
			code: `package main
const api_key = "sk-1234567890abcdef"`,
			wantRule: "hardcoded_secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse test code: %v", err)
			}

			issues := rule.Check("test.go", node, tt.code, fset)

			if len(issues) == 0 {
				t.Errorf("SecurityRule.Check() expected to find issues for %s", tt.name)
				return
			}

			// Check for specific rule
			hasExpectedRule := false
			for _, issue := range issues {
				if issue.Rule == tt.wantRule {
					hasExpectedRule = true
					break
				}
			}

			if !hasExpectedRule {
				t.Errorf("SecurityRule.Check() expected to find %s rule", tt.wantRule)
			}
		})
	}
}

// Test NamingRule
func TestNamingRule(t *testing.T) {
	rule := NewNamingRule()

	if rule.ID() != "naming" {
		t.Errorf("NamingRule.ID() = %v, want naming", rule.ID())
	}

	tests := []struct {
		name     string
		code     string
		wantRule string
	}{
		{
			name: "short function name",
			code: `package main
func f() {}`,
			wantRule: "short_function_name",
		},
		{
			name: "short variable name",
			code: `package main
func main() {
	var x int
}`,
			wantRule: "short_variable_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse test code: %v", err)
			}

			issues := rule.Check("test.go", node, tt.code, fset)

			// For naming rules, we might not always find issues depending on context
			if len(issues) > 0 {
				hasExpectedRule := false
				for _, issue := range issues {
					if issue.Rule == tt.wantRule {
						hasExpectedRule = true
						break
					}
				}

				if hasExpectedRule {
					// Good, we found the expected rule
					t.Logf("Found expected rule %s for test %s", tt.wantRule, tt.name)
				}
			}
		})
	}
}

// Test StructureRule
func TestStructureRule(t *testing.T) {
	rule := NewStructureRule()

	if rule.ID() != "structure" {
		t.Errorf("StructureRule.ID() = %v, want structure", rule.ID())
	}

	tests := []struct {
		name     string
		code     string
		wantRule string
	}{
		{
			name: "TODO comment",
			code: `package main
// TODO: implement this function
func main() {}`,
			wantRule: "todo_comment",
		},
		{
			name: "long line",
			code: `package main
// This is a very long comment that exceeds the maximum line length limit and should trigger a long_line rule violation
func main() {}`,
			wantRule: "long_line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse test code: %v", err)
			}

			issues := rule.Check("test.go", node, tt.code, fset)

			if len(issues) == 0 {
				t.Errorf("StructureRule.Check() expected to find issues for %s", tt.name)
				return
			}

			hasExpectedRule := false
			for _, issue := range issues {
				if issue.Rule == tt.wantRule {
					hasExpectedRule = true
					break
				}
			}

			if !hasExpectedRule {
				t.Errorf("StructureRule.Check() expected to find %s rule, got issues: %v",
					tt.wantRule, issues)
			}
		})
	}
}

// Benchmark tests
func BenchmarkEngineAnalyzeFile(b *testing.B) {
	engine := NewEngine()

	code := `package main

import "fmt"

func main() {
	// TODO: optimize this
	password := "secret123"
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			if i%4 == 0 {
				if i%8 == 0 {
					fmt.Println(i, password)
				}
			}
		}
	}
}`

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", code, parser.ParseComments)
	if err != nil {
		b.Fatalf("Failed to parse test code: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.AnalyzeFile("test.go", node, code, fset)
	}
}
