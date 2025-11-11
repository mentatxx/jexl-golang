package jexl

// Operator представляет оператор JEXL.
// Порт org.apache.commons.jexl3.JexlOperator.
type Operator struct {
	symbol     string
	methodName string
	arity      int
	base       *Operator
}

// Symbol возвращает символ оператора.
func (o *Operator) Symbol() string {
	return o.symbol
}

// MethodName возвращает имя метода, связанного с оператором.
func (o *Operator) MethodName() string {
	return o.methodName
}

// Arity возвращает количество аргументов оператора.
func (o *Operator) Arity() int {
	return o.arity
}

// BaseOperator возвращает базовый оператор (для side-effect операторов).
func (o *Operator) BaseOperator() *Operator {
	return o.base
}

// Определение операторов
var (
	// Арифметические операторы
	OpAdd      = &Operator{"+", "add", 2, nil}
	OpSubtract = &Operator{"-", "subtract", 2, nil}
	OpMultiply = &Operator{"*", "multiply", 2, nil}
	OpDivide   = &Operator{"/", "divide", 2, nil}
	OpMod      = &Operator{"%", "mod", 2, nil}

	// Побитовые операторы
	OpAnd         = &Operator{"&", "and", 2, nil}
	OpOr          = &Operator{"|", "or", 2, nil}
	OpXor         = &Operator{"^", "xor", 2, nil}
	OpShiftRight  = &Operator{">>", "shiftRight", 2, nil}
	OpShiftRightU = &Operator{">>>", "shiftRightUnsigned", 2, nil}
	OpShiftLeft   = &Operator{"<<", "shiftLeft", 2, nil}

	// Операторы сравнения
	OpEq      = &Operator{"==", "equals", 2, nil}
	OpEqStrict = &Operator{"===", "strictEquals", 2, nil}
	OpNe      = &Operator{"!=", "notEquals", 2, nil}
	OpLt      = &Operator{"<", "lessThan", 2, nil}
	OpLe      = &Operator{"<=", "lessThanOrEqual", 2, nil}
	OpGt      = &Operator{">", "greaterThan", 2, nil}
	OpGe      = &Operator{">=", "greaterThanOrEqual", 2, nil}

	// Строковые операторы
	OpContains   = &Operator{"=~", "contains", 2, nil}
	OpStartsWith = &Operator{"=^", "startsWith", 2, nil}
	OpEndsWith   = &Operator{"=$", "endsWith", 2, nil}

	// Унарные операторы
	OpNot        = &Operator{"!", "not", 1, nil}
	OpComplement = &Operator{"~", "complement", 1, nil}
	OpNegate     = &Operator{"-", "negate", 1, nil}
	OpPositivize = &Operator{"+", "positivize", 1, nil}
	OpEmpty      = &Operator{"empty", "empty", 1, nil}
	OpSize       = &Operator{"size", "size", 1, nil}

	// Side-effect операторы
	OpSelfAdd      = &Operator{"+=", "selfAdd", 2, OpAdd}
	OpSelfSubtract = &Operator{"-=", "selfSubtract", 2, OpSubtract}
	OpSelfMultiply = &Operator{"*=", "selfMultiply", 2, OpMultiply}
	OpSelfDivide   = &Operator{"/=", "selfDivide", 2, OpDivide}
	OpSelfMod      = &Operator{"%=", "selfMod", 2, OpMod}
	OpSelfAnd      = &Operator{"&=", "selfAnd", 2, OpAnd}
	OpSelfOr       = &Operator{"|=", "selfOr", 2, OpOr}
	OpSelfXor      = &Operator{"^=", "selfXor", 2, OpXor}
	OpSelfShiftRight  = &Operator{">>=", "selfShiftRight", 2, OpShiftRight}
	OpSelfShiftRightU = &Operator{">>>=", "selfShiftRightUnsigned", 2, OpShiftRightU}
	OpSelfShiftLeft   = &Operator{"<<=", "selfShiftLeft", 2, OpShiftLeft}

	// Инкремент/декремент
	OpIncrement      = &Operator{"+1", "increment", 1, nil}
	OpDecrement      = &Operator{"-1", "decrement", 1, nil}
	OpIncrementAndGet = &Operator{"++.", "incrementAndGet", 1, OpIncrement}
	OpGetAndIncrement = &Operator{".++", "getAndIncrement", 1, OpIncrement}
	OpDecrementAndGet = &Operator{"--.", "decrementAndGet", 1, OpDecrement}
	OpGetAndDecrement = &Operator{".--", "getAndDecrement", 1, OpDecrement}

	// Специальные операторы
	OpPropertyGet = &Operator{".", "propertyGet", 2, nil}
	OpPropertySet = &Operator{".=", "propertySet", 3, nil}
	OpArrayGet    = &Operator{"[]", "arrayGet", 2, nil}
	OpArraySet    = &Operator{"[]=", "arraySet", 3, nil}
	OpForEach     = &Operator{"for(...)", "forEach", 1, nil}
	OpCondition   = &Operator{"?", "testCondition", 1, nil}
	OpCompare     = &Operator{"<>", "compare", 2, nil}

	// Отрицательные операторы
	OpNotContains   = &Operator{"!~", "", 2, OpContains}
	OpNotStartsWith = &Operator{"!^", "", 2, OpStartsWith}
	OpNotEndsWith   = &Operator{"!$", "", 2, OpEndsWith}
)

// OperatorFromSymbol возвращает оператор по символу.
func OperatorFromSymbol(symbol string) *Operator {
	operators := map[string]*Operator{
		"+": OpAdd, "-": OpSubtract, "*": OpMultiply, "/": OpDivide, "%": OpMod,
		"&": OpAnd, "|": OpOr, "^": OpXor,
		">>": OpShiftRight, ">>>": OpShiftRightU, "<<": OpShiftLeft,
		"==": OpEq, "===": OpEqStrict, "!=": OpNe,
		"<": OpLt, "<=": OpLe, ">": OpGt, ">=": OpGe,
		"=~": OpContains, "=^": OpStartsWith, "=$": OpEndsWith,
		"!": OpNot, "~": OpComplement,
		"empty": OpEmpty, "size": OpSize,
		"+=": OpSelfAdd, "-=": OpSelfSubtract, "*=": OpSelfMultiply,
		"/=": OpSelfDivide, "%=": OpSelfMod,
		"&=": OpSelfAnd, "|=": OpSelfOr, "^=": OpSelfXor,
		">>=": OpSelfShiftRight, ">>>=": OpSelfShiftRightU, "<<=": OpSelfShiftLeft,
		"++.": OpIncrementAndGet, ".++": OpGetAndIncrement,
		"--.": OpDecrementAndGet, ".--": OpGetAndDecrement,
		".": OpPropertyGet, ".=": OpPropertySet,
		"[]": OpArrayGet, "[]=": OpArraySet,
		"for(...)": OpForEach, "?": OpCondition, "<>": OpCompare,
		"!~": OpNotContains, "!^": OpNotStartsWith, "!$": OpNotEndsWith,
	}
	return operators[symbol]
}

