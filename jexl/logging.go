package jexl

// Logger определяет минимальный интерфейс для логирования, аналог org.apache.commons.logging.Log.
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// NoopLogger используется по умолчанию.
type NoopLogger struct{}

// Debugf игнорирует сообщения.
func (NoopLogger) Debugf(string, ...any) {}

// Infof игнорирует сообщения.
func (NoopLogger) Infof(string, ...any) {}

// Warnf игнорирует сообщения.
func (NoopLogger) Warnf(string, ...any) {}

// Errorf игнорирует сообщения.
func (NoopLogger) Errorf(string, ...any) {}
