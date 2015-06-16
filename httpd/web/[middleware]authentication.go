package web

import (
	"github.com/d2g/controller"
	//"log"
	"net/http"
)

type Authentication struct {
	controller.HTTPController

	Sessions *ActiveSessions
}

/*
 * Override the Routes Function
 */
func (t *Authentication) Routes() (http.Handler, error) {

	routes, err := t.HTTPController.Routes()
	if err != nil {
		return nil, err
	} else {
		return &AuthenticationHandler{Handler: routes, Sessions: t.Sessions}, nil
	}
}

type AuthenticationHandler struct {
	http.Handler
	Sessions *ActiveSessions
}

func (s *AuthenticationHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	myuser := s.Sessions.CurrentUser(request)

	// If there isn't a current session user, redirect to login.
	if myuser.Username == "" {
		http.Redirect(response, request, "/login", http.StatusTemporaryRedirect)
		return
	}

	//The user must be logged in to get here.
	//Lets attach them to the device.

	s.Handler.ServeHTTP(response, request)
}
