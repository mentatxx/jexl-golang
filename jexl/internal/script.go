package internal

import (
	"strings"

	"github.com/mentatxx/jexl-golang/jexl"
)

// script реализует jexl.Script и jexl.Expression.
// Порт org.apache.commons.jexl3.internal.Script.
type script struct {
	engine    jexl.Engine
	source    string
	ast       *jexl.ScriptNode
	boundArgs []any
}

// NewScript создаёт новый script из AST.
func NewScript(engine jexl.Engine, source string, ast *jexl.ScriptNode) jexl.Script {
	return &script{
		engine: engine,
		source: source,
		ast:    ast,
	}
}

// Execute выполняет скрипт с контекстом и аргументами.
func (s *script) Execute(ctx jexl.Context, args ...any) (any, error) {
	execCtx := ctx
	if execCtx == nil {
		execCtx = jexl.NewMapContext()
	}

	allArgs := append([]any{}, s.boundArgs...)
	allArgs = append(allArgs, args...)

	paramNames := s.Parameters()
	if len(paramNames) > 0 {
		params := make(map[string]any, len(paramNames))
		for i, name := range paramNames {
			if i < len(allArgs) {
				params[name] = allArgs[i]
			} else {
				params[name] = nil
			}
		}
		execCtx = newArgumentContext(execCtx, params)
	}

	interp := newInterpreter(s.engine, execCtx)
	return interp.interpret(s.ast)
}

// Evaluate реализует Expression.Evaluate.
func (s *script) Evaluate(ctx jexl.Context) (any, error) {
	return s.Execute(ctx)
}

// SourceText возвращает исходный текст.
func (s *script) SourceText() string {
	return s.source
}

// ParsedText возвращает распарсенный текст.
func (s *script) ParsedText() string {
	return s.ast.String()
}

// ParsedTextWithIndent возвращает распарсенный текст с отступами.
func (s *script) ParsedTextWithIndent(indent int) string {
	if indent <= 0 {
		return s.ParsedText()
	}
	prefix := strings.Repeat(" ", indent)
	lines := strings.Split(s.ParsedText(), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// Parameters возвращает параметры скрипта.
func (s *script) Parameters() []string {
	return s.ast.Parameters()
}

// LocalVariables возвращает локальные переменные.
func (s *script) LocalVariables() []string {
	return s.ast.Variables()
}

// Pragmas возвращает pragma директивы.
func (s *script) Pragmas() map[string]any {
	return s.ast.Pragmas()
}

// UnboundParameters возвращает несвязанные параметры.
func (s *script) UnboundParameters() []string {
	params := s.Parameters()
	if len(s.boundArgs) >= len(params) {
		return nil
	}
	return params[len(s.boundArgs):]
}

// Variables возвращает переменные скрипта.
func (s *script) Variables() [][]string {
	return nil
}

// Curry создаёт новый скрипт с частично применёнными аргументами.
func (s *script) Curry(args ...any) jexl.Script {
	if len(args) == 0 {
		return s
	}
	allArgs := append([]any{}, s.boundArgs...)
	allArgs = append(allArgs, args...)

	return &script{
		engine:    s.engine,
		source:    s.source,
		ast:       s.ast,
		boundArgs: allArgs,
	}
}

// Callable создаёт Callable для асинхронного выполнения.
func (s *script) Callable(ctx jexl.Context) func() (any, error) {
	return func() (any, error) {
		return s.Execute(ctx)
	}
}

// CallableWithArgs создаёт Callable с аргументами.
func (s *script) CallableWithArgs(ctx jexl.Context, args ...any) func() (any, error) {
	return func() (any, error) {
		return s.Execute(ctx, args...)
	}
}

type argumentContext struct {
	base   jexl.Context
	params map[string]any
}

func newArgumentContext(base jexl.Context, params map[string]any) jexl.Context {
	if base == nil {
		base = jexl.NewMapContext()
	}
	return &argumentContext{
		base:   base,
		params: params,
	}
}

func (a *argumentContext) Get(name string) any {
	if val, ok := a.params[name]; ok {
		return val
	}
	return a.base.Get(name)
}

func (a *argumentContext) Has(name string) bool {
	if _, ok := a.params[name]; ok {
		return true
	}
	return a.base.Has(name)
}

func (a *argumentContext) Set(name string, value any) {
	if _, ok := a.params[name]; ok {
		a.params[name] = value
		return
	}
	a.base.Set(name, value)
}
