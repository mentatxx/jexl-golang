package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestForEachBreakBroken тестирует, что break вне цикла вызывает ошибку
func TestForEachBreakBroken(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	_, err = engine.CreateScript(nil, nil, "if (true) { break; }")
	if err == nil {
		t.Error("Expected parsing error for break outside loop")
	}
}

// TestForEachBreakMethod тестирует break в foreach цикле
func TestForEachBreakMethod(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil,
		"var rr = -1; for (var item : [1, 2, 3 ,4 ,5, 6]) { if (item == 3) { rr = item; break; }} rr")
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

	if actual != 3 {
		t.Errorf("Expected 3, got %d", actual)
	}
}

// TestForEachContinueBroken тестирует, что continue вне цикла вызывает ошибку
func TestForEachContinueBroken(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	_, err = engine.CreateScript(nil, nil, "var rr = 0; continue;")
	if err == nil {
		t.Error("Expected parsing error for continue outside loop")
	}
}

// TestForEachContinueMethod тестирует continue в foreach цикле
func TestForEachContinueMethod(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil,
		"var rr = 0; for (var item : [1, 2, 3 ,4 ,5, 6]) { if (item <= 3) continue; rr = rr + item;}")
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

	// 4 + 5 + 6 = 15
	if actual != 15 {
		t.Errorf("Expected 15, got %d", actual)
	}
}

// TestForEachWithArray тестирует foreach с массивом
func TestForEachWithArray(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	list := []any{"Hello", "World"}
	ctx.Set("list", list)
	script, err := engine.CreateScript(nil, nil, "for (item : list) item")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result != "World" {
		t.Errorf("Expected 'World', got %v", result)
	}
}

// TestForEachWithBlock тестирует foreach с блоком
func TestForEachWithBlock(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	list := []any{int64(2), int64(3)}
	ctx.Set("list", list)
	ctx.Set("x", int64(1))
	script, err := engine.CreateScript(nil, nil, "for (var in : list) { x = x + in; }")
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

	// 1 + 2 + 3 = 6
	if actual != 6 {
		t.Errorf("Expected 6, got %d", actual)
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

	if xVal != 6 {
		t.Errorf("Expected x to be 6, got %d", xVal)
	}
}

// TestForEachWithCollection тестирует foreach с коллекцией
func TestForEachWithCollection(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	list := []any{"Hello", "World"}
	ctx.Set("list", list)
	script, err := engine.CreateScript(nil, nil, "for (var item : list) item")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result != "World" {
		t.Errorf("Expected 'World', got %v", result)
	}
}

// TestForEachWithEmptyList тестирует foreach с пустым списком
func TestForEachWithEmptyList(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	list := []any{}
	ctx.Set("list", list)
	script, err := engine.CreateScript(nil, nil, "for (item : list) 1+1")
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

// TestForEachWithEmptyStatement тестирует foreach с пустым statement
func TestForEachWithEmptyStatement(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	list := []any{}
	ctx.Set("list", list)
	script, err := engine.CreateScript(nil, nil, "for (item : list) ;")
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

// TestForEachWithMap тестирует foreach с мапой
func TestForEachWithMap(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	m := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}
	ctx.Set("list", m)
	script, err := engine.CreateScript(nil, nil, "for(item : list) item")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Результат должен быть одним из значений мапы
	if result == nil {
		t.Fatal("Result is nil")
	}
}

