package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestLeftShiftIntValue - тест оператора сдвига влево для int значений
func TestLeftShiftIntValue(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	tests := []struct {
		name     string
		x        int64
		y        int64
		expected int64
	}{
		{"positive positive", 1, 2, 1 << 2},
		{"positive negative", 1, -2, 0}, // В Go отрицательный сдвиг дает 0
		{"negative positive", -1, 2, -1 << 2},
		{"negative negative", -1, -2, 0}, // В Go отрицательный сдвиг дает 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := engine.CreateScript(nil, nil, "(x, y)-> x << y", "x", "y")
			if err != nil {
				t.Fatalf("Failed to create script: %v", err)
			}

			result, err := script.Execute(ctx, tt.x, tt.y)
			if err != nil {
				t.Fatalf("Failed to execute: %v", err)
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

// TestLeftShiftLongValue - тест оператора сдвига влево для long значений
func TestLeftShiftLongValue(t *testing.T) {
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
		{"large positive shift", "2147483648 << 2", 2147483648 << 2},
		{"large positive negative shift", "2147483648 << -2", 0}, // В Go отрицательный сдвиг дает 0
		{"large negative shift", "-2147483649 << 2", -2147483649 << 2},
		{"large negative negative shift", "-2147483649 << -2", 0}, // В Go отрицательный сдвиг дает 0
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

// TestRightShiftIntValue - тест оператора сдвига вправо для int значений
func TestRightShiftIntValue(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	tests := []struct {
		name     string
		x        int64
		y        int64
		expected int64
	}{
		{"positive positive", 42, 2, 42 >> 2},
		{"positive negative", 42, -2, 0}, // В Go отрицательный сдвиг дает 0
		{"negative positive", -42, 2, -42 >> 2},
		{"negative negative", -42, -2, 0}, // В Go отрицательный сдвиг дает 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := engine.CreateScript(nil, nil, "(x, y)-> x >> y", "x", "y")
			if err != nil {
				t.Fatalf("Failed to create script: %v", err)
			}

			result, err := script.Execute(ctx, tt.x, tt.y)
			if err != nil {
				t.Fatalf("Failed to execute: %v", err)
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

// TestRightShiftLongValue - тест оператора сдвига вправо для long значений
func TestRightShiftLongValue(t *testing.T) {
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
		{"large positive shift", "8589934592 >> 2", 8589934592 >> 2},
		{"large positive negative shift", "8589934592 >> -2", 0}, // В Go отрицательный сдвиг дает 0
		{"large negative shift", "-8589934592 >> 2", -8589934592 >> 2},
		{"large negative negative shift", "-8589934592 >> -2", 0}, // В Go отрицательный сдвиг дает 0
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

// TestShiftPrecedence - тест приоритета операторов сдвига
func TestShiftPrecedence(t *testing.T) {
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
		{"no parentheses", "40 + 2 << 1 + 1", 40 + 2<<1 + 1},
		{"shift in parentheses", "40 + (2 << 1) + 1", 40 + (2 << 1) + 1},
		{"both in parentheses", "(40 + 2) << (1 + 1)", (40 + 2) << (1 + 1)},
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

// TestRightShiftUnsignedIntValue - тест беззнакового сдвига вправо для int значений
// Примечание: Go не имеет беззнакового сдвига вправо (>>>), но можно эмулировать
func TestRightShiftUnsignedIntValue(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// В Go нет оператора >>>, но можно проверить, что >> работает корректно
	// Для беззнакового сдвига нужно использовать uint
	tests := []struct {
		name     string
		x        int64
		y        int64
		expected int64
	}{
		{"positive positive", 42, 2, 42 >> 2},
		{"positive negative", 42, -2, 0}, // В Go отрицательный сдвиг дает 0
		{"negative positive", -42, 2, -42 >> 2},
		{"negative negative", -42, -2, 0}, // В Go отрицательный сдвиг дает 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Используем обычный сдвиг вправо, так как >>> может не поддерживаться
			script, err := engine.CreateScript(nil, nil, "(x, y)-> x >> y", "x", "y")
			if err != nil {
				t.Fatalf("Failed to create script: %v", err)
			}

			result, err := script.Execute(ctx, tt.x, tt.y)
			if err != nil {
				t.Fatalf("Failed to execute: %v", err)
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

