package jexl_test

import (
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestEmptyMapLiteral тестирует пустой литерал мапы
func TestEmptyMapLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "{}")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	if len(m) != 0 {
		t.Errorf("Expected empty map, got size %d", len(m))
	}
}

// TestMapLiteralWithStrings тестирует мапу со строками
func TestMapLiteralWithStrings(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "{'foo': 'bar'}")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	if len(m) != 1 {
		t.Fatalf("Expected map size 1, got %d", len(m))
	}

	if m["foo"] != "bar" {
		t.Errorf("Expected m['foo'] to be 'bar', got %v", m["foo"])
	}
}

// TestMapLiteralWithMultipleEntries тестирует мапу с несколькими записями
func TestMapLiteralWithMultipleEntries(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "{'foo': 'bar', 'eat': 'food'}")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	if len(m) != 2 {
		t.Fatalf("Expected map size 2, got %d", len(m))
	}

	if m["foo"] != "bar" {
		t.Errorf("Expected m['foo'] to be 'bar', got %v", m["foo"])
	}

	if m["eat"] != "food" {
		t.Errorf("Expected m['eat'] to be 'food', got %v", m["eat"])
	}
}

// TestMapLiteralWithNumbers тестирует мапу с числами
func TestMapLiteralWithNumbers(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "{5: 10}")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	// Ключи преобразуются в строки
	if m["5"] != 10 {
		t.Errorf("Expected m['5'] to be 10, got %v", m["5"])
	}
}

// TestMapArrayLiteral тестирует мапу с массивом
func TestMapArrayLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "{'foo': [1, 2, 3]}")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	arr, ok := m["foo"].([]any)
	if !ok {
		t.Fatalf("Expected []any in map, got %T", m["foo"])
	}

	if len(arr) != 3 {
		t.Fatalf("Expected array length 3, got %d", len(arr))
	}
}

// TestMapMapLiteral тестирует мапу с мапой
func TestMapMapLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "{'foo': {'inner': 'bar'}}")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	inner, ok := m["foo"].(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any in map, got %T", m["foo"])
	}

	if inner["inner"] != "bar" {
		t.Errorf("Expected inner['inner'] to be 'bar', got %v", inner["inner"])
	}
}

// TestMapAccess тестирует доступ к мапе
func TestMapAccess(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("m", map[string]any{"foo": "bar"})

	expr, err := engine.CreateExpression(nil, "m['foo']")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != "bar" {
		t.Errorf("Expected 'bar', got %v", result)
	}
}

