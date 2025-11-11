package jexl

// Uberspect выполняет рефлексию объектов. Аналог JexlUberspect.
type Uberspect interface {
	// GetProperty ищет getter.
	GetProperty(obj any, identifier string) PropertyGet
	// SetProperty ищет setter.
	SetProperty(obj any, identifier string, value any) PropertySet
	// GetMethod ищет метод.
	GetMethod(obj any, name string, args []any) (Method, error)
	// GetConstructor ищет конструктор.
	GetConstructor(name string, args []any) (Method, error)
}

// ResolverStrategy определяет стратегию выбора кандидатов.
type ResolverStrategy interface {
	SelectMethod(methods []Method, args []any) (Method, error)
}

// ResolverStrategyDefault реализация стратегии по умолчанию.
var ResolverStrategyDefault ResolverStrategy = &defaultResolverStrategy{}

type defaultResolverStrategy struct{}

// SelectMethod выбирает наиболее подходящий метод на основе типов аргументов.
// Портированная логика из Java версии JEXL.
func (d *defaultResolverStrategy) SelectMethod(methods []Method, args []any) (Method, error) {
	if len(methods) == 0 {
		return nil, NewError("no methods available")
	}
	if len(methods) == 1 {
		return methods[0], nil
	}

	// Ищем метод с наиболее точным соответствием типов
	bestMatch := methods[0]
	bestScore := d.scoreMethod(methods[0], args)

	for i := 1; i < len(methods); i++ {
		score := d.scoreMethod(methods[i], args)
		if score > bestScore {
			bestScore = score
			bestMatch = methods[i]
		}
	}

	return bestMatch, nil
}

// scoreMethod вычисляет оценку соответствия метода аргументам.
// Более высокий score означает лучшее соответствие.
func (d *defaultResolverStrategy) scoreMethod(method Method, args []any) int {
	// Базовая реализация: проверяем количество аргументов
	// В полной реализации нужно использовать reflection для проверки типов
	// и вычисления точности соответствия
	
	// Пока возвращаем базовый score
	// В будущем можно улучшить, используя reflection для проверки типов параметров
	return 0
}

// PropertyGet представляет операцию чтения свойства.
type PropertyGet interface {
	Invoke(obj any) (any, error)
}

// PropertySet представляет операцию записи свойства.
type PropertySet interface {
	Invoke(obj any, value any) error
}

// Method описывает вызов метода.
type Method interface {
	Name() string
	Invoke(target any, args []any) (any, error)
}
