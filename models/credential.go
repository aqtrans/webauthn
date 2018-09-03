package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/asdine/storm"

	"github.com/asdine/storm/q"
)

// Credential is the stored credential for Auth
type Credential struct {
	ID        int64 `json:"id" storm:"id,increment"`
	CreatedAt time.Time
	Counter   []byte `json:"sign_count"`

	RelyingParty   RelyingParty `json:"rp"`
	RelyingPartyID string       `json:"rp_id"`

	User   User  `json:"user"`
	UserID int64 `json:"user_id"`

	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
	Flags  []byte `json:"flags,omitempty"`

	CredID string `json:"credential_id,omitempty"`

	PublicKey PublicKey `json:"public_key,omitempty"`
}

// PublicKey is parsed from the credential creation response
type PublicKey struct {
	_struct      bool   `codec:",int"`
	KeyType      int8   `codec:"1"`
	Type         int8   `codec:"3"`
	XCoord       []byte `codec:"-2"`
	YCoord       []byte `codec:"-3"`
	Curve        int8   `codec:"-1"`
	CredentialID int64  `storm:"id" codec:"-,omitempty"`
}

// CreateCredential creates a new credential object
func CreateCredential(c *Credential) error {
	fmt.Println("Creating Credential")
	_, err := GetCredentialForUserAndRelyingParty(&c.User, &c.RelyingParty)
	if err == nil {
		log.Println("CreateCredential error: credentials already exist")
		return errors.New("Credentials already exist")
	}
	if err != storm.ErrNotFound {
		log.Println("CreateCredential other error:", err)
		return err
	}

	creds := db.From("credentials")

	err = creds.Save(c)
	if err != nil {
		log.Println("CreateCredential/db.Save error:", err)
		return err
	}
	return err
}

// UpdateCredential updates the credential with new attributes.
func UpdateCredential(c *Credential) error {
	creds := db.From("credentials")
	err = creds.Save(c)
	return err
}

// GetCredentialForUserAndRelyingParty retrieves the first credential for a provided user and relying party.
func GetCredentialForUserAndRelyingParty(user *User, rp *RelyingParty) (Credential, error) {
	cred := Credential{}
	//err := db.Where("user_id = ? AND relying_party_id = ?", user.ID, rp.ID).Preload("PublicKey").First(&cred).Error
	creds := db.From("credentials")
	err := creds.Select(q.Eq("UserID", user.ID), q.Eq("RelyingPartyID", rp.ID)).First(&cred)
	if err == storm.ErrNotFound {
		return Credential{}, err
	}
	if err != nil {
		log.Println("GetCredentialForUserAndRelyingParty err:", err)
		return Credential{}, err
	}
	cred.User = *user
	cred.RelyingParty = *rp

	return cred, err
}

// GetCredentialsForUserAndRelyingParty retrieves all credentials for a provided user for a relying party.
func GetCredentialsForUserAndRelyingParty(user *User, rp *RelyingParty) ([]Credential, error) {
	creds := []Credential{}
	//err := db.Where("user_id = ? AND relying_party_id = ?", user.ID, rp.ID).Preload("PublicKey").Find(&creds).Error
	credsDB := db.From("credentials")
	err := credsDB.All(&creds)
	//err := credsDB.Select(q.Eq("UserID", user.ID), q.Eq("RelyingPartyID", rp.ID)).Find(&creds)
	if err != nil {
		log.Println("GetCredentialsForUserAndRelyingParty err:", err)
		return []Credential{}, err
	}

	for _, cred := range creds {
		if user.ID == cred.User.ID && rp.ID == cred.RelyingPartyID {
			creds = append(creds, cred)
		}
	}

	return creds, nil
}

// GetCredentialsForUser retrieves all credentials for a provided user regardless of relying party.
func GetCredentialsForUser(user *User) ([]Credential, error) {
	creds := []Credential{}
	//err := db.Where("user_id = ?", user.ID).Preload("PublicKey").Find(&creds).Error
	credsDB := db.From("credentials")
	err := credsDB.Find("UserID", user.ID, &creds)
	return creds, err
}

// GetCredentialForUser retrieves a specific credential for a user.
func GetCredentialForUser(user *User, credentialID string) (Credential, error) {
	cred := Credential{}
	//err := db.Where("user_id = ? AND cred_id = ?", user.ID, credentialID).Preload("PublicKey").Find(&cred).Error
	credsDB := db.From("credentials")
	err := credsDB.Select(q.Eq("UserID", user.ID), q.Eq("CredID", credentialID)).First(&cred)
	return cred, err
}

// DeleteCredentialByID gets a credential by its ID. In practice, this would be a bad function without
// some other checks (like what user is logged in) because someone could hypothetically delete ANY credential.
func DeleteCredentialByID(credentialID string) error {
	//return db.Where("cred_id = ?", credentialID).Delete(&Credential{}).Error
	credsDB := db.From("credentials")
	return credsDB.DeleteStruct(&Credential{CredID: credentialID})
}

// GetUnformattedPublicKeyForCredential gives you the raw PublicKey model for a credential
func GetUnformattedPublicKeyForCredential(c *Credential) (PublicKey, error) {
	return c.PublicKey, err
}

// GetPublicKeyForCredential gets the formatted `models.PublicKey` for a provided credential
func GetPublicKeyForCredential(c *Credential) (ecdsa.PublicKey, error) {
	log.Println("PubKey:", c.PublicKey)
	return FormatPublicKey(c.PublicKey)
}

// FormatPublicKey formats a `models.PublicKey` into an `ecdsa.PublicKey`
func FormatPublicKey(pk PublicKey) (ecdsa.PublicKey, error) {
	ecPoint, err := AssembleUncompressedECPoint(pk.XCoord, pk.YCoord)
	if err != nil {
		return ecdsa.PublicKey{}, err
	}
	xInt, yInt := elliptic.Unmarshal(elliptic.P256(), ecPoint)
	return ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     xInt,
		Y:     yInt,
	}, err
}

// AssembleUncompressedECPoint will properly format the EC coordinates into
func AssembleUncompressedECPoint(xCoord []byte, yCoord []byte) ([]byte, error) {
	point := make([]byte, 65)
	if len(xCoord) != 32 || len(yCoord) != 32 {
		fmt.Println("X coord byte length : ", len(xCoord))
		fmt.Println("Y coord byte length : ", len(yCoord))
		err := errors.New("Coordinates are not 32 bytes long")
		return point, err
	}
	point[0] = 0x04
	copy(point[1:33], xCoord)
	copy(point[33:], yCoord)
	return point, nil
}
