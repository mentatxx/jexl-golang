package internal

import (
	"github.com/mentatxx/jexl-golang/jexl"
)

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
	
	baseScript := &script{
		engine: engine,
		source: lambda.SourceText(),
		ast:    scriptNode,
	}
	
	return &closure{
		script:          baseScript,
		capturedContext: capturedContext,
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
	
	// Используем захваченный контекст как базовый, если переданный контекст nil
	baseCtx := ctx
	if baseCtx == nil {
		baseCtx = c.capturedContext
	}
	if baseCtx == nil {
		baseCtx = jexl.NewMapContext()
	}
	
	// Создаём объединённый контекст: параметры -> захваченный контекст -> переданный контекст
	execCtx := newClosureContext(baseCtx, c.capturedContext, paramValues)
	
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
	// Затем захваченный контекст
	if c.capturedCtx != nil {
		if c.capturedCtx.Has(name) {
			return c.capturedCtx.Get(name)
		}
	}
	// Затем базовый контекст
	if c.baseCtx != nil {
		return c.baseCtx.Get(name)
	}
	return nil
}

func (c *closureContext) Has(name string) bool {
	if _, ok := c.paramValues[name]; ok {
		return true
	}
	if c.capturedCtx != nil && c.capturedCtx.Has(name) {
		return true
	}
	if c.baseCtx != nil {
		return c.baseCtx.Has(name)
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

