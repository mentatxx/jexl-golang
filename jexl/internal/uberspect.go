package internal

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mentatxx/jexl-golang/jexl"
)

// NewUberspect создаёт новый Uberspect с заданными параметрами.
func NewUberspect(logger jexl.Logger, strategy jexl.ResolverStrategy, permissions *jexl.Permissions) jexl.Uberspect {
	// TODO: реализовать полноценный Uberspect
	return &uberspectImpl{
		logger:      logger,
		strategy:    strategy,
		permissions: permissions,
	}
}

// NewSandboxUberspect создаёт Uberspect с песочницей.
func NewSandboxUberspect(base jexl.Uberspect, sandbox *jexl.Sandbox) jexl.Uberspect {
	if sandbox == nil {
		return base
	}
	return &sandboxUberspect{
		base:    base,
		sandbox: sandbox,
	}
}

type sandboxUberspect struct {
	base    jexl.Uberspect
	sandbox *jexl.Sandbox
}

func (s *sandboxUberspect) getClassName(obj any) string {
	if obj == nil {
		return ""
	}
	typ := reflect.TypeOf(obj)
	if typ == nil {
		return ""
	}
	// Для указателей получаем базовый тип
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.String()
}

func (s *sandboxUberspect) GetProperty(obj any, identifier string) jexl.PropertyGet {
	if obj == nil {
		return nil
	}
	className := s.getClassName(obj)
	if !s.sandbox.Allowed(className, identifier) {
		return nil
	}
	return s.base.GetProperty(obj, identifier)
}

func (s *sandboxUberspect) SetProperty(obj any, identifier string, value any) jexl.PropertySet {
	if obj == nil {
		return nil
	}
	className := s.getClassName(obj)
	if !s.sandbox.Allowed(className, identifier) {
		return nil
	}
	return s.base.SetProperty(obj, identifier, value)
}

func (s *sandboxUberspect) GetMethod(obj any, name string, args []any) (jexl.Method, error) {
	if obj == nil {
		return nil, jexl.NewError("cannot get method on nil")
	}
	className := s.getClassName(obj)
	if !s.sandbox.Allowed(className, name) {
		return nil, jexl.NewError(fmt.Sprintf("method %s not allowed for class %s", name, className))
	}
	return s.base.GetMethod(obj, name, args)
}

func (s *sandboxUberspect) GetConstructor(name string, args []any) (jexl.Method, error) {
	// Проверяем разрешения для конструктора
	if !s.sandbox.Allowed(name, "") {
		return nil, jexl.NewError(fmt.Sprintf("constructor for %s not allowed", name))
	}
	return s.base.GetConstructor(name, args)
}

type uberspectImpl struct {
	logger      jexl.Logger
	strategy    jexl.ResolverStrategy
	permissions *jexl.Permissions
}

func (u *uberspectImpl) GetProperty(obj any, identifier string) jexl.PropertyGet {
	if obj == nil {
		return nil
	}

	val := reflect.ValueOf(obj)
	if !val.IsValid() {
		return nil
	}

	// Для мапов
	if val.Kind() == reflect.Map {
		keyVal := reflect.ValueOf(identifier)
		if val.Type().Key() == keyVal.Type() {
			return &mapPropertyGet{mapVal: val, key: keyVal}
		}
	}

	// Для указателей, разыменовываем
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	// Пробуем поле
	if field := val.FieldByName(identifier); field.IsValid() && field.CanInterface() {
		return &fieldPropertyGet{field: field}
	}

	// Пробуем метод Getter (GetXxx)
	getterName := "Get" + strings.ToUpper(identifier[:1]) + identifier[1:]
	if method := val.MethodByName(getterName); method.IsValid() && method.Type().NumIn() == 0 {
		return &methodPropertyGet{method: method}
	}

	// Пробуем метод с именем идентификатора (без параметров)
	if method := val.MethodByName(identifier); method.IsValid() && method.Type().NumIn() == 0 {
		return &methodPropertyGet{method: method}
	}

	return nil
}

func (u *uberspectImpl) SetProperty(obj any, identifier string, value any) jexl.PropertySet {
	if obj == nil {
		return nil
	}

	val := reflect.ValueOf(obj)
	if !val.IsValid() {
		return nil
	}

	// Для мапов
	if val.Kind() == reflect.Map {
		keyVal := reflect.ValueOf(identifier)
		if val.Type().Key() == keyVal.Type() {
			return &mapPropertySet{mapVal: val, key: keyVal, valueType: val.Type().Elem()}
		}
	}

	// Для указателей, разыменовываем
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	// Пробуем поле
	if field := val.FieldByName(identifier); field.IsValid() && field.CanSet() {
		return &fieldPropertySet{field: field, valueType: field.Type()}
	}

	// Пробуем метод Setter (SetXxx)
	setterName := "Set" + strings.ToUpper(identifier[:1]) + identifier[1:]
	if method := val.MethodByName(setterName); method.IsValid() && method.Type().NumIn() == 1 {
		return &methodPropertySet{method: method, valueType: method.Type().In(0)}
	}

	return nil
}

