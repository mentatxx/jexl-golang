package jexl

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Info детализирует положение исходного кода для диагностики ошибок.
// Аналог org.apache.commons.jexl3.JexlInfo.
type Info struct {
	name   string
	line   int
	column int
	detail Detail
}

// Detail уточняет выделение участка кода.
type Detail interface {
	Start() int
	End() int
	fmt.Stringer
}

// NewInfo создаёт Info, используя стек вызовов для попытки определить точку происхождения.
func NewInfo() *Info {
	const skipInternalFrames = 2
	pcs := make([]uintptr, 32)
	n := runtime.Callers(skipInternalFrames, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	var (
		frame runtime.Frame
		more  bool
	)
	for {
		frame, more = frames.Next()
		if !more {
			break
		}
		if !isJexlInternal(frame.Function) {
			return &Info{
				name:   fmt.Sprintf("%s:%d", frame.Function, frame.Line),
				line:   1,
				column: 1,
			}
		}
	}

	return &Info{
		name:   "?",
		line:   1,
		column: 1,
	}
}

func isJexlInternal(function string) bool {
	if function == "" {
		return false
	}
	path := filepath.ToSlash(function)
	return containsAny(path,
		"/jexl.", "/jexl/internal/", "/jexl/parser/", "/jexl/Jexl")
}

func containsAny(s string, needles ...string) bool {
	for _, needle := range needles {
		if needle != "" && strings.Contains(s, needle) {
			return true
		}
	}
	return false
}

// NewInfoAt создаёт Info с заданными значениями.
func NewInfoAt(name string, line, column int) *Info {
	if line <= 0 {
		line = 1
	}
	if column <= 0 {
		column = 1
	}
	return &Info{name: name, line: line, column: column}
}

// WithDetail устанавливает Detail и возвращает Info.
func (i *Info) WithDetail(detail Detail) *Info {
	if i == nil {
		return nil
	}
	cp := *i
	cp.detail = detail
	return &cp
}

// At создаёт копию с новой позицией.
func (i *Info) At(line, column int) *Info {
	if i == nil {
		return nil
	}
	return NewInfoAt(i.name, line, column).WithDetail(i.detail)
}

// Name возвращает имя источника.
func (i *Info) Name() string {
	if i == nil {
		return ""
	}
	return i.name
}

// Line возвращает номер строки.
func (i *Info) Line() int {
	if i == nil {
		return 0
	}
	return i.line
}

// Column возвращает номер столбца.
func (i *Info) Column() int {
	if i == nil {
		return 0
	}
	return i.column
}

// Detail возвращает детальную информацию.
func (i *Info) Detail() Detail {
	if i == nil {
		return nil
	}
	return i.detail
}

// String форматирует в формате name@line:column.
func (i *Info) String() string {
	if i == nil {
		return ""
	}
	result := fmt.Sprintf("%s@%d:%d", i.name, i.line, i.column)
	if d := i.detail; d != nil {
		result += fmt.Sprintf("![%d,%d]: '%s'", d.Start(), d.End(), d.String())
	}
	return result
}
