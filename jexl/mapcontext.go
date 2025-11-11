package jexl

import "sync"

// MapContext оборачивает map в контекст JEXL.
// Каждая пара ключ-значение в map считается переменной.
type MapContext struct {
	mu   sync.RWMutex
	vars map[string]any
}

// NewMapContext создаёт MapContext с автоматически выделенной map.
func NewMapContext() *MapContext {
	return NewMapContextWithMap(nil)
}

// NewMapContextWithMap создаёт MapContext, оборачивающий существующую map.
func NewMapContextWithMap(vars map[string]any) *MapContext {
	if vars == nil {
		vars = make(map[string]any)
	}
	return &MapContext{vars: vars}
}

// Get возвращает значение переменной.
func (m *MapContext) Get(name string) any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.vars[name]
}

// Has проверяет, определена ли переменная.
func (m *MapContext) Has(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.vars[name]
	return ok
}

// Set устанавливает значение переменной.
func (m *MapContext) Set(name string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.vars[name] = value
}

// Clear очищает все переменные.
func (m *MapContext) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k := range m.vars {
		delete(m.vars, k)
	}
}

// Vars возвращает копию внутренней map переменных.
func (m *MapContext) Vars() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]any, len(m.vars))
	for k, v := range m.vars {
		result[k] = v
	}
	return result
}
