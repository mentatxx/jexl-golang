package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestArithmeticOperations тестирует базовые арифметические операции
func TestArithmeticOperations(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("a", 10)
	ctx.Set("b", 3)

	tests := []struct {
		name     string
		expr     string
		validate func(any) bool
	}{
		{"addition", "a + b", func(r any) bool {
			rat, ok := r.(*big.Rat)
			if !ok {
				return false
			}
			expected := big.NewRat(13, 1)
			return rat.Cmp(expected) == 0
		}},
		{"subtraction", "a - b", func(r any) bool {
			rat, ok := r.(*big.Rat)
			if !ok {
				return false
			}
			expected := big.NewRat(7, 1)
			return rat.Cmp(expected) == 0
		}},
		{"multiplication", "a * b", func(r any) bool {
			rat, ok := r.(*big.Rat)
			if !ok {
				return false
			}
			expected := big.NewRat(30, 1)
			return rat.Cmp(expected) == 0
		}},
		{"division", "a / b", func(r any) bool {
			rat, ok := r.(*big.Rat)
			if !ok {
				return false
			}
			// 10/3 = 3.333...
			expected := big.NewRat(10, 3)
			return rat.Cmp(expected) == 0
		}},
		{"modulo", "a % b", func(r any) bool {
			rat, ok := r.(*big.Rat)
			if !ok {
				return false
			}
			expected := big.NewRat(1, 1) // 10 % 3 = 1
			return rat.Cmp(expected) == 0
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

			if !tt.validate(result) {
				t.Errorf("Validation failed for result: %v", result)
			}
		})
	}
}

// TestArithmeticWithDifferentTypes тестирует арифметику с разными типами
func TestArithmeticWithDifferentTypes(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("intVal", int(5))
	ctx.Set("int64Val", int64(10))
	ctx.Set("floatVal", float64(2.5))

	tests := []struct {
		name string
		expr string
	}{
		{"int + int64", "intVal + int64Val"},
		{"int64 + float", "int64Val + floatVal"},
		{"float + int", "floatVal + intVal"},
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

			if result == nil {
				t.Error("Result is nil")
			}
		})
	}
}

// TestArithmeticNegate тестирует унарный минус
func TestArithmeticNegate(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 5)

	expr, err := engine.CreateExpression(nil, "-x")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	rat, ok := result.(*big.Rat)
	if !ok {
		t.Fatalf("Expected *big.Rat, got %T", result)
	}

	expected := big.NewRat(-5, 1)
	if rat.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, rat)
	}
}

// TestArithmeticPrecedence тестирует приоритет операций
func TestArithmeticPrecedence(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("a", 2)
	ctx.Set("b", 3)
	ctx.Set("c", 4)

	// 2 + 3 * 4 = 2 + 12 = 14
	expr, err := engine.CreateExpression(nil, "a + b * c")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	rat, ok := result.(*big.Rat)
	if !ok {
		t.Fatalf("Expected *big.Rat, got %T", result)
	}

	expected := big.NewRat(14, 1)
	if rat.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, rat)
	}
}

// TestArithmeticDivisionByZero тестирует деление на ноль
func TestArithmeticDivisionByZero(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("a", 10)
	ctx.Set("b", 0)

	expr, err := engine.CreateExpression(nil, "a / b")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	// Используем recover для обработки паники
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Паника ожидаема для деления на ноль
				t.Logf("Caught expected panic: %v", r)
			}
		}()
		_, err = expr.Evaluate(ctx)
		if err != nil {
			// Ошибка тоже приемлема
			return
		}
		// Если нет ни паники, ни ошибки - это проблема
		t.Error("Expected error or panic for division by zero")
	}()
}
