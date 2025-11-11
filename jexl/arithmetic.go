package jexl

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strings"
)

// Arithmetic определяет операции, используемые движком.
type Arithmetic interface {
	IsStrict() bool
	MathContext() *MathContext
	MathScale() int
	Compare(a, b any) (int, error)
	Add(a, b any) (any, error)
	Subtract(a, b any) (any, error)
	Multiply(a, b any) (any, error)
	Divide(a, b any) (any, error)
	Modulo(a, b any) (any, error)
	Negate(value any) (any, error)
	ToBoolean(value any) (bool, error)
	// Битовые операции
	BitwiseAnd(a, b any) (any, error)
	BitwiseOr(a, b any) (any, error)
	BitwiseXor(a, b any) (any, error)
	BitwiseComplement(value any) (any, error)
	ShiftLeft(a, b any) (any, error)
	ShiftRight(a, b any) (any, error)
	ShiftRightUnsigned(a, b any) (any, error)
	// Строковые операции
	Contains(a, b any) (any, error)
	ContainsAll(collection, elements any) (any, error) // Проверяет, содержится ли каждый элемент elements в collection
	StartsWith(a, b any) (any, error)
	EndsWith(a, b any) (any, error)
	// Range операция
	CreateRange(left, right any) (any, error)
}

// BaseArithmetic предоставляет простую реализацию с ограниченным функционалом.
type BaseArithmetic struct {
	strict  bool
	context *MathContext
	scale   int
}

// NewBaseArithmetic создаёт базовый Arithmetic.
func NewBaseArithmetic(strict bool, ctx *MathContext, scale int) *BaseArithmetic {
	return &BaseArithmetic{
		strict:  strict,
		context: ctx,
		scale:   scale,
	}
}

// IsStrict сообщает, включена ли строгая арифметика.
func (a *BaseArithmetic) IsStrict() bool {
	return a.strict
}

// MathContext возвращает контекст.
func (a *BaseArithmetic) MathContext() *MathContext {
	return a.context
}

// MathScale возвращает масштаб.
func (a *BaseArithmetic) MathScale() int {
	return a.scale
}

// Compare выполняет сравнение, используя big.Rat при необходимости.
func (a *BaseArithmetic) Compare(lhs, rhs any) (int, error) {
	// Специальная обработка для null
	if lhs == nil && rhs == nil {
		return 0, nil
	}
	if lhs == nil {
		return -1, nil // null меньше любого не-null значения
	}
	if rhs == nil {
		return 1, nil // любое не-null значение больше null
	}
	
	// Специальная обработка для bool
	if lb, ok := lhs.(bool); ok {
		if rb, ok := rhs.(bool); ok {
			if lb == rb {
				return 0, nil
			}
			if lb {
				return 1, nil
			}
			return -1, nil
		}
		// bool vs не-bool: преобразуем bool в число
		if lb {
			lhs = 1
		} else {
			lhs = 0
		}
	}
	if rb, ok := rhs.(bool); ok {
		if rb {
			rhs = 1
		} else {
			rhs = 0
		}
	}
	
	// Специальная обработка для строк
	if ls, ok := lhs.(string); ok {
		if rs, ok := rhs.(string); ok {
			if ls < rs {
				return -1, nil
			}
			if ls > rs {
				return 1, nil
			}
			return 0, nil
		}
	}
	
	al, ok := toBig(lhs)
	if !ok {
		return 0, ErrUnsupportedOperand
	}
	ar, ok := toBig(rhs)
	if !ok {
		return 0, ErrUnsupportedOperand
	}
	return al.Cmp(ar), nil
}

