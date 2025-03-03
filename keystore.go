package keystore

import (
	"time"
)

type Collection struct {
	Keys        map[any]Record
	Description string
}

type Record struct {
	Value any

	ContentType string
	Description string

	CreatedAt  time.Time
	ModifiedAt time.Time
	AccessedAt time.Time
	ExpiresAt  time.Time

	Writes int
	Reads  int
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
	value.Writes++
	c.Keys[key] = value
}

// Read a record from the collection
//
// ex: record, ok := (db.Read("Black Lotus")) { ... }
func (c *Collection) Read(key any) (bool, Record) {
	if record, ok := c.Keys[key]; ok {
		record.AccessedAt = time.Now()
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
