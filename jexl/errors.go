package jexl

import "fmt"

// Error представляет ошибку движка JEXL.
type Error struct {
	message string
	cause   error
	info    *Info
}

// Error реализует интерфейс error.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// Unwrap возвращает вложенную ошибку.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// Info возвращает связанный Info.
func (e *Error) Info() *Info {
	if e == nil {
		return nil
	}
	return e.info
}

// NewError создаёт ошибку.
func NewError(message string) *Error {
	return &Error{message: message}
}

// WrapError оборачивает существующую ошибку.
func WrapError(message string, cause error, info *Info) *Error {
	return &Error{
		message: message,
		cause:   cause,
		info:    info,
	}
}

// ParsingError представляет ошибку парсинга.
type ParsingError struct {
	*Error
	expression string
}

// NewParsingError создаёт ошибку парсинга.
func NewParsingError(message, expression string, info *Info) *ParsingError {
	return &ParsingError{
		Error:      WrapError(message, nil, info),
		expression: expression,
	}
}

// Expression возвращает выражение, вызвавшее ошибку.
func (e *ParsingError) Expression() string {
	return e.expression
}

// MethodError представляет ошибку вызова метода.
type MethodError struct {
	*Error
	method string
	args   []any
}

// NewMethodError создаёт ошибку вызова метода.
func NewMethodError(method string, args []any, info *Info, cause error) *MethodError {
	return &MethodError{
		Error:  WrapError("unsolvable function/method '"+methodSignature(method, args)+"'", cause, info),
		method: method,
		args:   args,
	}
}

// Method возвращает имя метода.
func (e *MethodError) Method() string {
	return e.method
}

// Args возвращает аргументы метода.
func (e *MethodError) Args() []any {
	return e.args
}

// OperatorError представляет ошибку оператора.
type OperatorError struct {
	*Error
	symbol string
}

// NewOperatorError создаёт ошибку оператора.
func NewOperatorError(symbol string, info *Info, cause error) *OperatorError {
	return &OperatorError{
		Error:  WrapError("error calling operator '"+symbol+"'", cause, info),
		symbol: symbol,
	}
}

// Symbol возвращает символ оператора.
func (e *OperatorError) Symbol() string {
	return e.symbol
}

// PropertyError представляет ошибку доступа к свойству.
type PropertyError struct {
	*Error
	property string
}

// NewPropertyError создаёт ошибку доступа к свойству.
func NewPropertyError(property string, info *Info, cause error) *PropertyError {
	return &PropertyError{
		Error:    WrapError("error accessing property '"+property+"'", cause, info),
		property: property,
	}
}

// Property возвращает имя свойства.
func (e *PropertyError) Property() string {
	return e.property
}

// methodSignature создаёт строку сигнатуры метода.
func methodSignature(method string, args []any) string {
	if len(args) == 0 {
		return method + "()"
	}
	// Упрощённая версия - в реальности нужно форматировать типы
	return method + "(...)"
}

var (
	// ErrNotImplemented используется для незавершённых элементов порта.
	ErrNotImplemented = NewError("not implemented")
)
