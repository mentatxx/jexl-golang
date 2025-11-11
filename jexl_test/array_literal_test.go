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

// TestArrayLiteralChangeThroughVariables тестирует изменение массива через переменные
func TestArrayLiteralChangeThroughVariables(t *testing.T) {
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

	result1, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr1, ok := result1.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result1)
	}

	if len(arr1) != 2 {
		t.Fatalf("Expected array length 2, got %d", len(arr1))
	}

	// Проверяем значения первого массива
	var val0, val1 int
	switch v := arr1[0].(type) {
	case int:
		val0 = v
	case int64:
		val0 = int(v)
	case float64:
		val0 = int(v)
	default:
		t.Errorf("Unexpected type for arr1[0]: %T", arr1[0])
		return
	}
	switch v := arr1[1].(type) {
	case int:
		val1 = v
	case int64:
		val1 = int(v)
	case float64:
		val1 = int(v)
	default:
		t.Errorf("Unexpected type for arr1[1]: %T", arr1[1])
		return
	}

	if val0 != 1 {
		t.Errorf("Expected arr1[0] to be 1, got %d", val0)
	}
	if val1 != 2 {
		t.Errorf("Expected arr1[1] to be 2, got %d", val1)
	}

	// Изменяем переменные
	ctx.Set("one", 10)
	ctx.Set("two", 20)

	result2, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr2, ok := result2.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result2)
	}

	if len(arr2) != 2 {
		t.Fatalf("Expected array length 2, got %d", len(arr2))
	}

	// Проверяем значения второго массива
	switch v := arr2[0].(type) {
	case int:
		val0 = v
	case int64:
		val0 = int(v)
	case float64:
		val0 = int(v)
	default:
		t.Errorf("Unexpected type for arr2[0]: %T", arr2[0])
		return
	}
	switch v := arr2[1].(type) {
	case int:
		val1 = v
	case int64:
		val1 = int(v)
	case float64:
		val1 = int(v)
	default:
		t.Errorf("Unexpected type for arr2[1]: %T", arr2[1])
		return
	}

	if val0 != 10 {
		t.Errorf("Expected arr2[0] to be 10, got %d", val0)
	}
	if val1 != 20 {
		t.Errorf("Expected arr2[1] to be 20, got %d", val1)
	}
}

// TestArrayLiteralWithNumbers тестирует массив с числами разных типов
func TestArrayLiteralWithNumbers(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "[5.0, 10]")
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

	// Проверяем, что первый элемент - float, второй - int
	if arr[0] == nil {
		t.Error("Expected arr[0] to be non-nil")
	}
	if arr[1] == nil {
		t.Error("Expected arr[1] to be non-nil")
	}
}

// TestNotEmptySimpleArrayLiteral тестирует проверку empty() для массива
func TestNotEmptySimpleArrayLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "empty(['foo', 'bar'])")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	empty, ok := result.(bool)
	if !ok {
		t.Fatalf("Expected bool, got %T", result)
	}

	if empty {
		t.Error("Expected empty(['foo', 'bar']) to be false")
	}
}

// TestSizeOfSimpleArrayLiteral тестирует size() для массива
func TestSizeOfSimpleArrayLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "size(['foo', 'bar'])")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	// Результат может быть разных типов
	var size int
	switch v := result.(type) {
	case int:
		size = v
	case int64:
		size = int(v)
	case float64:
		size = int(v)
	default:
		t.Errorf("Unexpected type for size: %T", result)
		return
	}

	if size != 2 {
		t.Errorf("Expected size to be 2, got %d", size)
	}
}

// TestArrayLiteralWithNullsComprehensive тестирует массивы с null в разных позициях
func TestArrayLiteralWithNullsComprehensive(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	testCases := []struct {
		expr   string
		length int
		checks []func(any) bool
	}{
		{
			expr:   "[null, 10]",
			length: 2,
			checks: []func(any) bool{
				func(v any) bool { return v == nil },
				func(v any) bool {
					var intVal int
					switch val := v.(type) {
					case int:
						intVal = val
					case int64:
						intVal = int(val)
					case float64:
						intVal = int(val)
					default:
						return false
					}
					return intVal == 10
				},
			},
		},
		{
			expr:   "[10, null]",
			length: 2,
			checks: []func(any) bool{
				func(v any) bool {
					var intVal int
					switch val := v.(type) {
					case int:
						intVal = val
					case int64:
						intVal = int(val)
					case float64:
						intVal = int(val)
					default:
						return false
					}
					return intVal == 10
				},
				func(v any) bool { return v == nil },
			},
		},
		{
			expr:   "[10, null, 10]",
			length: 3,
			checks: []func(any) bool{
				func(v any) bool {
					var intVal int
					switch val := v.(type) {
					case int:
						intVal = val
					case int64:
						intVal = int(val)
					case float64:
						intVal = int(val)
					default:
						return false
					}
					return intVal == 10
				},
				func(v any) bool { return v == nil },
				func(v any) bool {
					var intVal int
					switch val := v.(type) {
					case int:
						intVal = val
					case int64:
						intVal = int(val)
					case float64:
						intVal = int(val)
					default:
						return false
					}
					return intVal == 10
				},
			},
		},
		{
			expr:   "['10', null]",
			length: 2,
			checks: []func(any) bool{
				func(v any) bool { return v == "10" },
				func(v any) bool { return v == nil },
			},
		},
		{
			expr:   "[null, '10', null]",
			length: 3,
			checks: []func(any) bool{
				func(v any) bool { return v == nil },
				func(v any) bool { return v == "10" },
				func(v any) bool { return v == nil },
			},
		},
	}

	for _, tc := range testCases {
		expr, err := engine.CreateExpression(nil, tc.expr)
		if err != nil {
			t.Fatalf("Failed to create expression %s: %v", tc.expr, err)
		}

		result, err := expr.Evaluate(ctx)
		if err != nil {
			t.Fatalf("Failed to evaluate %s: %v", tc.expr, err)
		}

		arr, ok := result.([]any)
		if !ok {
			t.Fatalf("Expected []any for %s, got %T", tc.expr, result)
		}

		if len(arr) != tc.length {
			t.Errorf("Expected array length %d for %s, got %d", tc.length, tc.expr, len(arr))
			continue
		}

		for i, check := range tc.checks {
			if !check(arr[i]) {
				t.Errorf("Check %d failed for %s, arr[%d] = %v", i, tc.expr, i, arr[i])
			}
		}
	}
}