// Add складывает значения.
func (a *BaseArithmetic) Add(lhs, rhs any) (any, error) {
	// Специальная обработка для строковой конкатенации
	if ls, ok := lhs.(string); ok {
		rs := fmt.Sprintf("%v", rhs)
		return ls + rs, nil
	}
	if rs, ok := rhs.(string); ok {
		ls := fmt.Sprintf("%v", lhs)
		return ls + rs, nil
	}
	
	al, ok := toBig(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	ar, ok := toBig(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	return new(big.Rat).Add(al, ar), nil
}

// Subtract вычитает значения.
func (a *BaseArithmetic) Subtract(lhs, rhs any) (any, error) {
	al, ok := toBig(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	ar, ok := toBig(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	return new(big.Rat).Sub(al, ar), nil
}

// Multiply умножает значения.
func (a *BaseArithmetic) Multiply(lhs, rhs any) (any, error) {
	al, ok := toBig(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	ar, ok := toBig(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	return new(big.Rat).Mul(al, ar), nil
}

// Divide делит значения.
func (a *BaseArithmetic) Divide(lhs, rhs any) (any, error) {
	al, ok := toBig(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	ar, ok := toBig(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	// Проверяем деление на ноль
	if ar.Sign() == 0 {
		return nil, NewError("division by zero")
	}
	return new(big.Rat).Quo(al, ar), nil
}

// Modulo вычисляет остаток.
func (a *BaseArithmetic) Modulo(lhs, rhs any) (any, error) {
	al, ok := toBig(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	ar, ok := toBig(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	// Для big.Rat модуль не имеет прямого смысла, используем целочисленное деление
	// Преобразуем в int64 для вычисления модуля
	alInt := al.Num().Int64()
	arInt := ar.Num().Int64()
	if arInt == 0 {
		return nil, NewError("division by zero")
	}
	return big.NewRat(alInt%arInt, 1), nil
}

// Negate возвращает противоположное значение.
func (a *BaseArithmetic) Negate(value any) (any, error) {
	v, ok := toBig(value)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	return new(big.Rat).Neg(v), nil
}

// ToBoolean приводит значение к bool.
func (a *BaseArithmetic) ToBoolean(value any) (bool, error) {
	switch v := value.(type) {
	case nil:
		return false, nil
	case bool:
		return v, nil
	case *big.Rat:
		return v.Sign() != 0, nil
	case string:
		return len(v) > 0, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		// Преобразуем в big.Rat для проверки
		if rat, ok := toBig(v); ok {
			return rat.Sign() != 0, nil
		}
		return false, nil
	case float32, float64:
		if rat, ok := toBig(v); ok {
			return rat.Sign() != 0, nil
		}
		return false, nil
	case []any:
		return len(v) > 0, nil
	case map[string]any:
		return len(v) > 0, nil
	default:
		// Для неизвестных типов пробуем преобразовать через toBig
		if rat, ok := toBig(v); ok {
			return rat.Sign() != 0, nil
		}
		// Если не число, считаем что true (объект существует)
		return true, nil
	}
}

var (
	// ErrUnsupportedOperand сигнализирует о неподдерживаемом типе операнда.
	ErrUnsupportedOperand = NewError("unsupported operand type")
	// ErrUnsupportedOperation сигнализирует о неподдерживаемой операции.
	ErrUnsupportedOperation = NewError("unsupported operation")
)

// BitwiseAnd выполняет побитовую операцию AND.
func (a *BaseArithmetic) BitwiseAnd(lhs, rhs any) (any, error) {
	left, ok := toInt64(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	right, ok := toInt64(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	return int64(left & right), nil
}

// BitwiseOr выполняет побитовую операцию OR.
func (a *BaseArithmetic) BitwiseOr(lhs, rhs any) (any, error) {
	left, ok := toInt64(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	right, ok := toInt64(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	return int64(left | right), nil
}

// BitwiseXor выполняет побитовую операцию XOR.
func (a *BaseArithmetic) BitwiseXor(lhs, rhs any) (any, error) {
	left, ok := toInt64(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	right, ok := toInt64(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	return int64(left ^ right), nil
}

// BitwiseComplement выполняет побитовую операцию NOT (дополнение).
func (a *BaseArithmetic) BitwiseComplement(value any) (any, error) {
	v, ok := toInt64(value)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	return int64(^v), nil
}

// ShiftLeft выполняет побитовый сдвиг влево.
func (a *BaseArithmetic) ShiftLeft(lhs, rhs any) (any, error) {
	left, ok := toInt64(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	right, ok := toInt64(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	// В Go отрицательный сдвиг дает 0
	if right < 0 {
		return int64(0), nil
	}
	if right > 63 {
		return int64(0), nil
	}
	return int64(left << uint(right)), nil
}

// ShiftRight выполняет побитовый сдвиг вправо (арифметический).
func (a *BaseArithmetic) ShiftRight(lhs, rhs any) (any, error) {
	left, ok := toInt64(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	right, ok := toInt64(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	// В Go отрицательный сдвиг дает 0
	if right < 0 {
		return int64(0), nil
	}
	if right > 63 {
		// Для больших сдвигов результат зависит от знака
		if left < 0 {
			return int64(-1), nil
		}
		return int64(0), nil
	}
	return int64(left >> uint(right)), nil
}

// ShiftRightUnsigned выполняет побитовый сдвиг вправо (логический, беззнаковый).
func (a *BaseArithmetic) ShiftRightUnsigned(lhs, rhs any) (any, error) {
	// Для беззнакового сдвига преобразуем знаковое число в беззнаковое
	leftSigned, ok := toInt64(lhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	// Преобразуем в uint64 для беззнакового сдвига
	left := uint64(leftSigned)
	right, ok := toInt64(rhs)
	if !ok {
		return nil, ErrUnsupportedOperand
	}
	// В Go отрицательный сдвиг дает 0
	if right < 0 {
		return int64(0), nil
	}
	if right > 63 {
		return int64(0), nil
	}
	return int64(left >> uint(right)), nil
}

// toInt64 преобразует значение в int64 для битовых операций.
func toInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		if v <= uint64(^uint64(0)>>1) {
			return int64(v), true
		}
		return 0, false
	case *big.Rat:
		if v.IsInt() {
			return v.Num().Int64(), true
		}
		return 0, false
	case *big.Int:
		return v.Int64(), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	case nil:
		return 0, true // null coerced to 0
	default:
		return 0, false
	}
}

// toUint64 преобразует значение в uint64 для беззнаковых операций.
func toUint64(value any) (uint64, bool) {
	switch v := value.(type) {
	case int:
		if v >= 0 {
			return uint64(v), true
		}
		return 0, false
	case int8:
		if v >= 0 {
			return uint64(v), true
		}
		return 0, false
	case int16:
		if v >= 0 {
			return uint64(v), true
		}
		return 0, false
	case int32:
		if v >= 0 {
			return uint64(v), true
		}
		return 0, false
	case int64:
		if v >= 0 {
			return uint64(v), true
		}
		return 0, false
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case *big.Rat:
		if v.IsInt() && v.Sign() >= 0 {
			return v.Num().Uint64(), true
		}
		return 0, false
	case *big.Int:
		if v.Sign() >= 0 {
			return v.Uint64(), true
		}
		return 0, false
	case float32:
		if v >= 0 {
			return uint64(v), true
		}
		return 0, false
	case float64:
		if v >= 0 {
			return uint64(v), true
		}
		return 0, false
	case nil:
		return 0, true // null coerced to 0
	default:
		return 0, false
	}
}

// Contains проверяет, содержит ли lhs rhs.
// Для строк: проверяет соответствие регулярному выражению или подстроке
// Для коллекций: проверяет наличие элемента
func (a *BaseArithmetic) Contains(lhs, rhs any) (any, error) {
	// Если lhs - строка, проверяем соответствие rhs (регулярное выражение или подстрока)
	if ls, ok := lhs.(string); ok {
		rs := fmt.Sprintf("%v", rhs)
		// Пробуем как регулярное выражение
		if matched, err := regexp.MatchString(rs, ls); err == nil && matched {
			return true, nil
		}
		// Если не регулярное выражение, проверяем как подстроку
		return strings.Contains(ls, rs), nil
	}
	
	// Если rhs - строка, а lhs - нет, пробуем обратный порядок
	if rs, ok := rhs.(string); ok {
		ls := fmt.Sprintf("%v", lhs)
		// Пробуем как регулярное выражение
		if matched, err := regexp.MatchString(rs, ls); err == nil && matched {
			return true, nil
		}
		// Если не регулярное выражение, проверяем как подстроку
		return strings.Contains(ls, rs), nil
	}
	
	// Для коллекций проверяем наличие элемента
	rv := reflect.ValueOf(lhs)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i).Interface()
			// Используем улучшенное сравнение с приведением типов
			if equals(elem, rhs) {
				return true, nil
			}
			// Дополнительная проверка: пробуем преобразовать оба значения в числа
			if ar, ok1 := toBig(elem); ok1 {
				if br, ok2 := toBig(rhs); ok2 {
					if ar.Cmp(br) == 0 {
						return true, nil
					}
				}
			}
		}
		return false, nil
	case reflect.Map:
		// Для мапы проверяем ключи
		for _, key := range rv.MapKeys() {
			if equals(key.Interface(), rhs) {
				return true, nil
			}
		}
		return false, nil
	}
	
	// Для других типов пробуем преобразовать в строку
	ls := fmt.Sprintf("%v", lhs)
	rs := fmt.Sprintf("%v", rhs)
	return strings.Contains(ls, rs), nil
}

// ContainsAll проверяет, содержится ли каждый элемент elements в collection.
// Используется для проверки, является ли один массив подмножеством другого.
func (a *BaseArithmetic) ContainsAll(collection, elements any) (any, error) {
	// Получаем элементы для проверки
	elementsRv := reflect.ValueOf(elements)
	if elementsRv.Kind() != reflect.Array && elementsRv.Kind() != reflect.Slice {
		return false, NewError("ContainsAll: elements must be an array or slice")
	}
	
	// Если элементов нет, возвращаем true (пустое множество содержится в любом множестве)
	if elementsRv.Len() == 0 {
		return true, nil
	}
	
	// Проверяем каждый элемент
	for i := 0; i < elementsRv.Len(); i++ {
		elem := elementsRv.Index(i).Interface()
		// Проверяем, содержится ли элемент в collection
		contains, err := a.Contains(collection, elem)
		if err != nil {
			return false, err
		}
		if b, ok := contains.(bool); !ok || !b {
			return false, nil
		}
	}
	
	return true, nil
}

// StartsWith проверяет, начинается ли строка с подстроки.
func (a *BaseArithmetic) StartsWith(lhs, rhs any) (any, error) {
	ls := fmt.Sprintf("%v", lhs)
	rs := fmt.Sprintf("%v", rhs)
	return strings.HasPrefix(ls, rs), nil
}

// EndsWith проверяет, заканчивается ли строка подстрокой.
func (a *BaseArithmetic) EndsWith(lhs, rhs any) (any, error) {
	ls := fmt.Sprintf("%v", lhs)
	rs := fmt.Sprintf("%v", rhs)
	return strings.HasSuffix(ls, rs), nil
}

// CreateRange создаёт range от left до right включительно.
// Возвращает слайс чисел от left до right.
func (a *BaseArithmetic) CreateRange(left, right any) (any, error) {
	// Преобразуем left и right в числа
	// Сначала пробуем через toBig для поддержки всех числовых типов
	leftRat, ok := toBig(left)
	if !ok {
		// Если toBig не сработал, пробуем toInteger
		leftNum, err := a.toInteger(left)
		if err != nil {
			return nil, NewError("range left operand must be a number")
		}
		leftRat = big.NewRat(leftNum, 1)
	}
	rightRat, ok := toBig(right)
	if !ok {
		// Если toBig не сработал, пробуем toInteger
		rightNum, err := a.toInteger(right)
		if err != nil {
			return nil, NewError("range right operand must be a number")
		}
		rightRat = big.NewRat(rightNum, 1)
	}
	
	// Преобразуем в int64
	if !leftRat.IsInt() {
		return nil, NewError("range left operand must be an integer")
	}
	if !rightRat.IsInt() {
		return nil, NewError("range right operand must be an integer")
	}
	leftNum := leftRat.Num().Int64()
	rightNum := rightRat.Num().Int64()

	// Создаём range
	var result []int64
	if leftNum <= rightNum {
		// Возрастающий range
		for i := leftNum; i <= rightNum; i++ {
			result = append(result, i)
		}
	} else {
		// Убывающий range
		for i := leftNum; i >= rightNum; i-- {
			result = append(result, i)
		}
	}

	return result, nil
}

// toInteger преобразует значение в int64.
func (a *BaseArithmetic) toInteger(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > uint64(^uint64(0)>>1) {
			return 0, NewError("value too large for int64")
		}
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case *big.Rat:
		if !v.IsInt() {
			return 0, NewError("cannot convert non-integer to int64")
		}
		return v.Num().Int64(), nil
	case *big.Int:
		if !v.IsInt64() {
			return 0, NewError("value too large for int64")
		}
		return v.Int64(), nil
	case string:
		// Пробуем распарсить как число
		if r, ok := new(big.Rat).SetString(v); ok {
			if !r.IsInt() {
				return 0, NewError("cannot convert non-integer string to int64")
			}
			return r.Num().Int64(), nil
		}
		return 0, NewError("cannot convert string to int64")
	default:
		// Пробуем использовать toBig для преобразования
		if rat, ok := toBig(value); ok {
			if !rat.IsInt() {
				return 0, NewError("cannot convert non-integer to int64")
			}
			// Для отрицательных чисел Num() возвращает отрицательный big.Int
			// Int64() правильно обрабатывает отрицательные числа
			return rat.Num().Int64(), nil
		}
		// Если toBig не сработал, пробуем через reflection или другие способы
		// Это может быть необходимо для некоторых типов, которые не обрабатываются напрямую
		return 0, NewError("cannot convert value to int64: unsupported operand type")
	}
}

// equals сравнивает два значения на равенство.
func equals(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	// Прямое сравнение
	if a == b {
		return true
	}
	
	// Сравнение строк
	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			return sa == sb
		}
		// Пробуем преобразовать b в строку
		return sa == fmt.Sprintf("%v", b)
	}
	if sb, ok := b.(string); ok {
		return fmt.Sprintf("%v", a) == sb
	}
	
	// Сравнение чисел - используем toBig для приведения к общему типу
	ar, ok1 := toBig(a)
	br, ok2 := toBig(b)
	if ok1 && ok2 {
		return ar.Cmp(br) == 0
	}
	
	// Для других типов пробуем прямое сравнение через reflection
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	
	// Если типы совпадают, используем прямое сравнение
	if va.Type() == vb.Type() {
		return va.Interface() == vb.Interface()
	}
	
	// Пробуем преобразовать оба в строки и сравнить
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func toBig(value any) (*big.Rat, bool) {
	switch v := value.(type) {
	case bool:
		if v {
			return big.NewRat(1, 1), true
		}
		return big.NewRat(0, 1), true
	case int:
		return big.NewRat(int64(v), 1), true
	case int8:
		return big.NewRat(int64(v), 1), true
	case int16:
		return big.NewRat(int64(v), 1), true
	case int32:
		return big.NewRat(int64(v), 1), true
	case int64:
		return big.NewRat(v, 1), true
	case uint:
		return big.NewRat(int64(v), 1), true
	case uint8:
		return big.NewRat(int64(v), 1), true
	case uint16:
		return big.NewRat(int64(v), 1), true
	case uint32:
		return big.NewRat(int64(v), 1), true
	case uint64:
		// Осторожно с большими uint64
		if v <= uint64(^uint64(0)>>1) {
			return big.NewRat(int64(v), 1), true
		}
		return nil, false
	case float32:
		r := new(big.Rat)
		r.SetFloat64(float64(v))
		return r, true
	case float64:
		r := new(big.Rat)
		r.SetFloat64(v)
		return r, true
	case *big.Rat:
		return new(big.Rat).Set(v), true
	case *big.Int:
		return new(big.Rat).SetInt(v), true
	case string:
		// Попытка преобразовать строку в число
		if r, ok := new(big.Rat).SetString(v); ok {
			return r, true
		}
		return nil, false
	default:
		return nil, false
	}
}
