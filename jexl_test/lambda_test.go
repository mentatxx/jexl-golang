package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestLambda270 - тест из LambdaTest.java test270
func TestLambda270(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "(x, y, z)->{ x + y + z }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	text := base.SourceText()
	script := base.Curry(5, 15)
	if script.SourceText() != text {
		t.Errorf("Expected source text to match, got %q, want %q", script.SourceText(), text)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("s", base)
	script2, err := engine.CreateScript(nil, nil, "return s")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}
	result, err := script2.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
	if resultScript, ok := result.(jexl.Script); ok {
		if resultScript.SourceText() != text {
			t.Errorf("Expected source text to match, got %q, want %q", resultScript.SourceText(), text)
		}
	} else {
		t.Errorf("Expected Script, got %T", result)
	}

	script3, err := engine.CreateScript(nil, nil, "return s.curry(1)")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}
	result2, err := script3.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
	if resultScript2, ok := result2.(jexl.Script); ok {
		if resultScript2.SourceText() != text {
			t.Errorf("Expected source text to match, got %q, want %q", resultScript2.SourceText(), text)
		}
	} else {
		t.Errorf("Expected Script, got %T", result2)
	}
}

// TestLambda271a - тест из LambdaTest.java test271a
func TestLambda271a(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "var base = 1; var x = (a)->{ var y = (b) -> {base + b}; return base + y(a)}; x(40)")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := base.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambda271b - тест из LambdaTest.java test271b
func TestLambda271b(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "var base = 2; var sum = (x, y, z)->{ base + x + y + z }; var y = sum.curry(1); y(2,3)")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := base.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var expected int64 = 8
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambda271c - тест из LambdaTest.java test271c
func TestLambda271c(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "(x, y, z)->{ 2 + x + y + z }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	y := base.Curry(1)
	result, err := y.Execute(nil, 2, 3)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var expected int64 = 8
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambda271d - тест из LambdaTest.java test271d
func TestLambda271d(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "var base = 2; (x, y, z)->base + x + y + z;")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := base.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	lambdaScript, ok := result.(jexl.Script)
	if !ok {
		t.Fatalf("Expected Script, got %T", result)
	}

	y := lambdaScript.Curry(1)
	result2, err := y.Execute(nil, 2, 3)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var expected int64 = 8
	var actual int64
	switch v := result2.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result2, result2)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambda271e - тест из LambdaTest.java test271e
