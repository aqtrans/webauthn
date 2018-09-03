package models

// RelyingParty is the group the User is authenticating with
type RelyingParty struct {
	ID          string `json:"id" storm:"id"`
	DisplayName string `json:"display_name"`
	Icon        string `json:"icon,omitempty"`
	Users       []User `json:"users,omitempty" storm:"unique"`
}

// GetDefaultRelyingParty gets the RP associated with the configured hostname
func GetDefaultRelyingParty() (RelyingParty, error) {
	rp := RelyingParty{}
	//err := db.Where("id=?", config.Conf.HostAddress).First(&rp).Error
	err := db.One("ID", "localhost", &rp)
	if err != nil {
		return rp, err
	}
	return rp, nil
}

// GetRelyingPartyByHost gets the RP by hostname which in this case is the ID
func GetRelyingPartyByHost(hostname string) (RelyingParty, error) {
	rp := RelyingParty{}
	//err := db.Where("id = ?", hostname).First(&rp).Error
	err := db.One("ID", hostname, &rp)
	if err != nil {
		return rp, err
	}
	return rp, nil
}

// PutRelyingParty creates or updates a Relying Party
func PutRelyingParty(rp *RelyingParty) error {
	err := db.Save(rp)
	return err
}
