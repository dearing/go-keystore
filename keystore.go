package keystore

import (
	"math"
	"time"
)

// Collection is a key-value store
type Collection struct {
	// Keys is a map of keys to records
	Keys map[any]Record
	// Description is a human-readable description of the collection
	Description string
}

// Record is a key-value store record
type Record struct {
	// Value is anything worth storing
	Value any

	// ContentType is string representation of the type of Value
	ContentType string
	// Description is a human-readable description of the record
	Description string

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

// NewCollection creates a new collection with a description
//
// ex: db := NewCollection("my MTG cards")
func NewCollection(description string) *Collection {
	return &Collection{
		Keys:        make(map[any]Record),
		Description: description,
	}
}

// Create a new record in the collection
//
// ex: db.Create("Black Lotus", Record{Value: "Black Lotus"})
func (c *Collection) Create(key any, value Record) {
	value.CreatedAt = time.Now()
	value.ModifiedAt = time.Now()
	value.AccessedAt = time.Now()
	if value.Writes >= math.MaxInt {
		value.Writes = 0 // reset Writes to 0 if it reaches the maximum
	}
	value.Writes++
	c.Keys[key] = value
}

// Read a record from the collection
//
// ex: record, ok := (db.Read("Black Lotus")) { ... }
func (c *Collection) Read(key any) (bool, Record) {
	if record, ok := c.Keys[key]; ok {
		record.AccessedAt = time.Now()
		if record.Reads >= math.MaxInt {
			record.Reads = 0 // reset Reads to 0 if it reaches the maximum
		}
		record.Reads++
		return true, record
	}
	return false, Record{}
}

// Update a record in the collection
// returns true if the record was updated, false if the record does not exist
//
// ex: if (db.Update("Black Lotus", Record{Value: "Blacker Lotus"})) { ... }
func (c *Collection) Update(key any, value Record) bool {
	if _, ok := c.Keys[key]; ok {
		value.ModifiedAt = time.Now()
		value.AccessedAt = time.Now()
		if value.Writes >= math.MaxInt {
			value.Writes = 0 // reset Writes to 0 if it reaches the maximum
		}
		value.Writes++
		c.Keys[key] = value
		return true
	}
	return false
}

// Delete a record from the collection
//
// ex: db.Delete("Black Lotus")
func (c *Collection) Delete(key any) bool {
	if _, ok := c.Keys[key]; ok {
		delete(c.Keys, key)
		return true
	}
	return false
}
