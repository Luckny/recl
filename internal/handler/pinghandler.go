package handler

import (
	"log/slog"

	"github.com/Luckny/recl/internal/database"
	"github.com/Luckny/recl/internal/parser"
)

type PingHandler struct {
	db *database.DB
}

func NewPingHandler(db *database.DB) *PingHandler {
	return &PingHandler{db: db}
}

func (h *PingHandler) Handle(args []string, logger *slog.Logger) parser.Value {
	return parser.SimpleStringValue("PONG")
}
