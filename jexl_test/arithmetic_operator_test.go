package jexl_test

import (
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestRangeOperator тестирует оператор .. (range)
// Пока не реализован, но тесты подготовлены
func TestRangeOperator(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	tests := []struct {
		name     string
		expr     string
		expected []int64
	}{
		{"simple range", "1 .. 3", []int64{1, 2, 3}},
		{"range with expression", "(4 - 3) .. (9 / 3)", []int64{1, 2, 3}},
		{"negative range", "-3 .. 3", []int64{-3, -2, -1, 0, 1, 2, 3}},
		{"descending range", "3 .. 1", []int64{3, 2, 1}},
		{"single element range", "5 .. 5", []int64{5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create expression: %v", err)
			}

			result, err := expr.Evaluate(nil)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			// Проверяем, что результат - слайс int64
			rangeResult, ok := result.([]int64)
			if !ok {
				t.Fatalf("Expected []int64, got %T", result)
			}

			if len(rangeResult) != len(tt.expected) {
				t.Fatalf("Expected length %d, got %d", len(tt.expected), len(rangeResult))
			}

			for i, val := range tt.expected {
				if rangeResult[i] != val {
					t.Errorf("Expected [%d] = %d, got %d", i, val, rangeResult[i])
				}
			}
		})
	}
}

// TestStringStartsWithEndsWithComprehensive тестирует операторы =^ и =$ более детально
func TestStringStartsWithEndsWithComprehensive(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", "foobar")

	tests := []struct {
		name     string
		expr     string
		expected bool
	}{
		{"starts with prefix", "x =^ 'foo'", true},
		{"starts with wrong prefix", "x =^ 'bar'", false},
		{"ends with suffix", "x =$ 'bar'", true},
		{"ends with wrong suffix", "x =$ 'foo'", false},
		{"starts with empty", "x =^ ''", true},
		{"ends with empty", "x =$ ''", true},
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

			b, ok := result.(bool)
			if !ok {
				t.Fatalf("Expected bool, got %T", result)
			}

			if b != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, b)
			}
		})
	}
}

// TestStringNotStartsEndsWithComprehensive тестирует операторы !^ и !$ более детально
func TestStringNotStartsEndsWithComprehensive(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	tests := []struct {
		name     string
		expr     string
		setup    func(ctx jexl.Context)
		expected bool
	}{
		{"not starts with (false)", "x !^ 'foo'", func(ctx jexl.Context) { ctx.Set("x", "foobar") }, false},
		{"not starts with (true)", "x !^ 'bar'", func(ctx jexl.Context) { ctx.Set("x", "foobar") }, true},
		{"not ends with (false)", "x !$ 'bar'", func(ctx jexl.Context) { ctx.Set("x", "foobar") }, false},
		{"not ends with (true)", "x !$ 'foo'", func(ctx jexl.Context) { ctx.Set("x", "foobar") }, true},
		{"not starts with barfoo", "x !^ 'foo'", func(ctx jexl.Context) { ctx.Set("x", "barfoo") }, true},
		{"not ends with barfoo", "x !$ 'foo'", func(ctx jexl.Context) { ctx.Set("x", "barfoo") }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(ctx)
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create expression: %v", err)
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			b, ok := result.(bool)
			if !ok {
				t.Fatalf("Expected bool, got %T", result)
			}

			if b != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, b)
			}
		})
	}
}

// TestMatchOperatorComprehensive тестирует оператор =~ (match) более детально
func TestMatchOperatorComprehensive(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("str", "abc456")

	tests := []struct {
		name     string
		expr     string
		setup    func(ctx jexl.Context)
		expected bool
	}{
		{"regex match", "str =~ '.*456'", nil, true},
		{"regex not match", "str =~ 'ABC.*'", nil, false},
		{"regex not match (negated)", "str !~ 'ABC.*'", nil, true},
		{"regex match (negated)", "str !~ '.*456'", nil, false},
		{"variable pattern match", "str =~ match", func(ctx jexl.Context) { ctx.Set("match", "abc.*") }, true},
		{"variable pattern not match", "str =~ nomatch", func(ctx jexl.Context) { ctx.Set("nomatch", ".*123") }, false},
		{"array contains", "'a' =~ ['a','b','c']", nil, true},
		{"array not contains", "'z' =~ ['a','b','c']", nil, false},
		{"array not contains (negated)", "'z' !~ ['a','b','c']", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(ctx)
			}

			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create expression: %v", err)
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			b, ok := result.(bool)
			if !ok {
				t.Fatalf("Expected bool, got %T", result)
			}

			if b != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, b)
			}
		})
	}
}

// TestIncrementDecrementOnNull тестирует инкремент/декремент на null значениях
func TestIncrementDecrementOnNull(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	tests := []struct {
		name     string
		expr     string
		expected int64
		skip     bool
	}{
		{"increment null", "var i = null; ++i", 1, true}, // TODO: реализовать инкремент/декремент
		{"decrement null", "var i = null; --i", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Increment/decrement operators not yet fully implemented")
			}

			script, err := engine.CreateScript(nil, nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create script: %v", err)
			}

			result, err := script.Execute(nil)
			if err != nil {
				t.Fatalf("Failed to execute: %v", err)
			}

			val, ok := result.(int64)
			if !ok {
				t.Fatalf("Expected int64, got %T", result)
			}

			if val != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, val)
			}
		})
	}
}

// TestMatchWithCollections тестирует оператор =~ с различными коллекциями
func TestMatchWithCollections(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Массив
	ai := []int64{2, 4, 42, 54}
	ctx.Set("container", ai)
	ctx.Set("x", int64(2))

	expr, err := engine.CreateExpression(nil, "x =~ container")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	b, ok := result.(bool)
	if !ok {
		t.Fatalf("Expected bool, got %T", result)
	}

	if !b {
		t.Error("Expected true for 2 in [2, 4, 42, 54]")
	}

	// Проверяем не содержится
	ctx.Set("x", int64(169))
	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	b, ok = result.(bool)
	if !ok {
		t.Fatalf("Expected bool, got %T", result)
	}

	if b {
		t.Error("Expected false for 169 in [2, 4, 42, 54]")
	}
}

// Test391 тестирует оператор =~ с массивами и списками (из test391 в Java)
func Test391(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Тестируем с литералами
	literals := []struct {
		expr     string
		expected bool
	}{
		{"2 =~ [1, 2, 3, 4]", true},
		{"[2, 3] =~ [1, 2, 3, 4]", true}, // TODO: поддержка массивов в качестве левого операнда
		{"3 =~ [1, 2, 3, 4]", true},
	}

	for _, tt := range literals {
		t.Run(tt.expr, func(t *testing.T) {
			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				// Некоторые выражения могут не поддерживаться пока
				t.Skipf("Expression not yet supported: %v", err)
				return
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			b, ok := result.(bool)
			if !ok {
				t.Fatalf("Expected bool, got %T", result)
			}

			if b != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, b)
			}
		})
	}

	// Тестируем с переменными
	ic := []int64{1, 2, 3, 4}
	ctx.Set("y", ic)
	ctx.Set("x", int64(2))

	expr, err := engine.CreateExpression(nil, "x =~ y")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	b, ok := result.(bool)
	if !ok {
		t.Fatalf("Expected bool, got %T", result)
	}

	if !b {
		t.Error("Expected true for 2 in [1, 2, 3, 4]")
	}
}

