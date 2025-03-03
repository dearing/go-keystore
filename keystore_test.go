package keystore

import (
	"math"
	"testing"
)

func TestNewCollectionString(t *testing.T) {
	db := NewCollection("my MTG cards")
	if db.Description != "my MTG cards" {
		t.Errorf("expected %s, got %s", "my MTG cards", db.Description)
	}

	if len(db.Keys) != 0 {
		t.Errorf("expected %d, got %d", 0, len(db.Keys))
	}

	card := Record{Value: "Black Lotus"}
	db.Create("Black Lotus", card)

	if len(db.Keys) != 1 {
		t.Errorf("expected %d, got %d", 1, len(db.Keys))
	}

	ok, record := db.Read("Black Lotus")
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	if record.Value != "Black Lotus" {
		t.Errorf("expected %s, got %s", "Black Lotus", record.Value)
	}

	if record.Reads != 1 {
		t.Errorf("expected %d, got %d", 1, record.Reads)
	}

	if record.Writes != 1 {
		t.Errorf("expected %d, got %d", 1, record.Writes)
	}

	if record.CreatedAt.IsZero() {
		t.Errorf("expected %s, got %s", "not zero", record.CreatedAt)
	}

	if record.ModifiedAt.IsZero() {
		t.Errorf("expected %s, got %s", "not zero", record.ModifiedAt)
	}

	if record.AccessedAt.IsZero() {
		t.Errorf("expected %s, got %s", "not zero", record.AccessedAt)
	}

	ok, record = db.Read("Black Lotus")
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	if record.Reads != 1 {
		t.Errorf("expected %d, got %d", 1, record.Reads)
	}

	if record.Writes != 1 {
		t.Errorf("expected %d, got %d", 1, record.Writes)
	}

	if record.CreatedAt.IsZero() {
		t.Errorf("expected %s, got %s", "not zero", record.CreatedAt)
	}

	if record.ModifiedAt.IsZero() {
		t.Errorf("expected %s, got %s", "not zero", record.ModifiedAt)
	}

	if record.AccessedAt.IsZero() {
		t.Errorf("expected %s, got %s", "not zero", record.AccessedAt)
	}

	cardName := "Blacker Lotus"

	ok = db.Update("Black Lotus", Record{Value: cardName})
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	ok, record = db.Read("Black Lotus")
	if ok {
		if record.Value != cardName {
			t.Errorf("expected %s, got %s", cardName, record.Value)
		}
	}
}

func TestNewCollectionComplex(t *testing.T) {
	db := NewCollection("players")
	type player struct {
		Username string
		Password string
	}

	alice := &player{Username: "alice", Password: "password"}
	db.Create("alice", Record{Value: alice})

	ok, record := db.Read("alice")
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	if record.Value.(*player).Username != "alice" {
		t.Errorf("expected %s, got %s", "alice", record.Value.(*player).Username)
	}

	if record.Value.(*player).Password != "password" {
		t.Errorf("expected %s, got %s", "password", record.Value.(*player).Password)
	}

}

func TestNewCollectionOverFlow(t *testing.T) {
	db := NewCollection("stars")
	star := &Record{
		Value:  "Sun",
		Reads:  math.MaxInt,
		Writes: math.MaxInt,
	}

	db.Create("Sun", *star)

	ok, record := db.Read("Sun")
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	if record.Reads != 1 {
		t.Errorf("expected %d, got %d", 0, record.Reads)
	}

	if record.Writes != 1 {
		t.Errorf("expected %d, got %d", 1, record.Writes)
	}

	record.Reads = math.MaxInt
	record.Writes = math.MaxInt

	if db.Update("Sun", record); !ok {
		t.Errorf("expected %t, got %t", false, ok)
	}

	ok, record = db.Read("Sun")
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	if record.Reads != 1 {
		t.Errorf("expected %d, got %d", 0, record.Reads)
	}

	if record.Writes != 1 {
		t.Errorf("expected %d, got %d", 1, record.Writes)
	}

}
