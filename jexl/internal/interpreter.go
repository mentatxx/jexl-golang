package internal

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/mentatxx/jexl-golang/jexl"
)

// interpreter выполняет AST узлы.
// Порт org.apache.commons.jexl3.internal.Interpreter.
type interpreter struct {
	engine  jexl.Engine
	context jexl.Context
	options *jexl.Options
}

// newInterpreter создаёт новый интерпретатор.
func newInterpreter(engine jexl.Engine, ctx jexl.Context) *interpreter {
	opts := engine.Options()
	return &interpreter{
		engine:  engine,
		context: ctx,
		options: opts,
	}
}

// interpret выполняет AST узел.
func (i *interpreter) interpret(node jexl.Node) (any, error) {
	switch n := node.(type) {
	case *jexl.ScriptNode:
		return i.interpretScript(n)
	case *jexl.LiteralNode:
		return n.Value(), nil
	case *jexl.IdentifierNode:
		return i.interpretIdentifier(n)
	case *jexl.BinaryOpNode:
		return i.interpretBinaryOp(n)
	case *jexl.UnaryOpNode:
		return i.interpretUnaryOp(n)
	case *jexl.PropertyAccessNode:
		return i.interpretPropertyAccess(n)
	case *jexl.IndexAccessNode:
		return i.interpretIndexAccess(n)
	case *jexl.MethodCallNode:
		return i.interpretMethodCall(n)
	case *jexl.AssignmentNode:
		return i.interpretAssignment(n)
	case *jexl.TernaryNode:
		return i.interpretTernary(n)
	case *jexl.ElvisNode:
		return i.interpretElvis(n)
	case *jexl.RangeNode:
		return i.interpretRange(n)
	case *jexl.ArrayLiteralNode:
		return i.interpretArrayLiteral(n)
	case *jexl.MapLiteralNode:
		return i.interpretMapLiteral(n)
	case *jexl.SetLiteralNode:
		return i.interpretSetLiteral(n)
	case *jexl.IfNode:
		return i.interpretIf(n)
	case *jexl.ForNode:
		return i.interpretFor(n)
	case *jexl.ForeachNode:
		return i.interpretForeach(n)
	case *jexl.WhileNode:
		return i.interpretWhile(n)
	case *jexl.DoWhileNode:
		return i.interpretDoWhile(n)
	case *jexl.BlockNode:
		return i.interpretBlock(n)
	case *jexl.BreakNode:
		return nil, &BreakError{}
	case *jexl.ContinueNode:
		return nil, &ContinueError{}
	case *jexl.ReturnNode:
		return i.interpretReturn(n)
	default:
		return nil, jexl.NewError("unsupported node type")
	}
}

