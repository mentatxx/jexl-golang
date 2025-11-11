package jexl

// Expression представляет одно выражение JEXL.
// Аналог интерфейса org.apache.commons.jexl3.JexlExpression.
type Expression interface {
	// Callable возвращает функцию, которую можно выполнить асинхронно.
	Callable(ctx Context) func() (any, error)
	// Evaluate вычисляет выражение в заданном контексте.
	Evaluate(ctx Context) (any, error)
	// ParsedText возвращает восстановленный текст выражения из AST.
	ParsedText() string
	// SourceText возвращает исходный текст.
	SourceText() string
}
