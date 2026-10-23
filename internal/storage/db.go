// Package storage provides a simple in-memory key-value store.
package storage

import "sync"

type Val struct {
	Data string
}

type DB struct {
	store map[string]*Val
	mu    sync.RWMutex
}

func NewDB() *DB {
	return &DB{
		store: map[string]*Val{},
		mu:    sync.RWMutex{},
	}
}

func (d *DB) set(key string, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.store[key] = &Val{Data: value}
}

func (d *DB) get(key string) (*Val, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	val, ok := d.store[key]
	return val, ok
}
