package handler

import (
	"log/slog"

	"github.com/Luckny/recl/internal/parser"
	"github.com/Luckny/recl/internal/storage"
)

type PingHandler struct {
	store *storage.Store
}

func NewPingHandler(store *storage.Store) *PingHandler {
	return &PingHandler{store: store}
}

func (h *PingHandler) Handle(args []string, logger *slog.Logger) parser.Value {
	return parser.SimpleStringValue("PONG")
}
