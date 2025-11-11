package jexl

import (
	"context"
	"sync/atomic"
)

// Context описывает контейнер переменных, доступных во время выполнения выражений JEXL.
// Большинство методов и вложенных интерфейсов следуют структуре Java-интерфейса JexlContext.
type Context interface {
	// Get возвращает значение переменной.
	Get(name string) any
	// Has проверяет, определена ли переменная в контексте.
	Has(name string) bool
	// Set устанавливает значение переменной.
	Set(name string, value any)
}

// AnnotationProcessor обрабатывает пользовательские аннотации во время интерпретации.
type AnnotationProcessor interface {
	ProcessAnnotation(name string, args []any, statement func() (any, error)) (any, error)
}

// CancellationHandle предоставляет доступ к признаку отмены вычислений.
type CancellationHandle interface {
	// Cancellation возвращает atomic.Boolean аналог из Go.
	Cancellation() *atomic.Bool
}

// ClassNameResolver разрешает простые имена классов в fully-qualified имена.
type ClassNameResolver interface {
	ResolveClassName(name string) string
}

// ModuleProcessor отвечает за обработку определения модулей pragma.
type ModuleProcessor interface {
	ProcessModule(engine Engine, info *Info, name string, body string) any
}

// NamespaceFunctor создаёт объект-функтор для пространства имён.
type NamespaceFunctor interface {
	CreateFunctor(ctx Context) any
}

// NamespaceResolver разрешает имя пространства имён в объект.
type NamespaceResolver interface {
	ResolveNamespace(name string) any
}

// OptionsHandle предоставляет доступ к текущим опциям исполнения.
type OptionsHandle interface {
	EngineOptions() *Options
}

// PragmaProcessor обрабатывает pragma-директивы.
type PragmaProcessor interface {
	ProcessPragma(opts *Options, key string, value any)
}

// ThreadLocalContext определяет контекст, который должен быть помещён в thread-local область во время выполнения.
type ThreadLocalContext interface {
	Context
}

// ContextAware устанавливает контекст выполнения в Engine.
type ContextAware interface {
	context.Context
}
