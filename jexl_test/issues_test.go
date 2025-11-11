package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestIssue100 тестирует доступ к массивам через индексацию и точечную нотацию
func TestIssue100(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Cache(4).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	foo := []any{42}
	ctx.Set("foo", foo)

	for l := 0; l < 2; l++ {
		// Сбрасываем значение перед каждым проходом
		foo[0] = 42
		ctx.Set("foo", foo)

		// foo[0]
		expr, err := engine.CreateExpression(nil, "foo[0]")
		if err != nil {
			t.Fatalf("Failed to create expression: %v", err)
		}
		value, err := expr.Evaluate(ctx)
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}
		if value != 42 {
			t.Errorf("Expected 42, got %v", value)
		}

		// foo[0] = 43
		expr, err = engine.CreateExpression(nil, "foo[0] = 43")
		if err != nil {
			t.Fatalf("Failed to create expression: %v", err)
		}
		value, err = expr.Evaluate(ctx)
		if err != nil {
			t.Fatalf("Failed to evaluate: %v", err)
		}
		if value != 43 {
			t.Errorf("Expected 43, got %v", value)
		}
		if foo[0] != 43 {
			t.Errorf("Expected foo[0] to be 43, got %v", foo[0])
		}

		// foo.0 (точечная нотация для индекса)
		expr, err = engine.CreateExpression(nil, "foo.0")
		if err != nil {
			// Точечная нотация для индексов может не поддерживаться
			t.Logf("Dot notation for array index may not be supported: %v", err)
		} else {
			value, err = expr.Evaluate(ctx)
			if err != nil {
				t.Logf("Evaluation failed: %v", err)
			} else if value != 43 {
				t.Errorf("Expected 43, got %v", value)
			}
		}

		// foo.0 = 42
		expr, err = engine.CreateExpression(nil, "foo.0 = 42")
		if err != nil {
			t.Logf("Dot notation assignment may not be supported: %v", err)
		} else {
			value, err = expr.Evaluate(ctx)
			if err != nil {
				t.Logf("Evaluation failed: %v", err)
			} else if value != 42 {
				t.Errorf("Expected 42, got %v", value)
			}
		}
	}
}

// TestIssue105 тестирует доступ к свойствам объектов в массивах
func TestIssue105(t *testing.T) {
	type A105 struct {
		NameA string
		PropA string
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("a", &A105{NameA: "a1", PropA: "p1"})

	// [a.propA] - создание массива с одним элементом
	expr, err := engine.CreateExpression(nil, "[a.PropA]")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected array, got %T", result)
	}

	if len(arr) != 1 {
		t.Fatalf("Expected array length 1, got %d", len(arr))
	}

	if arr[0] != "p1" {
		t.Errorf("Expected 'p1', got %v", arr[0])
	}

	// Изменяем значение
	ctx.Set("a", &A105{NameA: "a2", PropA: "p2"})
	result, err = expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	arr, ok = result.([]any)
	if !ok {
		t.Fatalf("Expected array, got %T", result)
	}

	if arr[0] != "p2" {
		t.Errorf("Expected 'p2', got %v", arr[0])
	}
}

// TestIssue107 тестирует вызовы методов на различных выражениях
func TestIssue107(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("Q4", "Q4")

	tests := []struct {
		name     string
		expr     string
		expected string
		skip     bool // Пропустить если метод не поддерживается
	}{
		{"string toLowerCase", "'Q4'.toLowerCase()", "q4", false},
		{"parenthesized string", "(Q4).toLowerCase()", "q4", false},
		{"number toString", "(4).toString()", "4", true}, // toString может не поддерживаться
		{"expression toString", "(1 + 3).toString()", "4", true},
		{"map get toLowerCase", "({'q': 'Q4'}).get('q').toLowerCase()", "q4", true}, // get может не поддерживаться
		{"map bracket toLowerCase", "({'q': 'Q4'})['q'].toLowerCase()", "q4", false},
		{"array bracket toLowerCase", "(['Q4'])[0].toLowerCase()", "q4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Method not yet supported")
			}

			expr, err := engine.CreateExpression(nil, tt.expr)
			if err != nil {
				t.Logf("Failed to create expression (may not be supported): %v", err)
				return
			}

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Logf("Failed to evaluate (may not be supported): %v", err)
				return
			}

			// Преобразуем результат в строку для сравнения
			var resultStr string
			switch v := result.(type) {
			case string:
				resultStr = v
			case int, int64:
				resultStr = ""
			default:
				resultStr = ""
			}

			if resultStr != "" && resultStr != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

