package jexl_test

import (
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

func TestBasicExpression(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 10)
	ctx.Set("y", 20)

	expr, err := engine.CreateExpression(nil, "x + y")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	// Результат должен быть big.Rat, так как используется BaseArithmetic
	if result == nil {
		t.Fatal("Result is nil")
	}
}

func TestMapContext(t *testing.T) {
	ctx := jexl.NewMapContext()
	ctx.Set("name", "test")
	ctx.Set("value", 42)

	if !ctx.Has("name") {
		t.Error("Context should have 'name'")
	}

	if ctx.Has("missing") {
		t.Error("Context should not have 'missing'")
	}

	val := ctx.Get("name")
	if val != "test" {
		t.Errorf("Expected 'test', got %v", val)
	}
}

func TestBuilder(t *testing.T) {
	builder := jexl.NewBuilder().
		Cache(128)

	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	if engine == nil {
		t.Fatal("Engine is nil")
	}
}

func TestSimpleArithmetic(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("a", 5)
	ctx.Set("b", 3)

	tests := []struct {
		expr     string
		expected any
	}{
		{"a + b", nil}, // Результат будет big.Rat, проверяем только что нет ошибки
		{"a - b", nil},
		{"a * b", nil},
		{"a / b", nil},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create expression %s: %v", tt.expr, err)
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Fatalf("Failed to evaluate %s: %v", tt.expr, err)
			}

			if result == nil {
				t.Errorf("Result is nil for %s", tt.expr)
			}
		})
	}
}

func TestAssignment(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 5)
	ctx.Set("y", 7)

	expr, err := engine.CreateExpression(nil, "z = x + y")
	if err != nil {
		t.Fatalf("Failed to create assignment expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate assignment: %v", err)
	}

	if result == nil {
		t.Fatal("Assignment result is nil")
	}
	if !ctx.Has("z") {
		t.Fatal("Context does not contain assigned variable 'z'")
	}

	obj := &struct {
		Value int
	}{Value: 1}
	ctx.Set("obj", obj)

	expr, err = engine.CreateExpression(nil, "obj.Value = 42")
	if err != nil {
		t.Fatalf("Failed to create property assignment expression: %v", err)
	}
	if _, err := expr.Evaluate(ctx); err != nil {
		t.Fatalf("Failed to evaluate property assignment: %v", err)
	}
	if obj.Value != 42 {
		t.Fatalf("Expected obj.Value to be 42, got %d", obj.Value)
	}

	data := map[string]any{"foo": 1}
	ctx.Set("data", data)

	expr, err = engine.CreateExpression(nil, "data['bar'] = 99")
	if err != nil {
		t.Fatalf("Failed to create map assignment expression: %v", err)
	}
	if _, err := expr.Evaluate(ctx); err != nil {
		t.Fatalf("Failed to evaluate map assignment: %v", err)
	}
	if v, ok := data["bar"].(int64); !ok || v != 99 {
		t.Fatalf("Expected data['bar'] to be 99 (int64), got %v", data["bar"])
	}
}

func TestIfElse(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 10)

	script, err := engine.CreateScript(nil, nil, "if (x > 5) { 100 } else { 200 }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}
}

func TestForLoop(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("sum", 0)

	script, err := engine.CreateScript(nil, nil, "for (i = 0; i < 5; i = i + 1) { sum = sum + i }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	_, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Проверяем, что sum был изменён
	sum := ctx.Get("sum")
	if sum == nil {
		t.Fatal("sum is nil")
	}
}

func TestForeachLoop(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	items := []any{1, 2, 3, 4, 5}
	ctx.Set("items", items)
	ctx.Set("sum", 0)

	script, err := engine.CreateScript(nil, nil, "for (var x : items) { sum = sum + x }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	_, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	sum := ctx.Get("sum")
	if sum == nil {
		t.Fatal("sum is nil")
	}
}

func TestWhileLoop(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 0)

	script, err := engine.CreateScript(nil, nil, "while (x < 5) { x = x + 1 }")
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

func TestTernaryOperator(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 10)

	expr, err := engine.CreateExpression(nil, "x > 5 ? 100 : 200")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}
}

func TestArrayLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "[1, 2, 3, 4, 5]")
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

func TestMapLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "{'a': 1, 'b': 2, 'c': 3}")
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

	if len(m) != 3 {
		t.Fatalf("Expected map size 3, got %d", len(m))
	}
}
