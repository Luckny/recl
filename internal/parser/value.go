package parser

import (
	"fmt"
	"strings"
)

type ValueType rune

const (
	ArrayType       ValueType = '*'
	BulkType        ValueType = '$'
	SimpleStingType ValueType = '+'
	ErrorType       ValueType = '-'
)

type Value struct {
	Type         ValueType
	Bulk         string
	Array        []Value
	Error        string
	SimpleString string
}

func (v *Value) Print() string {
	switch v.Type {
	case ArrayType:
		return fmt.Sprintf("Array(%d) - %s", len(v.Array), v.PrintArray())
	case BulkType:
		return fmt.Sprintf("Bulk(%d) - %s", len(v.Bulk), v.Bulk)
	case SimpleStingType:
		return fmt.Sprintf("SimpleString(%d) - %s", len(v.SimpleString), v.SimpleString)
	case ErrorType:
		return fmt.Sprintf("Error(%d) - %s", len(v.Error), v.Error)
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

func SimpleStringValue(message string) Value {
	return Value{Type: SimpleStingType, SimpleString: message}
}

func ErrorValue(message string) Value {
	return Value{Type: ErrorType, Error: "ERR " + message}
}
