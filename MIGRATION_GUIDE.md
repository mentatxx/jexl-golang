# Руководство по портированию приложений с Java JEXL на Go JEXL

Этот документ описывает процесс портирования приложений, использующих Apache Commons JEXL на Java, на Go-версию библиотеки. Здесь описаны основные различия в API, синтаксисе и использовании.

## Содержание

1. [Обзор различий](#обзор-различий)
2. [Создание движка](#создание-движка)
3. [Работа с контекстом](#работа-с-контекстом)
4. [Создание и выполнение выражений](#создание-и-выполнение-выражений)
5. [Создание и выполнение скриптов](#создание-и-выполнение-скриптов)
6. [Конфигурация движка](#конфигурация-движка)
7. [Обработка ошибок](#обработка-ошибок)
8. [Типы данных и преобразования](#типы-данных-и-преобразования)
9. [Особенности Go](#особенности-go)
10. [Примеры портирования](#примеры-портирования)

## Обзор различий

### Основные архитектурные различия

| Аспект | Java JEXL | Go JEXL |
|--------|-----------|---------|
| **Интерфейсы** | Абстрактные классы и интерфейсы | Интерфейсы Go |
| **Ошибки** | Исключения (unchecked) | Возврат ошибок (error) |
| **Потокобезопасность** | Thread-safe по умолчанию | Thread-safe через sync.RWMutex |
| **Регистрация** | Автоматическая | Требует blank-import `jexl/impl` |
| **Типы** | Object, примитивы | any (interface{}), конкретные типы Go |

### Ключевые отличия в API

1. **Обработка ошибок**: В Java используются исключения (`JexlException`), в Go - возврат ошибок через `error`
2. **Создание движка**: В Go требуется явный blank-import для регистрации реализации
3. **Типизация**: Go использует строгую типизацию, но с `any` для динамических значений
4. **Контекст**: API контекста идентичен, но реализация использует Go-идиомы

## Создание движка

### Java

```java
import org.apache.commons.jexl3.*;

// Простое создание
JexlEngine jexl = new JexlBuilder().create();

// С конфигурацией
JexlEngine jexl = new JexlBuilder()
    .cache(512)
    .strict(true)
    .silent(false)
    .create();
```

### Go

```go
import (
    "github.com/mentatxx/jexl-golang/jexl"
    _ "github.com/mentatxx/jexl-golang/jexl/impl" // Важно: blank-import для регистрации
)

// Простое создание
builder := jexl.NewBuilder()
engine, err := builder.Build()
if err != nil {
    // обработка ошибки
}

// С конфигурацией
builder := jexl.NewBuilder().
    Cache(512).
    Strict(true).
    Silent(false)
engine, err := builder.Build()
if err != nil {
    // обработка ошибки
}
```

**Важно**: В Go версии необходимо добавить blank-import пакета `jexl/impl` для регистрации стандартной реализации движка. Без этого `Build()` вернёт ошибку.

## Работа с контекстом

### Java

```java
// Создание контекста
JexlContext context = new MapContext();
context.set("x", 10);
context.set("y", 20);
context.set("name", "JEXL");

// Использование существующей Map
Map<String, Object> vars = new HashMap<>();
vars.put("x", 10);
JexlContext context = new MapContext(vars);

// Проверка наличия переменной
if (context.has("x")) {
    Object value = context.get("x");
}
```

### Go

```go
// Создание контекста
ctx := jexl.NewMapContext()
ctx.Set("x", 10)
ctx.Set("y", 20)
ctx.Set("name", "JEXL")

// Использование существующей map
vars := map[string]any{
    "x": 10,
    "y": 20,
}
ctx := jexl.NewMapContextWithMap(vars)

// Проверка наличия переменной
if ctx.Has("x") {
    value := ctx.Get("x")
}
```

**Различия**:
- В Go используется `any` вместо `Object`
- Методы начинаются с заглавной буквы (экспортируемые)
- `NewMapContextWithMap` принимает `map[string]any` вместо `Map<String, Object>`

## Создание и выполнение выражений

### Java

```java
// Создание выражения
JexlExpression expr = jexl.createExpression("x + y");

// Выполнение
Object result = expr.evaluate(context);

// С обработкой ошибок
try {
    JexlExpression expr = jexl.createExpression("x + y");
    Object result = expr.evaluate(context);
} catch (JexlException e) {
    // обработка ошибки
}
```

### Go

```go
// Создание выражения
expr, err := engine.CreateExpression(nil, "x + y")
if err != nil {
    // обработка ошибки парсинга
}

// Выполнение
result, err := expr.Evaluate(ctx)
if err != nil {
    // обработка ошибки выполнения
}

// Полный пример
expr, err := engine.CreateExpression(nil, "x + y")
if err != nil {
    log.Fatalf("Failed to create expression: %v", err)
}

result, err := expr.Evaluate(ctx)
if err != nil {
    log.Fatalf("Failed to evaluate expression: %v", err)
}
fmt.Println("Result:", result)
```

**Различия**:
- В Go все операции возвращают `error` вместо выбрасывания исключений
- `CreateExpression` принимает `*Info` (может быть `nil`) и строку
- Результат имеет тип `any`, требует приведения типов при необходимости

## Создание и выполнение скриптов

### Java

```java
// Создание скрипта
JexlScript script = jexl.createScript("for (i = 0; i < 5; i = i + 1) { sum = sum + i }");

// Выполнение
Object result = script.execute(context);

// С параметрами
JexlScript script = jexl.createScript("x + y", "x", "y");
Object result = script.execute(context, 10, 20);
```

### Go

```go
// Создание скрипта
script, err := engine.CreateScript(nil, nil, "for (i = 0; i < 5; i = i + 1) { sum = sum + i }")
if err != nil {
    // обработка ошибки
}

// Выполнение
result, err := script.Execute(ctx)
if err != nil {
    // обработка ошибки
}

// С параметрами
script, err := engine.CreateScript(nil, nil, "x + y", "x", "y")
if err != nil {
    // обработка ошибки
}
result, err := script.Execute(ctx, 10, 20)
if err != nil {
    // обработка ошибки
}
```

**Различия**:
- Сигнатура: `CreateScript(features *Features, info *Info, source string, names ...string)`
- Все параметры могут быть `nil` для значений по умолчанию
- Имена параметров передаются как variadic аргументы

## Конфигурация движка

### Java

```java
JexlEngine jexl = new JexlBuilder()
    .cache(512)                    // Размер кэша
    .strict(true)                  // Строгий режим
    .silent(false)                 // Не silent режим
    .safe(true)                    // Safe navigation
    .debug(true)                   // Отладочный режим
    .cancellable(true)             // Поддержка отмены
    .lexical(true)                 // Лексическая область видимости
    .features(JexlFeatures.ANCHOR) // Особенности
    .permissions(permissions)      // Разрешения
    .sandbox(sandbox)              // Песочница
    .arithmetic(arithmetic)        // Арифметика
    .logger(logger)                // Логгер
    .create();
```

### Go

```go
builder := jexl.NewBuilder().
    Cache(512).                    // Размер кэша
    Strict(true).                  // Строгий режим
    Silent(false).                 // Не silent режим
    Safe(true).                    // Safe navigation
    Debug(true).                   // Отладочный режим
    Cancellable(true).             // Поддержка отмены
    Lexical(true).                 // Лексическая область видимости
    Features(features).            // Особенности
    Permissions(permissions).      // Разрешения
    Sandbox(sandbox).              // Песочница
    Arithmetic(arithmetic).        // Арифметика
    Logger(logger)                 // Логгер

engine, err := builder.Build()
if err != nil {
    // обработка ошибки
}
```

**Различия**:
- В Go используется цепочка методов (method chaining) как в Java
- Все методы возвращают `*Builder` для цепочки
- `Build()` возвращает `(Engine, error)` вместо просто `JexlEngine`
- Некоторые методы имеют другие названия (например, `Cancellable` вместо `cancellable`)

## Обработка ошибок

### Java

```java
try {
    JexlExpression expr = jexl.createExpression("x + y");
    Object result = expr.evaluate(context);
} catch (JexlException e) {
    // JexlException - unchecked exception
    System.err.println("JEXL Error: " + e.getMessage());
    e.printStackTrace();
} catch (Exception e) {
    // Другие исключения
    e.printStackTrace();
}
```

### Go

```go
expr, err := engine.CreateExpression(nil, "x + y")
if err != nil {
    // Проверка типа ошибки
    if parsingErr, ok := err.(*jexl.ParsingError); ok {
        log.Printf("Parsing error at line %d: %v", parsingErr.Line(), parsingErr)
    } else {
        log.Printf("Error creating expression: %v", err)
    }
    return
}

result, err := expr.Evaluate(ctx)
if err != nil {
    // Проверка типа ошибки выполнения
    if methodErr, ok := err.(*jexl.MethodError); ok {
        log.Printf("Method error: %v", methodErr)
    } else if propErr, ok := err.(*jexl.PropertyError); ok {
        log.Printf("Property error: %v", propErr)
    } else {
        log.Printf("Evaluation error: %v", err)
    }
    return
}
```

**Различия**:
- В Go ошибки возвращаются явно, не через исключения
- Доступны типизированные ошибки: `ParsingError`, `MethodError`, `OperatorError`, `PropertyError`
- Проверка типа ошибки через type assertion
- Нет автоматического stack trace, но можно использовать `fmt.Printf("%+v", err)`

## Типы данных и преобразования

### Java

```java
// Автоматическое приведение типов
JexlContext ctx = new MapContext();
ctx.set("x", 10);        // int
ctx.set("y", 20L);       // long
ctx.set("z", 30.5);      // double

JexlExpression expr = jexl.createExpression("x + y + z");
Number result = (Number) expr.evaluate(ctx);
```

### Go

```go
// Типы Go
ctx := jexl.NewMapContext()
ctx.Set("x", 10)         // int
ctx.Set("y", int64(20))  // int64
ctx.Set("z", 30.5)       // float64

expr, _ := engine.CreateExpression(nil, "x + y + z")
result, _ := expr.Evaluate(ctx)

// Приведение типов
if num, ok := result.(*big.Rat); ok {
    // Результат - big.Rat для точных вычислений
    fmt.Println(num)
} else if num, ok := result.(float64); ok {
    // Результат - float64
    fmt.Println(num)
}
```

**Различия**:
- Go использует конкретные числовые типы (`int`, `int64`, `float64`)
- Для точных вычислений используется `big.Rat` (аналог `BigDecimal`)
- Требуется явное приведение типов через type assertion
- `any` используется для динамических значений

## Особенности Go

### 1. Blank Import для регистрации

В Go версии необходимо явно импортировать пакет реализации:

```go
import (
    "github.com/mentatxx/jexl-golang/jexl"
    _ "github.com/mentatxx/jexl-golang/jexl/impl" // Blank import
)
```

### 2. Потокобезопасность

`MapContext` в Go версии потокобезопасен благодаря использованию `sync.RWMutex`:

```go
// Можно безопасно использовать из разных горутин
ctx := jexl.NewMapContext()
go func() {
    ctx.Set("x", 10)
}()
go func() {
    value := ctx.Get("x")
}()
```

### 3. Интерфейсы вместо классов

В Go используются интерфейсы, что позволяет легко создавать свои реализации:

```go
type MyContext struct{}

func (m *MyContext) Get(name string) any {
    // реализация
}

func (m *MyContext) Has(name string) bool {
    // реализация
}

func (m *MyContext) Set(name string, value any) {
    // реализация
}

// Использование
ctx := &MyContext{}
expr.Evaluate(ctx)
```

### 4. Отсутствие перегрузки методов

В Go нет перегрузки методов, поэтому используются разные имена:

```go
// Вместо перегрузки createScript(String) и createScript(String, String...)
script, err := engine.CreateScript(nil, nil, "code")
script, err := engine.CreateScript(nil, nil, "code", "param1", "param2")
```

## Примеры портирования

### Пример 1: Простое арифметическое выражение

**Java:**
```java
JexlEngine jexl = new JexlBuilder().create();
JexlContext ctx = new MapContext();
ctx.set("x", 10);
ctx.set("y", 20);

JexlExpression expr = jexl.createExpression("x + y");
Object result = expr.evaluate(ctx);
System.out.println("Result: " + result);
```

**Go:**
```go
import (
    "fmt"
    "github.com/mentatxx/jexl-golang/jexl"
    _ "github.com/mentatxx/jexl-golang/jexl/impl"
)

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
fmt.Printf("Result: %v\n", result)
```

### Пример 2: Скрипт с циклом

**Java:**
```java
JexlEngine jexl = new JexlBuilder().create();
JexlContext ctx = new MapContext();
ctx.set("sum", 0);

JexlScript script = jexl.createScript(
    "for (i = 0; i < 10; i = i + 1) { sum = sum + i }"
);
Object result = script.execute(ctx);
```

**Go:**
```go
builder := jexl.NewBuilder()
engine, err := builder.Build()
if err != nil {
    panic(err)
}

ctx := jexl.NewMapContext()
ctx.Set("sum", 0)

script, err := engine.CreateScript(nil, nil, 
    "for (i = 0; i < 10; i = i + 1) { sum = sum + i }")
if err != nil {
    panic(err)
}

result, err := script.Execute(ctx)
if err != nil {
    panic(err)
}
```

### Пример 3: Работа с объектами

**Java:**
```java
public class Person {
    private String name;
    private int age;
    
    public Person(String name, int age) {
        this.name = name;
        this.age = age;
    }
    
    public String getName() { return name; }
    public int getAge() { return age; }
}

JexlEngine jexl = new JexlBuilder().create();
JexlContext ctx = new MapContext();
ctx.set("person", new Person("John", 30));

JexlExpression expr = jexl.createExpression("person.name + ' is ' + person.age");
String result = (String) expr.evaluate(ctx);
```

**Go:**
```go
type Person struct {
    Name string
    Age  int
}

builder := jexl.NewBuilder()
engine, err := builder.Build()
if err != nil {
    panic(err)
}

ctx := jexl.NewMapContext()
ctx.Set("person", &Person{Name: "John", Age: 30})

expr, err := engine.CreateExpression(nil, "person.Name + ' is ' + person.Age")
if err != nil {
    panic(err)
}

result, err := expr.Evaluate(ctx)
if err != nil {
    panic(err)
}
// result будет строкой
```

**Примечание**: В Go поля структуры должны быть экспортируемыми (с заглавной буквы) для доступа через reflection.

### Пример 4: Конфигурация с кэшем и опциями

**Java:**
```java
JexlEngine jexl = new JexlBuilder()
    .cache(512)
    .strict(true)
    .silent(false)
    .debug(true)
    .create();
```

**Go:**
```go
builder := jexl.NewBuilder().
    Cache(512).
    Strict(true).
    Silent(false).
    Debug(true)

engine, err := builder.Build()
if err != nil {
    panic(err)
}
```

### Пример 5: Обработка ошибок

**Java:**
```java
try {
    JexlExpression expr = jexl.createExpression("x + y");
    Object result = expr.evaluate(context);
} catch (JexlException e) {
    logger.error("JEXL error: " + e.getMessage(), e);
    // обработка
}
```

**Go:**
```go
expr, err := engine.CreateExpression(nil, "x + y")
if err != nil {
    if parsingErr, ok := err.(*jexl.ParsingError); ok {
        log.Printf("Parsing error: %v", parsingErr)
    } else {
        log.Printf("Error: %v", err)
    }
    return
}

result, err := expr.Evaluate(ctx)
if err != nil {
    log.Printf("Evaluation error: %v", err)
    return
}
```

## Несовместимости и ограничения

### Что не поддерживается в Go версии

1. **Namespace и Import pragma**: Директивы `#pragma namespace` и `#pragma import` пока не реализованы
2. **Annotation обработка**: Пользовательские аннотации не поддерживаются
3. **JSR-223 Scripting API**: Не применимо к Go
4. **ClassLoader**: В Go нет аналога Java ClassLoader, используется другой механизм

### Ограничения Go

1. **Создание объектов по имени класса**: В Go нет прямого способа создать экземпляр по строковому имени типа (как `Class.forName()` в Java). Это влияет на оператор `new()` в скриптах.
2. **Reflection ограничения**: Go reflection имеет некоторые ограничения по сравнению с Java reflection.

## Рекомендации по портированию

1. **Начните с простых выражений**: Портируйте сначала базовые арифметические выражения и простые скрипты
2. **Обрабатывайте ошибки**: В Go все ошибки возвращаются явно, не забывайте их проверять
3. **Тестируйте типы**: Используйте type assertion для приведения типов результатов
4. **Используйте blank-import**: Не забудьте добавить `_ "github.com/mentatxx/jexl-golang/jexl/impl"`
5. **Проверяйте совместимость**: Сверяйте поддерживаемые функции с `PORTING_STATUS.md`

## Дополнительные ресурсы

- [PORTING_STATUS.md](./PORTING_STATUS.md) - Статус портирования функций
- [README.md](./README.md) - Общая документация
- [TEST_COVERAGE.md](./TEST_COVERAGE.md) - Покрытие тестами

## Заключение

Портирование с Java JEXL на Go JEXL в целом прямолинейно благодаря сохранению структуры API. Основные изменения связаны с:
- Обработкой ошибок (исключения → возврат ошибок)
- Типизацией (Object → any)
- Необходимостью blank-import для регистрации
- Синтаксисом Go вместо Java

Большинство базовых функций работают идентично, что упрощает миграцию существующего кода.