func TestLambda271e(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "var base = 1000; var f = (x, y)->{ var base = x + y + (base?:-1000); base; }; f(100, 20)")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := base.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var expected int64 = 1120
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaCurry1 - тест из LambdaTest.java testCurry1
func TestLambdaCurry1(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "(x, y, z)->{ x + y + z }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	parms := base.UnboundParameters()
	if len(parms) != 3 {
		t.Errorf("Expected 3 unbound parameters, got %d", len(parms))
	}

	script := base.Curry(5)
	parms = script.UnboundParameters()
	if len(parms) != 2 {
		t.Errorf("Expected 2 unbound parameters, got %d", len(parms))
	}

	script = script.Curry(15)
	parms = script.UnboundParameters()
	if len(parms) != 1 {
		t.Errorf("Expected 1 unbound parameter, got %d", len(parms))
	}

	script = script.Curry(22)
	parms = script.UnboundParameters()
	if len(parms) != 0 {
		t.Errorf("Expected 0 unbound parameters, got %d", len(parms))
	}

	result, err := script.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaCurry2 - тест из LambdaTest.java testCurry2
func TestLambdaCurry2(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "(x, y, z)->{ x + y + z }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	script := base.Curry(5, 15)
	parms := script.UnboundParameters()
	if len(parms) != 1 {
		t.Errorf("Expected 1 unbound parameter, got %d", len(parms))
	}

	script = script.Curry(22)
	result, err := script.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaCurry3 - тест из LambdaTest.java testCurry3
func TestLambdaCurry3(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "(x, y, z)->{ x + y + z }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	script := base.Curry(5, 15)
	result, err := script.Execute(nil, 22)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaCurry4 - тест из LambdaTest.java testCurry4
func TestLambdaCurry4(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "(x, y, z)->{ x + y + z }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	script := base.Curry(5)
	result, err := script.Execute(nil, 15, 22)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaCurry5 - тест из LambdaTest.java testCurry5
func TestLambdaCurry5(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	base, err := engine.CreateScript(nil, nil, "var t = x + y + z; return t", "x", "y", "z")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	script := base.Curry(5)
	result, err := script.Execute(nil, 15, 22)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaIdentity - тест из LambdaTest.java testIdentity
func TestLambdaIdentity(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	script, err := engine.CreateScript(nil, nil, "(x)->{ x }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	params := script.Parameters()
	if len(params) != 1 || params[0] != "x" {
		t.Errorf("Expected parameters [x], got %v", params)
	}

	result, err := script.Execute(nil, 42)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambda - тест из LambdaTest.java testLambda
func TestLambda(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	strs := "var s = function(x) { x + x }; s(21)"
	s42, err := engine.CreateScript(nil, nil, strs)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := s42.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}

	strs2 := "var s = function(x, y) { x + y }; s(15, 27)"
	s42_2, err := engine.CreateScript(nil, nil, strs2)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result2, err := s42_2.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var actual2 int64
	switch v := result2.(type) {
	case int:
		actual2 = int64(v)
	case int64:
		actual2 = v
	case *big.Rat:
		actual2 = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result2, result2)
	}

	if actual2 != expected {
		t.Errorf("Expected %d, got %d", expected, actual2)
	}
}

// TestLambdaClosure - тест из LambdaTest.java testLambdaClosure
func TestLambdaClosure(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	strs := "var t = 20; var s = function(x, y) { x + y + t}; s(15, 7)"
	s42, err := engine.CreateScript(nil, nil, strs)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := s42.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}

	strs2 := "var t = 19; var s = function(x, y) { var t = 20; x + y + t}; s(15, 7)"
	s42_2, err := engine.CreateScript(nil, nil, strs2)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result2, err := s42_2.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var actual2 int64
	switch v := result2.(type) {
	case int:
		actual2 = int64(v)
	case int64:
		actual2 = v
	case *big.Rat:
		actual2 = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result2, result2)
	}

	if actual2 != expected {
		t.Errorf("Expected %d, got %d", expected, actual2)
	}
}

// TestLambdaExpr0 - тест из LambdaTest.java testLambdaExpr0
func TestLambdaExpr0(t *testing.T) {
	src := "(x, y) -> x + y"
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	script, err := engine.CreateScript(nil, nil, src)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(nil, 11, 31)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaExpr1 - тест из LambdaTest.java testLambdaExpr1
func TestLambdaExpr1(t *testing.T) {
	src := "x -> x + x"
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	script, err := engine.CreateScript(nil, nil, src)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(nil, 21)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaRecurse - тест из LambdaTest.java testRecurse
func TestLambdaRecurse(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	jc := jexl.NewMapContext()
	script, err := engine.CreateScript(nil, nil, "var fact = (x)->{ if (x <= 1) 1; else x * fact(x - 1) }; fact(5)")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(jc)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var expected int64 = 120
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaRecurse1 - тест из LambdaTest.java testRecurse1
func TestLambdaRecurse1(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	jc := jexl.NewMapContext()
	src := "var fact = (x)-> x <= 1? 1 : x * fact(x - 1);\nfact(5);\n"
	script, err := engine.CreateScript(nil, nil, src)
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(jc)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var expected int64 = 120
	var actual int64
	switch v := result.(type) {
	case int:
		actual = int64(v)
	case int64:
		actual = v
	case *big.Rat:
		actual = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaScriptArguments - тест из LambdaTest.java testScriptArguments
func TestLambdaScriptArguments(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	s, err := engine.CreateScript(nil, nil, " x + x ", "x")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	s42, err := engine.CreateScript(nil, nil, "s(21)", "s")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := s42.Execute(nil, s)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// TestLambdaScriptContext - тест из LambdaTest.java testScriptContext
func TestLambdaScriptContext(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	s, err := engine.CreateScript(nil, nil, "function(x) { x + x }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := s.Execute(nil, 21)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
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
		t.Fatalf("Unexpected result type: %T, value: %v", result, result)
	}

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}

	s42, err := engine.CreateScript(nil, nil, "s(21)")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	ctxt := jexl.NewMapContext()
	ctxt.Set("s", s)
	result2, err := s42.Execute(ctxt)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var actual2 int64
	switch v := result2.(type) {
	case int:
		actual2 = int64(v)
	case int64:
		actual2 = v
	case *big.Rat:
		actual2 = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result2, result2)
	}

	if actual2 != expected {
		t.Errorf("Expected %d, got %d", expected, actual2)
	}

	s42_2, err := engine.CreateScript(nil, nil, "x-> { x + x }")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result3, err := s42_2.Execute(ctxt, 21)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var actual3 int64
	switch v := result3.(type) {
	case int:
		actual3 = int64(v)
	case int64:
		actual3 = v
	case *big.Rat:
		actual3 = v.Num().Int64()
	default:
		t.Fatalf("Unexpected result type: %T, value: %v", result3, result3)
	}

	if actual3 != expected {
		t.Errorf("Expected %d, got %d", expected, actual3)
	}
}

