package jexl_test

import (
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestSimpleAssignment тестирует простое присваивание
func TestSimpleAssignment(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("aString", "Hello")

	expr, err := engine.CreateExpression(nil, "hello = 'world'")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != "world" {
		t.Errorf("Expected 'world', got %v", result)
	}

	if ctx.Get("hello") != "world" {
		t.Error("Variable 'hello' not set in context")
	}
}

// TestPropertyAssignment тестирует присваивание свойств
func TestPropertyAssignment(t *testing.T) {
	type Foo struct {
		Property1 string
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	foo := &Foo{Property1: "initial"}
	ctx.Set("foo", foo)

	expr, err := engine.CreateExpression(nil, "foo.Property1 = '99'")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != "99" {
		t.Errorf("Expected '99', got %v", result)
	}

	if foo.Property1 != "99" {
		t.Errorf("Expected foo.Property1 to be '99', got %s", foo.Property1)
	}
}

// TestMapAssignment тестирует присваивание в мапу
func TestMapAssignment(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	data := map[string]any{"foo": 1}
	ctx.Set("data", data)

	expr, err := engine.CreateExpression(nil, "data['bar'] = 99")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	_, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	barVal := data["bar"]
	// Проверяем значение (может быть разных типов)
	var intVal int
	switch v := barVal.(type) {
	case int:
		intVal = v
	case int64:
		intVal = int(v)
	case float64:
		intVal = int(v)
	default:
		t.Errorf("Unexpected type for data['bar']: %T, value: %v", barVal, barVal)
		return
	}
	
	if intVal != 99 {
		t.Errorf("Expected data['bar'] to be 99, got %d", intVal)
	}
}

// TestNestedPropertyAssignment тестирует вложенное присваивание свойств
func TestNestedPropertyAssignment(t *testing.T) {
	type Froboz struct {
		Value int
	}
	type Quux struct {
		Froboz *Froboz
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	quux := &Quux{Froboz: &Froboz{Value: 0}}
	ctx.Set("quux", quux)

	expr, err := engine.CreateExpression(nil, "quux.Froboz.Value = 10")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != 10 {
		t.Errorf("Expected 10, got %v", result)
	}

	if quux.Froboz.Value != 10 {
		t.Errorf("Expected quux.Froboz.Value to be 10, got %d", quux.Froboz.Value)
	}
}

// TestArrayAssignment тестирует присваивание в массив
func TestArrayAssignment(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	array := []any{100, 101, 102}
	ctx.Set("array", array)

	expr, err := engine.CreateExpression(nil, "array[1] = 1010")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != 1010 {
		t.Errorf("Expected 1010, got %v", result)
	}

	if array[1] != 1010 {
		t.Errorf("Expected array[1] to be 1010, got %v", array[1])
	}
}

// TestExpressionAssignment тестирует присваивание результата выражения
func TestExpressionAssignment(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("a", 5)
	ctx.Set("b", 3)

	expr, err := engine.CreateExpression(nil, "result = a + b")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	_, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	result := ctx.Get("result")
	if result == nil {
		t.Fatal("Variable 'result' not set in context")
	}
}

