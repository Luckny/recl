// Package client represents a single client connection to the Recl server.
package client

import (
	"bytes"
	"fmt"
	"io"
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

func (c *Client) WriteValue(value parser.Value) error {
	resp, err := parser.RESPValue(value)
	if err != nil {
		return err
	}
	_, err = io.Copy(c.Conn, bytes.NewReader(resp))
	if err != nil {
		c.logger.Error("failed to write value", "error", err)
	}

	return nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
