// Package handler
package handler

import (
	"strings"

	"github.com/Luckny/recl/internal/parser"
)

type Handler func(args []string) parser.Value

var Handlers = map[string]Handler{
	"PING": PingHandler,
	"SET":  SetHandler,
}

func HandleValue(v parser.Value) parser.Value {
	// should arrive as array
	if v.Type != parser.ArrayType {
		return parser.ErrorValue("expected an array")
	}

	if len(v.Array) == 0 {
		return parser.ErrorValue("empty command")
	}

	args := make([]string, len(v.Array))
	for i, arg := range v.Array {
		if arg.Type != parser.BulkType {
			return parser.ErrorValue("expected a bulk string")
		}
		args[i] = arg.Bulk
	}

	command := strings.ToUpper(args[0])
	handler, ok := Handlers[command]
	if !ok {
		return parser.ErrorValue("invalid command")
	}

	return handler(args[1:])
}
