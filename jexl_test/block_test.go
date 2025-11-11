package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestBlockExecutesAll тестирует, что блок выполняет все выражения
func TestBlockExecutesAll(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (true) { x = 'Hello'; y = 'World';}")
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

	x := ctx.Get("x")
	if x != "Hello" {
		t.Errorf("Expected x to be 'Hello', got %v", x)
	}

	y := ctx.Get("y")
	if y != "World" {
		t.Errorf("Expected y to be 'World', got %v", y)
	}
}

// TestBlockLastExecuted01 тестирует, что блок возвращает последнее выполненное выражение
func TestBlockLastExecuted01(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (true) { x = 1; } else { x = 2; }")
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

// TestBlockLastExecuted02 тестирует, что блок возвращает последнее выполненное выражение (else ветка)
func TestBlockLastExecuted02(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (false) { x = 1; } else { x = 2; }")
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

// TestBlockSimple тестирует простой блок
func TestBlockSimple(t *testing.T) {
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

// TestEmptyBlock тестирует пустой блок
func TestEmptyBlock(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (true) { }")
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

// TestNestedBlock тестирует вложенный блок
func TestNestedBlock(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "if (true) { x = 'hello'; y = 'world'; if (true) { x; } y; }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result != "world" {
		t.Errorf("Expected 'world', got %v", result)
	}
}

