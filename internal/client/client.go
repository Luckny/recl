// Package client represents a single client connection to the Recl server.
package client

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Luckny/recl/internal/parser"
)

type Client struct {
	Conn   net.Conn
	Name   string
	id     uint64 // private for now
	Parser *parser.Parser
	logger *slog.Logger
}

func NewClient(conn net.Conn, id uint64, logger *slog.Logger) *Client {
	name := fmt.Sprintf("connection-%d", id)
	return &Client{
		Conn:   conn,
		Name:   name,
		id:     id,
		Parser: parser.NewParser(logger),
		logger: logger,
	}
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
