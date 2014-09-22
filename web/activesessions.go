package web

import (
	"encoding/gob"
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/sessions"
	"log"
	"net/http"
)

type ActiveSessions struct {
	sessions.SessionInfo
}

func (t *ActiveSessions) CurrentUser(response http.ResponseWriter, request *http.Request) datastore.User {

	session, err := t.GetSession(request)
	if err != nil {
		log.Printf("Error: Getting Current Session:%s\n", err.Error())
		return datastore.User{}
	}

	user, err := session.Get("User")
	if err != nil {
		log.Printf("Error Getting Current Session User:%s\n", err.Error())
		return datastore.User{}
	}

	if user != nil {
		assertedUser, ok := user.(*datastore.User)
		if ok {
			return *assertedUser
		} else {
			log.Printf("Error: Unable to Asset Session User")
			return datastore.User{}
		}
	} else {
		return datastore.User{}
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

func init() {
	//We need to register types we use in sessions as gob has to encode them :(
	gob.Register(&datastore.User{})
}
