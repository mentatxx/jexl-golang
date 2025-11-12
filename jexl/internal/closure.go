package internal

import (
	"github.com/mentatxx/jexl-golang/jexl"
)

// getBaseContext извлекает базовый контекст из argumentContext, если это возможно
func getBaseContext(ctx jexl.Context) jexl.Context {
	if ctx == nil {
		return nil
	}
	// Проверяем, является ли контекст argumentContext
	// argumentContext находится в том же пакете, поэтому можем использовать type assertion
	// Но нам нужно проверить через интерфейс или рефлексию
	// Проще всего - использовать интерфейс для получения базового контекста
	// Но для этого нужно добавить метод в Context или использовать другой подход
	// Временно просто возвращаем контекст как есть
	return ctx
}

// closure реализует jexl.Script для lambda функций с захватом контекста.
// Порт org.apache.commons.jexl3.internal.Closure.
type closure struct {
	*script
	capturedContext jexl.Context
}

// NewClosure создаёт новый closure из lambda узла.
func NewClosure(engine jexl.Engine, lambda *jexl.LambdaNode, capturedContext jexl.Context) jexl.Script {
	// Создаём ScriptNode из lambda для выполнения
	scriptNode := jexl.NewScriptNode(nil, lambda.SourceText(), nil)
	scriptNode.AddChild(lambda.Body())

	// Получаем имена параметров
	params := make([]string, 0, len(lambda.Parameters()))
	for _, param := range lambda.Parameters() {
		params = append(params, param.Name())
	}
	scriptNode.SetParameters(params)

	// Создаём snapshot контекста для захвата переменных
	// В Java версии closure захватывает ссылку на контекст, а не копию
	// Это позволяет видеть изменения переменных после создания closure
	// Но для правильной работы нужно использовать ссылку на контекст, а не копию
	var snapshot jexl.Context = capturedContext

	baseScript := &script{
		engine: engine,
		source: lambda.SourceText(),
		ast:    scriptNode,
	}

	return &closure{
		script:          baseScript,
		capturedContext: snapshot,
	}
}

// Execute выполняет closure с аргументами, используя захваченный контекст.
func (c *closure) Execute(ctx jexl.Context, args ...any) (any, error) {
	// Создаём контекст для параметров
	paramNames := c.Parameters()
	paramValues := make(map[string]any, len(paramNames))
	for i, name := range paramNames {
		if i < len(args) {
			paramValues[name] = args[i]
		} else {
			paramValues[name] = nil
		}
	}

	// Используем переданный контекст как базовый, если он есть
	// Если переданный контекст nil, используем захваченный контекст
	baseCtx := ctx
	if baseCtx == nil {
		baseCtx = c.capturedContext
	}
	if baseCtx == nil {
		baseCtx = jexl.NewMapContext()
	}

	// Если базовый контекст - это argumentContext, извлекаем его базовый контекст
	// Это позволяет видеть переменные, установленные в скрипте (например, через var)
	// Но также проверяем захваченный контекст, чтобы видеть переменные, установленные до создания closure
	if argCtx, ok := baseCtx.(*argumentContext); ok {
		// Если базовый контекст - это argumentContext, используем его базовый контекст
		// Но также добавляем захваченный контекст, чтобы видеть переменные из скрипта
		baseCtx = argCtx.base
		// Если базовый контекст и захваченный контекст - это один и тот же объект, используем только базовый
		if baseCtx == c.capturedContext {
			// Не нужно добавлять capturedContext, так как это тот же объект
		} else {
			// Используем базовый контекст как есть - он уже содержит переменные из скрипта
		}
	}

	// Создаём объединённый контекст: параметры -> базовый контекст (где могут быть переменные после создания closure) -> захваченный контекст
	// Важно: базовый контекст проверяется первым, чтобы видеть переменные, установленные после создания closure
	// Если базовый контекст и захваченный контекст - это один и тот же объект, используем только базовый
	var capturedCtx jexl.Context = c.capturedContext
	if baseCtx == capturedCtx {
		capturedCtx = nil // Избегаем двойной проверки одного и того же контекста
	}
	execCtx := newClosureContext(baseCtx, capturedCtx, paramValues)

	// Выполняем тело lambda напрямую через интерпретатор
	interp := newInterpreter(c.engine, execCtx)
	return interp.interpret(c.ast.Children()[0]) // Тело lambda - первый (и единственный) дочерний узел
}

// closureContext объединяет несколько контекстов для closure.
type closureContext struct {
	paramValues map[string]any
	capturedCtx jexl.Context
	baseCtx     jexl.Context
}

func newClosureContext(base jexl.Context, captured jexl.Context, params map[string]any) jexl.Context {
	return &closureContext{
		baseCtx:     base,
		capturedCtx: captured,
		paramValues: params,
	}
}

func (c *closureContext) Get(name string) any {
	// Сначала проверяем параметры
	if val, ok := c.paramValues[name]; ok {
		return val
	}
	// Затем базовый контекст (где могут быть установлены переменные после создания closure)
	// Базовый контекст может быть argumentContext, который содержит параметры скрипта и базовый контекст
	if c.baseCtx != nil {
		// Проверяем Has перед Get, чтобы избежать ошибок
		if c.baseCtx.Has(name) {
			return c.baseCtx.Get(name)
		}
	}
	// Затем захваченный контекст (для переменных, установленных до создания closure)
	if c.capturedCtx != nil {
		if c.capturedCtx.Has(name) {
			return c.capturedCtx.Get(name)
		}
	}
	return nil
}

func (c *closureContext) Has(name string) bool {
	if _, ok := c.paramValues[name]; ok {
		return true
	}
	// Сначала проверяем базовый контекст
	if c.baseCtx != nil && c.baseCtx.Has(name) {
		return true
	}
	// Затем захваченный контекст
	if c.capturedCtx != nil && c.capturedCtx.Has(name) {
		return true
	}
	return false
}

func (c *closureContext) Set(name string, value any) {
	// Если это параметр, обновляем его
	if _, ok := c.paramValues[name]; ok {
		c.paramValues[name] = value
		return
	}
	// Иначе устанавливаем в базовый контекст
	if c.baseCtx != nil {
		c.baseCtx.Set(name, value)
	}
}
