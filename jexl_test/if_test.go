package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestBlockElse тестирует, что if statement правильно обрабатывает блоки в else statement
func TestBlockElse(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (false) {1} else {2 ; 3}")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Результат должен быть 3 (последнее выражение в блоке)
	if result == nil {
		t.Fatal("Result is nil")
	}

	// Проверяем, что результат равен 3
	var expected int64 = 3
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestBlockIfTrue тестирует, что if statement правильно обрабатывает блоки
func TestBlockIfTrue(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (true) { 'hello'; }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}
}

// TestIfElseIfExpression тестирует if/else if/else выражения
func TestIfElseIfExpression(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (x == 1) { 10; } else if (x == 2) 20  else 30", "x")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	tests := []struct {
		x        int64
		expected int64
	}{
		{1, 10},
		{2, 20},
		{4, 30},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, err := script.Execute(ctx, tt.x)
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
				t.Fatalf("Unexpected result type: %T, value: %v", result, result)
			}

			if actual != tt.expected {
				t.Errorf("For x=%d, expected %d, got %d", tt.x, tt.expected, actual)
			}
		})
	}
}

// TestIfElseIfReturnExpression тестирует if/else if с return выражениями
func TestIfElseIfReturnExpression(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil,
		"if (x == 1) return 10;  if (x == 2) return 20  else if (x == 3) return 30; else return 40;", "x")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	tests := []struct {
		x        int64
		expected int64
	}{
		{1, 10},
		{2, 20},
		{3, 30},
		{4, 40},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, err := script.Execute(ctx, tt.x)
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
				t.Fatalf("Unexpected result type: %T, value: %v", result, result)
			}

			if actual != tt.expected {
				t.Errorf("For x=%d, expected %d, got %d", tt.x, tt.expected, actual)
			}
		})
	}
}

// TestIfElseIfReturnExpression0 тестирует другой вариант if/else if с return
func TestIfElseIfReturnExpression0(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil,
		"if (x == 1) return 10; if (x == 2)  return 20; else if (x == 3) return 30  else { return 40 }", "x")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	tests := []struct {
		x        int64
		expected int64
	}{
		{1, 10},
		{2, 20},
		{3, 30},
		{4, 40},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result, err := script.Execute(ctx, tt.x)
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
				t.Fatalf("Unexpected result type: %T, value: %v", result, result)
			}

			if actual != tt.expected {
				t.Errorf("For x=%d, expected %d, got %d", tt.x, tt.expected, actual)
			}
		})
	}
}

// TestIfWithArithmeticExpression тестирует if statement с арифметическими выражениями
func TestIfWithArithmeticExpression(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", int64(2))
	script, err := engine.CreateScript(nil, nil, "if ((x * 2) + 1 == 5) true;")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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

	if !boolResult {
		t.Error("Expected true, got false")
	}
}

// TestIfWithAssignment тестирует if statement с присваиванием
func TestIfWithAssignment(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 2.5)
	script, err := engine.CreateScript(nil, nil, "if ((x * 2) == 5) {y = 1} else {y = 2;}")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	_, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	result := ctx.Get("y")
	if result == nil {
		t.Fatal("y is nil")
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

	if actual != 1 {
		t.Errorf("Expected y to be 1, got %d", actual)
	}
}

// TestIfWithDecimalArithmeticExpression тестирует if statement с десятичными арифметическими выражениями
func TestIfWithDecimalArithmeticExpression(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 2.5)
	script, err := engine.CreateScript(nil, nil, "if ((x * 2) == 5) true")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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

	if !boolResult {
		t.Error("Expected true, got false")
	}
}

// TestIfWithSimpleExpression тестирует if statement с простым выражением
func TestIfWithSimpleExpression(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", int64(1))
	script, err := engine.CreateScript(nil, nil, "if (x == 1) true;")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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

	if !boolResult {
		t.Error("Expected true, got false")
	}
}

// TestNullCoalescing тестирует null coalescing оператор (??)
func TestNullCoalescing(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	expr, err := engine.CreateExpression(nil, "x??true")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	// x не определен, должен вернуть true
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

	if !boolResult {
		t.Error("Expected true, got false")
	}

	// x = false, должен вернуть false
	ctx.Set("x", false)
	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	switch v := result.(type) {
	case bool:
		boolResult = v
	case *big.Rat:
		boolResult = v.Sign() != 0
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if boolResult {
		t.Error("Expected false, got true")
	}

	// y не определен, должен вернуть 1
	expr2, err := engine.CreateExpression(nil, "y??1")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr2.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	var intResult int64
	switch v := result.(type) {
	case int64:
		intResult = v
	case *big.Rat:
		if v.IsInt() {
			intResult = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if intResult != 1 {
		t.Errorf("Expected 1, got %d", intResult)
	}

	// y = 0, должен вернуть 0
	ctx.Set("y", int64(0))
	result, err = expr2.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	switch v := result.(type) {
	case int64:
		intResult = v
	case *big.Rat:
		if v.IsInt() {
			intResult = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if intResult != 0 {
		t.Errorf("Expected 0, got %d", intResult)
	}
}

// TestNullCoalescingScript тестирует null coalescing оператор в скрипте
func TestNullCoalescingScript(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "x??true")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	// x не определен, должен вернуть true
	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
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

	if !boolResult {
		t.Error("Expected true, got false")
	}

	// x = false, должен вернуть false
	ctx.Set("x", false)
	result, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	switch v := result.(type) {
	case bool:
		boolResult = v
	case *big.Rat:
		boolResult = v.Sign() != 0
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if boolResult {
		t.Error("Expected false, got true")
	}

	// y не определен, должен вернуть 1
	script2, err := engine.CreateScript(nil, nil, "y??1")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err = script2.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	var intResult int64
	switch v := result.(type) {
	case int64:
		intResult = v
	case *big.Rat:
		if v.IsInt() {
			intResult = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if intResult != 1 {
		t.Errorf("Expected 1, got %d", intResult)
	}

	// y = 0, должен вернуть 0
	ctx.Set("y", int64(0))
	result, err = script2.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	switch v := result.(type) {
	case int64:
		intResult = v
	case *big.Rat:
		if v.IsInt() {
			intResult = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if intResult != 0 {
		t.Errorf("Expected 0, got %d", intResult)
	}
}

// TestSimpleElse тестирует простой else
func TestSimpleElse(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (false) 1 else 2;")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
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

	if actual != 2 {
		t.Errorf("Expected 2, got %d", actual)
	}
}

// TestSimpleIfFalse тестирует, что if false не выполняет true statement
func TestSimpleIfFalse(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (false) 1")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

// TestSimpleIfTrue тестирует, что if true выполняет true statement
func TestSimpleIfTrue(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (true) 1")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
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

	if actual != 1 {
		t.Errorf("Expected 1, got %d", actual)
	}
}

