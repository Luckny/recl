// Package handler
package handler

import (
	"log/slog"

	"github.com/Luckny/recl/internal/database"
	"github.com/Luckny/recl/internal/parser"
)

type Handler interface {
	Handle(args []string, logger *slog.Logger) parser.Value
}

func InitHandlers(db *database.DB) map[string]Handler {
	return map[string]Handler{
		"PING": NewPingHandler(db),
		"SET":  NewSetHandler(db),
	}
}
