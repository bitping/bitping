package database

import (
	"log"

	"github.com/tidwall/buntdb"
)

// OpenOrCreateDatabase creates a buntdb at a specific path
func OpenOrCreateDatabase(path string) *buntdb.DB {
	db, err := buntdb.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// CloseDatabase will close the buntdb database if it exists
func CloseDatabase(db *buntdb.DB) {
	if db != nil {
		db.Close()
	}
}
