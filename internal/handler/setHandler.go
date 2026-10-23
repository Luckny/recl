package handler

import "github.com/Luckny/recl/internal/parser"

func SetHandler(args []string) parser.Value {
	return parser.SimpleStringValue("PONG")
}
