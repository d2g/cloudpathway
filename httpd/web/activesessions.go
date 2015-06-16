package web

import (
	"bytes"
	"encoding/gob"
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/sessions"
	"log"
	"net/http"
	"time"
)

type ActiveSessions struct {
	sessions.SessionInfo
}

func (t *ActiveSessions) CurrentUser(request *http.Request) datastore.User {

	session, err := t.GetSession(request)
	if err != nil {
		log.Printf("Error: Getting Current Session:%s\n", err.Error())
		return datastore.User{}
	}

	user, err := t.currentUser(session)
	if err != nil {
		log.Printf("Error: Getting Session User:%s\n", err.Error())
		return datastore.User{}
	}

	return user
}

func (t *ActiveSessions) currentUser(s sessions.Session) (datastore.User, error) {
	user, err := s.Get("User")
	if err != nil {
		return datastore.User{}, err
	}

	if user != nil {
		assertedUser, ok := user.(*datastore.User)
		if ok {
			return *assertedUser, nil
		} else {
			log.Printf("Error: Unable to Asset Session User")
			return datastore.User{}, nil
		}
	} else {
		return datastore.User{}, nil
	}
}

func (t *ActiveSessions) SetCurrentUser(response http.ResponseWriter, request *http.Request, user datastore.User) error {
	session, err := t.GetSession(request)
	if err != nil {
		http.Error(response, err.Error(), 500)
		return err
	}

	err = session.Set("User", user)
	if err != nil {
		http.Error(response, err.Error(), 500)
		return err
	}

	t.SetSession(request, session)
	return nil
}

func (t *ActiveSessions) CurrentDevice(request *http.Request) datastore.Device {
	s, err := t.GetSession(request)
	if err != nil {
		log.Printf("Error: Getting Current Session:%s\n", err.Error())
		return datastore.Device{}
	}

	d, err := t.currentDevice(s)
	if err != nil {
		log.Printf("Error: Getting Device for Current Session:%s\n", err.Error())
		return datastore.Device{}
	}

	return d
}

func (t *ActiveSessions) currentDevice(s sessions.Session) (datastore.Device, error) {
	device, err := s.Get("Device")
	if err != nil {
		return datastore.Device{}, err
	}

	if device != nil {
		assertedDevice, ok := device.(*datastore.Device)
		if ok {
			return *assertedDevice, nil
		} else {
			log.Printf("Error: Unable to Asset Session Device")
			return datastore.Device{}, nil
		}
	}

	return datastore.Device{}, nil
}

func (t *ActiveSessions) SetCurrentDevice(response http.ResponseWriter, request *http.Request, device datastore.Device) error {
	session, err := t.GetSession(request)
	if err != nil {
		http.Error(response, err.Error(), 500)
		return err
	}

	err = session.Set("Device", device)
	if err != nil {
		http.Error(response, err.Error(), 500)
		return err
	}

	t.SetSession(request, session)
	return nil
}

func (t *ActiveSessions) GCSessions() {

	//Device Datastore
	helper, err := datastore.GetDeviceHelper()
	if err != nil {
		log.Printf("Error: GC Sessions Unable to get Device Helper \"%s\"\n", err.Error())
		return
	}

	for {
		select {
		case <-time.After(t.SessionInfo.Timeout):
		}

		log.Printf("Debug: Running Session GC\n")
		sessions, err := t.SessionInfo.Store.All()
		if err != nil {
			log.Printf("Error: getting all sessions \"%v\"\n", err.Error())
			continue
		}

		for i := range sessions {
			if sessions[i].Expiry().Before(time.Now()) {
				id, err := sessions[i].ID()
				if err != nil {
					log.Printf("Error: GC Failed to get session ID to delete \"%v\"\n", err.Error())
					continue
				}

				//If there was a device with this session?
				device, err := t.currentDevice(sessions[i])
				if err != nil {
					log.Printf("Warning: GC Failed to find device for session \"%s\" because of \"%v\"\n", id, err.Error())
				}

				if !bytes.Equal(device.MACAddress, []byte{}) {
					//Ok we were logged onto a device.
					//Lets get the device as it currentlky is (Not how it was when the session started.
					currentdevice, err := helper.GetDevice(device.MACAddress)
					if err != nil {
						log.Printf("Error: GC Failed to get device\"%s\" because of \"%v\"\n", device.MACAddress, err.Error())
					}

					if !bytes.Equal(device.MACAddress, []byte{}) && bytes.Equal(device.MACAddress, currentdevice.MACAddress) {
						//We have a current device that was in the session.
						user, err := t.currentUser(sessions[i])
						if err != nil {
							log.Printf("Error: GC Failed to get User from session \"%s\" because of \"%v\"\n", id, err.Error())
						}

						if user.Username != "" && currentdevice.CurrentUser != nil && currentdevice.CurrentUser.Username == user.Username {
							currentdevice.CurrentUser = nil
							err := helper.SetDevice(&currentdevice)
							if err != nil {
								log.Printf("Error: GC Failed to get update device \"%s\" because of \"%v\"\n", currentdevice.MACAddress, err.Error())
							}
						}
					}
				}

				err = t.SessionInfo.Store.Delete(id)
				if err != nil {
					log.Printf("Error: GC Failed to delete session %s because \"%v\"\n", id, err.Error())
					continue
				}

				log.Printf("Debug: Deleted Session %s\n", id)
			}
		}

	}
}

func init() {
	//We need to register types we use in sessions as gob has to encode them :(
	gob.Register(&datastore.User{})
	gob.Register(&datastore.Device{})
}
