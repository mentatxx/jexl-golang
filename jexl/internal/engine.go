package internal

import (
	"strings"
	"sync"

	"github.com/mentatxx/jexl-golang/jexl"
)

// engine реализует интерфейс jexl.Engine.
// Это порт org.apache.commons.jexl3.internal.Engine.
type engine struct {
	mu                 sync.RWMutex
	uberspect          jexl.Uberspect
	arithmetic         jexl.Arithmetic
	logger             jexl.Logger
	options            *jexl.Options
	features           *jexl.Features
	scriptFeatures     *jexl.Features
	expressionFeatures *jexl.Features
	cache              jexl.Cache[string, any]
	cacheThreshold     int
	stackOverflow      int
	collectMode        int
	strict             bool
	safe               bool
	silent             bool
	cancellable        bool
	debug              bool
	charset            string
	classLoader        any
	parserFactory      jexl.ParserFactory
	threadContext      jexl.ThreadLocalContext
	parser             Parser
}

// NewEngine создаёт новый экземпляр движка из Builder.
func NewEngine(builder *jexl.Builder) (jexl.Engine, error) {
	opts := builder.Options()
	if opts == nil {
		opts = jexl.NewOptions()
	} else {
		opts = opts.Copy()
	}

	eng := &engine{
		options:        opts,
		strict:         opts.Strict(),
		safe:           opts.Safe(),
		silent:         opts.Silent(),
		cancellable:    opts.Cancellable(),
		cacheThreshold: 64, // default
		stackOverflow:  maxInt,
		collectMode:    1,
		debug:          true,
		arithmetic:     jexl.NewBaseArithmetic(opts.StrictArithmetic(), opts.MathContext(), opts.MathScale()),
		logger:         jexl.NoopLogger{},
	}

	// Настройка features
	features := builder.FeaturesValue()
	if features == nil {
		features = jexl.FeaturesDefault()
	}
	eng.features = features
	eng.scriptFeatures = features
	eng.expressionFeatures = features

	// Настройка uberspect
	uberspect := builder.UberspectValue()
	if uberspect == nil {
		permissions := builder.PermissionsValue()
		if permissions == nil {
			permissions = jexl.PermissionsRestricted.Clone()
		}
		strategy := builder.StrategyValue()
		if strategy == nil {
			strategy = jexl.ResolverStrategyDefault
		}
		uberspect = NewUberspect(eng.logger, strategy, permissions)
	}
	eng.uberspect = uberspect

	// Настройка sandbox
	sandbox := builder.SandboxValue()
	if sandbox != nil {
		eng.uberspect = NewSandboxUberspect(eng.uberspect, sandbox)
	}

	// Настройка кэша
	cacheSize := builder.CacheSize()
	if cacheSize > 0 {
		cacheFactory := builder.CacheFactoryValue()
		if cacheFactory == nil {
			cacheFactory = jexl.DefaultCacheFactory
		}
		eng.cache = cacheFactory(cacheSize)
	}

	// Настройка parser factory
	eng.parserFactory = builder.ParserFactoryValue()

	// Инициализация парсера
	eng.parser = eng.getParser()

	// Настройка других параметров
	if builder.ArithmeticValue() != nil {
		eng.arithmetic = builder.ArithmeticValue()
	}
	if builder.LoggerValue() != nil {
		eng.logger = builder.LoggerValue()
	}
	if builder.DebugValue() != nil {
		eng.debug = *builder.DebugValue()
	}
	eng.collectMode = builder.CollectModeValue()
	eng.stackOverflow = builder.StackOverflowValue()
	eng.cacheThreshold = builder.CacheThresholdValue()

	// Настройка charset
	charset := builder.CharsetValue()
	if charset == "" {
		charset = "UTF-8" // default
	}
	eng.charset = charset

	return eng, nil
}

// ClearCache очищает кэш выражений.
func (e *engine) ClearCache() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cache != nil {
		e.cache.Clear()
	}
}