func (u *uberspectImpl) GetMethod(obj any, name string, args []any) (jexl.Method, error) {
	if obj == nil {
		return nil, jexl.NewError("cannot get method on nil")
	}

	val := reflect.ValueOf(obj)
	if !val.IsValid() {
		return nil, jexl.NewError("invalid object")
	}

	// Для указателей, разыменовываем для получения типа
	typ := val.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Собираем все кандидаты методов
	var candidates []reflect.Method
	
	// Ищем методы на значении
	for i := 0; i < val.Type().NumMethod(); i++ {
		m := val.Type().Method(i)
		if m.Name == name {
			candidates = append(candidates, m)
		}
	}
	
	// Ищем методы на типе (для указателей)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		elemType := val.Elem().Type()
		for i := 0; i < elemType.NumMethod(); i++ {
			m := elemType.Method(i)
			if m.Name == name {
				// Проверяем, что метод не был уже добавлен
				found := false
				for _, c := range candidates {
					if c.Name == m.Name && c.Type == m.Type {
						found = true
						break
					}
				}
				if !found {
					candidates = append(candidates, m)
				}
			}
		}
	}

	if len(candidates) == 0 {
		return nil, jexl.NewError(fmt.Sprintf("method %s not found", name))
	}

	// Выбираем наиболее подходящий метод
	bestMethod, err := u.selectBestMethod(candidates, val, args)
	if err != nil {
		return nil, err
	}

	return &reflectionMethod{
		method: bestMethod,
		name:   name,
	}, nil
}

// selectBestMethod выбирает наиболее подходящий метод из кандидатов.
func (u *uberspectImpl) selectBestMethod(candidates []reflect.Method, val reflect.Value, args []any) (reflect.Value, error) {
	var bestMethod reflect.Method
	bestScore := -1

	for _, candidate := range candidates {
		score := u.scoreMethod(candidate, args)
		if score > bestScore {
			bestScore = score
			bestMethod = candidate
		}
	}

	if bestScore < 0 {
		return reflect.Value{}, jexl.NewError("no compatible method found")
	}

	// Получаем reflect.Value метода
	methodVal := val.MethodByName(bestMethod.Name)
	if !methodVal.IsValid() && val.Kind() == reflect.Ptr && !val.IsNil() {
		methodVal = val.Elem().MethodByName(bestMethod.Name)
	}

	return methodVal, nil
}

// scoreMethod вычисляет оценку соответствия метода аргументам.
// Возвращает -1 если метод несовместим, иначе положительное число.
func (u *uberspectImpl) scoreMethod(method reflect.Method, args []any) int {
	methodType := method.Type
	
	// Проверяем количество аргументов (первый аргумент - receiver)
	numIn := methodType.NumIn() - 1 // -1 для receiver
	if numIn != len(args) {
		// Проверяем variadic методы
		if !methodType.IsVariadic() || numIn > len(args)+1 {
			return -1
		}
	}

	score := 0
	argIdx := 0
	
	// Пропускаем receiver (индекс 0)
	for i := 1; i < methodType.NumIn(); i++ {
		paramType := methodType.In(i)
		
		if methodType.IsVariadic() && i == methodType.NumIn()-1 {
			// Variadic параметр
			elemType := paramType.Elem()
			for argIdx < len(args) {
				if !u.isAssignable(args[argIdx], elemType) {
					return -1
				}
				score += u.typeMatchScore(args[argIdx], elemType)
				argIdx++
			}
		} else {
			if argIdx >= len(args) {
				return -1
			}
			if !u.isAssignable(args[argIdx], paramType) {
				return -1
			}
			score += u.typeMatchScore(args[argIdx], paramType)
			argIdx++
		}
	}

	return score
}

// isAssignable проверяет, можно ли присвоить значение типу.
func (u *uberspectImpl) isAssignable(value any, targetType reflect.Type) bool {
	if value == nil {
		// nil можно присвоить любому указателю, интерфейсу, slice, map, func, channel
		switch targetType.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
			return true
		default:
			return false
		}
	}

	valType := reflect.TypeOf(value)
	
	// Прямое присваивание
	if valType.AssignableTo(targetType) {
		return true
	}
	
	// Преобразование типов
	if valType.ConvertibleTo(targetType) {
		return true
	}
	
	// Специальные случаи для числовых типов
	if u.isNumeric(valType) && u.isNumeric(targetType) {
		return true
	}
	
	return false
}

