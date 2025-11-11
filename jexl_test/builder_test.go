package jexl_test

import (
	"testing"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

// TestBuilderFlags - тест флагов Builder
func TestBuilderFlags(t *testing.T) {
	// Тест safe через Options
	builder := jexl.NewBuilder()
	builder.Options().SetSafe(true)
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}
	// Проверяем, что движок создан
	if engine == nil {
		t.Error("Engine should not be nil")
	}

	// Тест strict через Options
	builder = jexl.NewBuilder()
	builder.Options().SetStrict(true)
	engine, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}
	if engine == nil {
		t.Error("Engine should not be nil")
	}

	// Тест silent через Options
	builder = jexl.NewBuilder()
	builder.Options().SetSilent(true)
	engine, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}
	if engine == nil {
		t.Error("Engine should not be nil")
	}
}

// TestBuilderCache - тест настройки кэша
func TestBuilderCache(t *testing.T) {
	builder := jexl.NewBuilder()
	builder = builder.Cache(128)
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}
	if engine == nil {
		t.Error("Engine should not be nil")
	}
}

// TestBuilderValues - тест установки значений
func TestBuilderValues(t *testing.T) {
	builder := jexl.NewBuilder()
	
	// Тест cache threshold
	builder = builder.Cache(32)
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}
	if engine == nil {
		t.Error("Engine should not be nil")
	}
}

// TestBuilderBuild - тест создания движка
func TestBuilderBuild(t *testing.T) {
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}
	if engine == nil {
		t.Error("Engine should not be nil")
	}

	// Проверяем, что движок работает
	ctx := jexl.NewMapContext()
	ctx.Set("x", 10)
	ctx.Set("y", 20)

	expr, err := engine.CreateExpression(nil, "x + y")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result == nil {
		t.Error("Result should not be nil")
	}
}

// TestBuilderChaining - тест цепочки вызовов Builder
func TestBuilderChaining(t *testing.T) {
	builder := jexl.NewBuilder().
		Cache(64)
	builder.Options().SetSafe(true)
	builder.Options().SetStrict(false)
	builder.Options().SetSilent(true)

	engine, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build engine: %v", err)
	}
	if engine == nil {
		t.Error("Engine should not be nil")
	}

	// Проверяем, что движок работает
	ctx := jexl.NewMapContext()
	ctx.Set("x", 10)

	expr, err := engine.CreateExpression(nil, "x + 5")
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	result, err := expr.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate: %v", err)
	}

	if result == nil {
		t.Error("Result should not be nil")
	}
}

