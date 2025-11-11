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