// interpretScript выполняет ScriptNode.
func (i *interpreter) interpretScript(node *jexl.ScriptNode) (any, error) {
	var result any
	var err error

	for _, child := range node.Children() {
		result, err = i.interpret(child)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// interpretIdentifier выполняет IdentifierNode.
func (i *interpreter) interpretIdentifier(node *jexl.IdentifierNode) (any, error) {
	if i.context == nil {
		return nil, jexl.NewError("context is nil")
	}

	if !i.context.Has(node.Name()) {
		if i.options != nil && i.options.Strict() {
			return nil, jexl.WrapError("variable not found", nil, node.Info())
		}
		return nil, nil
	}

	return i.context.Get(node.Name()), nil
}

// interpretBinaryOp выполняет BinaryOpNode.
func (i *interpreter) interpretBinaryOp(node *jexl.BinaryOpNode) (any, error) {
	left, err := i.interpret(node.Left())
	if err != nil {
		return nil, err
	}

	right, err := i.interpret(node.Right())
	if err != nil {
		return nil, err
	}

	// Получаем арифметику из движка
	arithmetic := i.engine.Arithmetic()
	if arithmetic == nil {
		// Используем базовую арифметику
		arithmetic = jexl.NewBaseArithmetic(true, nil, 0)
	}

	// Выполняем операцию в зависимости от типа
	switch node.Op() {
	case "+":
		return arithmetic.Add(left, right)
	case "-":
		return arithmetic.Subtract(left, right)
	case "*":
		return arithmetic.Multiply(left, right)
	case "/":
		return arithmetic.Divide(left, right)
	case "%":
		return arithmetic.Modulo(left, right)
	case "==", "eq":
		cmp, err := arithmetic.Compare(left, right)
		if err != nil {
			return nil, err
		}
		return cmp == 0, nil
	case "!=", "ne":
		cmp, err := arithmetic.Compare(left, right)
		if err != nil {
			return nil, err
		}
		return cmp != 0, nil
	case "<", "lt":
		cmp, err := arithmetic.Compare(left, right)
		if err != nil {
			return nil, err
		}
		return cmp < 0, nil
	case ">", "gt":
		cmp, err := arithmetic.Compare(left, right)
		if err != nil {
			return nil, err
		}
		return cmp > 0, nil
	case "<=", "le":
		cmp, err := arithmetic.Compare(left, right)
		if err != nil {
			return nil, err
		}
		return cmp <= 0, nil
	case ">=", "ge":
		cmp, err := arithmetic.Compare(left, right)
		if err != nil {
			return nil, err
		}
		return cmp >= 0, nil
	case "&&", "and":
		leftBool, err := arithmetic.ToBoolean(left)
		if err != nil {
			return nil, err
		}
		if !leftBool {
			return false, nil
		}
		return arithmetic.ToBoolean(right)
	case "||", "or":
		leftBool, err := arithmetic.ToBoolean(left)
		if err != nil {
			return nil, err
		}
		if leftBool {
			return true, nil
		}
		return arithmetic.ToBoolean(right)
	case "&":
		return arithmetic.BitwiseAnd(left, right)
	case "|":
		return arithmetic.BitwiseOr(left, right)
	case "^":
		return arithmetic.BitwiseXor(left, right)
	case "<<":
		return arithmetic.ShiftLeft(left, right)
	case ">>":
		return arithmetic.ShiftRight(left, right)
	case ">>>":
		return arithmetic.ShiftRightUnsigned(left, right)
	case "=~":
		// В JEXL x =~ y:
		// - Если y - строка: проверяем соответствие x регулярному выражению y
		// - Если y - коллекция: проверяем, содержится ли x в y
		// Для строк: left =~ right означает "left соответствует right"
		// Для коллекций: left =~ right означает "left содержится в right"
		if _, ok := right.(string); ok {
			// Строковый паттерн - проверяем соответствие left паттерну right
			return arithmetic.Contains(left, right)
		}
		// Коллекция - проверяем, содержится ли left в right
		return arithmetic.Contains(right, left)
	case "=^":
		return arithmetic.StartsWith(left, right)
	case "=$":
		return arithmetic.EndsWith(left, right)
	case "!~":
		// В JEXL x !~ y - отрицание =~
		var result any
		var err error
		if _, ok := right.(string); ok {
			result, err = arithmetic.Contains(left, right)
		} else {
			result, err = arithmetic.Contains(right, left)
		}
		if err != nil {
			return nil, err
		}
		if b, ok := result.(bool); ok {
			return !b, nil
		}
		return nil, jexl.NewError("contains operation did not return boolean")
	case "!^":
		result, err := arithmetic.StartsWith(left, right)
		if err != nil {
			return nil, err
		}
		if b, ok := result.(bool); ok {
			return !b, nil
		}
		return nil, jexl.NewError("startsWith operation did not return boolean")
	case "!$":
		result, err := arithmetic.EndsWith(left, right)
		if err != nil {
			return nil, err
		}
		if b, ok := result.(bool); ok {
			return !b, nil
		}
		return nil, jexl.NewError("endsWith operation did not return boolean")
	default:
		return nil, jexl.WrapError("unsupported binary operation: "+node.Op(), nil, nil)
	}
}

// interpretUnaryOp выполняет UnaryOpNode.
func (i *interpreter) interpretUnaryOp(node *jexl.UnaryOpNode) (any, error) {
	value, err := i.interpret(node.Operand())
	if err != nil {
		return nil, err
	}

	arithmetic := i.engine.Arithmetic()
	if arithmetic == nil {
		arithmetic = jexl.NewBaseArithmetic(true, nil, 0)
	}

	switch node.Op() {
	case "+":
		return value, nil
	case "-":
		return arithmetic.Negate(value)
	case "!":
		b, err := arithmetic.ToBoolean(value)
		if err != nil {
			return nil, err
		}
		return !b, nil
	case "~":
		return arithmetic.BitwiseComplement(value)
	default:
		return nil, jexl.WrapError("unsupported unary operation: "+node.Op(), nil, nil)
	}
}

// interpretPropertyAccess выполняет PropertyAccessNode.
func (i *interpreter) interpretPropertyAccess(node *jexl.PropertyAccessNode) (any, error) {
	obj, err := i.interpret(node.Object())
	if err != nil {
		return nil, err
	}

	if obj == nil {
		if i.options != nil && i.options.Safe() {
			return nil, nil
		}
		return nil, jexl.NewError("cannot access property on nil")
	}

	propNode := node.Property()
	propIdent, ok := propNode.(*jexl.IdentifierNode)
	if !ok {
		return nil, jexl.NewError("property must be an identifier")
	}

	propName := propIdent.Name()

	// Используем Uberspect для получения свойства
	uberspect := i.engine.Uberspect()
	if uberspect == nil {
		return nil, jexl.NewError("uberspect not available")
	}

	propGet := uberspect.GetProperty(obj, propName)
	if propGet == nil {
		if i.options != nil && i.options.Strict() {
			return nil, jexl.NewError("property not found: " + propName)
		}
		return nil, nil
	}

	return propGet.Invoke(obj)
}

// interpretIndexAccess выполняет IndexAccessNode.
func (i *interpreter) interpretIndexAccess(node *jexl.IndexAccessNode) (any, error) {
	obj, err := i.interpret(node.Object())
	if err != nil {
		return nil, err
	}

	if obj == nil {
		if i.options != nil && i.options.Safe() {
			return nil, nil
		}
		return nil, jexl.NewError("cannot access index on nil")
	}

	index, err := i.interpret(node.Index())
	if err != nil {
		return nil, err
	}

	// Сначала проверяем, является ли объект мапой
	// Для мап индекс может быть строкой или числом (преобразуется в строку)
	if m, ok := obj.(map[string]any); ok {
		key, ok := index.(string)
		if !ok {
			// Преобразуем ключ в строку
			key = fmt.Sprintf("%v", index)
		}
		val, ok := m[key]
		if !ok {
			if i.options != nil && i.options.Strict() {
				return nil, jexl.NewError("map key not found: " + key)
			}
			return nil, nil
		}
		return val, nil
	}

	// Используем reflection для проверки других типов мап
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Map {
		keyType := val.Type().Key()
		
		// Специальная обработка для map[interface{}]interface{}
		if keyType.Kind() == reflect.Interface {
			// Для interface{} ключей нужно найти ключ по значению
			// Пробуем найти ключ, сравнивая значения
			for _, mapKey := range val.MapKeys() {
				mapKeyVal := mapKey.Interface()
				// Сравниваем значения через арифметику
				if arith := i.engine.Arithmetic(); arith != nil {
					cmp, err := arith.Compare(mapKeyVal, index)
					if err == nil && cmp == 0 {
						return val.MapIndex(mapKey).Interface(), nil
					}
				} else {
					// Простое сравнение
					if mapKeyVal == index {
						return val.MapIndex(mapKey).Interface(), nil
					}
				}
			}
			// Если не нашли, пробуем преобразовать индекс и поискать снова
			keyVal := reflect.ValueOf(index)
			mapVal := val.MapIndex(keyVal)
			if mapVal.IsValid() {
				return mapVal.Interface(), nil
			}
			// Если все еще не нашли, пробуем преобразовать индекс в разные типы
			if intIdx, err := toIntIndex(index); err == nil {
				// Пробуем как int
				if mapVal = val.MapIndex(reflect.ValueOf(intIdx)); mapVal.IsValid() {
					return mapVal.Interface(), nil
				}
				// Пробуем как int64
				if mapVal = val.MapIndex(reflect.ValueOf(int64(intIdx))); mapVal.IsValid() {
					return mapVal.Interface(), nil
				}
			}
			// Пробуем как строку
			keyStr := fmt.Sprintf("%v", index)
			if mapVal = val.MapIndex(reflect.ValueOf(keyStr)); mapVal.IsValid() {
				return mapVal.Interface(), nil
			}
			if i.options != nil && i.options.Strict() {
				return nil, jexl.NewError("map key not found")
			}
			return nil, nil
		}
		
		keyVal := reflect.ValueOf(index)
		if !keyVal.Type().AssignableTo(keyType) {
			if keyVal.Type().ConvertibleTo(keyType) {
				keyVal = keyVal.Convert(keyType)
			} else if keyType.Kind() == reflect.String {
				// Преобразуем ключ в строку
				keyStr := fmt.Sprintf("%v", index)
				keyVal = reflect.ValueOf(keyStr)
			} else {
				return nil, jexl.NewError("map key type mismatch")
			}
		}
		mapVal := val.MapIndex(keyVal)
		if !mapVal.IsValid() {
			if i.options != nil && i.options.Strict() {
				return nil, jexl.NewError("map key not found")
			}
			return nil, nil
		}
		return mapVal.Interface(), nil
	}

	// Преобразуем индекс в int для массивов и слайсов
	intIndex, err := toIntIndex(index)
	if err != nil {
		return nil, jexl.NewError("array index must be integer")
	}

	// Поддержка массивов и слайсов
	switch v := obj.(type) {
	case []any:
		if intIndex < 0 || intIndex >= len(v) {
			return nil, jexl.NewError("array index out of bounds")
		}
		return v[intIndex], nil
	default:
		// Используем reflection для других типов слайсов и массивов
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			if intIndex < 0 || intIndex >= val.Len() {
				return nil, jexl.NewError("array index out of bounds")
			}
			return val.Index(intIndex).Interface(), nil
		}
		
		// Попытка использовать reflection через Uberspect
		uberspect := i.engine.Uberspect()
		if uberspect != nil {
			// Для объектов с методом Get или индексацией
			if propGet := uberspect.GetProperty(obj, fmt.Sprintf("[%v]", index)); propGet != nil {
				return propGet.Invoke(obj)
			}
		}
		return nil, jexl.NewError("unsupported index access type")
	}
}

// interpretMethodCall выполняет MethodCallNode.
func (i *interpreter) interpretMethodCall(node *jexl.MethodCallNode) (any, error) {
	// Вычисляем аргументы
	args := make([]any, len(node.Args()))
	for j, argNode := range node.Args() {
		arg, err := i.interpret(argNode)
		if err != nil {
			return nil, err
		}
		args[j] = arg
	}

	methodNode := node.Method()
	methodIdent, ok := methodNode.(*jexl.IdentifierNode)
	if !ok {
		return nil, jexl.NewError("method name must be an identifier")
	}
	methodName := methodIdent.Name()

	// Если есть target, это вызов метода объекта
	if targetNode := node.Target(); targetNode != nil {
		obj, err := i.interpret(targetNode)
		if err != nil {
			return nil, err
		}

		if obj == nil {
			if i.options != nil && i.options.Safe() {
				return nil, nil
			}
			return nil, jexl.NewError("cannot call method on nil")
		}

		// Используем Uberspect для вызова метода
		uberspect := i.engine.Uberspect()
		if uberspect == nil {
			return nil, jexl.NewError("uberspect not available")
		}

		method, err := uberspect.GetMethod(obj, methodName, args)
		if err != nil || method == nil {
			if i.options != nil && i.options.Strict() {
				return nil, jexl.NewError("method not found: " + methodName)
			}
			return nil, nil
		}

		return method.Invoke(obj, args)
	}

	// Проверяем встроенные функции
	switch methodName {
	case "empty":
		if len(args) != 1 {
			return nil, jexl.NewError("empty() requires exactly 1 argument")
		}
		return i.interpretEmpty(args[0])
	case "size":
		if len(args) != 1 {
			return nil, jexl.NewError("size() requires exactly 1 argument")
		}
		return i.interpretSize(args[0])
	}

	// Функция верхнего уровня - ищем в контексте
	if i.context == nil {
		return nil, jexl.NewError("context is nil")
	}

	if !i.context.Has(methodName) {
		if i.options != nil && i.options.Strict() {
			return nil, jexl.NewError("function not found: " + methodName)
		}
		return nil, nil
	}

	funcValue := i.context.Get(methodName)
	if funcValue == nil {
		return nil, nil
	}

	// Попытка вызвать как функцию
	if fn, ok := funcValue.(func(...any) (any, error)); ok {
		return fn(args...)
	}

	return nil, jexl.NewError("value is not callable: " + methodName)
}

func (i *interpreter) interpretAssignment(node *jexl.AssignmentNode) (any, error) {
	if node == nil {
		return nil, jexl.NewError("assignment node is nil")
	}

	// Проверяем, является ли это постфиксным инкрементом/декрементом
	// По source text определяем: если заканчивается на ++ или --, это постфиксный
	source := node.SourceText()
	isPostfix := (len(source) >= 2 && source[len(source)-2:] == "++") ||
		(len(source) >= 2 && source[len(source)-2:] == "--")
	
	var oldValue any
	if isPostfix {
		// Для постфиксного инкремента/декремента сохраняем старое значение ДО вычисления нового
		switch target := node.Target().(type) {
		case *jexl.IdentifierNode:
			if i.context != nil {
				oldValue = i.context.Get(target.Name())
			}
		case *jexl.PropertyAccessNode:
			obj, err := i.interpret(target.Object())
			if err == nil && obj != nil {
				propIdent, ok := target.Property().(*jexl.IdentifierNode)
				if ok {
					uberspect := i.engine.Uberspect()
					if uberspect != nil {
						if propGet := uberspect.GetProperty(obj, propIdent.Name()); propGet != nil {
							oldValue, _ = propGet.Invoke(obj)
						}
					}
				}
			}
		case *jexl.IndexAccessNode:
			obj, err := i.interpret(target.Object())
			if err == nil && obj != nil {
				_, err := i.interpret(target.Index())
				if err == nil {
					oldValue, _ = i.interpretIndexAccess(target)
				}
			}
		}
	}

	// Вычисляем новое значение (это может изменить значение в контексте, если используется x + 1)
	value, err := i.interpret(node.Value())
	if err != nil {
		return nil, err
	}

	switch target := node.Target().(type) {
	case *jexl.IdentifierNode:
		if i.context == nil {
			return nil, jexl.NewError("context is nil")
		}
		// Для постфиксного инкремента/декремента нужно вернуть старое значение
		if isPostfix {
			// Сохраняем старое значение перед присваиванием
			if oldValue == nil {
				// Если старое значение не было сохранено, пробуем восстановить из нового
				arithmetic := i.engine.Arithmetic()
				if arithmetic != nil {
					one := int64(1)
					if len(source) >= 2 && source[len(source)-2:] == "++" {
						oldValue, _ = arithmetic.Subtract(value, one)
					} else if len(source) >= 2 && source[len(source)-2:] == "--" {
						oldValue, _ = arithmetic.Add(value, one)
					}
				}
			}
			// Устанавливаем новое значение
			i.context.Set(target.Name(), value)
			// Возвращаем старое значение
			return oldValue, nil
		}
		// Для префиксного и обычного присваивания устанавливаем и возвращаем новое значение
		i.context.Set(target.Name(), value)
		return value, nil
	case *jexl.PropertyAccessNode:
		result, err := i.assignProperty(target, value)
		if err != nil {
			return nil, err
		}
		// Для постфиксного возвращаем старое значение
		if isPostfix && oldValue != nil {
			return oldValue, nil
		}
		return result, nil
	case *jexl.IndexAccessNode:
		result, err := i.assignIndex(target, value)
		if err != nil {
			return nil, err
		}
		// Для постфиксного возвращаем старое значение
		if isPostfix && oldValue != nil {
			return oldValue, nil
		}
		return result, nil
	default:
		return nil, jexl.NewError("unsupported assignment target")
	}
}

func (i *interpreter) assignProperty(target *jexl.PropertyAccessNode, value any) (any, error) {
	obj, err := i.interpret(target.Object())
	if err != nil {
		return nil, err
	}
	if obj == nil {
		if i.options != nil && i.options.Safe() {
			return nil, nil
		}
		return nil, jexl.NewError("cannot assign property on nil")
	}

	propIdent, ok := target.Property().(*jexl.IdentifierNode)
	if !ok {
		return nil, jexl.NewError("property must be an identifier")
	}

	uberspect := i.engine.Uberspect()
	if uberspect == nil {
		return nil, jexl.NewError("uberspect not available")
	}

	propSet := uberspect.SetProperty(obj, propIdent.Name(), value)
	if propSet == nil {
		if i.options != nil && i.options.Strict() {
			return nil, jexl.NewError("property setter not found: " + propIdent.Name())
		}
		return value, nil
	}

	if err := propSet.Invoke(obj, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (i *interpreter) assignIndex(target *jexl.IndexAccessNode, value any) (any, error) {
	obj, err := i.interpret(target.Object())
	if err != nil {
		return nil, err
	}

	if obj == nil {
		if i.options != nil && i.options.Safe() {
			return nil, nil
		}
		return nil, jexl.NewError("cannot assign index on nil")
	}

	index, err := i.interpret(target.Index())
	if err != nil {
		return nil, err
	}

	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		if objValue.IsNil() {
			return nil, jexl.NewError("cannot assign index on nil pointer")
		}
		objValue = objValue.Elem()
	}

	switch objValue.Kind() {
	case reflect.Map:
		return i.assignMapIndex(objValue, index, value)
	case reflect.Slice, reflect.Array:
		return i.assignSliceIndex(objValue, index, value)
	default:
		// Попытка использовать uberspect для объекта с индексатором
		uberspect := i.engine.Uberspect()
		if uberspect != nil {
			idxExpr := fmt.Sprintf("[%v]", index)
			if propSet := uberspect.SetProperty(obj, idxExpr, value); propSet != nil {
				if err := propSet.Invoke(obj, value); err != nil {
					return nil, err
				}
				return value, nil
			}
		}
		return nil, jexl.NewError("unsupported index assignment type")
	}
}

func (i *interpreter) assignMapIndex(mapValue reflect.Value, index any, value any) (any, error) {
	// Для map[string]any используем прямое присваивание
	if mapValue.Type().Key().Kind() == reflect.String {
		keyStr, ok := index.(string)
		if !ok {
			keyStr = fmt.Sprintf("%v", index)
		}
		
		// Если это map[string]any, можем напрямую установить значение
		if mapValue.Type().Elem().Kind() == reflect.Interface {
			mapValue.SetMapIndex(reflect.ValueOf(keyStr), reflect.ValueOf(value))
			return value, nil
		}
	}
	
	keyVal := reflect.ValueOf(index)
	keyType := mapValue.Type().Key()
	
	// Специальная обработка для map[interface{}]interface{}
	if keyType.Kind() == reflect.Interface {
		// Для interface{} ключей можем использовать любой тип
		keyVal = reflect.ValueOf(index)
	} else if !keyVal.Type().AssignableTo(keyType) {
		if keyVal.Type().ConvertibleTo(keyType) {
			keyVal = keyVal.Convert(keyType)
		} else {
			// Для map[string]any пробуем преобразовать ключ в строку
			if keyType.Kind() == reflect.String {
				keyStr := fmt.Sprintf("%v", index)
				keyVal = reflect.ValueOf(keyStr)
			} else {
				return nil, jexl.NewError("map key type mismatch")
			}
		}
	}

	valVal := reflect.ValueOf(value)
	elemType := mapValue.Type().Elem()
	if !valVal.IsValid() {
		valVal = reflect.Zero(elemType)
	} else if !valVal.Type().AssignableTo(elemType) {
		if elemType.Kind() == reflect.Interface {
			// Для interface{} значений можем использовать любое значение
			valVal = reflect.ValueOf(value)
		} else if valVal.Type().ConvertibleTo(elemType) {
			valVal = valVal.Convert(elemType)
		} else {
			return nil, jexl.NewError("map value type mismatch")
		}
	}

	mapValue.SetMapIndex(keyVal, valVal)
	return value, nil
}

func (i *interpreter) assignSliceIndex(sliceValue reflect.Value, index any, value any) (any, error) {
	intIndex, err := toIntIndex(index)
	if err != nil {
		return nil, err
	}

	if intIndex < 0 || intIndex >= sliceValue.Len() {
		return nil, jexl.NewError("slice index out of bounds")
	}

	// Для []any используем прямое присваивание
	if sliceValue.Type().Elem().Kind() == reflect.Interface {
		sliceValue.Index(intIndex).Set(reflect.ValueOf(value))
		return value, nil
	}

	valVal := reflect.ValueOf(value)
	if !valVal.IsValid() {
		valVal = reflect.Zero(sliceValue.Type().Elem())
	} else if !valVal.Type().AssignableTo(sliceValue.Type().Elem()) {
		if valVal.Type().ConvertibleTo(sliceValue.Type().Elem()) {
			valVal = valVal.Convert(sliceValue.Type().Elem())
		} else {
			return nil, jexl.NewError("slice value type mismatch")
		}
	}

	sliceValue.Index(intIndex).Set(valVal)
	return value, nil
}

func toIntIndex(index any) (int, error) {
	switch v := index.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case *big.Rat:
		if !v.IsInt() {
			return 0, jexl.NewError("index must be integer")
		}
		return int(v.Num().Int64()), nil
	case *big.Int:
		if !v.IsInt64() {
			return 0, jexl.NewError("index must be integer")
		}
		return int(v.Int64()), nil
	default:
		return 0, jexl.NewError("index must be integer")
	}
}

// interpretTernary выполняет тернарный оператор.
func (i *interpreter) interpretTernary(node *jexl.TernaryNode) (any, error) {
	condition, err := i.interpret(node.Condition())
	if err != nil {
		return nil, err
	}

	arithmetic := i.engine.Arithmetic()
	if arithmetic == nil {
		arithmetic = jexl.NewBaseArithmetic(true, nil, 0)
	}

	test, err := arithmetic.ToBoolean(condition)
	if err != nil {
		if i.options != nil && i.options.Strict() {
			return nil, err
		}
		test = false
	}

	if test {
		return i.interpret(node.TrueExpr())
	}
	return i.interpret(node.FalseExpr())
}

// interpretElvis выполняет Elvis оператор.
// Elvis оператор (??) возвращает default только если expr == nil или undefined.
// Если expr имеет значение (даже false, 0, пустая строка), возвращается expr.
func (i *interpreter) interpretElvis(node *jexl.ElvisNode) (any, error) {
	expr, err := i.interpret(node.Expr())
	if err != nil {
		// Если ошибка и это не строгий режим, возвращаем default
		if i.options != nil && i.options.Safe() {
			return i.interpret(node.DefaultExpr())
		}
		// В строгом режиме проверяем, является ли это ошибкой "variable not found"
		if i.options == nil || !i.options.Strict() {
			// В нестрогом режиме для ошибок переменных возвращаем default
			return i.interpret(node.DefaultExpr())
		}
		return nil, err
	}

	// Elvis оператор возвращает default только если expr == nil
	// Если expr имеет любое другое значение (включая false, 0, ""), возвращаем expr
	if expr == nil {
		return i.interpret(node.DefaultExpr())
	}

	// Возвращаем expr, даже если он false, 0, пустая строка и т.д.
	return expr, nil
}

// interpretRange выполняет RangeNode.
func (i *interpreter) interpretRange(node *jexl.RangeNode) (any, error) {
	left, err := i.interpret(node.Left())
	if err != nil {
		return nil, err
	}

	right, err := i.interpret(node.Right())
	if err != nil {
		return nil, err
	}

	// Получаем арифметику из движка
	arithmetic := i.engine.Arithmetic()
	if arithmetic == nil {
		// Используем базовую арифметику
		arithmetic = jexl.NewBaseArithmetic(true, nil, 0)
	}

	// Создаём range через арифметику
	return arithmetic.CreateRange(left, right)
}

// interpretArrayLiteral выполняет литерал массива.
func (i *interpreter) interpretArrayLiteral(node *jexl.ArrayLiteralNode) (any, error) {
	elements := node.Elements()
	result := make([]any, len(elements))

	for j, elem := range elements {
		val, err := i.interpret(elem)
		if err != nil {
			return nil, err
		}
		result[j] = val
	}

	return result, nil
}

// interpretMapLiteral выполняет литерал мапы.
func (i *interpreter) interpretMapLiteral(node *jexl.MapLiteralNode) (any, error) {
	entries := node.Entries()
	result := make(map[string]any, len(entries))

	for _, entry := range entries {
		key, err := i.interpret(entry.Key)
		if err != nil {
			return nil, err
		}
		value, err := i.interpret(entry.Value)
		if err != nil {
			return nil, err
		}

		keyStr, ok := key.(string)
		if !ok {
			// Пробуем преобразовать в строку
			keyStr = fmt.Sprintf("%v", key)
		}
		result[keyStr] = value
	}

	return result, nil
}

// interpretSetLiteral выполняет литерал множества.
func (i *interpreter) interpretSetLiteral(node *jexl.SetLiteralNode) (any, error) {
	elements := node.Elements()
	if elements == nil {
		return []any{}, nil
	}

	result := make([]any, 0, len(elements))
	seen := make(map[any]bool)

	for _, elem := range elements {
		val, err := i.interpret(elem)
		if err != nil {
			return nil, err
		}

		// Проверяем на дубликаты (простая реализация)
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}

	return result, nil
}

// BreakError используется для выхода из цикла.
type BreakError struct{}

func (e *BreakError) Error() string {
	return "break"
}

// ContinueError используется для продолжения цикла.
type ContinueError struct{}

func (e *ContinueError) Error() string {
	return "continue"
}

// interpretIf выполняет if/else statement.
func (i *interpreter) interpretIf(node *jexl.IfNode) (any, error) {
	condition, err := i.interpret(node.Condition())
	if err != nil {
		return nil, err
	}

	arithmetic := i.engine.Arithmetic()
	if arithmetic == nil {
		arithmetic = jexl.NewBaseArithmetic(true, nil, 0)
	}

	test, err := arithmetic.ToBoolean(condition)
	if err != nil {
		if i.options != nil && i.options.Strict() {
			return nil, err
		}
		test = false
	}

	if test {
		return i.interpret(node.ThenBranch())
	}

	if node.ElseBranch() != nil {
		return i.interpret(node.ElseBranch())
	}

	return nil, nil
}

// interpretFor выполняет цикл for (init; condition; step) body.
func (i *interpreter) interpretFor(node *jexl.ForNode) (any, error) {
	// Инициализация
	if node.Init() != nil {
		_, err := i.interpret(node.Init())
		if err != nil {
			return nil, err
		}
	}

	arithmetic := i.engine.Arithmetic()
	if arithmetic == nil {
		arithmetic = jexl.NewBaseArithmetic(true, nil, 0)
	}

	var result any
	for {
		// Проверка условия
		if node.Condition() != nil {
			condition, err := i.interpret(node.Condition())
			if err != nil {
				return nil, err
			}
			test, err := arithmetic.ToBoolean(condition)
			if err != nil {
				if i.options != nil && i.options.Strict() {
					return nil, err
				}
				test = false
			}
			if !test {
				break
			}
		}

		// Выполнение тела
		if node.Body() != nil {
			bodyResult, err := i.interpret(node.Body())
			if err != nil {
				if _, isBreak := err.(*BreakError); isBreak {
					break
				}
				if _, isContinue := err.(*ContinueError); isContinue {
					// Продолжаем цикл
					if node.Step() != nil {
						_, err := i.interpret(node.Step())
						if err != nil {
							return nil, err
						}
					}
					continue
				}
				return nil, err
			}
			result = bodyResult
		}

		// Шаг
		if node.Step() != nil {
			_, err := i.interpret(node.Step())
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// interpretForeach выполняет цикл foreach (var x : items) body.
func (i *interpreter) interpretForeach(node *jexl.ForeachNode) (any, error) {
	items, err := i.interpret(node.Items())
	if err != nil {
		return nil, err
	}

	varName := node.Variable()
	var ident *jexl.IdentifierNode
	if id, ok := varName.(*jexl.IdentifierNode); ok {
		ident = id
	} else {
		return nil, jexl.NewError("foreach variable must be an identifier")
	}

	// Преобразуем items в итерируемую коллекцию
	var iterable []any
	switch v := items.(type) {
	case []any:
		iterable = v
	case []int:
		iterable = make([]any, len(v))
		for j, val := range v {
			iterable[j] = val
		}
	case []string:
		iterable = make([]any, len(v))
		for j, val := range v {
			iterable[j] = val
		}
	case map[string]any:
		// Для мап итерируем по значениям
		iterable = make([]any, 0, len(v))
		for _, val := range v {
			iterable = append(iterable, val)
		}
	default:
		return nil, jexl.NewError("foreach items must be iterable")
	}

	var result any
	for _, item := range iterable {
		// Устанавливаем переменную цикла
		if i.context != nil {
			i.context.Set(ident.Name(), item)
		}

		// Выполняем тело
		if node.Body() != nil {
			bodyResult, err := i.interpret(node.Body())
			if err != nil {
				if _, isBreak := err.(*BreakError); isBreak {
					break
				}
				if _, isContinue := err.(*ContinueError); isContinue {
					continue
				}
				return nil, err
			}
			result = bodyResult
		}
	}

	return result, nil
}

// interpretWhile выполняет цикл while (condition) body.
func (i *interpreter) interpretWhile(node *jexl.WhileNode) (any, error) {
	arithmetic := i.engine.Arithmetic()
	if arithmetic == nil {
		arithmetic = jexl.NewBaseArithmetic(true, nil, 0)
	}

	var result any
	for {
		condition, err := i.interpret(node.Condition())
		if err != nil {
			return nil, err
		}

		test, err := arithmetic.ToBoolean(condition)
		if err != nil {
			if i.options != nil && i.options.Strict() {
				return nil, err
			}
			test = false
		}

		if !test {
			break
		}

		if node.Body() != nil {
			bodyResult, err := i.interpret(node.Body())
			if err != nil {
				if _, isBreak := err.(*BreakError); isBreak {
					break
				}
				if _, isContinue := err.(*ContinueError); isContinue {
					continue
				}
				return nil, err
			}
			result = bodyResult
		}
	}

	return result, nil
}

// interpretDoWhile выполняет цикл do body while (condition).
func (i *interpreter) interpretDoWhile(node *jexl.DoWhileNode) (any, error) {
	arithmetic := i.engine.Arithmetic()
	if arithmetic == nil {
		arithmetic = jexl.NewBaseArithmetic(true, nil, 0)
	}

	var result any
	for {
		if node.Body() != nil {
			bodyResult, err := i.interpret(node.Body())
			if err != nil {
				if _, isBreak := err.(*BreakError); isBreak {
					break
				}
				if _, isContinue := err.(*ContinueError); isContinue {
					// Проверяем условие и продолжаем
					condition, err := i.interpret(node.Condition())
					if err != nil {
						return nil, err
					}
					test, err := arithmetic.ToBoolean(condition)
					if err != nil {
						if i.options != nil && i.options.Strict() {
							return nil, err
						}
						test = false
					}
					if !test {
						break
					}
					continue
				}
				return nil, err
			}
			result = bodyResult
		}

		condition, err := i.interpret(node.Condition())
		if err != nil {
			return nil, err
		}

		test, err := arithmetic.ToBoolean(condition)
		if err != nil {
			if i.options != nil && i.options.Strict() {
				return nil, err
			}
			test = false
		}

		if !test {
			break
		}
	}

	return result, nil
}

// interpretBlock выполняет блок кода.
func (i *interpreter) interpretBlock(node *jexl.BlockNode) (any, error) {
	var result any
	for _, stmt := range node.Statements() {
		var err error
		result, err = i.interpret(stmt)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// interpretReturn выполняет return statement.
// interpretEmpty проверяет, является ли значение пустым.
func (i *interpreter) interpretEmpty(value any) (any, error) {
	if value == nil {
		return true, nil
	}
	
	switch v := value.(type) {
	case string:
		return len(v) == 0, nil
	case []any:
		return len(v) == 0, nil
	case map[string]any:
		return len(v) == 0, nil
	case bool:
		return !v, nil
	default:
		// Для чисел проверяем, равно ли нулю
		arithmetic := i.engine.Arithmetic()
		if arithmetic == nil {
			arithmetic = jexl.NewBaseArithmetic(true, nil, 0)
		}
		b, err := arithmetic.ToBoolean(value)
		if err != nil {
			return false, nil
		}
		return !b, nil
	}
}

// interpretSize возвращает размер коллекции или строки.
func (i *interpreter) interpretSize(value any) (any, error) {
	if value == nil {
		return int64(0), nil
	}
	
	switch v := value.(type) {
	case string:
		return int64(len(v)), nil
	case []any:
		return int64(len(v)), nil
	case map[string]any:
		return int64(len(v)), nil
	default:
		// Для других типов пробуем использовать reflection
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
			return int64(rv.Len()), nil
		default:
			return nil, jexl.NewError("size() is not applicable to this type")
		}
	}
}

func (i *interpreter) interpretReturn(node *jexl.ReturnNode) (any, error) {
	if node.Value() != nil {
		return i.interpret(node.Value())
	}
	return nil, nil
}
