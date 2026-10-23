package handler

import (
	"log/slog"

	"github.com/Luckny/recl/internal/parser"
	"github.com/Luckny/recl/internal/storage"
)

type SetHandler struct {
	store *storage.Store
}

func NewSetHandler(store *storage.Store) *SetHandler {
	return &SetHandler{store}
}

func (h *SetHandler) Handle(args []string, logger *slog.Logger) parser.Value {
	if len(args) != 2 {
		return parser.ErrorValue("wrong number of arguments for 'set' command")
	}

	key := args[0]
	value := args[1]

	h.store.Set(key, value)

	return parser.SimpleStringValue("OK")
}
