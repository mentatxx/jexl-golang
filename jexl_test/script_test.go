package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// Tester - тестовый класс для testScriptUpdatesContext
type Tester struct {
	code string
}

func (t *Tester) GetCode() string {
	return t.code
}

func (t *Tester) SetCode(c string) {
	t.code = c
}

// TestScriptUpdatesContext тестирует, что скрипт обновляет контекст
func TestScriptUpdatesContext(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	jexlCode := "resultat.setCode('OK')"
	expr, err := engine.CreateExpression(nil, jexlCode)
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	script, err := engine.CreateScript(nil, nil, jexlCode)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	resultatJexl := &Tester{}
	ctx := jexl.NewMapContext()
	ctx.Set("resultat", resultatJexl)

	resultatJexl.SetCode("")
	_, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}
	if resultatJexl.GetCode() != "OK" {
		t.Errorf("Expected 'OK', got '%s'", resultatJexl.GetCode())
	}

	resultatJexl.SetCode("")
	_, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
	if resultatJexl.GetCode() != "OK" {
		t.Errorf("Expected 'OK', got '%s'", resultatJexl.GetCode())
	}
}

// TestSimpleScript тестирует создание скрипта из строки
func TestSimpleScript(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	code := "while (x < 10) x = x + 1;"
	script, err := engine.CreateScript(nil, nil, code)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", int64(1))

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

	if script.SourceText() != code {
		t.Errorf("Expected source text '%s', got '%s'", code, script.SourceText())
	}
}

// TestSpacesScript тестирует создание скрипта из пробелов
func TestSpacesScript(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	code := " "
	script, err := engine.CreateScript(nil, nil, code)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	if script == nil {
		t.Error("Script should not be nil")
	}
}

// TestScriptWithParameters тестирует скрипт с параметрами
func TestScriptWithParameters(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Создаем скрипт с параметрами x и y
	script, err := engine.CreateScript(nil, nil, "x + y", "x", "y")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	ctx := jexl.NewMapContext()
	result, err := script.Execute(ctx, int64(13), int64(29))
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

	// 13 + 29 = 42
	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}

	params := script.Parameters()
	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}
	if len(params) > 0 && params[0] != "x" {
		t.Errorf("Expected first parameter 'x', got '%s'", params[0])
	}
	if len(params) > 1 && params[1] != "y" {
		t.Errorf("Expected second parameter 'y', got '%s'", params[1])
	}
}

// TestScriptCurry тестирует каррирование скрипта
func TestScriptCurry(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Создаем скрипт с параметрами x и y
	script, err := engine.CreateScript(nil, nil, "x + y", "x", "y")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	// Каррируем первый аргумент
	curried := script.Curry(int64(10))

	ctx := jexl.NewMapContext()
	result, err := curried.Execute(ctx, int64(32))
	if err != nil {
		t.Fatalf("Failed to execute curried script: %v", err)
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

	// 10 + 32 = 42
	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}
}

// TestScriptWithSingleParameter тестирует скрипт с одним параметром
func TestScriptWithSingleParameter(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	script, err := engine.CreateScript(nil, nil, "x * 2", "x")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	ctx := jexl.NewMapContext()
	result, err := script.Execute(ctx, int64(21))
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

	// 21 * 2 = 42
	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}
}

// TestScriptWithoutParameters тестирует скрипт без параметров
func TestScriptWithoutParameters(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", int64(10))
	ctx.Set("y", int64(32))

	script, err := engine.CreateScript(nil, nil, "x + y")
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

	// 10 + 32 = 42
	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}

	params := script.Parameters()
	if len(params) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(params))
	}
}

// TestScriptLocalVariables тестирует локальные переменные в скрипте
func TestScriptLocalVariables(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	script, err := engine.CreateScript(nil, nil, "var x = 21; var y = 21; x + y")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	ctx := jexl.NewMapContext()
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

	// 21 + 21 = 42
	if actual != 42 {
		t.Errorf("Expected 42, got %d", actual)
	}
}

