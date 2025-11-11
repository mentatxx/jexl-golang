package jexl

import "slices"

// Feature представляет отдельный синтаксический флаг.
type Feature int

const (
	FeatureRegister Feature = iota
	FeatureReserved
	FeatureLocalVar
	FeatureSideEffect
	FeatureSideEffectGlobal
	FeatureArrayReferenceExpr
	FeatureNewInstance
	FeatureLoop
	FeatureLambda
	FeatureMethodCall
	FeatureStructuredLiteral
	FeaturePragma
	FeatureAnnotation
	FeatureScript
	FeatureLexical
	FeatureLexicalShade
	FeatureThinArrow
	FeatureFatArrow
	FeatureNamespacePragma
	FeatureNamespaceIdentifier
	FeatureImportPragma
	FeatureComparatorNames
	FeaturePragmaAnywhere
	FeatureConstCapture
	FeatureReferenceCapture
	FeatureAmbiguousStatement
)

// Features задаёт набор разрешённых синтаксических конструкций.
type Features struct {
	flags        uint64
	reserved     []string
	namespaceSet []string
}

// NewFeatures создаёт Features с заданными флагами.
func NewFeatures(enabled ...Feature) *Features {
	f := &Features{}
	for _, feat := range enabled {
		f.Enable(feat, true)
	}
	return f
}

// Enable включает или выключает конкретный флаг.
func (f *Features) Enable(feature Feature, value bool) {
	if f == nil {
		return
	}
	mask := uint64(1) << uint(feature)
	if value {
		f.flags |= mask
	} else {
		f.flags &^= mask
	}
}

// Enabled проверяет включён ли флаг.
func (f *Features) Enabled(feature Feature) bool {
	if f == nil {
		return false
	}
	mask := uint64(1) << uint(feature)
	return f.flags&mask != 0
}

// SetReservedNames задаёт зарезервированные идентификаторы.
func (f *Features) SetReservedNames(names []string) {
	if f == nil {
		return
	}
	f.reserved = slices.Clone(names)
}

// ReservedNames возвращает список зарезервированных имён.
func (f *Features) ReservedNames() []string {
	if f == nil {
		return nil
	}
	return slices.Clone(f.reserved)
}

// FeaturesDefault создаёт набор features по умолчанию.
func FeaturesDefault() *Features {
	f := NewFeatures(
		FeatureLocalVar,
		FeatureSideEffect,
		FeatureSideEffectGlobal,
		FeatureArrayReferenceExpr,
		FeatureNewInstance,
		FeatureLoop,
		FeatureLambda,
		FeatureMethodCall,
		FeatureStructuredLiteral,
		FeaturePragma,
		FeatureAnnotation,
		FeatureScript,
		FeatureThinArrow,
		FeatureNamespacePragma,
		FeatureImportPragma,
		FeatureComparatorNames,
		FeaturePragmaAnywhere,
	)
	return f
}

// With включает указанные features и возвращает новый объект.
func (f *Features) With(features ...Feature) *Features {
	if f == nil {
		f = &Features{}
	}
	result := &Features{
		flags:        f.flags,
		reserved:     slices.Clone(f.reserved),
		namespaceSet: slices.Clone(f.namespaceSet),
	}
	for _, feat := range features {
		result.Enable(feat, true)
	}
	return result
}

// Without выключает указанные features и возвращает новый объект.
func (f *Features) Without(features ...Feature) *Features {
	if f == nil {
		f = &Features{}
	}
	result := &Features{
		flags:        f.flags,
		reserved:     slices.Clone(f.reserved),
		namespaceSet: slices.Clone(f.namespaceSet),
	}
	for _, feat := range features {
		result.Enable(feat, false)
	}
	return result
}

// SupportsExpression проверяет поддержку выражений (всегда true).
func (f *Features) SupportsExpression() bool {
	return true // Выражения всегда поддерживаются
}

// SupportsScript проверяет поддержку скриптов.
func (f *Features) SupportsScript() bool {
	return f.Enabled(FeatureScript)
}

// SupportsLoops проверяет поддержку циклов.
func (f *Features) SupportsLoops() bool {
	return f.Enabled(FeatureLoop)
}

// SupportsLocalVar проверяет поддержку локальных переменных.
func (f *Features) SupportsLocalVar() bool {
	return f.Enabled(FeatureLocalVar)
}

// SupportsLambda проверяет поддержку lambda функций.
func (f *Features) SupportsLambda() bool {
	return f.Enabled(FeatureLambda)
}

// SupportsMethodCall проверяет поддержку вызовов методов.
func (f *Features) SupportsMethodCall() bool {
	return f.Enabled(FeatureMethodCall)
}

// SupportsNewInstance проверяет поддержку создания экземпляров.
func (f *Features) SupportsNewInstance() bool {
	return f.Enabled(FeatureNewInstance)
}

// SupportsStructuredLiteral проверяет поддержку структурированных литералов.
func (f *Features) SupportsStructuredLiteral() bool {
	return f.Enabled(FeatureStructuredLiteral)
}

// SupportsAnnotation проверяет поддержку аннотаций.
func (f *Features) SupportsAnnotation() bool {
	return f.Enabled(FeatureAnnotation)
}

// SupportsPragma проверяет поддержку pragma директив.
func (f *Features) SupportsPragma() bool {
	return f.Enabled(FeaturePragma)
}

// IsLexical проверяет включен ли lexical scope.
func (f *Features) IsLexical() bool {
	return f.Enabled(FeatureLexical)
}

// IsLexicalShade проверяет включен ли lexical shade.
func (f *Features) IsLexicalShade() bool {
	return f.Enabled(FeatureLexicalShade)
}

// SupportsConstCapture проверяет поддержку const capture.
func (f *Features) SupportsConstCapture() bool {
	return f.Enabled(FeatureConstCapture)
}

// SupportsComparatorNames проверяет поддержку именованных операторов сравнения (eq, ne, lt, и т.д.).
func (f *Features) SupportsComparatorNames() bool {
	return f.Enabled(FeatureComparatorNames)
}

// SupportsPragmaAnywhere проверяет поддержку pragma везде.
func (f *Features) SupportsPragmaAnywhere() bool {
	return f.Enabled(FeaturePragmaAnywhere)
}
