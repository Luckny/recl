package handler

import "github.com/Luckny/recl/internal/parser"

func PingHandler(args []string) parser.Value {
	return parser.SimpleStringValue("PONG")
}