// CreateExpression компилирует строку в выражение.
func (e *engine) CreateExpression(info *jexl.Info, source string) (jexl.Expression, error) {
	if info == nil {
		info = e.createInfo()
	}

	// Используем парсер для создания AST
	parser := e.getParser()
	if parser == nil {
		return nil, jexl.NewError("parser not available")
	}

	ast, err := parser.ParseExpression(info, source, e.expressionFeatures)
	if err != nil {
		return nil, err
	}

	// Создаём script (который также реализует Expression)
	return NewScript(e, source, ast), nil
}

// CreateScript компилирует строку в исполняемый скрипт.
func (e *engine) CreateScript(features *jexl.Features, info *jexl.Info, source string, names ...string) (jexl.Script, error) {
	if info == nil {
		info = e.createInfo()
	}
	if features == nil {
		features = e.scriptFeatures
	}

	// Используем парсер для создания AST
	parser := e.getParser()
	if parser == nil {
		return nil, jexl.NewError("parser not available")
	}

	ast, err := parser.ParseScript(info, source, features, names)
	if err != nil {
		return nil, err
	}

	return NewScript(e, source, ast), nil
}

// getParser возвращает парсер, создавая его при необходимости.
func (e *engine) getParser() Parser {
	e.mu.RLock()
	p := e.parser
	e.mu.RUnlock()

	if p != nil {
		return p
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Двойная проверка
	if e.parser != nil {
		return e.parser
	}

	if e.parserFactory != nil {
		// Парсер из фабрики должен реализовывать наш интерфейс
		// Пока используем дефолтный
		e.parser = NewDefaultParser()
	} else {
		e.parser = NewDefaultParser()
	}

	return e.parser
}

// CreateTemplateEngine создаёт движок шаблонов JXLT.
func (e *engine) CreateTemplateEngine(opts ...jexl.TemplateOption) (*jexl.TemplateEngine, error) {
	cfg := &jexl.TemplateConfig{
		NoScript:      false,
		CacheSize:     256,
		ImmediateRune: '$',
		DeferredRune:  '#',
	}

	for _, opt := range opts {
		opt.Apply(cfg)
	}

	return jexl.NewTemplateEngine(e, cfg), nil
}

// ThreadContext возвращает текущий thread-local контекст.
func (e *engine) ThreadContext() jexl.ThreadLocalContext {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.threadContext
}

// SetThreadContext устанавливает thread-local контекст.
func (e *engine) SetThreadContext(ctx jexl.ThreadLocalContext) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.threadContext = ctx
}

// Uberspect возвращает introspection-объект.
func (e *engine) Uberspect() jexl.Uberspect {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.uberspect
}

// Options возвращает текущие опции движка.
func (e *engine) Options() *jexl.Options {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.options.Copy()
}

// IsStrict сообщает, включён ли строгий режим.
func (e *engine) IsStrict() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.strict
}

// Arithmetic возвращает арифметику движка.
func (e *engine) Arithmetic() jexl.Arithmetic {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.arithmetic
}

func (e *engine) createInfo() *jexl.Info {
	return jexl.NewInfo()
}

// CreateInfo создаёт Info структуру из текущего стека вызовов.
func (e *engine) CreateInfo() *jexl.Info {
	return jexl.NewInfo()
}

// CreateInfoAt создаёт Info структуру с заданными параметрами.
func (e *engine) CreateInfoAt(name string, line, column int) *jexl.Info {
	return jexl.NewInfoAt(name, line, column)
}

// GetCharset возвращает кодировку, используемую для парсинга.
func (e *engine) GetCharset() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.charset == "" {
		return "UTF-8"
	}
	return e.charset
}

// IsDebug сообщает, включён ли режим отладки.
func (e *engine) IsDebug() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.debug
}

// IsSilent сообщает, включён ли silent режим (ошибки не выбрасываются).
func (e *engine) IsSilent() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.silent
}

// IsCancellable сообщает, будет ли движок выбрасывать исключение при прерывании.
func (e *engine) IsCancellable() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.cancellable
}

// SetClassLoader устанавливает class loader для создания экземпляров по имени класса.
// В Go это не применимо напрямую, но добавлено для совместимости API.
func (e *engine) SetClassLoader(loader any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.classLoader = loader
}

