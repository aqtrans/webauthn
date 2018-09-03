package models

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// Open opens the boltDB, closing after use
func openDB() {
	db, err = bolt.Open("webauthn.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		log.Fatalln("openDB() error:", err)
	}
}
