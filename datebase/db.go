// Package db implements simple method of treatment to bbolt db.
package db

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	bolt "go.etcd.io/bbolt"
)

// itoa convet int to string.
var itoa = strconv.Itoa

// DB implements interface access to bbolt db.
type DB struct {
	db *bolt.DB
}

// Bucket implements interface simple access to read/write in bbolt db.
type Bucket struct {
	db   *bolt.DB
	name string
}

// DB returns pointer to bbolt.DB.
func (b *Bucket) DB() *bolt.DB {
	return b.db
}

// Name returns string, name of Bucket.
func (b *Bucket) Name() string {
	return b.name
}

// Open return pointer to DB,
// If DB does not exist then error.
func Open(path string) (*DB, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Close implements access to close DB.
func (this *DB) Close() error {
	return this.db.Close()
}

// Bucket returns pointer to Bucket in db,
// Returns error if name is blank, or name is too long.
func (this *DB) Bucket(name string) (*Bucket, error) {
	err := this.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
	if err != nil {
		return nil, err
	}
	return &Bucket{this.db, name}, nil
}

// ExistsBucket returns true if bucket exists.
func (this *DB) ExistsBucket(name string) bool {
	var exists bool
	this.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		exists = (bucket != nil)
		return nil
	})
	return exists
}

// Set implements setting value of key in bucket.
func (this *Bucket) Set(key int, value string) error {
	return this.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(this.name))
		return bucket.Put([]byte(itoa(key)), []byte(value))
	})
}

// Get implements getting value of key in bucket.
func (this *Bucket) Get(key int) (string, error) {
	var value string
	err := this.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(this.name))
		value = string(bucket.Get([]byte(itoa(key))))
		if value == "" {
			return fmt.Errorf("Key `%v` is not exists", key)
		}
		return nil
	})
	return value, err
}

// GetOfField returns json-string of field in bucket.
func (this *Bucket) GetOfField(field string, value string) (string, error) {
	for inc := 1; true; inc++ {
		v, err := this.Get(inc)
		if err != nil {
			return "", err
		}

		var data fiber.Map
		err = json.Unmarshal([]byte(v), &data)
		if err != nil {
			return "ErrorJSON", err
		}

		if data[field] == nil {
			continue
		}

		if data[field].(string) == value {
			return v, nil
		}
	}
	return "", nil
}
