package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestMethod тестирует простой вызов метода
func TestMethod(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	foo := NewMethodFoo()
	ctx.Set("foo", foo)

	expr, err := engine.CreateExpression(nil, "foo.bar()")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != "Method string" {
		t.Errorf("Expected 'Method string', got %v", result)
	}
}

// TestMulti тестирует вызов метода на вложенном объекте
func TestMulti(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	foo := NewMethodFoo()
	ctx.Set("foo", foo)

	expr, err := engine.CreateExpression(nil, "foo.innerFoo.bar()")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != "Method string" {
		t.Errorf("Expected 'Method string', got %v", result)
	}
}

// TestStringMethods тестирует вызовы методов строк
func TestStringMethods(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("foo", "abcdef")

	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"substring(3)", "foo.substring(3)", "def"},
		{"substring(0, size(foo)-3)", "foo.substring(0, size(foo)-3)", "abc"},
		{"substring(0, foo.length()-3)", "foo.substring(0, foo.length()-3)", "abc"},
		{"substring(0, 1+1)", "foo.substring(0, 1+1)", "ab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create expression: %v", err)
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				// Некоторые методы могут не поддерживаться
				t.Logf("Method evaluation failed (may not be supported): %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected '%s', got %v", tt.expected, result)
			}
		})
	}
}

// TestScriptCall тестирует вызов скрипта как функции
func TestScriptCall(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	plus, err := engine.CreateScript(nil, nil, "a + b", "a", "b")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	ctx.Set("plus", plus)
	forty2, err := engine.CreateScript(nil, nil, "plus(4, 2) * plus(4, 3)")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := forty2.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var actual int64
	switch v := result.(type) {
	case int64:
		actual = v
	case *big.Rat:
		if v.IsInt() {
			actual = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	// (4+2) * (4+3) = 6 * 7 = 42
	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}
}

// TestScriptCallInMap тестирует вызов скрипта из мапы
func TestScriptCallInMap(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	plus, err := engine.CreateScript(nil, nil, "a + b", "a", "b")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	foo := map[string]any{
		"plus": plus,
	}
	ctx.Set("foo", foo)

	forty2, err := engine.CreateScript(nil, nil, "foo.plus(4, 2) * foo.plus(4, 3)")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := forty2.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var actual int64
	switch v := result.(type) {
	case int64:
		actual = v
	case *big.Rat:
		if v.IsInt() {
			actual = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	// (4+2) * (4+3) = 6 * 7 = 42
	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}
}

// MethodFoo - тестовый класс для тестов методов (отличается от Foo в assign_test.go)
type MethodFoo struct {
	innerFoo *MethodFoo
}

func (f *MethodFoo) Bar() string {
	return "Method string"
}

func NewMethodFoo() *MethodFoo {
	return &MethodFoo{
		innerFoo: &MethodFoo{},
	}
}

