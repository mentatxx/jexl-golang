package jexl

import (
	"fmt"
	"math/big"
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
