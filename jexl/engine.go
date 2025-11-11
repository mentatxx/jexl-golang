package jexl

// Engine описывает поведение движка JEXL.
// Это прямой аналог абстрактного класса org.apache.commons.jexl3.JexlEngine.
type Engine interface {
	// ClearCache очищает кэш выражений.
	ClearCache()

	// CreateExpression компилирует строку в выражение.
	CreateExpression(info *Info, source string) (Expression, error)

	// CreateScript компилирует строку в исполняемый скрипт.
	CreateScript(features *Features, info *Info, source string, names ...string) (Script, error)

	// CreateTemplateEngine создаёт движок шаблонов JXLT.
	CreateTemplateEngine(opts ...TemplateOption) (*TemplateEngine, error)

	// ThreadContext возвращает текущий thread-local контекст.
	ThreadContext() ThreadLocalContext

	// SetThreadContext устанавливает thread-local контекст.
	SetThreadContext(ctx ThreadLocalContext)

	// Uberspect возвращает introspection-объект.
	Uberspect() Uberspect

	// Options возвращает текущие опции движка.
	Options() *Options

	// IsStrict сообщает, включён ли строгий режим.
	IsStrict() bool
	
	// Arithmetic возвращает арифметику движка.
	Arithmetic() Arithmetic

	// GetProperty получает свойство объекта по выражению.
	GetProperty(ctx Context, bean any, expr string) (any, error)

	// GetPropertyWithoutContext получает свойство объекта по выражению без контекста.
	GetPropertyWithoutContext(bean any, expr string) (any, error)

	// SetProperty устанавливает свойство объекта по выражению.
	SetProperty(ctx Context, bean any, expr string, value any) error

	// SetPropertyWithoutContext устанавливает свойство объекта по выражению без контекста.
	SetPropertyWithoutContext(bean any, expr string, value any) error

	// InvokeMethod вызывает метод объекта.
	InvokeMethod(obj any, method string, args ...any) (any, error)

	// NewInstance создаёт новый экземпляр объекта.
	NewInstance(className string, args ...any) (any, error)
}
