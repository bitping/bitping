package database

import (
	"log"

	"github.com/tidwall/buntdb"
)

var defaultDatabasePath = "/tmp/database/bitping.db"

// GetDatabasePath returns the default database path
func GetDatabasePath() string {
	return defaultDatabasePath
}

// OpenOrCreateDatabase creates a buntdb at a specific path
func OpenOrCreateDatabase(path *string) *buntdb.DB {
	if path == nil {
		path = &defaultDatabasePath
	}

	db, err := buntdb.Open(*path)
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