// typeMatchScore вычисляет оценку соответствия типов (чем выше, тем лучше).
func (u *uberspectImpl) typeMatchScore(value any, targetType reflect.Type) int {
	if value == nil {
		return 0
	}

	valType := reflect.TypeOf(value)
	
	// Точное совпадение
	if valType == targetType {
		return 100
	}
	
	// Прямое присваивание
	if valType.AssignableTo(targetType) {
		return 80
	}
	
	// Преобразование типов
	if valType.ConvertibleTo(targetType) {
		return 60
	}
	
	// Числовые типы
	if u.isNumeric(valType) && u.isNumeric(targetType) {
		return 40
	}
	
	return 20
}

// isNumeric проверяет, является ли тип числовым.
func (u *uberspectImpl) isNumeric(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func (u *uberspectImpl) GetConstructor(name string, args []any) (jexl.Method, error) {
	// В Go нет прямого способа получить тип по имени строки без использования реестра типов
	// Это ограничение Go - нужно использовать реестр или другие механизмы
	return nil, jexl.NewError("constructor lookup by name not supported in Go")
}

// Реализации PropertyGet

type fieldPropertyGet struct {
	field reflect.Value
}

func (f *fieldPropertyGet) Invoke(obj any) (any, error) {
	if !f.field.IsValid() {
		return nil, jexl.NewError("invalid field")
	}
	return f.field.Interface(), nil
}

type methodPropertyGet struct {
	method reflect.Value
}

func (m *methodPropertyGet) Invoke(obj any) (any, error) {
	results := m.method.Call(nil)
	if len(results) == 0 {
		return nil, nil
	}
	if len(results) == 1 {
		return results[0].Interface(), nil
	}
	// Множественные возвращаемые значения
	values := make([]any, len(results))
	for i, r := range results {
		values[i] = r.Interface()
	}
	return values, nil
}

type mapPropertyGet struct {
	mapVal reflect.Value
	key    reflect.Value
}

func (m *mapPropertyGet) Invoke(obj any) (any, error) {
	result := m.mapVal.MapIndex(m.key)
	if !result.IsValid() {
		return nil, nil
	}
	return result.Interface(), nil
}

// Реализации PropertySet

type fieldPropertySet struct {
	field     reflect.Value
	valueType reflect.Type
}

func (f *fieldPropertySet) Invoke(obj any, value any) error {
	if !f.field.IsValid() || !f.field.CanSet() {
		return jexl.NewError("field is not settable")
	}

	val := reflect.ValueOf(value)
	if !val.IsValid() {
		f.field.Set(reflect.Zero(f.field.Type()))
		return nil
	}

	if val.Type().AssignableTo(f.field.Type()) {
		f.field.Set(val)
		return nil
	}

	if val.Type().ConvertibleTo(f.field.Type()) {
		f.field.Set(val.Convert(f.field.Type()))
		return nil
	}

	return jexl.NewError(fmt.Sprintf("cannot assign %s to %s", val.Type(), f.field.Type()))
}

type methodPropertySet struct {
	method    reflect.Value
	valueType reflect.Type
}

func (m *methodPropertySet) Invoke(obj any, value any) error {
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		val = reflect.Zero(m.valueType)
	}

	if !val.Type().AssignableTo(m.valueType) {
		if val.Type().ConvertibleTo(m.valueType) {
			val = val.Convert(m.valueType)
		} else {
			return jexl.NewError(fmt.Sprintf("cannot assign %s to %s", val.Type(), m.valueType))
		}
	}

	results := m.method.Call([]reflect.Value{val})
	if len(results) > 0 {
		if err, ok := results[0].Interface().(error); ok && err != nil {
			return err
		}
	}
	return nil
}

type mapPropertySet struct {
	mapVal    reflect.Value
	key       reflect.Value
	valueType reflect.Type
}

func (m *mapPropertySet) Invoke(obj any, value any) error {
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		val = reflect.Zero(m.valueType)
	}

	if !val.Type().AssignableTo(m.valueType) {
		if val.Type().ConvertibleTo(m.valueType) {
			val = val.Convert(m.valueType)
		} else {
			return jexl.NewError(fmt.Sprintf("cannot assign %s to %s", val.Type(), m.valueType))
		}
	}

	m.mapVal.SetMapIndex(m.key, val)
	return nil
}

// Реализация Method

type reflectionMethod struct {
	method reflect.Value
	name   string
}

func (r *reflectionMethod) Name() string {
	return r.name
}

func (r *reflectionMethod) Invoke(target any, args []any) (any, error) {
	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValues[i] = reflect.ValueOf(arg)
	}

	results := r.method.Call(argValues)
	if len(results) == 0 {
		return nil, nil
	}
	if len(results) == 1 {
		result := results[0].Interface()
		if err, ok := result.(error); ok {
			return nil, err
		}
		return result, nil
	}

	// Множественные возвращаемые значения
	values := make([]any, len(results))
	for i, r := range results {
		values[i] = r.Interface()
	}
	// Последнее значение может быть error
	if len(values) > 0 {
		if err, ok := values[len(values)-1].(error); ok && err != nil {
			return nil, err
		}
	}
	return values, nil
}
