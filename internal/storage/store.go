package storage

import (
	"github.com/Luckny/recl/internal/aof"
)

type Store struct {
	db  *DB
	aof *aof.AOF
}

func NewStore(db *DB, aof *aof.AOF) *Store {
	return &Store{db: db, aof: aof}
}

func (s *Store) Set(key string, value string) {
	s.db.set(key, value)
}

func (s *Store) Get(key string) (*Val, bool) {
	return s.db.get(key)
}
