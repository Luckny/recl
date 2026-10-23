// Package database provides a simple in-memory key-value store.
package database

type DB struct{}

func NewDB() *DB { return &DB{} }
