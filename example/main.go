package main

import (
	"fmt"
	"log"

	"github.com/mentatxx/jexl-golang/jexl"
	_ "github.com/mentatxx/jexl-golang/jexl/impl"
)

func main() {
	// Создаём движок
	builder := jexl.NewBuilder()
	engine, err := builder.Build()
	if err != nil {
		log.Fatalf("Failed to build engine: %v", err)
	}

	// Создаём контекст с переменными
	ctx := jexl.NewMapContext()
	ctx.Set("x", 10)
	ctx.Set("y", 20)
	ctx.Set("name", "JEXL")

	// Простое арифметическое выражение
	expr1, err := engine.CreateExpression(nil, "x + y")
	if err != nil {
		log.Fatalf("Failed to create expression: %v", err)
	}

	result1, err := expr1.Evaluate(ctx)
	if err != nil {
		log.Fatalf("Failed to evaluate expression: %v", err)
	}

	fmt.Printf("x + y = %v\n", result1)

	// Выражение со строкой
	expr2, err := engine.CreateExpression(nil, "name + ' для Go'")
	if err != nil {
		log.Fatalf("Failed to create expression: %v", err)
	}

	result2, err := expr2.Evaluate(ctx)
	if err != nil {
		log.Fatalf("Failed to evaluate expression: %v", err)
	}

	fmt.Printf("name + ' для Go' = %v\n", result2)
}
