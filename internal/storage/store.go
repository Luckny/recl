package storage

import (
	"github.com/Luckny/recl/config"
	"github.com/Luckny/recl/internal/aof"
	"github.com/Luckny/recl/internal/parser"
)

type Store struct {
	db   *DB
	aof  *aof.AOF
	conf *config.ReclConfig
}

func NewStore(db *DB, aof *aof.AOF, conf *config.ReclConfig) *Store {
	return &Store{db: db, aof: aof, conf: conf}
}

func (s *Store) Set(key string, value string) {
	if s.aof != nil {
		command := parser.CommandValue("SET", key, value)
		s.aof.Write(command)

		if s.conf.AOF.Sync == config.AOFSyncModeAlways {
			s.aof.Flush()
		}
	}
	s.db.set(key, value)
}

func (s *Store) Get(key string) (*Val, bool) {
	return s.db.get(key)
}
