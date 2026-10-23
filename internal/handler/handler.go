// Package handler
package handler

import (
	"github.com/Luckny/recl/internal/parser"
)

func HandleValue(v parser.Value) parser.Value {
	// should arrive as array
	if v.Type != parser.ArrayPrefix {
		return parser.ErrorValue("expected an array")
	}

	return parser.ErrorValue("test error")
}
