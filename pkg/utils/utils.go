package utils

import (
	"encoding/base32"
	"time"

	bolt "go.etcd.io/bbolt"
)

type ID []byte

var Generator func() string

func (id ID) String() string {
	return base32.HexEncoding.EncodeToString(id)
}

func NextID() ID {
	return ID(Generator())
}
const initialMmapSize = 10 * 1 << 30

func Open(path string) (*bolt.DB, error) {
	opts := &bolt.Options{
		Timeout:         10 * time.Second,
		InitialMmapSize: initialMmapSize,
	}
	db, err := bolt.Open(path, 0600, opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}
