package datastore

import (
	"encoding/json"
	//"log"
	"time"
)

type User struct {
	Username    string
	Password    string
	DisplayName string
	IsAdmin     bool      //So we're starting with Admin & User only (I'm sure this will change one day)
	DOB         time.Time //DOB YYYY-MM-DD
}

/**
 * Identifies if this user is the default user for the given device.
 */
func (this *User) IsDefaultUserForDevice(device Device) bool {
	return this.Username == device.DefaultUser.Username
}

/**
 * Returns the display name for this user.
 * If there is a DisplayName present, returns that, otherwise return the Username.
 */
func (this *User) GetDisplayName() string {
	if this.DisplayName != "" {
		return this.DisplayName
	}

	return this.Username
}

/**
 * Returns the DOB field in YYYY-MM-DD format.
 */
func (this *User) ShortDOB() string {
	if this.DOB.IsZero() {
		return ""
	} else {
		return this.DOB.Format("2006-01-02")
	}
}

/**
 * Returns the number of filter collections assigned to this user.
 */
func (this *User) NumberOfFilterCollections() (int, error) {
	userFilterCollectionsHelper, err := GetUserFilterCollectionsHelper()
	if err != nil {
		return 0, err
	}

	userFilterCollections, err := userFilterCollectionsHelper.GetUserFilterCollections(this.Username)
	if err != nil {
		return 0, err
	}

	return userFilterCollections.NumberOfCollections(), nil
}

/**
 * Sets the DOB field using the given YYYY-MM-DD format.
 */
func (this *User) SetShortDOB(shortFormat string) error {
	time, err := time.Parse("2006-01-02", shortFormat)
	if err != nil {
		return err
	}

	this.DOB = time
	return nil
}

/**
 * Gets and returns the UserFilterCollections related to this user.
 */
func (this *User) UserFilterCollections() (UserFilterCollections, error) {
	userFilterCollectionsHelper, err := GetUserFilterCollectionsHelper()
	if err != nil {
		return UserFilterCollections{}, err
	}

	userFilterCollection, err := userFilterCollectionsHelper.GetUserFilterCollections(this.Username)
	if err != nil {
		return UserFilterCollections{}, err
	}

	return userFilterCollection, nil
}

func (t User) MarshalJSON() ([]byte, error) {
	s := struct {
		Username    string
		Password    string
		DisplayName string
		IsAdmin     bool
		DOB         string
	}{
		t.Username,
		t.Password,
		t.DisplayName,
		t.IsAdmin,
		t.DOB.Format("2006-01-02"),
	}

	return json.Marshal(s)
}

func (t *User) UnmarshalJSON(data []byte) error {
	s := struct {
		Username    string
		Password    string
		DisplayName string
		IsAdmin     bool
		DOB         string
	}{}

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	t.Username = s.Username
	t.Password = s.Password
	t.DisplayName = s.DisplayName
	t.IsAdmin = s.IsAdmin
	t.DOB, _ = time.Parse("2006-01-02", s.DOB)

	return nil
}
