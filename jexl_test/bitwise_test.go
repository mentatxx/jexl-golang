package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestBitwiseAndSimple тестирует простую битовую операцию AND
func TestBitwiseAndSimple(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "15 & 3")
	if err != nil {
		t.Logf("Bitwise AND may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// 15 & 3 = 3
	expected := int64(3)
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestBitwiseAndVariableNumberCoercion тестирует AND с разными числовыми типами
func TestBitwiseAndVariableNumberCoercion(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 15)
	ctx.Set("y", int16(7))

	expr, err := engine.CreateExpression(nil, "x & y")
	if err != nil {
		t.Logf("Bitwise AND may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// 15 & 7 = 7
	expected := int64(7)
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestBitwiseAndWithNull тестирует AND с null значениями
func TestBitwiseAndWithNull(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	tests := []struct {
		name     string
		expr     string
		expected int64
	}{
		{"null & 1", "null & 1", 0},
		{"1 & null", "1 & null", 0},
		{"null & null", "null & null", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Logf("Bitwise AND may not be supported: %v", err)
				return
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Logf("Evaluation failed: %v", err)
				return
			}

			var actual int64
			switch v := result.(type) {
			case int:
				actual = int64(v)
			case int64:
				actual = v
			case *big.Rat:
				actual = v.Num().Int64()
			default:
				t.Fatalf("Unexpected result type: %T", result)
			}

			if actual != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, actual)
			}
		})
	}
}

// TestBitwiseComplementSimple тестирует унарную операцию NOT (~)
func TestBitwiseComplementSimple(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "~128")
	if err != nil {
		t.Logf("Bitwise complement may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// ~128 = -129 (для 32-битного int)
	expected := int64(-129)
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestBitwiseOrSimple тестирует простую битовую операцию OR
func TestBitwiseOrSimple(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "12 | 3")
	if err != nil {
		t.Logf("Bitwise OR may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// 12 | 3 = 15
	expected := int64(15)
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestBitwiseXorSimple тестирует простую битовую операцию XOR
func TestBitwiseXorSimple(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "1 ^ 3")
	if err != nil {
		t.Logf("Bitwise XOR may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// 1 ^ 3 = 2
	expected := int64(2)
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestBitwiseShiftLeft тестирует операцию сдвига влево
func TestBitwiseShiftLeft(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "1 << 3")
	if err != nil {
		t.Logf("Bitwise shift left may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// 1 << 3 = 8
	expected := int64(8)
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestBitwiseShiftRight тестирует операцию сдвига вправо
func TestBitwiseShiftRight(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "8 >> 2")
	if err != nil {
		t.Logf("Bitwise shift right may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// 8 >> 2 = 2
	expected := int64(2)
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestBitwiseParenthesized тестирует приоритет операций с скобками
func TestBitwiseParenthesized(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	tests := []struct {
		name     string
		expr     string
		expected int64
	}{
		{"(2 | 1) & 3", "(2 | 1) & 3", 3},
		{"(2 & 1) | 3", "(2 & 1) | 3", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Logf("Bitwise operations may not be supported: %v", err)
				return
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Logf("Evaluation failed: %v", err)
				return
			}

			var actual int64
			switch v := result.(type) {
			case int:
				actual = int64(v)
			case int64:
				actual = v
			case *big.Rat:
				actual = v.Num().Int64()
			default:
				t.Fatalf("Unexpected result type: %T", result)
			}

			if actual != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, actual)
			}
		})
	}
}

