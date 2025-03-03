package keystore

import (
	"encoding/gob"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"time"
)

// Collection is a key-value store
type Collection[K comparable, T any] struct {
	// Keys is a map of keys to records
	Keys map[K]Record[T]
	// Description is a human-readable description of the collection
	Description string
}

// Record is a key-value store record
type Record[T any] struct {
	// Value is anything worth storing
	Value T

	// ContentType is string representation of the type of Value
	ContentType string

	// CreatedAt is the time the record was created
	CreatedAt time.Time
	// ModifiedAt is the time the record was last modified
	ModifiedAt time.Time
	// AccessedAt is the time the record was last accessed
	AccessedAt time.Time

	// Writes is the number of times the record was written
	Writes int
	// Reads is the number of times the record was read
	Reads int
}

// NewRecord creates a new record with a value, content type, and description
//
// ex: record := NewRecord(mtgCard, "Card", "A powerful card")
func NewRecord[T any](value T) Record[T] {
	now := time.Now()
	return Record[T]{
		Value:       value,
		ContentType: reflect.TypeOf(value).String(),
		CreatedAt:   now,
		ModifiedAt:  now,
		AccessedAt:  now,
		Writes:      0,
		Reads:       0,
	}
}

// NewCollection creates a new collection with a description
//
// ex: db := NewCollection("my MTG cards")
func NewCollection[K comparable, T any](description string) Collection[K, T] {
	//slog.Info("Creating new collection", "description", description)
	return Collection[K, T]{
		Keys:        make(map[K]Record[T]),
		Description: description,
	}
}

// Set a record in the collection
//
// If the record already exists, it will be updated with write plus one,
// if the record does not exist, it will be created with write at zero.
//
// ex: db.Set("Black Lotus", "Black Lotus")
func (c *Collection[K, T]) Set(key K, value T) {
	if record, exists := c.Keys[key]; exists {
		record.Value = value
		record.ModifiedAt = time.Now()
		record.Writes++
		c.Keys[key] = record
		//slog.Info("Updating record", "value", value)
	} else {
		//slog.Info("Setting new record", "value", value)
		c.Keys[key] = NewRecord(value)
	}
}

// Read a record from the collection
//
// ex: record, ok := db.Read("Black Lotus")
func (c *Collection[K, T]) Get(key K) (T, bool) {
	//slog.Info("Reading value", "key", key)
	if record, exists := c.Keys[key]; exists {
		record.AccessedAt = time.Now()
		record.Reads++
		c.Keys[key] = record
		return record.Value, true
	}

	var zero T
	return zero, false
}

// GetRecord retrieves the full record from the collection
//
// ex: record, ok := db.GetRecord("Black Lotus")
func (c *Collection[K, T]) GetRecord(key K) (Record[T], bool) {
	//slog.Info("Getting record", "key", key)
	record, exists := c.Keys[key]
	if exists {
		record.AccessedAt = time.Now()
		record.Reads++
		c.Keys[key] = record
	}
	return record, exists
}

// Delete removes a record from the collection
//
// ex: db.Delete("Black Lotus")
func (c *Collection[K, T]) Delete(key K) {
	//slog.Info("Deleting record", "key", key)
	delete(c.Keys, key)
}

// Clear the collection
//
// ex: db.Clear()
func (c Collection[K, T]) Clear() {
	//slog.Info("Clearing collection")
	clear(c.Keys)
}

// Len returns the number of records in the collection
//
// ex: size := db.Len()
func (c Collection[K, T]) Len() int {
	//slog.Info("Getting collection size", "len", len(c.Keys))
	return len(c.Keys)
}

// // All returns a sequence of all records in the collection
// //
// // ex: for record := range db.All() { ... }
// func (c Collection[K, T]) All() iter.Seq[T] {
// 	return func(yield func(T) bool) {
// 		for _, record := range c.Keys {
// 			if !yield(record) {
// 				return
// 			}
// 		}
// 	}
// }

func (c Collection[K, T]) Save(fileName string) error {
	slog.Info("Saving collection")

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(c)
	if err != nil {
		return fmt.Errorf("error encoding collection: %w", err)
	}
	return nil
}

func (c Collection[K, T]) Load(fileName string) error {
	slog.Info("Loading collection")

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()

	dec := gob.NewDecoder(file)
	err = dec.Decode(&c)
	if err != nil {
		return fmt.Errorf("error decoding collection: %w", err)
	}
	return nil

}

func (c Collection[K, T]) RecordsCreatedSince(time time.Time) []Record[T] {
	var records []Record[T]
	for _, record := range c.Keys {
		if record.CreatedAt.After(time) {
			records = append(records, record)
		}
	}
	return records
}

func (c Collection[K, T]) RecordsModifiedSince(time time.Time) []Record[T] {
	var records []Record[T]
	for _, record := range c.Keys {
		if record.ModifiedAt.After(time) {
			records = append(records, record)
		}
	}
	return records
}

func (c Collection[K, T]) RecordsAccessedSince(time time.Time) []Record[T] {
	var records []Record[T]
	for _, record := range c.Keys {
		if record.AccessedAt.After(time) {
			records = append(records, record)
		}
	}
	return records
}
