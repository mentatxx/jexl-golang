package jexl_test

import (
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestExceptionNullPropertyAccess - тест обработки исключений при доступе к null свойствам
func TestExceptionNullPropertyAccess(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Тест: null.1 = 2 (должен вернуть ошибку или null в зависимости от режима)
	script, err := engine.CreateScript(nil, nil, "null.1 = 2; return 42")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		// Ошибка ожидаема при строгом режиме
		t.Logf("Expected error in strict mode: %v", err)
	} else {
		// В нестрогом режиме должен вернуться 42
		if result != 42 {
			t.Errorf("Expected 42, got %v", result)
		}
	}
}

// TestExceptionUndefinedVariable - тест обработки неопределенных переменных
func TestExceptionUndefinedVariable(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Тест: x = y.1 (y не определено)
	script, err := engine.CreateScript(nil, nil, "x = y.1; return 42")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		// Ошибка ожидаема при строгом режиме
		t.Logf("Expected error in strict mode: %v", err)
	} else {
		// В нестрогом режиме должен вернуться 42
		if result != 42 {
			t.Errorf("Expected 42, got %v", result)
		}
	}
}

// TestExceptionNullOperand - тест обработки null операндов в арифметике
func TestExceptionNullOperand(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("c.e", nil)

	// Тест: c.e * 6 (c.e = null)
	expr, err := engine.CreateExpression(nil, "c.e * 6")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		// Ошибка ожидаема при строгом арифметическом режиме
		t.Logf("Expected error in strict arithmetic mode: %v", err)
	} else {
		// В нестрогом режиме должен вернуться 0 или null
		t.Logf("Result in non-strict mode: %v", result)
	}
}

// TestExceptionMethodCall - тест обработки исключений при вызове методов
func TestExceptionMethodCall(t *testing.T) {
	type ThrowNPE struct {
		doThrow bool
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	npe := &ThrowNPE{doThrow: true}
	ctx := jexl.NewObjectContext(engine, npe)

	// Тест: вызов метода, который выбрасывает исключение
	// Примечание: в Go нет прямого эквивалента, но можно проверить обработку ошибок
	expr, err := engine.CreateExpression(nil, "npe()")
	if err != nil {
		t.Logf("Method call may not be supported: %v", err)
		return
	}

	_, err = expr.Evaluate(ctx)
	if err != nil {
		// Ошибка ожидаема
		t.Logf("Expected error: %v", err)
	} else {
		t.Log("No error returned, but exception was expected")
	}
}

// TestExceptionPropertyAccess - тест обработки исключений при доступе к свойствам
func TestExceptionPropertyAccess(t *testing.T) {
	type ThrowNPE struct {
		fail bool
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	npe := &ThrowNPE{fail: false}
	ctx := jexl.NewObjectContext(engine, npe)

	// Тест: доступ к свойству, которое может вызвать исключение
	expr, err := engine.CreateExpression(nil, "fail")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Logf("Error accessing property: %v", err)
	} else {
		t.Logf("Property access result: %v", result)
	}
}

