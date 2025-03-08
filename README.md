# go-keystore

A lightweight, in-memory key-value store for Go applications with type-safety through generics.

## features

- Generic type support - store any type in your collections
- Simple API for basic CRUD operations
- Advanced key filtering with regex pattern matching and wildcards
- Record metadata tracking (creation time, modification time, access time)
- Usage statistics (read/write counts)
- Persistence with gob encoding
- Time-based filtering for records

## installation

```bash
go get github.com/dearing/go-keystore
```

## coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## usage

```go
// Create a new string collection
db := keystore.NewCollection[string]("my string database")

// Set values
db.Set("key1", "value1")
db.Set("key2", "value2")

// Get values
value, exists := db.Get("key1")
if exists {
    fmt.Println(value) // "value1"
}

// Get full record with metadata
record, exists := db.GetRecord("key1")
if exists {
    fmt.Println(record.Value)
    fmt.Println(record.CreatedAt)
    fmt.Println(record.ModifiedAt)
    fmt.Println(record.Reads)
}

// Find keys matching a regex pattern
matches, _ := db.MatchKeys("key[0-9]+")

// Find values matching a regex pattern
values, _ := db.MatchValues("value[0-9]+")

// Get keys with a wildcard prefix
keys, _ := db.Prefix("user_*")

// Stream keys with a wildcard prefix
for key := range db.PrefixChan("user_*") {
    fmt.Println(key)
}

// Time-based filtering
recentlyCreated := db.CreatedSince(time.Now().Add(-24 * time.Hour))
recentlyModified := db.ModifiedSince(time.Now().Add(-1 * time.Hour))
recentlyAccessed := db.AccessedSince(time.Now().Add(-30 * time.Minute))

// Save collection to file
err := db.Save("database.gob")
if err != nil {
    log.Fatal(err)
}

// Load collection from file
newDB := keystore.NewCollection[string]("loaded database")
err = newDB.Load("database.gob")
if err != nil {
    log.Fatal(err)
}

// Delete a value
db.Delete("key2")

// Clear the entire collection
db.Clear()

type Card struct {
    Name      string
    ManaCost  string
    Type      string
    Power     int
    Toughness int
}

// Create a collection for your custom type
cardDB := keystore.NewCollection[Card]("my card collection")

// Add cards to the collection
cardDB.Set("black-lotus", Card{
    Name:     "Black Lotus",
    ManaCost: "0",
    Type:     "Artifact",
})

// Retrieve a card
card, exists := cardDB.Get("black-lotus")
```

