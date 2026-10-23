// Package server implements the server.
package server

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"os"

	"github.com/Luckny/recl/config"
	"github.com/Luckny/recl/internal/aof"
	"github.com/Luckny/recl/internal/client"
	"github.com/Luckny/recl/internal/parser"
)

type Server struct {
	config *config.ReclConfig
	logger *slog.Logger
	aof    *aof.AOF
}

func NewServer(conf *config.ReclConfig, logger *slog.Logger) (*Server, error) {
	// create data dir if not created
	err := os.MkdirAll(conf.DataDir, 0o755)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config: conf,
		logger: logger,
	}

	if conf.AOF.Enabled {
		aof, err := aof.NewAOF(conf.DataDir, &conf.AOF)
		if err != nil {
			return nil, err
		}

		aof.SyncFromDisk(logger)
		server.aof = aof
	}
	return server, nil
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		s.logger.Error("failed to listen", "error", err)
		return err
	}

	defer listener.Close()
	s.logger.Info("server started", "address", listener.Addr())

	var nextConnectionID uint64
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("failed to accept connection", "error", err)
			return err
		}

		nextConnectionID++
		go s.handleConnection(conn, nextConnectionID)
	}
}

func (s *Server) handleConnection(conn net.Conn, connectionID uint64) {
	client := client.NewClient(conn, connectionID, s.logger)
	s.logger.Info("new client connected", "id", client.Name)
	defer client.Close()

	buf := make([]byte, 4096)

	for {

		consumedLen, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				s.logger.Error("error reading from client", "name", client.Name, "error", err)
			}
			return
		}
		client.Parser.Append(buf[:consumedLen])

		for {
			_, err := client.Parser.ParseNext()
			if errors.Is(err, parser.ErrIncomplete) {
				break
			}
			if err != nil {
				s.logger.Error("error parsing command", "error", err)
				return
			}

		}
	}
}
