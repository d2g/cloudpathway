package datastore

import (
	"encoding/json"
	"net"
)

//MACAddress Is the Primary Key
//IPAddress Is Needed as a Lookup [IPAddress -> MACAddress] (Things from the Kernel have the IP)
type Device struct {
	MACAddress  net.HardwareAddr
	IPAddress   net.IP
	Hostname    string
	Nickname    string
	CurrentUser *User
	DefaultUser *User
}

func (t Device) MarshalJSON() ([]byte, error) {
	s := struct {
		MACAddress  string
		IPAddress   string
		Hostname    string
		Nickname    string
		CurrentUser string
		DefaultUser string
	}{
		(t.MACAddress.String()),
		(t.IPAddress.String()),
		t.Hostname,
		t.Nickname,
		"",
		"",
	}

	if t.CurrentUser != nil {
		s.CurrentUser = t.CurrentUser.Username
	}

	if t.DefaultUser != nil {
		s.DefaultUser = t.DefaultUser.Username
	}

	return json.Marshal(s)
}

func (t *Device) UnmarshalJSON(data []byte) error {
	s := struct {
		MACAddress  string
		IPAddress   string
		Hostname    string
		Nickname    string
		CurrentUser string
		DefaultUser string
	}{}

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	tmpMACAddress, err := net.ParseMAC(s.MACAddress)
	if err != nil {
		return err
	}

	t.MACAddress = tmpMACAddress
	t.IPAddress = net.ParseIP(s.IPAddress)
	t.Hostname = s.Hostname
	t.Nickname = s.Nickname

	if s.CurrentUser != "" || s.DefaultUser != "" {
		userHelper, err := GetUserHelper()
		if err != nil {
			return err
		}

		if s.CurrentUser != "" {
			//Lookup Current User
			tmpCurrentUser, err := userHelper.GetUser(s.CurrentUser)
			if err != nil {
				return err
			}
			t.CurrentUser = &tmpCurrentUser
		}

		if s.DefaultUser != "" {
			//Lookup Current User
			tmpDefaultUser, err := userHelper.GetUser(s.DefaultUser)
			if err != nil {
				return err
			}
			t.DefaultUser = &tmpDefaultUser
		}
	}

	return nil
}

/**
 * Returns the display name for this device.
 * If there is a Nickname present, returns that, otherwise return the hostname.
 */
func (this *Device) GetDisplayName() string {
	if this.Nickname != "" {
		return this.Nickname
	}

	return this.Hostname
}

func (t Device) DisplayName() string {
	return t.GetDisplayName()
}

/*
 * Get the Current Username (CurrentUser, Default User if Blank and blank if both blank)
 */
func (this *Device) GetActiveUser() *User {
	if this.CurrentUser != nil {
		return this.CurrentUser
	}
	if this.DefaultUser != nil {
		return this.DefaultUser
	}
	return nil
}
