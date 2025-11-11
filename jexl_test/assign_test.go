package jexl_test

import (
	"math/big"
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

	if resultInt != 10 {
		t.Errorf("Expected 10, got %d", resultInt)
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

	if resultInt != 1010 {
		t.Errorf("Expected 1010, got %d", resultInt)
	}

	// Проверяем значение в массиве
	var arrayVal int
	switch v := array[1].(type) {
	case int:
		arrayVal = v
	case int64:
		arrayVal = int(v)
	case *big.Rat:
		if !v.IsInt() {
			t.Errorf("Expected integer in array, got %v", array[1])
		} else {
			arrayVal = int(v.Num().Int64())
		}
	default:
		t.Errorf("Unexpected array value type: %T, value: %v", array[1], array[1])
	}

	if arrayVal != 1010 {
		t.Errorf("Expected array[1] to be 1010, got %d", arrayVal)
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

// TestAntish тестирует присваивание через точку (ant-style)
func TestAntish(t *testing.T) {
	type Froboz struct {
		Value int
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	froboz := &Froboz{Value: 0}
	ctx.Set("froboz", froboz)

	assignExpr, err := engine.CreateExpression(nil, "froboz.Value = 10")
	if err != nil {
		t.Fatalf("Failed to create assign expression: %v", err)
	}

	checkExpr, err := engine.CreateExpression(nil, "froboz.Value")
	if err != nil {
		t.Fatalf("Failed to create check expression: %v", err)
	}

	result, err := assignExpr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate assign: %v", err)
	}

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

	if resultInt != 10 {
		t.Errorf("Expected 10, got %d", resultInt)
	}

	checkResult, err := checkExpr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate check: %v", err)
	}

	var checkInt int
	switch v := checkResult.(type) {
	case int:
		checkInt = v
	case int64:
		checkInt = int(v)
	case *big.Rat:
		if !v.IsInt() {
			t.Errorf("Expected integer result, got %v", checkResult)
		} else {
			checkInt = int(v.Num().Int64())
		}
	default:
		t.Errorf("Unexpected check result type: %T, value: %v", checkResult, checkResult)
	}

	if checkInt != 10 {
		t.Errorf("Expected check to return 10, got %d", checkInt)
	}
}

// FrobozBeanish - структура для теста TestBeanish
type FrobozBeanish struct {
	value int
}

// GetValue возвращает значение
func (f *FrobozBeanish) GetValue() int {
	return f.value
}

// SetValue устанавливает значение
func (f *FrobozBeanish) SetValue(v int) {
	f.value = v
}

// TestBeanish тестирует присваивание через getter/setter
func TestBeanish(t *testing.T) {

	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	froboz := &FrobozBeanish{value: -169}
	ctx.Set("froboz", froboz)

	assignExpr, err := engine.CreateExpression(nil, "froboz.value = 10")
	if err != nil {
		t.Fatalf("Failed to create assign expression: %v", err)
	}

	checkExpr, err := engine.CreateExpression(nil, "froboz.value")
	if err != nil {
		t.Fatalf("Failed to create check expression: %v", err)
	}

	result, err := assignExpr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate assign: %v", err)
	}

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

	if resultInt != 10 {
		t.Errorf("Expected 10, got %d", resultInt)
	}

	checkResult, err := checkExpr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate check: %v", err)
	}

	var checkInt int
	switch v := checkResult.(type) {
	case int:
		checkInt = v
	case int64:
		checkInt = int(v)
	case *big.Rat:
		if !v.IsInt() {
			t.Errorf("Expected integer result, got %v", checkResult)
		} else {
			checkInt = int(v.Num().Int64())
		}
	default:
		t.Errorf("Unexpected check result type: %T, value: %v", checkResult, checkResult)
	}

	if checkInt != 10 {
		t.Errorf("Expected check to return 10, got %d", checkInt)
	}
}

// TestArrayAssignmentIndex тестирует присваивание через индекс массива
func TestArrayAssignmentIndex(t *testing.T) {
	type Froboz struct {
		data map[string]any
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	froboz := &Froboz{data: make(map[string]any)}
	ctx.Set("froboz", froboz)

	// Используем мапу для хранения данных
	data := make(map[string]any)
	ctx.Set("frobozData", data)

	assignExpr, err := engine.CreateExpression(nil, "frobozData[\"value\"] = 10")
	if err != nil {
		t.Fatalf("Failed to create assign expression: %v", err)
	}

	checkExpr, err := engine.CreateExpression(nil, "frobozData[\"value\"]")
	if err != nil {
		t.Fatalf("Failed to create check expression: %v", err)
	}

	result, err := assignExpr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate assign: %v", err)
	}

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

	if resultInt != 10 {
		t.Errorf("Expected 10, got %d", resultInt)
	}

	checkResult, err := checkExpr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate check: %v", err)
	}

	var checkInt int
	switch v := checkResult.(type) {
	case int:
		checkInt = v
	case int64:
		checkInt = int(v)
	case *big.Rat:
		if !v.IsInt() {
			t.Errorf("Expected integer result, got %v", checkResult)
		} else {
			checkInt = int(v.Num().Int64())
		}
	default:
		t.Errorf("Unexpected check result type: %T, value: %v", checkResult, checkResult)
	}

	if checkInt != 10 {
		t.Errorf("Expected check to return 10, got %d", checkInt)
	}
}

// TestMiniAssignment тестирует простое присваивание переменной
func TestMiniAssignment(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "quux = 10")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

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

	if resultInt != 10 {
		t.Errorf("Expected 10, got %d", resultInt)
	}
}

