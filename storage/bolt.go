package storage

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"os"
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

type TestDB struct {
	*BoltDB
}

// NewTestDB returns a TestDB using a temporary path.
func NewTestDB() *TestDB {
	// Retrieve a temporary path.
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(fmt.Sprintf("temp file: %s", err))
	}
	path := f.Name()
	f.Close()
	os.Remove(path)
	// Open the database.
	db, err := NewBoltDB(path)
	if err != nil {
		panic(fmt.Sprintf("temp db: %s", err))
	}
	// Return wrapped type.
	return &TestDB{db}
}

// Close and delete Bolt database.
func (db *TestDB) Close() {
	defer os.Remove(db.Path())
	db.DB.Close()
}