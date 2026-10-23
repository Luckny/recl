package handler

import (
	"log/slog"

	"github.com/Luckny/recl/internal/database"
	"github.com/Luckny/recl/internal/parser"
)

type SetHandler struct {
	db *database.DB
}

func NewSetHandler(db *database.DB) *SetHandler {
	return &SetHandler{db: db}
}

func (h *SetHandler) Handle(args []string, logger *slog.Logger) parser.Value {
	if len(args) != 2 {
		return parser.ErrorValue("wrong number of arguments for 'set' command")
	}

	return parser.NullValue()
}
