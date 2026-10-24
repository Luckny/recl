package storage

import (
	"github.com/Luckny/recl/internal/aof"
	"github.com/Luckny/recl/internal/parser"
)

type Store struct {
	db  *DB
	aof *aof.AOF
}

func NewStore(db *DB, aof *aof.AOF) *Store {
	return &Store{db: db, aof: aof}
}

func (s *Store) SetAofStore(aof *aof.AOF) (ok bool) {
	s.aof = aof
	return true
}

func (s *Store) Set(key string, value string) {
	if s.aof != nil {
		command := parser.CommandValue("SET", key, value)
		s.aof.Write(command)
	}
	s.db.set(key, value)
}

func (s *Store) Get(key string) (*Val, bool) {
	return s.db.get(key)
}
