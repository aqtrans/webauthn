package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/asdine/storm"

	"github.com/asdine/storm/q"
	"github.com/gorilla/sessions"
)

// SessionData is the Model
type SessionData struct {
	ID          uint   `json:"id" storm:"increment"`
	Challenge   []byte `json:"challenge"`
	Origin      string `json:"origin"`
	SessionType string `json:"session_type"`

	User   User `json:"user"`
	UserID uint `json:"user_id"`

	RelyingParty   RelyingParty `json:"rp"`
	RelyingPartyID string       `json:"rp_id"`
}

// ErrInvalidSessionType is thrown when an invalid session type is created
var ErrInvalidSessionType = errors.New("SessionType needs to be 'reg' or 'att'")

// CreateNewSession - Create new user/rp session
func CreateNewSession(u *User, rp *RelyingParty, st string) (SessionData, error) {
	ch, err := CreateChallenge(16)
	if err != nil {
		fmt.Println("Error Creating Challenge")
		return SessionData{}, err
	}

	if !(st == "reg" || st == "att") {
		return SessionData{}, ErrInvalidSessionType
	}

	sd := SessionData{
		Challenge:      ch,
		Origin:         rp.ID,
		UserID:         u.ID,
		RelyingPartyID: rp.ID,
		SessionType:    st,
	}

	err = PutSession(&sd)
	if err != nil {
		log.Println("PutSession error")
		return SessionData{}, err
	}

	return sd, nil
}

// GetSessionsByUsernameAndRelyingParty - Get the last recorded SessionData for a user/rp
func GetSessionsByUsernameAndRelyingParty(uid uint, rpid string) (SessionData, error) {
	sd := SessionData{}

	//err := db.Where("user_id = ? AND relying_party_id = ?", uid, rpid).Last(&sd).Error
	err := db.Select(q.Eq("user_id", uid), q.Eq("relying_party_id", rpid)).Reverse().First(&sd)
	return sd, err
}

// GetSessionData returns the SessionData that the given id corresponds to. If no user is found, an
// error is thrown.
func GetSessionData(id uint) (SessionData, error) {
	sd := SessionData{}
	//err := db.Where("id = ?", id).First(&sd).Error
	sds := []SessionData{}
	err := db.All(&sds)
	if err != nil {
		log.Println("GetSessionData:", err)
	}
	err = db.Select(q.Eq("id", id)).First(&sd)
	//err = db.One("id", id, &sd)
	if err != nil {
		return sd, err
	}
	/*
		//err = db.Model(&sd).Related(&sd.User).Error
		err = db.Select()
		if err != nil {
			fmt.Println("Error retrieving User data for session")
			return sd, err
		}
		err = db.Model(&sd).Related(&sd.RelyingParty).Error
		if err != nil {
			fmt.Println("Error retrieving RP data for session")
			return sd, err
		}
	*/
	return sd, nil
}

// GetSessionByUsername returns the user that the given username corresponds to. If no user is found, an
// error is thrown.
func GetSessionByUsername(username string) (User, error) {
	u := User{}
	//err := db.Where("username = ?", username).First(&u).Error
	err := db.One("username", username, &u)
	if err == storm.ErrNotFound {
		return u, nil
	}
	// No issue if we don't find a record
	/*
		if err == gorm.ErrRecordNotFound {
			return u, nil
		} else if err == nil {
			return u, ErrUsernameTaken
		}
	*/
	return u, err
}

// GetSessionForRequest gets the stored session data for a provided request.
func GetSessionForRequest(r *http.Request, store *sessions.CookieStore) (SessionData, error) {
	session, _ := store.Get(r, "registration-session")
	sessionID := session.Values["session_id"].(uint)
	sessionData, err := GetSessionData(sessionID)
	return sessionData, err
}

// PutSession - Update or Create SessionData
func PutSession(sd *SessionData) error {
	err := db.Save(sd)
	return err
}

// CreateChallenge - Create a new challenge to be sent to the authenticator
func CreateChallenge(len int) ([]byte, error) {
	challenge := make([]byte, len)
	_, err := rand.Read(challenge)
	if err != nil {
		return nil, err
	}
	return challenge, nil
}
