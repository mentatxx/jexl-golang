package jexl

// ScriptParser соответствует org.apache.commons.jexl3.parser.JexlScriptParser.
type ScriptParser interface {
	Parse(info *Info, source string, features *Features, names []string) (Script, error)
}

// ParserFactory создаёт новый парсер.
type ParserFactory interface {
	New() ScriptParser
}

// ParserFactoryFunc адаптер функций.
type ParserFactoryFunc func() ScriptParser

// New создаёт парсер.
func (f ParserFactoryFunc) New() ScriptParser {
	return f()
}
