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

var (
	// ErrNotImplemented используется для незавершённых элементов порта.
	ErrNotImplemented = NewError("not implemented")
)