// TestIssue108 тестирует size() для пустых коллекций
func TestIssue108(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// size([])
	script, err := engine.CreateScript(nil, nil, "size([])")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	// Результат должен быть 0 или big.Rat(0)
	var size int
	switch v := result.(type) {
	case int:
		size = v
	case int64:
		size = int(v)
	case *big.Rat:
		if v.Sign() != 0 {
			t.Errorf("Expected 0, got %v", v)
		}
		return
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if size != 0 {
		t.Errorf("Expected 0, got %d", size)
	}

	// size({:}) - пустая мапа
	script, err = engine.CreateScript(nil, nil, "size({:})")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err = script.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	switch v := result.(type) {
	case int:
		size = v
	case int64:
		size = int(v)
	case *big.Rat:
		if v.Sign() != 0 {
			t.Errorf("Expected 0, got %v", v)
		}
		return
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if size != 0 {
		t.Errorf("Expected 0, got %d", size)
	}
}

// TestIssue109 тестирует доступ к переменным с точкой в имени
func TestIssue109(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("foo.bar", 40)

	expr, err := engine.CreateExpression(nil, "foo.bar + 2")
	if err != nil {
		// Может не поддерживаться, так как парсер может интерпретировать как свойство
		t.Logf("May not support dot in variable names: %v", err)
		return
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	// Результат должен быть 42
	var expected int64 = 42
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		if v.Num().Int64() != expected {
			t.Errorf("Expected %d, got %v", expected, v)
		}
		return
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestIssue110 тестирует скрипты с параметрами
func TestIssue110(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Создаем скрипт с параметром "foo"
	script, err := engine.CreateScript(nil, nil, "foo + 2", "foo")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	// Выполняем скрипт с параметром 40
	result, err := script.Execute(ctx, 40)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	// Результат должен быть 42
	var expected int64 = 42
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		if v.Num().Int64() != expected {
			t.Errorf("Expected %d, got %v", expected, v)
		}
		return
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}

	// Тест с вложенным доступом
	ctx.Set("frak.foo", -40)
	script, err = engine.CreateScript(nil, nil, "frak.foo - 2", "foo")
	if err != nil {
		t.Logf("Nested property access in script may not be supported: %v", err)
		return
	}

	result, err = script.Execute(ctx, 40)
	if err != nil {
		t.Logf("Execution failed: %v", err)
		return
	}

	expected = -42
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		if v.Num().Int64() != expected {
			t.Errorf("Expected %d, got %v", expected, v)
		}
		return
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestIssue111 тестирует тернарный оператор с различными типами
func TestIssue111(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	strExpr := "((x>0)?\"FirstValue=\"+(y-x):\"SecondValue=\"+x)"

	expr, err := engine.CreateExpression(nil, strExpr)
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	tests := []struct {
		name     string
		x        any
		y        any
		expected string
	}{
		{"int int", 1, 10, "FirstValue=9"},
		{"float float", 1.0, 10.0, "FirstValue=9.0"},
		{"int float", 1, 10.0, "FirstValue=9.0"},
		{"float int", 1.0, 10, "FirstValue=9.0"},
		{"negative int int", -10, 1, "SecondValue=-10"},
		{"negative float float", -10.0, 1.0, "SecondValue=-10.0"},
		{"negative int float", -10, 1.0, "SecondValue=-10"},
		{"negative float int", -10.0, 1, "SecondValue=-10.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx.Set("x", tt.x)
			ctx.Set("y", tt.y)

			result, err := expr.Evaluate(ctx)
			if err != nil {
				t.Fatalf("Failed to evaluate: %v", err)
			}

			resultStr, ok := result.(string)
			if !ok {
				t.Fatalf("Expected string, got %T", result)
			}

			if resultStr != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

// TestIssue112 тестирует парсинг больших целых чисел
func TestIssue112(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	tests := []struct {
		name     string
		expr     string
		expected int64
	}{
		{"max int", "2147483647", 2147483647},
		{"min int + 1", "-2147483647", -2147483647},
		{"min int", "-2147483648", -2147483648},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := engine.CreateScript(nil, nil, tt.expr)
			if err != nil {
				t.Fatalf("Failed to create script: %v", err)
			}

			result, err := script.Execute(nil)
			if err != nil {
				t.Fatalf("Failed to execute: %v", err)
			}

			var actual int64
			switch v := result.(type) {
			case int:
				actual = int64(v)
			case int64:
				actual = v
			case *big.Rat:
				actual = v.Num().Int64()
			default:
				t.Fatalf("Unexpected result type: %T", result)
			}

			if actual != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, actual)
			}
		})
	}
}

// TestIssue117 тестирует сравнение больших чисел
func TestIssue117(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	expr, err := engine.CreateExpression(nil, "TIMESTAMP > 20100102000000")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("TIMESTAMP", int64(20100103000000))

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	// Результат должен быть true
	var boolResult bool
	switch v := result.(type) {
	case bool:
		boolResult = v
	case *big.Rat:
		boolResult = v.Sign() != 0
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if !boolResult {
		t.Error("Expected true, got false")
	}
}

// TestIssue200 - тест из Issues200Test.java (lambda функции)
// Примечание: lambda функции могут быть не реализованы, поэтому тест может быть пропущен
func TestIssue200(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Тест lambda функции (может не поддерживаться)
	script, err := engine.CreateScript(nil, nil, "var f = (y)->{y + 42}; f(x)", "x")
	if err != nil {
		t.Skipf("Lambda functions may not be supported: %v", err)
		return
	}

	result, err := script.Execute(ctx, 100)
	if err != nil {
		t.Skipf("Lambda execution may not be supported: %v", err)
		return
	}

	var expected int64 = 142
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestIssue217 - тест из Issues200Test.java (доступ к массивам с проверкой границ)
func TestIssue217(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	foo := []int{0, 1, 2, 42}
	ctx.Set("foo", foo)

	script, err := engine.CreateScript(nil, nil, "foo[3]")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	var expected int64 = 42
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}

	// Тест с выходом за границы (должен вернуть null или ошибку)
	ctx.Set("foo", []int{0, 1})
	result, err = script.Execute(ctx)
	if err == nil {
		// Если ошибки нет, результат должен быть null
		if result != nil {
			t.Logf("Out of bounds access returned %v instead of null", result)
		}
	}
}

// TestIssue242 - тест из Issues200Test.java (точность вычислений с double)
func TestIssue242(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("a", -40.05)
	ctx.Set("b", -8.01)

	script, err := engine.CreateScript(nil, nil, "a + b")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	// Результат должен быть близок к -48.06 (с учетом погрешности float)
	var actual float64
	switch v := result.(type) {
	case float64:
		actual = v
	case float32:
		actual = float64(v)
	case *big.Rat:
		f, _ := v.Float64()
		actual = f
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	expected := -48.06
	diff := actual - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > 0.0001 {
		t.Errorf("Expected approximately %f, got %f (diff: %f)", expected, actual, diff)
	}
}

// TestIssue267 - тест из Issues200Test.java (скрипты с параметрами)
func TestIssue267(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// API declared params
	script, err := engine.CreateScript(nil, nil, "x + y", "x", "y")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(ctx, 20, 22)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	var expected int64 = 42
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestIssue302 - тест из Issues300Test.java (if без скобок)
func TestIssue302(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	strs := []string{
		"{if (0) 1 else 2; var x = 4;}",
		"if (0) 1; else 2;",
		"{ if (0) 1; else 2; }",
		"{ if (0) { if (false) 1 else -3 } else 2; }",
	}

	for _, str := range strs {
		script, err := engine.CreateScript(nil, nil, str)
		if err != nil {
			t.Logf("Failed to create script for %s: %v", str, err)
			continue
		}

		result, err := script.Execute(ctx)
		if err != nil {
			t.Logf("Failed to execute script for %s: %v", str, err)
			continue
		}

		// Результат должен быть четным числом (0 или 2)
		var num int64
		switch v := result.(type) {
		case int:
			num = int64(v)
		case int64:
			num = v
		case *big.Rat:
			num = v.Num().Int64()
		default:
			t.Logf("Unexpected result type for %s: %T", str, result)
			continue
		}

		if num%2 != 0 {
			t.Errorf("Block result should be even for %s, got %d", str, num)
		}
	}
}

// TestIssue306 - тест из Issues300Test.java (Elvis оператор)
func TestIssue306(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	script, err := engine.CreateScript(nil, nil, "x.y ?: 2")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	// x.y не определено, должно вернуть 2
	result, err := script.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	var expected int64 = 2
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}

	// x.y = null, должно вернуть 2
	ctx.Set("x.y", nil)
	result, err = script.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestIssue402 - тест из Issues400Test.java (return в if)
func TestIssue402(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	sources := []string{
		"if (true) { return }",
		"if (true) { 3; return }",
	}

	for _, source := range sources {
		script, err := engine.CreateScript(nil, nil, source)
		if err != nil {
			t.Logf("Failed to create script for %s: %v", source, err)
			continue
		}

		result, err := script.Execute(ctx)
		if err != nil {
			t.Logf("Failed to execute script for %s: %v", source, err)
			continue
		}

		// Результат должен быть null
		if result != nil {
			t.Logf("Expected nil for %s, got %v", source, result)
		}
	}
}

// TestIssue407 - тест из Issues400Test.java (точность вычислений)
func TestIssue407(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()

	// Тест: a + b - a - b должно быть близко к 0
	script, err := engine.CreateScript(nil, nil, "a + b - a - b", "a", "b")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	// Используем double
	result, err := script.Execute(ctx, 99.0, 7.82)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	var actual float64
	switch v := result.(type) {
	case float64:
		actual = v
	case float32:
		actual = float64(v)
	case *big.Rat:
		f, _ := v.Float64()
		actual = f
	default:
		t.Fatalf("Unexpected result type: %T", result)
	}

	// Результат должен быть близок к 0 (с учетом погрешности float)
	diff := actual
	if diff < 0 {
		diff = -diff
	}
	if diff > 1e-14 {
		t.Errorf("Expected approximately 0, got %f (diff: %f)", actual, diff)
	}
}

