# JEXL для Go

Порт библиотеки Apache Commons JEXL на язык Go.

## Текущий статус

Портирование находится в процессе. Реализовано:

- ✅ Базовые интерфейсы (Context, Engine, Expression, Script, Builder)
- ✅ MapContext и ObjectContext
- ✅ Парсер с поддержкой основных конструкций JEXL
- ✅ Интерпретатор для базовых операций и control flow
- ✅ Арифметика (базовая)
- ✅ AST узлы (включая условные операторы и циклы)
- ✅ Options, Features, Permissions
- ✅ Uberspect (интроспекция через reflection)
- ✅ Sandbox для безопасности
- ✅ Cache
- ✅ Template Engine (JXLT) - базовая функциональность
- ✅ Обработка ошибок и логирование
- ✅ Условные операторы (if/else, тернарный оператор ?:, Elvis оператор ??)
- ✅ Циклы (for, while, do-while, foreach)
- ✅ Структурные литералы (массивы, мапы, множества)
- ✅ Операторы break, continue, return
- ✅ Блоки кода { }

## Использование

**Важно**: Для включения стандартной реализации движка необходимо добавить blank-import пакета `github.com/mentatxx/jexl-golang/jexl/impl`. Это обеспечит регистрацию движка через `init()`.

```go
package main

import (
    "fmt"
    "github.com/mentatxx/jexl-golang/jexl"
    _ "github.com/mentatxx/jexl-golang/jexl/impl"
)

func main() {
    builder := jexl.NewBuilder()
    engine, err := builder.Build()
    if err != nil {
        panic(err)
    }

    ctx := jexl.NewMapContext()
    ctx.Set("x", 10)
    ctx.Set("y", 20)

    expr, err := engine.CreateExpression(nil, "x + y")
    if err != nil {
        panic(err)
    }

    result, err := expr.Evaluate(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Println("Result:", result)
}
```

## Архитектура

Библиотека использует архитектуру с разделением на публичный API (пакет `jexl`) и внутреннюю реализацию (пакет `jexl/internal`). Это позволяет избежать экспозиции внутренних деталей реализации, но требует явной регистрации движка в текущей версии.

Пакет `jexl/impl` предоставляет стандартную реализацию, регистрирующую `internal`-движок при его импорте.

## Примеры использования

### Базовые выражения
```go
engine, _ := builder.Build()
ctx := jexl.NewMapContext()
ctx.Set("x", 10)
ctx.Set("y", 20)

expr, _ := engine.CreateExpression(nil, "x + y")
result, _ := expr.Evaluate(ctx)
```

### Условные операторы
```go
expr, _ := engine.CreateExpression(nil, "x > 5 ? 100 : 200")
result, _ := expr.Evaluate(ctx)

script, _ := engine.CreateScript(nil, nil, "if (x > 5) { 100 } else { 200 }")
result, _ := script.Execute(ctx)
```

### Циклы
```go
// For loop
script, _ := engine.CreateScript(nil, nil, "for (i = 0; i < 5; i = i + 1) { sum = sum + i }")

// Foreach loop
script, _ := engine.CreateScript(nil, nil, "for (var x : items) { sum = sum + x }")

// While loop
script, _ := engine.CreateScript(nil, nil, "while (x < 5) { x = x + 1 }")
```

### Структурные литералы
```go
// Массивы
expr, _ := engine.CreateExpression(nil, "[1, 2, 3, 4, 5]")

// Мапы
expr, _ := engine.CreateExpression(nil, "{'a': 1, 'b': 2, 'c': 3}")

// Множества
expr, _ := engine.CreateExpression(nil, "{1, 2, 3}")
```

### Template Engine
```go
jxlt, _ := engine.CreateTemplateEngine()
expr, _ := jxlt.CreateExpression("Hello ${name}, value is ${value}")
result, _ := expr.Evaluate(ctx)
```

## TODO

- [ ] Лямбда-функции и замыкания
- [ ] Расширенная поддержка типов и преобразований
- [ ] Улучшенная обработка ошибок с более детальной информацией
- [ ] Портировать основные тесты из Java версии
- [ ] Оптимизация производительности парсера и интерпретатора
