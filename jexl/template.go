package jexl

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// TemplateEngine представляет движок JXLT.
type TemplateEngine struct {
	engine        Engine
	config        *TemplateConfig
	expressionCache Cache[string, TemplateExpression]
}

// NewTemplateEngine создаёт новый TemplateEngine.
func NewTemplateEngine(engine Engine, config *TemplateConfig) *TemplateEngine {
	var cache Cache[string, TemplateExpression]
	if config.CacheSize > 0 {
		// Создаём кэш с правильным типом
		baseCache := DefaultCacheFactory(config.CacheSize)
		cache = &templateExpressionCacheAdapter{cache: baseCache}
	}
	return &TemplateEngine{
		engine:          engine,
		config:          config,
		expressionCache: cache,
	}
}

// templateExpressionCacheAdapter адаптирует Cache[string, any] к Cache[string, TemplateExpression]
type templateExpressionCacheAdapter struct {
	cache Cache[string, any]
}

func (a *templateExpressionCacheAdapter) Get(key string) (TemplateExpression, bool) {
	val, ok := a.cache.Get(key)
	if !ok {
		return nil, false
	}
	expr, ok := val.(TemplateExpression)
	return expr, ok
}

func (a *templateExpressionCacheAdapter) Put(key string, value TemplateExpression) {
	a.cache.Put(key, value)
}

func (a *templateExpressionCacheAdapter) Clear() {
	a.cache.Clear()
}

func (a *templateExpressionCacheAdapter) Size() int {
	return a.cache.Size()
}

// TemplateOption задаёт опцию для TemplateEngine.
type TemplateOption interface {
	Apply(*TemplateConfig)
}

// TemplateConfig агрегирует значения опций.
type TemplateConfig struct {
	NoScript      bool
	CacheSize     int
	ImmediateRune rune
	DeferredRune  rune
}

// TemplateOptionFunc облегчает создание опций на функциях.
type TemplateOptionFunc func(*TemplateConfig)

func (f TemplateOptionFunc) Apply(cfg *TemplateConfig) {
	f(cfg)
}

// WithNoScript запрещает использование скриптов в шаблонах.
func WithNoScript(value bool) TemplateOption {
	return TemplateOptionFunc(func(cfg *TemplateConfig) {
		cfg.NoScript = value
	})
}

// WithCacheSize задаёт размер кэша.
func WithCacheSize(size int) TemplateOption {
	return TemplateOptionFunc(func(cfg *TemplateConfig) {
		cfg.CacheSize = size
	})
}

// WithImmediateRune задаёт символ немедленной вставки.
func WithImmediateRune(r rune) TemplateOption {
	return TemplateOptionFunc(func(cfg *TemplateConfig) {
		cfg.ImmediateRune = r
	})
}

// WithDeferredRune задаёт символ отложенной вставки.
func WithDeferredRune(r rune) TemplateOption {
	return TemplateOptionFunc(func(cfg *TemplateConfig) {
		cfg.DeferredRune = r
	})
}

// TemplateExpression представляет выражение шаблона.
type TemplateExpression interface {
	Evaluate(ctx Context) (string, error)
	AsString() string
	IsDeferred() bool
	IsImmediate() bool
}

// CreateExpression создаёт выражение шаблона из строки.
func (te *TemplateEngine) CreateExpression(source string) (TemplateExpression, error) {
	if te.expressionCache != nil {
		if cached, ok := te.expressionCache.Get(source); ok {
			return cached, nil
		}
	}

	expr, err := te.parseTemplateExpression(source)
	if err != nil {
		return nil, err
	}

	if te.expressionCache != nil {
		te.expressionCache.Put(source, expr)
	}

	return expr, nil
}

// parseTemplateExpression парсит выражение шаблона.
func (te *TemplateEngine) parseTemplateExpression(source string) (TemplateExpression, error) {
	immediateRune := te.config.ImmediateRune
	deferredRune := te.config.DeferredRune

	var parts []templatePart
	var current strings.Builder
	var inExpression bool
	var isDeferred bool
	var exprStart int

	for i, r := range source {
		if !inExpression {
			if r == immediateRune && i+1 < len(source) && source[i+1] == '{' {
				// Начало immediate выражения
				if current.Len() > 0 {
					parts = append(parts, templatePart{text: current.String(), isText: true})
					current.Reset()
				}
				inExpression = true
				isDeferred = false
				exprStart = i + 2
				i++ // Пропускаем '{'
				continue
			} else if r == deferredRune && i+1 < len(source) && source[i+1] == '{' {
				// Начало deferred выражения
				if current.Len() > 0 {
					parts = append(parts, templatePart{text: current.String(), isText: true})
					current.Reset()
				}
				inExpression = true
				isDeferred = true
				exprStart = i + 2
				i++ // Пропускаем '{'
				continue
			}
			current.WriteRune(r)
		} else {
			if r == '}' {
				// Конец выражения
				expr := source[exprStart:i]
				parts = append(parts, templatePart{
					expression: expr,
					isText:     false,
					isDeferred: isDeferred,
				})
				current.Reset()
				inExpression = false
			}
		}
	}

	if current.Len() > 0 {
		parts = append(parts, templatePart{text: current.String(), isText: true})
	}

	return &templateExpressionImpl{
		engine: te.engine,
		parts:  parts,
		source: source,
	}, nil
}

