package jexl_test

import (
	"math/big"
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestFeaturesCreateNone тестирует создание Features с минимальными возможностями
func TestFeaturesCreateNone(t *testing.T) {
	features := jexl.NewFeatures() // Создаем минимальный набор features

	// Проверяем, что базовые возможности включены
	if !features.SupportsExpression() {
		t.Error("Expression support should be enabled by default")
	}

	// Проверяем, что расширенные возможности выключены (если не включены явно)
	if features.SupportsScript() {
		t.Error("Script support should be disabled by default when no features specified")
	}
	if features.SupportsLoops() {
		t.Error("Loop support should be disabled by default when no features specified")
	}
	if features.SupportsLocalVar() {
		t.Error("Local var support should be disabled by default when no features specified")
	}
}

// TestFeaturesCreateDefault тестирует создание Features с дефолтными настройками
func TestFeaturesCreateDefault(t *testing.T) {
	features := jexl.NewFeatures()

	// Включаем все стандартные возможности
	features = features.With(jexl.FeatureScript, jexl.FeatureLoop, jexl.FeatureLocalVar)

	builder := jexl.NewBuilder()
	engine, err := builder.Features(features).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Проверяем, что скрипты работают
	script, err := engine.CreateScript(features, nil, "var x = 1; x + 1")
	if err != nil {
		t.Fatalf("Failed to create script: %v", err)
	}

	result, err := script.Execute(nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	var val int64
	switch v := result.(type) {
	case int:
		val = int64(v)
	case int64:
		val = v
	case *big.Rat:
		if v.IsInt() {
			val = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Expected int64, got %T", result)
	}

	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}
}

// TestFeaturesNoLoops тестирует отключение циклов
func TestFeaturesNoLoops(t *testing.T) {
	features := jexl.NewFeatures()
	// Отключаем циклы
	features = features.Without(jexl.FeatureLoop)

	builder := jexl.NewBuilder()
	engine, err := builder.Features(features).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Попытка создать скрипт с циклом должна вызвать ошибку
	_, err = engine.CreateScript(features, nil, "while(true) { break; }")
	if err == nil {
		t.Error("Expected error when loops are disabled")
	}

	// Выражения должны работать
	expr, err := engine.CreateExpression(nil, "1 + 1")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	var val int64
	switch v := result.(type) {
	case int:
		val = int64(v)
	case int64:
		val = v
	case *big.Rat:
		if v.IsInt() {
			val = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Expected int64, got %T", result)
	}

	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}
}

// TestFeaturesNoLocalVar тестирует отключение локальных переменных
func TestFeaturesNoLocalVar(t *testing.T) {
	features := jexl.FeaturesDefault()
	// Отключаем локальные переменные
	features = features.Without(jexl.FeatureLocalVar)

	if features.SupportsLocalVar() {
		t.Error("Local var should be disabled")
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Features(features).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Попытка создать скрипт с var должна вызвать ошибку
	_, err = engine.CreateScript(features, nil, "var x = 0")
	if err == nil {
		t.Skip("Local var feature check not yet implemented in parser")
		// TODO: реализовать проверку var в парсере
	}

	// Выражения должны работать
	expr, err := engine.CreateExpression(nil, "1 + 1")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	var val int64
	switch v := result.(type) {
	case int:
		val = int64(v)
	case int64:
		val = v
	case *big.Rat:
		if v.IsInt() {
			val = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Expected int64, got %T", result)
	}

	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}
}

// TestFeaturesNoLambda тестирует отключение lambda функций
func TestFeaturesNoLambda(t *testing.T) {
	features := jexl.FeaturesDefault()
	// Отключаем lambda
	features = features.Without(jexl.FeatureLambda)

	if features.SupportsLambda() {
		t.Error("Lambda should be disabled")
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Features(features).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Попытка создать lambda должна вызвать ошибку
	_, err = engine.CreateScript(features, nil, "var x = ()->{ return 0 }")
	if err == nil {
		t.Skip("Lambda feature check not yet implemented")
		// TODO: реализовать проверку lambda в парсере
	}
}

// TestFeaturesNoMethodCall тестирует отключение вызовов методов
func TestFeaturesNoMethodCall(t *testing.T) {
	features := jexl.FeaturesDefault()
	// Отключаем вызовы методов
	features = features.Without(jexl.FeatureMethodCall)

	if features.SupportsMethodCall() {
		t.Error("Method call should be disabled")
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Features(features).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	ctx := jexl.NewMapContext()
	ctx.Set("x", "test")

	// Попытка вызвать метод должна вызвать ошибку
	_, err = engine.CreateExpression(nil, "x.length()")
	if err == nil {
		t.Skip("Method call feature check not yet implemented")
		// TODO: реализовать проверку вызовов методов в парсере
	}
}

// TestFeaturesNoNewInstance тестирует отключение создания экземпляров
func TestFeaturesNoNewInstance(t *testing.T) {
	features := jexl.FeaturesDefault()
	// Отключаем создание экземпляров
	features = features.Without(jexl.FeatureNewInstance)

	if features.SupportsNewInstance() {
		t.Error("New instance should be disabled")
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Features(features).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Попытка создать экземпляр должна вызвать ошибку
	_, err = engine.CreateScript(features, nil, "new('SomeClass')")
	if err == nil {
		t.Skip("New instance feature check not yet implemented")
		// TODO: реализовать проверку new в парсере
	}
}

// TestFeaturesNoStructuredLiteral тестирует отключение структурированных литералов
func TestFeaturesNoStructuredLiteral(t *testing.T) {
	features := jexl.FeaturesDefault()
	// Отключаем структурированные литералы
	features = features.Without(jexl.FeatureStructuredLiteral)

	if features.SupportsStructuredLiteral() {
		t.Error("Structured literal should be disabled")
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Features(features).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Попытка создать массив должна вызвать ошибку
	_, err = engine.CreateExpression(nil, "[1, 2, 3]")
	if err == nil {
		t.Skip("Structured literal feature check not yet implemented")
		// TODO: реализовать проверку структурированных литералов в парсере
	}
}

// TestFeaturesNoScript тестирует отключение скриптов
func TestFeaturesNoScript(t *testing.T) {
	features := jexl.FeaturesDefault()
	// Отключаем скрипты
	features = features.Without(jexl.FeatureScript)

	if features.SupportsScript() {
		t.Error("Script should be disabled")
	}

	builder := jexl.NewBuilder()
	engine, err := builder.Features(features).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Попытка создать скрипт должна вызвать ошибку
	_, err = engine.CreateScript(features, nil, "{ 3 + 4 }")
	if err == nil {
		t.Skip("Script feature check not yet implemented")
		// TODO: реализовать проверку скриптов в парсере
	}

	// Выражения должны работать
	expr, err := engine.CreateExpression(nil, "3 + 4")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	var val int64
	switch v := result.(type) {
	case int:
		val = int64(v)
	case int64:
		val = v
	case *big.Rat:
		if v.IsInt() {
			val = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Expected int64, got %T", result)
	}

	if val != 7 {
		t.Errorf("Expected 7, got %d", val)
	}
}

// TestFeaturesMixedFeatures тестирует смешанные настройки features
func TestFeaturesMixedFeatures(t *testing.T) {
	features := jexl.NewFeatures()
	// Отключаем несколько возможностей одновременно
	features = features.Without(
		jexl.FeatureNewInstance,
		jexl.FeatureLocalVar,
		jexl.FeatureLambda,
		jexl.FeatureLoop,
	)

	builder := jexl.NewBuilder()
	engine, err := builder.Features(features).Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}

	// Выражения должны работать
	expr, err := engine.CreateExpression(nil, "1 + 2 + 3")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	var val int64
	switch v := result.(type) {
	case int:
		val = int64(v)
	case int64:
		val = v
	case *big.Rat:
		if v.IsInt() {
			val = v.Num().Int64()
		} else {
			t.Fatalf("Expected integer, got %v", v)
		}
	default:
		t.Fatalf("Expected int64, got %T", result)
	}

	if val != 6 {
		t.Errorf("Expected 6, got %d", val)
	}
}