// GetProperty получает свойство объекта по выражению.
func (e *engine) GetProperty(ctx jexl.Context, bean any, expr string) (any, error) {
	if ctx == nil {
		ctx = jexl.EmptyContext{}
	}
	expression, err := e.CreateExpression(nil, expr)
	if err != nil {
		return nil, err
	}

	// Создаём временный контекст с объектом
	tempCtx := &tempContext{
		base:  ctx,
		bean:  bean,
		value: bean,
	}

	return expression.Evaluate(tempCtx)
}

// GetPropertyWithoutContext получает свойство объекта по выражению без контекста.
func (e *engine) GetPropertyWithoutContext(bean any, expr string) (any, error) {
	return e.GetProperty(nil, bean, expr)
}

// SetProperty устанавливает свойство объекта по выражению.
func (e *engine) SetProperty(ctx jexl.Context, bean any, expr string, value any) error {
	// Парсим выражение как присваивание
	// Для простоты поддерживаем только простые пути типа "prop" или "obj.prop"
	if ctx == nil {
		ctx = jexl.EmptyContext{}
	}

	uberspect := e.Uberspect()
	if uberspect == nil {
		return jexl.NewError("uberspect not available")
	}

	// Разбираем выражение на части
	parts := splitPropertyPath(expr)
	if len(parts) == 0 {
		return jexl.NewError("invalid property expression")
	}

	// Получаем объект для установки свойства
	var target any = bean
	for i := 0; i < len(parts)-1; i++ {
		propGet := uberspect.GetProperty(target, parts[i])
		if propGet == nil {
			return jexl.NewError("property not found: " + parts[i])
		}
		var err error
		target, err = propGet.Invoke(target)
		if err != nil {
			return err
		}
		if target == nil {
			return jexl.NewError("intermediate property is nil: " + parts[i])
		}
	}

	// Устанавливаем последнее свойство
	propSet := uberspect.SetProperty(target, parts[len(parts)-1], value)
	if propSet == nil {
		return jexl.NewError("property setter not found: " + parts[len(parts)-1])
	}

	return propSet.Invoke(target, value)
}

// SetPropertyWithoutContext устанавливает свойство объекта по выражению без контекста.
func (e *engine) SetPropertyWithoutContext(bean any, expr string, value any) error {
	return e.SetProperty(nil, bean, expr, value)
}

// InvokeMethod вызывает метод объекта.
func (e *engine) InvokeMethod(obj any, method string, args ...any) (any, error) {
	uberspect := e.Uberspect()
	if uberspect == nil {
		return nil, jexl.NewError("uberspect not available")
	}

	m, err := uberspect.GetMethod(obj, method, args)
	if err != nil {
		if e.strict {
			return nil, err
		}
		return nil, nil
	}

	return m.Invoke(obj, args)
}

// NewInstance создаёт новый экземпляр объекта.
func (e *engine) NewInstance(className string, args ...any) (any, error) {
	// В Go нет прямого способа создать экземпляр по имени класса
	// Это ограничение Go - нужно использовать реестр типов или другие механизмы
	return nil, jexl.NewError("constructor lookup by name not supported in Go")
}

// tempContext временный контекст для вычисления выражений свойств.
type tempContext struct {
	base  jexl.Context
	bean  any
	value any
}

func (t *tempContext) Get(name string) any {
	if name == "this" || name == "it" {
		return t.bean
	}
	if t.base != nil {
		return t.base.Get(name)
	}
	return nil
}

func (t *tempContext) Has(name string) bool {
	if name == "this" || name == "it" {
		return true
	}
	if t.base != nil {
		return t.base.Has(name)
	}
	return false
}

func (t *tempContext) Set(name string, value any) {
	if t.base != nil {
		t.base.Set(name, value)
	}
}

// splitPropertyPath разбивает путь к свойству на части.
func splitPropertyPath(path string) []string {
	return strings.Split(path, ".")
}

// emptyContext больше не используется, используем jexl.EmptyContext

const maxInt = int(^uint(0) >> 1)
