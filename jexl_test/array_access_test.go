package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestArrayAccessBasic - тест простого доступа к массивам
func TestArrayAccessBasic(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Тест доступа к списку
	list := []interface{}{1, 2, 3}
	ctx.Set("list", list)

	expr, err := engine.CreateExpression(nil, "list[1]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != 2 {
		t.Errorf("Expected 2, got %v", result)
	}

	// Тест доступа с выражением
	expr, err = engine.CreateExpression(nil, "list[1+1]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != 3 {
		t.Errorf("Expected 3, got %v", result)
	}

	// Тест доступа с переменной
	ctx.Set("loc", 1)
	expr, err = engine.CreateExpression(nil, "list[loc+1]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != 3 {
		t.Errorf("Expected 3, got %v", result)
	}

	// Тест доступа к массиву строк
	array := []string{"hello", "there"}
	ctx.Set("array", array)

	expr, err = engine.CreateExpression(nil, "array[0]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}

	// Тест доступа к мапе
	m := map[string]string{
		"foo": "bar",
	}
	ctx.Set("map", m)
	ctx.Set("key", "foo")

	expr, err = engine.CreateExpression(nil, `map["foo"]`)
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != "bar" {
		t.Errorf("Expected 'bar', got %v", result)
	}

	expr, err = engine.CreateExpression(nil, "map[key]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != "bar" {
		t.Errorf("Expected 'bar', got %v", result)
	}
}

// TestArrayArray - тест доступа к многомерным массивам
func TestArrayArray(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Создаем массив массивов
	foo := make([]interface{}, 3)
	foo[0] = foo // самореференция
	foo[1] = 42
	foo[2] = "fourty-two"

	ctx.Set("foo", foo)
	ctx.Set("zero", 0)
	ctx.Set("one", 1)
	ctx.Set("two", 2)

	// Тест foo[0]
	expr, err := engine.CreateExpression(nil, "foo[0]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	// Проверяем, что результат - это слайс
	resultSlice, ok := result.([]interface{})
	if !ok {
		t.Errorf("Expected slice, got %T", result)
	} else if len(resultSlice) != len(foo) {
		t.Errorf("Expected slice length %d, got %d", len(foo), len(resultSlice))
	}

	// Тест foo[1]
	expr, err = engine.CreateExpression(nil, "foo[1]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != 42 {
		t.Errorf("Expected 42, got %v", result)
	}

	// Тест foo[0][1]
	expr, err = engine.CreateExpression(nil, "foo[0][1]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != 42 {
		t.Errorf("Expected 42, got %v", result)
	}

	// Тест присваивания foo[0][1] = 43
	script, err := engine.CreateScript(nil, nil, "foo[0][1] = 43")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Результат может быть *big.Rat или int, проверяем значение
	var resultInt int
	switch v := result.(type) {
	case int:
		resultInt = v
	case int64:
		resultInt = int(v)
	case *big.Rat:
		if !v.IsInt() {
			t.Errorf("Expected integer result, got %v", result)
		} else {
			resultInt = int(v.Num().Int64())
		}
	default:
		t.Errorf("Unexpected result type: %T, value: %v", result, result)
	}

	if resultInt != 43 {
		t.Errorf("Expected 43, got %d", resultInt)
	}

	// Проверяем, что значение изменилось
	expr, err = engine.CreateExpression(nil, "foo[0][1]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	// Результат может быть *big.Rat или int, проверяем значение
	resultInt = 0
	switch v := result.(type) {
	case int:
		resultInt = v
	case int64:
		resultInt = int(v)
	case *big.Rat:
		if !v.IsInt() {
			t.Errorf("Expected integer result, got %v", result)
		} else {
			resultInt = int(v.Num().Int64())
		}
	default:
		t.Errorf("Unexpected result type: %T, value: %v", result, result)
	}

	if resultInt != 43 {
		t.Errorf("Expected 43, got %d", resultInt)
	}
}

// TestDoubleArrays - тест доступа к двумерным массивам
func TestDoubleArrays(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Создаем двумерный массив
	foo := make([][]interface{}, 2)
	foo[0] = make([]interface{}, 2)
	foo[0][0] = "one"
	foo[0][1] = "two"

	ctx.Set("foo", foo)

	// Тест foo[0][1]
	expr, err := engine.CreateExpression(nil, "foo[0][1]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != "two" {
		t.Errorf("Expected 'two', got %v", result)
	}

	// Тест присваивания
	script, err := engine.CreateScript(nil, nil, "foo[0][1] = 'three'")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result != "three" {
		t.Errorf("Expected 'three', got %v", result)
	}

	// Проверяем, что значение изменилось
	expr, err = engine.CreateExpression(nil, "foo[0][1]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != "three" {
		t.Errorf("Expected 'three', got %v", result)
	}
}

// TestDoubleMaps - тест доступа к вложенным мапам
func TestDoubleMaps(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Создаем вложенные мапы
	foo := make(map[interface{}]interface{})
	foo0 := make(map[interface{}]interface{})
	foo0[0] = "one"
	foo0[1] = "two"
	foo0["3.0"] = "three"
	foo[0] = foo0

	ctx.Set("foo", foo)

	// Тест foo[0][1]
	expr, err := engine.CreateExpression(nil, "foo[0][1]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != "two" {
		t.Errorf("Expected 'two', got %v", result)
	}

	// Тест присваивания
	script, err := engine.CreateScript(nil, nil, "foo[0][1] = 'three'")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result != "three" {
		t.Errorf("Expected 'three', got %v", result)
	}

	// Проверяем доступ к строковому ключу
	expr, err = engine.CreateExpression(nil, `foo[0]["3.0"]`)
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result != "three" {
		t.Errorf("Expected 'three', got %v", result)
	}
}

