package keystore

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

// errorWriter is a simulated writer that always returns an error
type errorWriter struct{}

// Write always returns an error
func (ew errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("simulated write error")
}

func TestNewCollectionTimeGroups(t *testing.T) {
	db := NewCollection[string]("my MTG cards")
	db.Set("c1", "Black Lotus")

	mark := time.Now()

	db.Set("c2", "Mox Pearl")

	if db.Len() != 2 {
		t.Errorf("expected %d, got %d", 2, db.Len())
	}

	// should only get back the Mox Pearl record
	for _, record := range db.CreatedSince(mark) {
		if record != "Mox Pearl" {
			t.Errorf("expected %s, got %s", "Mox Pearl", record)
		}
	}
	// should only get back the Mox Pearl record
	for _, record := range db.ModifiedSince(mark) {
		if record != "Mox Pearl" {
			t.Errorf("expected %s, got %s", "Mox Pearl", record)
		}
	}
	// should only get back the Mox Pearl record
	for _, record := range db.AccessedSince(mark) {
		if record != "Mox Pearl" {
			t.Errorf("expected %s, got %s", "Mox Pearl", record)
		}
	}

}

func TestNewCollectionPersist(t *testing.T) {
	db := NewCollection[string]("my MTG cards")

	if db.Len() != 0 {
		t.Errorf("expected %d, got %d", 0, db.Len())
	}

	db.Set("c1", "Black Lotus")
	db.Set("c2", "Mox Pearl")

	dataStore := bytes.Buffer{}

	err := db.Load(&dataStore)
	if err == nil {
		t.Errorf("encoding should throw and error on empty data store")
	}

	err = db.Save(&dataStore)
	if err != nil {
		t.Error(err)
	}

	db.Clear()

	if db.Len() != 0 {
		t.Errorf("expected %d, got %d", 0, db.Len())
	}

	err = db.Load(&dataStore)
	if err != nil {
		t.Error(err)
	}

	ew := errorWriter{}

	err = db.Save(ew)
	if err == nil {
		t.Error(err)
	}
}

func TestNewCollectionBasic(t *testing.T) {

	card1 := "Black Lotus"
	card2 := "Mox Pearl"

	db := NewCollection[string]("my MTG cards")

	// did the description get set?
	if db.Description != "my MTG cards" {
		t.Errorf("expected %s, got %s", "my MTG cards", db.Description)
	}

	// database should be empty
	if len(db.Keys) != 0 {
		t.Errorf("expected %d, got %d", 0, len(db.Keys))
	}

	// create a record
	db.Set("c1", card1)

	// database should have 1 record
	if len(db.Keys) != 1 {
		t.Errorf("expected %d, got %d", 1, len(db.Keys))
	}

	record, ok := db.Get("c1")

	// did we get the record back?
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	// did we get the right record back?
	if record != card1 {
		t.Errorf("expected %s, got %s", card1, record)
	}

	// attempt to get a non-existent record
	record, ok = db.Get("should not exist")
	if ok {
		t.Errorf("expected %t, got %t", false, ok)
	}

	// we should have an empty record on no match
	if record != "" {
		t.Errorf("expected %s, got %s", "", record)
	}

	// get a full record
	fullRecord, ok := db.GetRecord("c1")
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	// did we get the right record back?
	if fullRecord.Value != card1 {
		t.Errorf("expected %s, got %s", card1, fullRecord.Value)
	}

	// writes should be zero because it was a new record and not re-set
	if fullRecord.Writes != 0 {
		t.Errorf("expected %d, got %d", 1, fullRecord.Writes)
	}

	db.Set("c1", "Black Lotus (reprint)")

	// get a full record
	fullRecord, ok = db.GetRecord("c1")
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	// created should be set
	if fullRecord.CreatedAt.IsZero() {
		t.Errorf("expected %t, got %t", false, fullRecord.CreatedAt.IsZero())
	}

	// modified should be set
	if fullRecord.ModifiedAt.IsZero() {
		t.Errorf("expected %t, got %t", false, fullRecord.ModifiedAt.IsZero())
	}

	// accessed should be set
	if fullRecord.AccessedAt.IsZero() {
		t.Errorf("expected %t, got %t", false, fullRecord.AccessedAt.IsZero())
	}

	// create a new record
	db.Set("c2", card2)

	// database should have 2 records
	if len(db.Keys) != 2 {
		t.Errorf("expected %d, got %d", 2, len(db.Keys))
	}

	// delete a record
	db.Delete("c2")

	// database should have 1 record
	if len(db.Keys) != 1 {
		t.Errorf("expected %d, got %d", 1, len(db.Keys))
	}

	// attempt to get a deleted record
	record, ok = db.Get("c2")
	if ok {
		t.Errorf("expected %t, got %t", false, ok)
	}

	// we should have an empty record on no match
	if record != "" {
		t.Errorf("expected %s, got %s", "", record)
	}

	// clear the database
	db.Clear()

	// database should be empty
	if len(db.Keys) != 0 {
		t.Errorf("expected %d, got %d", 0, len(db.Keys))
	}

	db.Set("c1", card1)

	// database should have 1 record
	if len(db.Keys) != 1 {
		t.Errorf("expected %d, got %d", 1, len(db.Keys))
	}

}

// benchmark the Set method
func BenchmarkSetString(b *testing.B) {

	type Card struct {
		Name      string
		ManaCost  string
		Type      string
		SubType   string
		Power     int
		Toughness int
	}
	card := Card{
		Name:      "Black Lotus",
		ManaCost:  "0",
		Type:      "Artifact",
		SubType:   "Lotus",
		Power:     0,
		Toughness: 0,
	}

	db := NewCollection[Card]("my MTG cards")

	for b.Loop() {
		name := "Black Lotus"
		db.Set(name, card)
		_, _ = db.Get(name)
		db.Delete(name)
	}

	db.Clear()

}
