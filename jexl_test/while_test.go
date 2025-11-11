package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestSimpleWhileFalse тестирует простой while с false условием
func TestSimpleWhileFalse(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "while (false) ;")
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

// TestWhileExecutesExpressionWhenLooping тестирует, что while выполняет выражение при цикле
func TestWhileExecutesExpressionWhenLooping(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", int64(1))
	script, err := engine.CreateScript(nil, nil, "while (x < 10) x = x + 1;")
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

	if actual != 10 {
		t.Errorf("Expected 10, got %d", actual)
	}
}

// TestWhileWithBlock тестирует while с блоком
func TestWhileWithBlock(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", int64(1))
	ctx.Set("y", int64(1))
	script, err := engine.CreateScript(nil, nil, "while (x < 10) { x = x + 1; y = y * 2; }")
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

	// y = 1 * 2^9 = 512
	if actual != 512 {
		t.Errorf("Expected 512, got %d", actual)
	}

	x := ctx.Get("x")
	var xVal int64
	switch v := x.(type) {
	case int64:
		xVal = v
	case *big.Rat:
		if v.IsInt() {
			xVal = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected x type: %T", x)
	}

	if xVal != 10 {
		t.Errorf("Expected x to be 10, got %d", xVal)
	}

	y := ctx.Get("y")
	var yVal int64
	switch v := y.(type) {
	case int64:
		yVal = v
	case *big.Rat:
		if v.IsInt() {
			yVal = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected y type: %T", y)
	}

	if yVal != 512 {
		t.Errorf("Expected y to be 512, got %d", yVal)
	}
}

