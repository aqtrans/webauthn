package models

import (
	"log"
)

// User represents the user model.
type User struct {
	ID             uint         `json:"id" storm:"id,increment"`
	Name           string       `json:"name"`
	DisplayName    string       `json:"display_name"`
	Icon           string       `json:"icon,omitempty"`
	Credentials    []Credential `json:"credentials,omitempty"`
	RelyingParties []RelyingParty
}

// GetUser returns the user that the given id corresponds to. If no user is found, an
// error is thrown.
func GetUser(id int64) (User, error) {
	u := User{}
	//err := db.Where("id=?", id).Preload("Credential").Find(&u).Error
	err := db.One("id", id, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetUserByUsername returns the user that the given username corresponds to. If no user is found, an
// error is thrown.
func GetUserByUsername(username string) (User, error) {
	u := User{}
	//err := db.Where("name = ?", username).Preload("Credentials").Find(&u).Error
	err := db.One("name", username, &u)

	if err == nil {
		return u, err
	}

	return User{}, err
}

// PutUser updates the given user
func PutUser(u *User) error {
	log.Println(u)
	err := db.Save(u)
	if err != nil {
		log.Println("PutUser error:", err)
	}
	return err
}
