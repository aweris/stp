package storage

import (
	bolt "go.etcd.io/bbolt"
)

// DB wrapper for bolt.DB
type BoltDB struct {
	*bolt.DB
}

// NewBoltDB returns a BoltDB wrapper
func NewBoltDB(path string) (*BoltDB, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &BoltDB{db}, nil
}
