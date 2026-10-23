package handler

import (
	"log/slog"

	"github.com/Luckny/recl/internal/database"
	"github.com/Luckny/recl/internal/parser"
)

type GetHandler struct {
	db *database.DB
}

func NewGetHandler(db *database.DB) *GetHandler {
	return &GetHandler{db: db}
}

func (h *GetHandler) Handle(args []string, logger *slog.Logger) parser.Value {
	if len(args) != 1 {
		return parser.ErrorValue("wrong number of arguments for 'get' command")
	}

	key := args[0]

	val, ok := h.db.Get(key)
	if !ok {
		return parser.NullValue()
	}

	return parser.BulkStringValue(val.Data)
}
