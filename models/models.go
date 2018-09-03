package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	sessionBucket = "sessions"
	userBucket    = "users"
	rpBucket      = "rp"
)

var db *bolt.DB
var err error

// ErrUsernameTaken is thrown when a user attempts to register a username that is taken.
var ErrUsernameTaken = errors.New("username already taken")

// Logger is a global logger used to show informational, warning, and error messages
var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

// Copy of auth.GenerateSecureKey to prevent cyclic import with auth library
func generateSecureKey() string {
	k := make([]byte, 32)
	io.ReadFull(rand.Reader, k)
	return fmt.Sprintf("%x", k)
}

// Setup initializes the Conn object
// It also populates the Config object
func Setup() error {
	openDB()

	err := db.Update(func(tx *bolt.Tx) error {
		registerKeyBucket, err := tx.CreateBucketIfNotExists([]byte(registerKeysBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		userBucket, err := tx.CreateBucketIfNotExists([]byte(userInfoBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
	})

	// Create the default user
	initUser := User{
		ID:          1,
		Name:        "admin",
		DisplayName: "Mr. Admin Face",
	}
	// Create the default relying party
	initRP := RelyingParty{
		ID:          "localhost",
		DisplayName: "Acme, Inc",
		Icon:        "lol.catpics.png",
		Users:       []User{initUser},
	}

	err = db.Save(&initRP)
	if err != nil {
		Logger.Println(err)
		return err
	}

	err = db.Save(&initUser)
	if err != nil {
		Logger.Println(err)
		return err
	}

	db.Init(&SessionData{})

	return nil
}
