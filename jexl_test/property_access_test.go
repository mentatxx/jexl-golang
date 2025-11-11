package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestPropertyAccessSimple тестирует простой доступ к свойству
func TestPropertyAccessSimple(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	type TestStruct struct {
		Value int
	}

	ctx := jexl.NewMapContext()
	obj := &TestStruct{Value: 42}
	ctx.Set("obj", obj)

	expr, err := engine.CreateExpression(nil, "obj.Value")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	var actual int64
	switch v := result.(type) {
	case int64:
		actual = v
	case int:
		actual = int64(v)
	case *big.Rat:
		if v.IsInt() {
			actual = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}
}

// TestPropertyAccessNested тестирует вложенный доступ к свойствам
func TestPropertyAccessNested(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	type Inner struct {
		Value int
	}
	type Outer struct {
		Inner *Inner
	}

	ctx := jexl.NewMapContext()
	obj := &Outer{Inner: &Inner{Value: 42}}
	ctx.Set("obj", obj)

	expr, err := engine.CreateExpression(nil, "obj.Inner.Value")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	var actual int64
	switch v := result.(type) {
	case int64:
		actual = v
	case int:
		actual = int64(v)
	case *big.Rat:
		if v.IsInt() {
			actual = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}
}

// TestStructWithGetter - структура с getter/setter для тестов
type TestStructWithGetter struct {
	value int
}

func (t *TestStructWithGetter) GetValue() int {
	return t.value
}

func (t *TestStructWithGetter) SetValue(v int) {
	t.value = v
}

// TestPropertyAccessWithGetter тестирует доступ к свойству через getter
func TestPropertyAccessWithGetter(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	obj := &TestStructWithGetter{}
	obj.SetValue(42)
	ctx.Set("obj", obj)

	expr, err := engine.CreateExpression(nil, "obj.Value")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	var actual int64
	switch v := result.(type) {
	case int64:
		actual = v
	case int:
		actual = int64(v)
	case *big.Rat:
		if v.IsInt() {
			actual = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}
}

// TestPropertyAccessMap тестирует доступ к свойству мапы
func TestPropertyAccessMap(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	m := map[string]any{
		"key1": "value1",
		"key2": 42,
	}
	ctx.Set("m", m)

	expr, err := engine.CreateExpression(nil, "m.key1")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != "value1" {
		t.Errorf("Expected 'value1', got %v", result)
	}

	expr2, err := engine.CreateExpression(nil, "m.key2")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result2, err := expr2.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	var actual int64
	switch v := result2.(type) {
	case int64:
		actual = v
	case int:
		actual = int64(v)
	case *big.Rat:
		if v.IsInt() {
			actual = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Unexpected result type: %T", result2)
	}

	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}
}

// TestPropertyAccessWithIndex тестирует доступ к свойству через индекс
func TestPropertyAccessWithIndex(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	m := map[string]any{
		"key": "value",
	}
	ctx.Set("m", m)

	expr, err := engine.CreateExpression(nil, "m['key']")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != "value" {
		t.Errorf("Expected 'value', got %v", result)
	}
}

// TestPropertyAccessChaining тестирует цепочку доступа к свойствам
func TestPropertyAccessChaining(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	type Inner struct {
		Value string
	}
	type Middle struct {
		Inner *Inner
	}
	type Outer struct {
		Middle *Middle
	}

	ctx := jexl.NewMapContext()
	obj := &Outer{
		Middle: &Middle{
			Inner: &Inner{
				Value: "test",
			},
		},
	}
	ctx.Set("obj", obj)

	expr, err := engine.CreateExpression(nil, "obj.Middle.Inner.Value")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result != "test" {
		t.Errorf("Expected 'test', got %v", result)
	}
}

