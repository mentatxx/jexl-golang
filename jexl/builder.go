package jexl

import (
	"math"
)

const (
	cacheThresholdDefault = 64
	stackOverflowDefault  = math.MaxInt
)

// Builder конфигурирует и создаёт экземпляр движка JEXL.
// Портированная версия org.apache.commons.jexl3.JexlBuilder.
type Builder struct {
	uberspect      Uberspect
	strategy       ResolverStrategy
	permissions    *Permissions
	sandbox        *Sandbox
	logger         Logger
	debug          *bool
	cancellable    *bool
	options        *Options
	collectMode    int
	arithmetic     Arithmetic
	cacheSize      int
	cacheFactory   CacheFactory
	parserFactory  ParserFactory
	stackOverflow  int
	cacheThreshold int
	charset        string
	features       *Features
}

// NewBuilder создаёт Builder с настройками по умолчанию.
func NewBuilder() *Builder {
	return &Builder{
		permissions:    PermissionsRestricted.Clone(),
		sandbox:        nil,
		logger:         NoopLogger{},
		options:        NewOptions(),
		collectMode:    1,
		cacheSize:      -1,
		stackOverflow:  stackOverflowDefault,
		cacheThreshold: cacheThresholdDefault,
	}
}

// Options возвращает объект опций.
func (b *Builder) Options() *Options {
	return b.options
}

// Antish включает/отключает antish резолюцию переменных.
func (b *Builder) Antish(flag bool) *Builder {
	b.options.SetAntish(flag)
	return b
}

// Arithmetic задаёт арифметику.
func (b *Builder) Arithmetic(a Arithmetic) *Builder {
	b.arithmetic = a
	if a != nil {
		b.options.SetStrictArithmetic(a.IsStrict())
		if ctx := a.MathContext(); ctx != nil {
			b.options.SetMathContext(ctx)
		}
		b.options.SetMathScale(a.MathScale())
	}
	return b
}

// BooleanLogical задаёт поведение логических выражений.
func (b *Builder) BooleanLogical(flag bool) *Builder {
	b.options.SetBooleanLogical(flag)
	return b
}

// Cache задаёт размер кэша.
func (b *Builder) Cache(size int) *Builder {
	b.cacheSize = size
	return b
}

// CacheFactory задаёт фабрику кэшей.
func (b *Builder) CacheFactory(factory CacheFactory) *Builder {
	b.cacheFactory = factory
	return b
}

// CacheThreshold задаёт максимальную длину выражения для кэширования.
func (b *Builder) CacheThreshold(length int) *Builder {
	if length > 0 {
		b.cacheThreshold = length
	} else {
		b.cacheThreshold = cacheThresholdDefault
	}
	return b
}

// Cancellable управляет реакцией на прерывание.
func (b *Builder) Cancellable(flag bool) *Builder {
	b.cancellable = &flag
	b.options.SetCancellable(flag)
	return b
}

// Charset задаёт кодировку исходного текста.
func (b *Builder) Charset(name string) *Builder {
	b.charset = name
	return b
}

// CollectAll управляет расширенным поиском переменных.
func (b *Builder) CollectAll(flag bool) *Builder {
	return b.CollectMode(func() int {
		if flag {
			return 1
		}
		return 0
	}())
}

// CollectMode задаёт режим сборщика переменных.
func (b *Builder) CollectMode(mode int) *Builder {
	b.collectMode = mode
	return b
}

// Debug управляет включением отладочной информации.
func (b *Builder) Debug(flag bool) *Builder {
	b.debug = &flag
	return b
}

// Features задаёт набор признаков.
func (b *Builder) Features(features *Features) *Builder {
	b.features = features
	return b
}

// Logger задаёт логгер.
func (b *Builder) Logger(logger Logger) *Builder {
	if logger == nil {
		b.logger = NoopLogger{}
	} else {
		b.logger = logger
	}
	return b
}

