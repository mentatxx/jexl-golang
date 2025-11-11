package jexl_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestSetLiteralWithOneEntry - тест литерала множества с одним элементом
func TestSetLiteralWithOneEntry(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	sources := []string{"{ 'foo' }", "{ 'foo', }"}
	for _, src := range sources {
		expr, err := engine.CreateExpression(nil, src)
		if err != nil {
			t.Fatalf("Failed to create expression for %s: %v", src, err)
		}

		result, err := expr.Evaluate(ctx)
		if err != nil {
			t.Fatalf("Failed to evaluate expression for %s: %v", src, err)
		}

		// Проверяем, что это множество
		resultType := reflect.TypeOf(result)
		if resultType.Kind() != reflect.Map {
			t.Errorf("Expected map (set), got %v for %s", resultType, src)
		}

		// Проверяем размер
		resultValue := reflect.ValueOf(result)
		if resultValue.Len() != 1 {
			t.Errorf("Expected set size 1, got %d for %s", resultValue.Len(), src)
		}
	}
}

// TestNotEmptySimpleSetLiteral - тест что множество не пустое
func TestNotEmptySimpleSetLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "empty({ 'foo' , 'bar' })")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	if result == true {
		t.Error("Expected false (set is not empty), got true")
	}
}

// TestSetLiteralWithNulls - тест литерала множества с null значениями
func TestSetLiteralWithNulls(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	exprs := []string{
		"{  }",
		"{ 10 }",
		"{ 10 , null }",
		"{ 10 , null , 20}",
		"{ '10' , null }",
		"{ null, '10' , 20 }",
	}

	for _, exprStr := range exprs {
		script, err := engine.CreateScript(nil, nil, exprStr)
		if err != nil {
			t.Fatalf("Failed to create script for %s: %v", exprStr, err)
		}

		result, err := script.Execute(ctx)
		if err != nil {
			t.Fatalf("Failed to execute script for %s: %v", exprStr, err)
		}

		// Проверяем, что это множество (map в Go)
		if result == nil {
			// Пустое множество может быть nil
			if exprStr != "{  }" {
				t.Errorf("Unexpected nil result for %s", exprStr)
			}
			continue
		}
		resultType := reflect.TypeOf(result)
		if resultType == nil {
			t.Errorf("Result type is nil for %s", exprStr)
			continue
		}
		if resultType.Kind() != reflect.Map {
			t.Errorf("Expected map (set), got %v for %s", resultType, exprStr)
		}
	}
}

// TestSetLiteralWithNumbers - тест литерала множества с числами
func TestSetLiteralWithNumbers(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "{ 5.0 , 10 }")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	// Проверяем, что это множество
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Map {
		t.Errorf("Expected map (set), got %v", resultType)
	}

	// Проверяем размер
	resultValue := reflect.ValueOf(result)
	if resultValue.Len() != 2 {
		t.Errorf("Expected set size 2, got %d", resultValue.Len())
	}
}

// TestSetLiteralWithOneEntryScript - тест литерала множества в скрипте
func TestSetLiteralWithOneEntryScript(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	script, err := engine.CreateScript(nil, nil, "{ 'foo' }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Проверяем, что это множество
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Map {
		t.Errorf("Expected map (set), got %v", resultType)
	}

	// Проверяем размер
	resultValue := reflect.ValueOf(result)
	if resultValue.Len() != 1 {
		t.Errorf("Expected set size 1, got %d", resultValue.Len())
	}
}

// TestSetLiteralWithStrings - тест литерала множества со строками
func TestSetLiteralWithStrings(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	sources := []string{"{ 'foo', 'bar' }", "{ 'foo', 'bar', }"}
	for _, src := range sources {
		expr, err := engine.CreateExpression(nil, src)
		if err != nil {
			t.Fatalf("Failed to create expression for %s: %v", src, err)
		}

		result, err := expr.Evaluate(ctx)
		if err != nil {
			t.Fatalf("Failed to evaluate expression for %s: %v", src, err)
		}

		// Проверяем, что это множество
		resultType := reflect.TypeOf(result)
		if resultType.Kind() != reflect.Map {
			t.Errorf("Expected map (set), got %v for %s", resultType, src)
		}

		// Проверяем размер
		resultValue := reflect.ValueOf(result)
		if resultValue.Len() != 2 {
			t.Errorf("Expected set size 2, got %d for %s", resultValue.Len(), src)
		}
	}

	// Тест синтаксической ошибки
	_, err = engine.CreateExpression(nil, "{ , }")
	if err == nil {
		t.Error("Expected parsing error for '{ , }'")
	}
}

// TestSetLiteralWithStringsScript - тест литерала множества со строками в скрипте
func TestSetLiteralWithStringsScript(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	script, err := engine.CreateScript(nil, nil, "{ 'foo' , 'bar' }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Проверяем, что это множество
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Map {
		t.Errorf("Expected map (set), got %v", resultType)
	}

	// Проверяем размер
	resultValue := reflect.ValueOf(result)
	if resultValue.Len() != 2 {
		t.Errorf("Expected set size 2, got %d", resultValue.Len())
	}
}

// TestSizeOfSimpleSetLiteral - тест размера простого литерала множества
func TestSizeOfSimpleSetLiteral(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	expr, err := engine.CreateExpression(nil, "size({ 'foo' , 'bar'})")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate expression: %v", err)
	}

	// Проверяем, что размер равен 2
	var size int64
	switch v := result.(type) {
	case int:
		size = int64(v)
	case int64:
		size = v
	case *big.Rat:
		size = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}
	if size != 2 {
		t.Errorf("Expected size 2, got %d", size)
	}
}

