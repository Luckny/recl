// Package server implements the server.
package server

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/Luckny/recl/config"
	"github.com/Luckny/recl/internal/aof"
	"github.com/Luckny/recl/internal/client"
	"github.com/Luckny/recl/internal/handler"
	"github.com/Luckny/recl/internal/parser"
	"github.com/Luckny/recl/internal/storage"
)

type Server struct {
	config   *config.ReclConfig
	logger   *slog.Logger
	aof      *aof.AOF
	Handlers map[string]handler.Handler
	store    *storage.Store
	db       *storage.DB
}

func NewServer(conf *config.ReclConfig, logger *slog.Logger) (*Server, error) {
	// create data dir if not created
	err := os.MkdirAll(conf.DataDir, 0o755)
	if err != nil {
		return nil, err
	}

	db := storage.NewDB()
	store := storage.NewStore(db, nil)

	server := &Server{
		config:   conf,
		logger:   logger,
		store:    store,
		Handlers: handler.InitHandlers(store),
	}

	if conf.AOF.Enabled {
		aofStore, err := aof.NewAOF(conf.DataDir, &conf.AOF)
		if err != nil {
			return nil, err
		}

		if ok := store.SetAofStore(aofStore); !ok {
			return nil, errors.New("failed to set aof store")
		}

		commands, err := aofStore.GetSyncValueFromDisk(logger)
		if err != nil {
			return nil, err
		}

		for _, command := range commands {
			server.handleValue(command)
		}
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
			value, err := client.Parser.ParseNext()
			if errors.Is(err, parser.ErrIncomplete) {
				break
			}
			if err != nil {
				s.logger.Error("error parsing command", "error", err)
				return
			}

			response := s.handleValue(value)
			if err := client.WriteValue(response); err != nil {
				return
			}

		}
	}
}

func (s *Server) handleValue(v parser.Value) parser.Value {
	// should arrive as array
	if v.Type != parser.ArrayType {
		return parser.ErrorValue("expected an array")
	}

	if len(v.Array) == 0 {
		return parser.ErrorValue("empty command")
	}

	args := make([]string, len(v.Array))
	for i, arg := range v.Array {
		if arg.Type != parser.BulkType {
			return parser.ErrorValue("expected a bulk string")
		}
		args[i] = arg.Bulk
	}

	command := strings.ToUpper(args[0])
	handler, ok := s.Handlers[command]
	if !ok {
		return parser.ErrorValue("invalid command")
	}

	return handler.Handle(args[1:], s.logger)
}