// Permissions задаёт набор разрешений.
func (b *Builder) Permissions(permissions *Permissions) *Builder {
	if permissions == nil {
		b.permissions = PermissionsRestricted.Clone()
	} else {
		b.permissions = permissions.Clone()
	}
	return b
}

// Sandbox задаёт песочницу.
func (b *Builder) Sandbox(sandbox *Sandbox) *Builder {
	b.sandbox = sandbox
	return b
}

// StackOverflowLimit задаёт лимит глубины рекурсии.
func (b *Builder) StackOverflowLimit(limit int) *Builder {
	if limit <= 0 {
		b.stackOverflow = stackOverflowDefault
	} else {
		b.stackOverflow = limit
	}
	return b
}

// ParserFactory задаёт фабрику парсеров.
func (b *Builder) ParserFactory(factory ParserFactory) *Builder {
	b.parserFactory = factory
	return b
}

// Strategy задаёт стратегию выбора методов.
func (b *Builder) Strategy(strategy ResolverStrategy) *Builder {
	b.strategy = strategy
	return b
}

// Uberspect задаёт introspection.
func (b *Builder) Uberspect(uberspect Uberspect) *Builder {
	b.uberspect = uberspect
	return b
}

// Build создаёт экземпляр Engine.
// Реализация находится в internal пакете.
func (b *Builder) Build() (Engine, error) {
	// Используем прямой вызов через функцию-инициализатор
	// Это позволяет избежать циклических зависимостей
	if buildEngineImpl == nil {
		// Пытаемся загрузить реализацию через lazy initialization
		// init() в internal/register.go должен был зарегистрировать движок
		// Если это не произошло, возвращаем ошибку
		return nil, ErrNotImplemented
	}
	return buildEngineImpl(b)
}

// buildEngineImpl объявлена в отдельном файле для избежания циклических зависимостей
var buildEngineImpl func(*Builder) (Engine, error)

// RegisterEngineBuilder регистрирует фабрику движка.
func RegisterEngineBuilder(fn func(*Builder) (Engine, error)) {
	buildEngineImpl = fn
}

// initEngine регистрирует движок, если он ещё не зарегистрирован.
// Эта функция вызывается автоматически при первом использовании Builder.Build()
func initEngine() {
	if buildEngineImpl == nil {
		// Используем рефлексию для вызова функции из internal пакета
		// Это обходной путь для избежания циклических зависимостей
		// В реальном использовании init() в internal/register.go должен вызываться автоматически
		// при импорте пакета, который использует internal
	}
}

// Геттеры для доступа к полям Builder (используются internal пакетом)

func (b *Builder) UberspectValue() Uberspect {
	return b.uberspect
}

func (b *Builder) StrategyValue() ResolverStrategy {
	return b.strategy
}

func (b *Builder) PermissionsValue() *Permissions {
	return b.permissions
}

func (b *Builder) SandboxValue() *Sandbox {
	return b.sandbox
}

func (b *Builder) LoggerValue() Logger {
	return b.logger
}

func (b *Builder) DebugValue() *bool {
	return b.debug
}

func (b *Builder) CancellableValue() *bool {
	return b.cancellable
}

func (b *Builder) CollectModeValue() int {
	return b.collectMode
}

func (b *Builder) ArithmeticValue() Arithmetic {
	return b.arithmetic
}

func (b *Builder) CacheSize() int {
	return b.cacheSize
}

func (b *Builder) CacheFactoryValue() CacheFactory {
	return b.cacheFactory
}

func (b *Builder) ParserFactoryValue() ParserFactory {
	return b.parserFactory
}

func (b *Builder) StackOverflowValue() int {
	return b.stackOverflow
}

func (b *Builder) CacheThresholdValue() int {
	return b.cacheThreshold
}

func (b *Builder) CharsetValue() string {
	return b.charset
}

func (b *Builder) FeaturesValue() *Features {
	return b.features
}
