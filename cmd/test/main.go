package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/dearing/go-keystore"
)

var initDB = flag.Bool("init", false, "initialize the database")
var database = flag.String("database", "mtg-cards.gob", "database file")

func main() {

	flag.Parse()

	db := keystore.NewCollection[MTGCard]("Foundations")

	if *initDB {
		db.Clear()

		for _, set := range []string{"FDN.json", "FIN.json", "DFT.json"} {
			slog.Info("Loading set", "set", set)
			// load json from 'FIN.json' into a slice of MTGCard structs
			var cards []MTGCard
			file, err := os.Open(set)
			if err != nil {
				slog.Error("error opening file", "error", err)
				return
			}
			defer file.Close()

			decoder := json.NewDecoder(file)
			err = decoder.Decode(&cards)
			if err != nil {
				slog.Error("error decoding json", "set", set, "error", err)
				return
			}

			for _, card := range cards {
				key := fmt.Sprintf("/card/%s/%s", strings.ToLower(strings.TrimSuffix(set, ".json")), card.Name)
				db.Set(key, card)
				//slog.Info("Added card", "set", set, "key", key)
			}
		}

		db.Save(*database)
	}

	err := db.Load(*database)
	if err != nil {
		fmt.Println(err)
	}

	slog.Info("Database loaded", "size", db.Len())

	input := ""
	for {

		fmt.Printf("(%s) pattern:", input)
		fmt.Scanln(&input)

		start := time.Now()

		matches, error := db.Prefix(input)
		if error != nil {
			fmt.Println(error)
			continue
		}
		duration := time.Since(start)

		for i, card := range matches {
			fmt.Printf("item %d: '%s'\n", i, card)
		}

		slog.Info("results", "input", input, "len", len(matches), "duration", duration)

		start = time.Now()

		for key := range db.PrefixChan(input) {
			fmt.Printf("key: %s\n", key)
		}
		slog.Info("results(chan)", "input", input, "duration", time.Since(start))

	}

	// for i, r := range db.ModifiedSince(time.Now().Add(7 * -time.Minute)) {
	// 	slog.Info("Modified since", "record", r.ModifiedAt, "index", i)
	// }

	//db.Save(*database)

}

// pretty print MTGCard
func (c MTGCard) JSONString() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}

	return string(b)
}
