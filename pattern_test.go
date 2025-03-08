package keystore

import (
	"testing"
)

func TestRegexp(t *testing.T) {
	db := NewCollection[string]("my MTG cards")
	db.Set("card:1:black_lotus", "Black Lotus")
	db.Set("card:2:mox_pearl", "Mox Pearl")

	badRegexp := "card[0-9"
	_, err := db.MatchKeys(badRegexp)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	_, err = db.MatchValues(badRegexp)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	keys, err := db.MatchKeys("card:[0-9]:.")
	if err != nil {
		t.Error(err)
	}
	if len(keys) != 2 {
		t.Errorf("expected %d, got %d", 2, len(keys))
	}

	values, err := db.MatchValues("card:[0-9]:.")
	if err != nil {
		t.Error(err)
	}
	if len(values) != 2 {
		t.Errorf("expected %d, got %d", 2, len(values))
	}
}

func TestPrefix(t *testing.T) {
	db := NewCollection[string]("my MTG cards")
	db.Set("card:1:black_lotus", "Black Lotus")
	db.Set("card:2:mox_pearl", "Mox Pearl")

	keys, err := db.Prefix("card:2:mox_pearl")
	if err != nil {
		t.Error(err)
	}

	if len(keys) != 1 {
		t.Errorf("expected %d, got %d", 1, len(keys))
	}

	keys, err = db.Prefix("*:?:*")
	if err != nil {
		t.Error(err)
	}

	if len(keys) != 2 {
		t.Errorf("expected %d, got %d", 2, len(keys))
	}

	keysChan := db.PrefixChan("c?")
	count := 0
	for range keysChan {
		count++
	}
	if count != 2 {
		t.Errorf("expected %d, got %d", 2, count)
	}

}
