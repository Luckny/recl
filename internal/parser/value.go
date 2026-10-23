package parser

import (
	"fmt"
	"strings"
)

type ValueType rune

const (
	ArrayPrefix ValueType = '*'
	BulkPrefix  ValueType = '$'
	Error       ValueType = '-'
)

type Value struct {
	Type  ValueType
	Bulk  string
	Array []Value
	Error string
}

func (v *Value) Print() string {
	switch v.Type {
	case ArrayPrefix:
		return fmt.Sprintf("Array(%d) - %s", len(v.Array), v.PrintArray())
	case BulkPrefix:
		return fmt.Sprintf("Bulk(%d) - %s", len(v.Bulk), v.Bulk)
	default:
		return fmt.Sprintf("Unknown(%c)", v.Type)
	}
}

func (v *Value) PrintArray() string {
	values := make([]string, len(v.Array))
	for i, element := range v.Array {
		values[i] = element.Print()
	}
	return "[" + strings.Join(values, ", ") + "]"
}

func ErrorValue(message string) Value {
	return Value{Type: Error, Error: "ERR " + message}
}
