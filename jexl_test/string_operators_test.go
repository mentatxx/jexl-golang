package jexl_test

import (
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestStringContainsOperator тестирует оператор =~ (contains)
func TestStringContainsOperator(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("str", "abc456")
	ctx.Set("items", []any{int64(1), int64(2), int64(3), int64(4), int64(5)})

	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{"string contains substring", "str =~ '456'", true},
		{"string contains regex", "str =~ '.*456'", true},
		{"string not contains", "str =~ 'ABC'", false},
		{"array contains element", "2 =~ items", true},
		{"array not contains", "10 =~ items", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create expression: %v", err)
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			b, ok := result.(bool)
			if !ok {
				t.Fatalf("Expected bool, got %T", result)
			}

			if b != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, b)
			}
		})
	}
}

// TestStringStartsWithOperator тестирует оператор =^ (startsWith)
func TestStringStartsWithOperator(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("str", "abcdef")

	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{"starts with prefix", "str =^ 'abc'", true},
		{"not starts with", "str =^ 'def'", false},
		{"empty prefix", "str =^ ''", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create expression: %v", err)
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			b, ok := result.(bool)
			if !ok {
				t.Fatalf("Expected bool, got %T", result)
			}

			if b != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, b)
			}
		})
	}
}

// TestStringEndsWithOperator тестирует оператор =$ (endsWith)
func TestStringEndsWithOperator(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("str", "abcdef")

	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{"ends with suffix", "str =$ 'def'", true},
		{"not ends with", "str =$ 'abc'", false},
		{"empty suffix", "str =$ ''", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create expression: %v", err)
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			b, ok := result.(bool)
			if !ok {
				t.Fatalf("Expected bool, got %T", result)
			}

			if b != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, b)
			}
		})
	}
}

// TestStringNotOperators тестирует отрицательные операторы !~, !^, !$
func TestStringNotOperators(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("str", "abcdef")

	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{"not contains", "str !~ 'xyz'", true},
		{"not starts with", "str !^ 'def'", true},
		{"not ends with", "str !$ 'abc'", true},
		{"contains (negated false)", "str !~ 'abc'", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create expression: %v", err)
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			b, ok := result.(bool)
			if !ok {
				t.Fatalf("Expected bool, got %T", result)
			}

			if b != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, b)
			}
		})
	}
}

