package main

import (
	"encoding/binary"
	"errors"
	"github.com/cockroachdb/pebble"
	"io"
	"testing"
)

func BigEndian(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

// todo: 用于保存session的历史信息
func Test_pebble(t *testing.T) {
	db, err := pebble.Open("test-db", &pebble.Options{})
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func Get(db *pebble.DB, key []byte) []byte {
	retrievedValue, closer, err := db.Get(key)
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			panic(err)
		}
	}(closer)
	if errors.Is(err, pebble.ErrNotFound) {
		return nil
	}
	panic(err)
	return retrievedValue
}

func Set(db *pebble.DB, key []byte, val []byte) {
	if err := db.Set(key, val, pebble.Sync); err != nil {
		panic(err)
	}
}

func Delete(db *pebble.DB, key []byte) {
	if err := db.Delete(key, pebble.Sync); err != nil {
		panic(err)
	}
}
