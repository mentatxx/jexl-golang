package jexl

// Script представляет скрипт JEXL.
// Аналог org.apache.commons.jexl3.JexlScript.
// Script также реализует Expression, как в Java версии.
type Script interface {
	Expression // Script включает все методы Expression
	
	CallableWithArgs(ctx Context, args ...any) func() (any, error)
	Curry(args ...any) Script
	Execute(ctx Context, args ...any) (any, error)
	LocalVariables() []string
	Parameters() []string
	ParsedTextWithIndent(indent int) string
	Pragmas() map[string]any
	UnboundParameters() []string
	Variables() [][]string
}
