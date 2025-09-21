package rules

import (
	"go/ast"
	"go/token"

	"github.com/chaksack/nada/internal/types"
)

// Rule defines the interface for analysis rules
type Rule interface {
	ID() string
	Name() string
	Description() string
	Check(file string, node ast.Node, content string, fset *token.FileSet) []types.Issue
}

// Engine manages and executes analysis rules
type Engine struct {
	rules []Rule
}

// NewEngine creates a new rule engine with default rules
func NewEngine() *Engine {
	engine := &Engine{
		rules: make([]Rule, 0),
	}

	// Register default rules
	engine.RegisterRule(NewComplexityRule())
	engine.RegisterRule(NewSecurityRule())
	engine.RegisterRule(NewNamingRule())
	engine.RegisterRule(NewStructureRule())
	engine.RegisterRule(NewDocumentationRule())
	engine.RegisterRule(NewErrorHandlingRule())

	return engine
}

// RegisterRule adds a new rule to the engine
func (e *Engine) RegisterRule(rule Rule) {
	e.rules = append(e.rules, rule)
}

// GetRules returns all registered rules
func (e *Engine) GetRules() []Rule {
	return e.rules
}

// AnalyzeFile runs all rules against a file
func (e *Engine) AnalyzeFile(filePath string, node *ast.File, content string, fset *token.FileSet) []types.Issue {
	var allIssues []types.Issue

	for _, rule := range e.rules {
		issues := rule.Check(filePath, node, content, fset)
		allIssues = append(allIssues, issues...)
	}

	return allIssues
}
