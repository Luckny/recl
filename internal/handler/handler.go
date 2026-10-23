// Package handler
package handler

import (
	"log/slog"

	"github.com/Luckny/recl/internal/parser"
	"github.com/Luckny/recl/internal/storage"
)

type Handler interface {
	Handle(args []string, logger *slog.Logger) parser.Value
}

func InitHandlers(store *storage.Store) map[string]Handler {
	return map[string]Handler{
		"PING": NewPingHandler(store),
		"SET":  NewSetHandler(store),
		"GET":  NewGetHandler(store),
	}
}
