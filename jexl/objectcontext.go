package jexl

// ObjectContext оборачивает объект как контекст JEXL и NamespaceResolver.
type ObjectContext struct {
	engine Engine
	object any
}

// NewObjectContext создаёт новый ObjectContext.
func NewObjectContext(engine Engine, wrapped any) *ObjectContext {
	return &ObjectContext{
		engine: engine,
		object: wrapped,
	}
}

// Get возвращает значение свойства через introspection.
func (o *ObjectContext) Get(name string) any {
	uberspect := o.engine.Uberspect()
	if uberspect == nil {
		return nil
	}
	propGet := uberspect.GetProperty(o.object, name)
	if propGet == nil {
		return nil
	}
	value, err := propGet.Invoke(o.object)
	if err != nil {
		opts := o.engine.Options()
		if opts != nil && opts.Strict() {
			// В строгом режиме возвращаем ошибку, но это нарушает интерфейс Context
			// В Java это выбрасывает исключение, в Go мы возвращаем nil
			return nil
		}
		return nil
	}
	return value
}

// Has проверяет, доступно ли свойство.
func (o *ObjectContext) Has(name string) bool {
	uberspect := o.engine.Uberspect()
	if uberspect == nil {
		return false
	}
	return uberspect.GetProperty(o.object, name) != nil
}

// Set устанавливает значение свойства через introspection.
func (o *ObjectContext) Set(name string, value any) {
	uberspect := o.engine.Uberspect()
	if uberspect == nil {
		return
	}
	propSet := uberspect.SetProperty(o.object, name, value)
	if propSet == nil {
		return
	}
	if err := propSet.Invoke(o.object, value); err != nil {
		opts := o.engine.Options()
		if opts != nil && opts.Strict() {
			// В строгом режиме это должно выбрасывать исключение
			// В Go мы просто игнорируем ошибку, так как Set не возвращает error
		}
	}
}

// ResolveNamespace реализует NamespaceResolver.
func (o *ObjectContext) ResolveNamespace(name string) any {
	if name == "" {
		return o.object
	}
	return nil
}

// Object возвращает обёрнутый объект.
func (o *ObjectContext) Object() any {
	return o.object
}
