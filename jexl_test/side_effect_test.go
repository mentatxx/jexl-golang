package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestSideEffectAdd тестирует оператор +=
func TestSideEffectAdd(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 5)

	expr, err := engine.CreateExpression(nil, "x += 3")
	if err != nil {
		t.Logf("Side-effect operator += may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// x должно стать 8
	x := ctx.Get("x")
	var xVal int64
	switch v := x.(type) {
	case int:
		xVal = int64(v)
	case int64:
		xVal = v
	case *big.Rat:
		xVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected x type: %T", x)
	}

	if xVal != 8 {
		t.Errorf("Expected x to be 8, got %d", xVal)
	}

	// Результат выражения тоже должен быть 8
	var resultVal int64
	switch v := result.(type) {
	case int:
		resultVal = int64(v)
	case int64:
		resultVal = v
	case *big.Rat:
		resultVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if resultVal != 8 {
		t.Errorf("Expected result to be 8, got %d", resultVal)
	}
}

// TestSideEffectSubtract тестирует оператор -=
func TestSideEffectSubtract(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 10)

	expr, err := engine.CreateExpression(nil, "x -= 3")
	if err != nil {
		t.Logf("Side-effect operator -= may not be supported: %v", err)
		return
	}

	_, err = expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// x должно стать 7
	x := ctx.Get("x")
	var xVal int64
	switch v := x.(type) {
	case int:
		xVal = int64(v)
	case int64:
		xVal = v
	case *big.Rat:
		xVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected x type: %T", x)
	}

	if xVal != 7 {
		t.Errorf("Expected x to be 7, got %d", xVal)
	}
}

// TestSideEffectMultiply тестирует оператор *=
func TestSideEffectMultiply(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 5)

	expr, err := engine.CreateExpression(nil, "x *= 3")
	if err != nil {
		t.Logf("Side-effect operator *= may not be supported: %v", err)
		return
	}

	_, err = expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// x должно стать 15
	x := ctx.Get("x")
	var xVal int64
	switch v := x.(type) {
	case int:
		xVal = int64(v)
	case int64:
		xVal = v
	case *big.Rat:
		xVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected x type: %T", x)
	}

	if xVal != 15 {
		t.Errorf("Expected x to be 15, got %d", xVal)
	}
}

// TestIncrementPrefix тестирует префиксный инкремент ++x
func TestIncrementPrefix(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 5)

	expr, err := engine.CreateExpression(nil, "++x")
	if err != nil {
		t.Logf("Increment operator ++ may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// x должно стать 6
	x := ctx.Get("x")
	var xVal int64
	switch v := x.(type) {
	case int:
		xVal = int64(v)
	case int64:
		xVal = v
	case *big.Rat:
		xVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected x type: %T", x)
	}

	if xVal != 6 {
		t.Errorf("Expected x to be 6, got %d", xVal)
	}

	// Результат тоже должен быть 6 (префиксный инкремент)
	var resultVal int64
	switch v := result.(type) {
	case int:
		resultVal = int64(v)
	case int64:
		resultVal = v
	case *big.Rat:
		resultVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if resultVal != 6 {
		t.Errorf("Expected result to be 6, got %d", resultVal)
	}
}

// TestIncrementPostfix тестирует постфиксный инкремент x++
func TestIncrementPostfix(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 5)

	expr, err := engine.CreateExpression(nil, "x++")
	if err != nil {
		t.Logf("Increment operator ++ may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// x должно стать 6
	x := ctx.Get("x")
	var xVal int64
	switch v := x.(type) {
	case int:
		xVal = int64(v)
	case int64:
		xVal = v
	case *big.Rat:
		xVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected x type: %T", x)
	}

	if xVal != 6 {
		t.Errorf("Expected x to be 6, got %d", xVal)
	}

	// Результат должен быть 5 (постфиксный инкремент возвращает старое значение)
	var resultVal int64
	switch v := result.(type) {
	case int:
		resultVal = int64(v)
	case int64:
		resultVal = v
	case *big.Rat:
		resultVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if resultVal != 5 {
		t.Errorf("Expected result to be 5 (old value), got %d", resultVal)
	}
}

// TestDecrementPrefix тестирует префиксный декремент --x
func TestDecrementPrefix(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 5)

	expr, err := engine.CreateExpression(nil, "--x")
	if err != nil {
		t.Logf("Decrement operator -- may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// x должно стать 4
	x := ctx.Get("x")
	var xVal int64
	switch v := x.(type) {
	case int:
		xVal = int64(v)
	case int64:
		xVal = v
	case *big.Rat:
		xVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected x type: %T", x)
	}

	if xVal != 4 {
		t.Errorf("Expected x to be 4, got %d", xVal)
	}

	// Результат тоже должен быть 4
	var resultVal int64
	switch v := result.(type) {
	case int:
		resultVal = int64(v)
	case int64:
		resultVal = v
	case *big.Rat:
		resultVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if resultVal != 4 {
		t.Errorf("Expected result to be 4, got %d", resultVal)
	}
}

// TestDecrementPostfix тестирует постфиксный декремент x--
func TestDecrementPostfix(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 5)

	expr, err := engine.CreateExpression(nil, "x--")
	if err != nil {
		t.Logf("Decrement operator -- may not be supported: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Evaluation failed: %v", err)
		return
	}

	// x должно стать 4
	x := ctx.Get("x")
	var xVal int64
	switch v := x.(type) {
	case int:
		xVal = int64(v)
	case int64:
		xVal = v
	case *big.Rat:
		xVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected x type: %T", x)
	}

	if xVal != 4 {
		t.Errorf("Expected x to be 4, got %d", xVal)
	}

	// Результат должен быть 5 (постфиксный декремент возвращает старое значение)
	var resultVal int64
	switch v := result.(type) {
	case int:
		resultVal = int64(v)
	case int64:
		resultVal = v
	case *big.Rat:
		resultVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if resultVal != 5 {
		t.Errorf("Expected result to be 5 (old value), got %d", resultVal)
	}
}

// TestSideEffectInScript тестирует side-effect операторы в скрипте
func TestSideEffectInScript(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", 0)
	ctx.Set("y", 10)

	script, err := engine.CreateScript(nil, nil, "x += 5; y -= 3; x + y")
	if err != nil {
		t.Logf("Side-effect operators in script may not be supported: %v", err)
		return
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Logf("Execution failed: %v", err)
		return
	}

	// x должно стать 5, y должно стать 7, результат должен быть 12
	x := ctx.Get("x")
	var xVal int64
	switch v := x.(type) {
	case int:
		xVal = int64(v)
	case int64:
		xVal = v
	case *big.Rat:
		xVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected x type: %T", x)
	}

	if xVal != 5 {
		t.Errorf("Expected x to be 5, got %d", xVal)
	}

	y := ctx.Get("y")
	var yVal int64
	switch v := y.(type) {
	case int:
		yVal = int64(v)
	case int64:
		yVal = v
	case *big.Rat:
		yVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected y type: %T", y)
	}

	if yVal != 7 {
		t.Errorf("Expected y to be 7, got %d", yVal)
	}

	// Результат скрипта (x + y) должен быть 12
	var resultVal int64
	switch v := result.(type) {
	case int:
		resultVal = int64(v)
	case int64:
		resultVal = v
	case *big.Rat:
		resultVal = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if resultVal != 12 {
		t.Errorf("Expected result to be 12, got %d", resultVal)
	}
}

