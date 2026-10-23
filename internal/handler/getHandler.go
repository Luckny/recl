package handler

import (
	"log/slog"

	"github.com/Luckny/recl/internal/parser"
	"github.com/Luckny/recl/internal/storage"
)

type GetHandler struct {
	store *storage.Store
}

func NewGetHandler(store *storage.Store) *GetHandler {
	return &GetHandler{store}
}

func (h *GetHandler) Handle(args []string, logger *slog.Logger) parser.Value {
	if len(args) != 1 {
		return parser.ErrorValue("wrong number of arguments for 'get' command")
	}

	key := args[0]

	val, ok := h.store.Get(key)
	if !ok {
		return parser.NullValue()
	}

	return parser.BulkStringValue(val.Data)
}
