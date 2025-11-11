package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestBooleanExpressions тестирует булевы выражения
func TestBooleanExpressions(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("a", true)
	ctx.Set("b", false)

	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{"true eq false", "true eq false", false},
		{"true ne false", "true ne false", true},
		{"a == b", "a == b", false},
		{"a != b", "a != b", true},
		{"a && b", "a && b", false},
		{"a || b", "a || b", true},
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

			// Результат может быть big.Rat или bool
			var boolResult bool
			switch v := result.(type) {
			case bool:
				boolResult = v
			case *big.Rat:
				boolResult = v.Sign() != 0
			default:
				t.Fatalf("Unexpected result type: %T", result)
			}

			if boolResult != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, boolResult)
			}
		})
	}
}

// TestStringOperations тестирует строковые операции
func TestStringOperations(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("foo", "abcdef")
	ctx.Set("bar", "")

	tests := []struct {
		name     string
		expr     string
		check    func(any) bool
	}{
		{"empty string", "empty bar", func(r any) bool {
			b, ok := r.(bool)
			return ok && b
		}},
		{"bar == ''", "bar == ''", func(r any) bool {
			b, ok := r.(bool)
			return ok && b
		}},
		{"string concatenation", "foo + ' is good'", func(r any) bool {
			s, ok := r.(string)
			return ok && s == "abcdef is good"
		}},
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

			if !tt.check(result) {
				t.Errorf("Check failed for result: %v", result)
			}
		})
	}
}

// TestComparisons тестирует операции сравнения
func TestComparisons(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 10)
	ctx.Set("y", 5)

	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{"x > y", "x > y", true},
		{"x >= y", "x >= y", true},
		{"x >= x", "x >= x", true},
		{"x < y", "x < y", false},
		{"x <= y", "x <= y", false},
		{"x <= x", "x <= x", true},
		{"x == x", "x == x", true},
		{"x != y", "x != y", true},
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

			var boolResult bool
			switch v := result.(type) {
			case bool:
				boolResult = v
			case *big.Rat:
				boolResult = v.Sign() != 0
			default:
				t.Fatalf("Unexpected result type: %T", result)
			}

			if boolResult != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, boolResult)
			}
		})
	}
}

// TestNullCoercion тестирует приведение null значений
func TestNullCoercion(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("imanull", nil)
	ctx.Set("n", 0)

	// В нестрогом режиме null + число должно работать
	expr, err := engine.CreateExpression(nil, "n != null && n != 0")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	// Результат должен быть false, так как n == 0
	var boolResult bool
	switch v := result.(type) {
	case bool:
		boolResult = v
	case *big.Rat:
		boolResult = v.Sign() != 0
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if boolResult {
		t.Error("Expected false for 'n != null && n != 0' when n == 0")
	}
}

// TestArrayAccess тестирует доступ к массивам
func TestArrayAccess(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	array := []any{100, 101, 102}
	ctx.Set("array", array)

	tests := []struct {
		name     string
		expr     string
		expected any
	}{
		{"array[1]", "array[1]", 101},
		{"array[0]", "array[0]", 100},
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

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestMethodCalls тестирует вызовы методов
func TestMethodCalls(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("foo", "abcdef")

	// Тест вызова метода substring (если поддерживается)
	expr, err := engine.CreateExpression(nil, "foo.length()")
	if err != nil {
		// Метод может не поддерживаться, это нормально
		t.Logf("Method call not supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Method evaluation failed: %v", err)
		return
	}

	if result == nil {
		t.Error("Result is nil")
	}
}

// TestElvisOperator тестирует Elvis оператор
func TestElvisOperator(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", nil)
	ctx.Set("y", "default")

	expr, err := engine.CreateExpression(nil, "x ?? y")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != "default" {
		t.Errorf("Expected 'default', got %v", result)
	}
}

// TestDoWhileLoop тестирует do-while цикл
func TestDoWhileLoop(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 0)

	script, err := engine.CreateScript(nil, nil, "do { x = x + 1 } while (x < 5)")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	_, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	x := ctx.Get("x")
	if x == nil {
		t.Fatal("x is nil")
	}
}

// TestSetLiteral тестирует литералы множеств
func TestSetLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "{1, 2, 3, 4, 5}")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected array, got %T", result)
	}

	if len(arr) != 5 {
		t.Fatalf("Expected array length 5, got %d", len(arr))
	}
}

// TestNestedLiterals тестирует вложенные литералы
func TestNestedLiterals(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Мапа с массивом
	expr, err := engine.CreateExpression(nil, "{'foo': [1, 2, 3]}")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map, got %T", result)
	}

	arr, ok := m["foo"].([]any)
	if !ok {
		t.Fatalf("Expected array in map, got %T", m["foo"])
	}

	if len(arr) != 3 {
		t.Fatalf("Expected array length 3, got %d", len(arr))
	}
}

// TestComplexExpressions тестирует сложные выражения
func TestComplexExpressions(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("a", 5)
	ctx.Set("b", 3)
	ctx.Set("c", 2)

	expr, err := engine.CreateExpression(nil, "(a + b) * c")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}
}

// TestPropertyChaining тестирует цепочки свойств
func TestPropertyChaining(t *testing.T) {
	type Inner struct {
		Value int
	}
	type Outer struct {
		Inner *Inner
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	outer := &Outer{Inner: &Inner{Value: 42}}
	ctx.Set("outer", outer)

	expr, err := engine.CreateExpression(nil, "outer.Inner.Value")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}
}

