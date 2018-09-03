package models

import (
	"log"
)

// User represents the user model.
type User struct {
	ID             int64        `json:"id" storm:"id,increment"`
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
	userDB := db.From("users")
	err := userDB.One("ID", id, &u)
	/*
		users := []User{}
		err := db.All(&users)
		if err == nil {
			for _, v := range users {
				if v.ID == id {
					u = v
				}
			}
			return u, err
		}
	*/
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
	userDB := db.From("users")
	err := userDB.One("Name", username, &u)
	/*
		users := []User{}
		err := db.All(&users)
		if err == nil {
			for _, v := range users {
				if v.Name == username {
					u = v
				}
			}
			return u, err
		}
	*/
	log.Println("User:", u)

	return u, err
}

// PutUser updates the given user
func PutUser(u *User) error {
	log.Println(u)
	userDB := db.From("users")
	err := userDB.Save(u)
	if err != nil {
		log.Println("PutUser error:", err)
	}
	return err
}
