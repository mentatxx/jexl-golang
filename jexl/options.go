package jexl

import (
	"fmt"
	"math"
	"math/big"
	"slices"
)

// MathContext описывает настройки точности и округления для операций с big.Decimal, аналог java.math.MathContext.
type MathContext struct {
	Precision uint
	Rounding  big.RoundingMode
}

// Options инкапсулирует флаги и параметры, влияющие на поведение движка JEXL.
// Структура повторяет семантику org.apache.commons.jexl3.JexlOptions.
type Options struct {
	mathContext      *MathContext
	mathScale        int
	strictArithmetic bool
	flags            uint32
	namespaces       map[string]any
	imports          []string
}

const (
	flagCancellable uint32 = 1 << iota
	flagStrict
	flagSilent
	flagSafe
	flagLexical
	flagAntish
	flagLexicalShade
	flagSharedInstance
	flagConstCapture
	flagStrictInterpolation
	flagBooleanLogical
)

var optionFlagNames = []string{
	"cancellable",
	"strict",
	"silent",
	"safe",
	"lexical",
	"antish",
	"lexicalShade",
	"sharedInstance",
	"constCapture",
	"strictInterpolation",
	"booleanShortCircuit",
}

var defaultOptionFlags uint32 = flagCancellable | flagStrict | flagAntish | flagSafe

// NewOptions создаёт Options с настройками по умолчанию.
func NewOptions() *Options {
	return &Options{
		mathScale:        math.MinInt,
		strictArithmetic: true,
		flags:            defaultOptionFlags,
		namespaces:       map[string]any{},
		imports:          []string{},
	}
}

// Copy создаёт глубокую копию настроек.
func (o *Options) Copy() *Options {
	if o == nil {
		return nil
	}

	cp := *o
	if o.mathContext != nil {
		ctx := *o.mathContext
		cp.mathContext = &ctx
	}
	if len(o.namespaces) > 0 {
		cp.namespaces = make(map[string]any, len(o.namespaces))
		for k, v := range o.namespaces {
			cp.namespaces[k] = v
		}
	} else {
		cp.namespaces = map[string]any{}
	}
	cp.imports = slices.Clone(o.imports)
	return &cp
}

// Set копирует значения из src.
func (o *Options) Set(src *Options) *Options {
	if o == nil || src == nil {
		return o
	}
	*o = *src.Copy()
	return o
}

// Flags возвращает текущее значение битовой маски.
func (o *Options) Flags() uint32 {
	return o.flags
}

// SetFlags применяет набор флагов в формате "+flag" или "-flag".
func (o *Options) SetFlags(values ...string) error {
	mask := o.flags
	for _, value := range values {
		if value == "" {
			continue
		}
		desired := true
		name := value
		switch value[0] {
		case '+':
			name = value[1:]
			desired = true
		case '-':
			name = value[1:]
			desired = false
		}
		index := slices.Index(optionFlagNames, name)
		if index == -1 {
			return fmt.Errorf("jexl: неизвестный флаг опции %q", name)
		}
		bit := uint32(1) << uint(index)
		if desired {
			mask |= bit
		} else {
			mask &^= bit
		}
	}
	o.flags = mask
	return nil
}

// MathContext возвращает текущий math контекст.
func (o *Options) MathContext() *MathContext {
	return o.mathContext
}

// SetMathContext задаёт math контекст.
func (o *Options) SetMathContext(ctx *MathContext) {
	if ctx == nil {
		o.mathContext = nil
		return
	}
	cp := *ctx
	o.mathContext = &cp
}

// MathScale возвращает текущий масштаб.
func (o *Options) MathScale() int {
	return o.mathScale
}

// SetMathScale задаёт масштаб.
func (o *Options) SetMathScale(scale int) {
	o.mathScale = scale
}

// SetStrictArithmetic включает/выключает строгую арифметику.
func (o *Options) SetStrictArithmetic(strict bool) {
	o.strictArithmetic = strict
}

// StrictArithmetic сообщает включена ли строгая арифметика.
func (o *Options) StrictArithmetic() bool {
	return o.strictArithmetic
}

// Imports возвращает список импортов.
func (o *Options) Imports() []string {
	return slices.Clone(o.imports)
}

// SetImports задаёт список импортов.
func (o *Options) SetImports(values []string) {
	o.imports = slices.Clone(values)
}

// Namespaces возвращает карту пространств имён.
func (o *Options) Namespaces() map[string]any {
	if len(o.namespaces) == 0 {
		return map[string]any{}
	}
	res := make(map[string]any, len(o.namespaces))
	for k, v := range o.namespaces {
		res[k] = v
	}
	return res
}

// SetNamespaces задаёт пространства имён.
func (o *Options) SetNamespaces(values map[string]any) {
	if len(values) == 0 {
		o.namespaces = map[string]any{}
		return
	}
	o.namespaces = make(map[string]any, len(values))
	for k, v := range values {
		o.namespaces[k] = v
	}
}

// Utility методы проверки отдельных флагов.

func (o *Options) isSet(flag uint32) bool {
	return o.flags&flag != 0
}

func (o *Options) set(flag uint32, value bool) {
	if value {
		o.flags |= flag
	} else {
		o.flags &^= flag
	}
}

func (o *Options) Antish() bool                { return o.isSet(flagAntish) }
func (o *Options) SetAntish(flag bool)         { o.set(flagAntish, flag) }
func (o *Options) BooleanLogical() bool        { return o.isSet(flagBooleanLogical) }
func (o *Options) SetBooleanLogical(flag bool) { o.set(flagBooleanLogical, flag) }
func (o *Options) Cancellable() bool           { return o.isSet(flagCancellable) }
func (o *Options) SetCancellable(flag bool)    { o.set(flagCancellable, flag) }
func (o *Options) ConstCapture() bool          { return o.isSet(flagConstCapture) }
func (o *Options) SetConstCapture(flag bool)   { o.set(flagConstCapture, flag) }
func (o *Options) Lexical() bool               { return o.isSet(flagLexical) }
func (o *Options) SetLexical(flag bool)        { o.set(flagLexical, flag) }
func (o *Options) LexicalShade() bool          { return o.isSet(flagLexicalShade) }
func (o *Options) SetLexicalShade(flag bool)   { o.set(flagLexicalShade, flag) }
func (o *Options) Safe() bool                  { return o.isSet(flagSafe) }
func (o *Options) SetSafe(flag bool)           { o.set(flagSafe, flag) }
func (o *Options) SharedInstance() bool        { return o.isSet(flagSharedInstance) }
func (o *Options) SetSharedInstance(flag bool) { o.set(flagSharedInstance, flag) }
func (o *Options) Silent() bool                { return o.isSet(flagSilent) }
func (o *Options) SetSilent(flag bool)         { o.set(flagSilent, flag) }
func (o *Options) Strict() bool                { return o.isSet(flagStrict) }
func (o *Options) SetStrict(flag bool)         { o.set(flagStrict, flag) }
func (o *Options) StrictInterpolation() bool   { return o.isSet(flagStrictInterpolation) }
func (o *Options) SetStrictInterpolation(flag bool) {
	o.set(flagStrictInterpolation, flag)
}