type templatePart struct {
	text       string
	expression string
	isText     bool
	isDeferred bool
}

type templateExpressionImpl struct {
	engine Engine
	parts  []templatePart
	source string
}

func (t *templateExpressionImpl) Evaluate(ctx Context) (string, error) {
	var result strings.Builder

	for _, part := range t.parts {
		if part.isText {
			result.WriteString(part.text)
		} else {
			expr, err := t.engine.CreateExpression(nil, part.expression)
			if err != nil {
				return "", err
			}
			val, err := expr.Evaluate(ctx)
			if err != nil {
				return "", err
			}
			if val != nil {
				result.WriteString(fmt.Sprintf("%v", val))
			}
		}
	}

	return result.String(), nil
}

func (t *templateExpressionImpl) AsString() string {
	return t.source
}

func (t *templateExpressionImpl) IsDeferred() bool {
	for _, part := range t.parts {
		if !part.isText && part.isDeferred {
			return true
		}
	}
	return false
}

func (t *templateExpressionImpl) IsImmediate() bool {
	for _, part := range t.parts {
		if !part.isText && !part.isDeferred {
			return true
		}
	}
	return false
}

// ClearCache очищает кэш выражений.
func (te *TemplateEngine) ClearCache() {
	if te.expressionCache != nil {
		te.expressionCache.Clear()
	}
}

// Template представляет шаблон JXLT.
type Template interface {
	Evaluate(ctx Context, writer io.Writer) error
	EvaluateWithArgs(ctx Context, writer io.Writer, args ...any) error
	AsString() string
}

// CreateTemplate создаёт шаблон из строки.
func (te *TemplateEngine) CreateTemplate(source string) (Template, error) {
	// Упрощённая реализация: парсим шаблон, где строки начинающиеся с '$$' - это код
	lines := strings.Split(source, "\n")
	var scriptParts []string
	var templateParts []templatePart
	var currentText strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "$$") {
			// Это JEXL код
			if currentText.Len() > 0 {
				templateParts = append(templateParts, templatePart{
					text:   currentText.String(),
					isText: true,
				})
				currentText.Reset()
			}
			scriptParts = append(scriptParts, strings.TrimPrefix(trimmed, "$$"))
		} else {
			// Это текст с возможными выражениями
			currentText.WriteString(line)
			currentText.WriteString("\n")
		}
	}

	if currentText.Len() > 0 {
		templateParts = append(templateParts, templatePart{
			text:   currentText.String(),
			isText: true,
		})
	}

	// Парсим выражения в тексте
	var finalParts []templatePart
	for _, part := range templateParts {
		if part.isText {
			expr, err := te.parseTemplateExpression(part.text)
			if err == nil {
				if exprImpl, ok := expr.(*templateExpressionImpl); ok {
					finalParts = append(finalParts, exprImpl.parts...)
					continue
				}
			}
		}
		finalParts = append(finalParts, part)
	}

	// Создаём скрипт из частей кода
	var scriptSource strings.Builder
	for i, part := range scriptParts {
		if i > 0 {
			scriptSource.WriteString("; ")
		}
		scriptSource.WriteString(part)
	}

	var script Script
	var err error
	if scriptSource.Len() > 0 {
		script, err = te.engine.CreateScript(nil, nil, scriptSource.String())
		if err != nil {
			return nil, err
		}
	}

	return &templateImpl{
		engine:        te.engine,
		script:        script,
		templateParts: finalParts,
		source:        source,
	}, nil
}

type templateImpl struct {
	engine        Engine
	script        Script
	templateParts []templatePart
	source        string
}

func (t *templateImpl) Evaluate(ctx Context, writer io.Writer) error {
	return t.EvaluateWithArgs(ctx, writer)
}

func (t *templateImpl) EvaluateWithArgs(ctx Context, writer io.Writer, args ...any) error {
	// Если есть скрипт, выполняем его
	if t.script != nil {
		_, err := t.script.Execute(ctx, args...)
		if err != nil {
			return err
		}
	}

	// Выводим части шаблона
	for _, part := range t.templateParts {
		if part.isText {
			writer.Write([]byte(part.text))
		} else {
			expr, err := t.engine.CreateExpression(nil, part.expression)
			if err != nil {
				return err
			}
			val, err := expr.Evaluate(ctx)
			if err != nil {
				return err
			}
			if val != nil {
				fmt.Fprint(writer, val)
			}
		}
	}

	return nil
}

func (t *templateImpl) AsString() string {
	return t.source
}

// EvaluateString выполняет выражение шаблона и возвращает строку.
func (te *TemplateEngine) EvaluateString(ctx Context, source string) (string, error) {
	expr, err := te.CreateExpression(source)
	if err != nil {
		return "", err
	}
	return expr.Evaluate(ctx)
}

// EvaluateTemplate выполняет шаблон и записывает результат в writer.
func (te *TemplateEngine) EvaluateTemplate(ctx Context, writer io.Writer, source string) error {
	tmpl, err := te.CreateTemplate(source)
	if err != nil {
		return err
	}
	return tmpl.Evaluate(ctx, writer)
}

// EvaluateTemplateToString выполняет шаблон и возвращает результат как строку.
func (te *TemplateEngine) EvaluateTemplateToString(ctx Context, source string) (string, error) {
	var buf bytes.Buffer
	err := te.EvaluateTemplate(ctx, &buf, source)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
