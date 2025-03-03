package main

import (
	"flag"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/dearing/go-keystore"
)

var initDB = flag.Bool("init", false, "initialize the database")
var database = flag.String("database", "mtg-cards.gob", "database file")

func main() {

	flag.Parse()

	db := keystore.NewCollection[int, Card]("my MTG cards")

	if *initDB {
		db.Clear()

		for i := range 1000 {
			db.Set(i, Card{
				Name:      fmt.Sprintf("Black Lotus %d", i),
				ManaCost:  "0",
				Type:      "Artifact",
				SubType:   "Lotus",
				Power:     i,
				Toughness: 0,
			})
		}

		for i := range 10000 {
			db.Set(rand.IntN(1000), Card{
				Name:      fmt.Sprintf("Black Lotus %d", i),
				ManaCost:  "0",
				Type:      "Artifact",
				SubType:   "Lotus",
				Power:     i,
				Toughness: 0,
			})
			db.Get(rand.IntN(1000))
		}

		db.Save(*database)
	}

	db.Load(*database)

	// for _, record := range db.Keys {
	// 	card := record.Value
	// 	fmt.Println(card.Name, card.Power, record.CreatedAt, record.ModifiedAt, record.AccessedAt, record.Writes, record.Reads)
	// }

	for i, r := range db.RecordsModifiedSince(time.Now().Add(7 * -time.Minute)) {
		slog.Info("Modified since", "record", r.ModifiedAt, "index", i)
	}

	//db.Save(*database)

}

type Card struct {
	Name      string
	ManaCost  string
	Type      string
	SubType   string
	Power     int
	Toughness int
}
