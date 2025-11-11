package jexl_test

import (
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestEmptyArrayLiteral тестирует пустой литерал массива
func TestEmptyArrayLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "[]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result)
	}

	if len(arr) != 0 {
		t.Errorf("Expected empty array, got length %d", len(arr))
	}
}

// TestArrayLiteralWithIntegers тестирует массив с целыми числами
func TestArrayLiteralWithIntegers(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "[5, 10]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result)
	}

	if len(arr) != 2 {
		t.Fatalf("Expected array length 2, got %d", len(arr))
	}

	// Проверяем значения (могут быть разных типов)
	val0 := arr[0]
	val1 := arr[1]
	
	// Преобразуем к int для сравнения
	var int0, int1 int
	switch v := val0.(type) {
	case int:
		int0 = v
	case int64:
		int0 = int(v)
	case float64:
		int0 = int(v)
	default:
		t.Errorf("Unexpected type for arr[0]: %T", val0)
		return
	}
	
	switch v := val1.(type) {
	case int:
		int1 = v
	case int64:
		int1 = int(v)
	case float64:
		int1 = int(v)
	default:
		t.Errorf("Unexpected type for arr[1]: %T", val1)
		return
	}
	
	if int0 != 5 {
		t.Errorf("Expected arr[0] to be 5, got %d", int0)
	}
	
	if int1 != 10 {
		t.Errorf("Expected arr[1] to be 10, got %d", int1)
	}
}

// TestArrayLiteralWithStrings тестирует массив со строками
func TestArrayLiteralWithStrings(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "['foo', 'bar']")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result)
	}

	if len(arr) != 2 {
		t.Fatalf("Expected array length 2, got %d", len(arr))
	}

	if arr[0] != "foo" {
		t.Errorf("Expected arr[0] to be 'foo', got %v", arr[0])
	}

	if arr[1] != "bar" {
		t.Errorf("Expected arr[1] to be 'bar', got %v", arr[1])
	}
}

// TestArrayLiteralWithOneEntry тестирует массив с одним элементом
func TestArrayLiteralWithOneEntry(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "['foo']")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result)
	}

	if len(arr) != 1 {
		t.Fatalf("Expected array length 1, got %d", len(arr))
	}

	if arr[0] != "foo" {
		t.Errorf("Expected arr[0] to be 'foo', got %v", arr[0])
	}
}

// TestArrayLiteralWithVariables тестирует массив с переменными
func TestArrayLiteralWithVariables(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("one", 1)
	ctx.Set("two", 2)

	expr, err := engine.CreateExpression(nil, "quux = [one, two]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result)
	}

	if len(arr) != 2 {
		t.Fatalf("Expected array length 2, got %d", len(arr))
	}

	if arr[0] != 1 {
		t.Errorf("Expected arr[0] to be 1, got %v", arr[0])
	}

	if arr[1] != 2 {
		t.Errorf("Expected arr[1] to be 2, got %v", arr[1])
	}

	// Проверяем, что переменная установлена
	quux := ctx.Get("quux")
	if quux == nil {
		t.Fatal("Variable 'quux' not set in context")
	}
}

// TestArrayLiteralWithNulls тестирует массив с null значениями
func TestArrayLiteralWithNulls(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "[null, 10]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result)
	}

	if len(arr) != 2 {
		t.Fatalf("Expected array length 2, got %d", len(arr))
	}

	if arr[0] != nil {
		t.Errorf("Expected arr[0] to be nil, got %v", arr[0])
	}

	// Проверяем значение (может быть разных типов)
	val1 := arr[1]
	var int1 int
	switch v := val1.(type) {
	case int:
		int1 = v
	case int64:
		int1 = int(v)
	case float64:
		int1 = int(v)
	default:
		t.Errorf("Unexpected type for arr[1]: %T", val1)
		return
	}
	
	if int1 != 10 {
		t.Errorf("Expected arr[1] to be 10, got %d", int1)
	}
}

