package web

import (
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/controller"
	template "github.com/d2g/goti/html"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type User struct {
	Sessions *ActiveSessions

	base string
}

func (t *User) Routes() (http.Handler, error) {
	router := mux.NewRouter()

	/*
	 * Setup my Routes.
	 */
	router.HandleFunc(t.Base()+"edit", t.newuser)
	router.HandleFunc(t.Base()+"{username}/edit", t.edit)
	router.HandleFunc(t.Base()+"save", t.save).Methods("POST")
	router.HandleFunc(t.Base()+"{username}/delete", t.remove) //Deletes a Keyword...
	router.HandleFunc(t.Base(), t.index)

	return router, nil
}

func (t *User) SetBase(base string) controller.HTTPController {
	t.base = base
	return t
}

func (t *User) Base() string {
	return t.base
}

func (t *User) index(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	user := t.Sessions.CurrentUser(response, request)

	// Get all users.
	userHelper, err := datastore.GetUserHelper()

	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	allUsers, err := userHelper.GetUsers()

	// Check for error when loading users.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Retrieve the savecomplete parameter from the URL and determine if it is true or not.
	saveCompleteParam := request.URL.Query().Get("savecomplete")
	saveComplete := saveCompleteParam == "true"

	// Retrieve the saveerror parameter from the URL and determine if it is true or not.
	saveErrorParam := request.URL.Query().Get("saveerror")
	saveError := saveErrorParam == "true"

	// Setup the data structure to pass to the page.
	data := struct {
		Action       string
		User         datastore.User
		AllUsers     []datastore.User
		SaveComplete bool
		SaveError    bool
	}{
		"userSettings",
		user,
		allUsers,
		saveComplete,
		saveError,
	}

	// Parse the page and execute the template.
	tpl, _ := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/users/index.tpl")
	tpl.Execute(response, data)
}

/**
 * Handles editing a user.
 */
func (t *User) edit(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	user := t.Sessions.CurrentUser(response, request)

	// Get the user.
	username := mux.Vars(request)["username"]
	userDataStoreHelper, err := datastore.GetUserHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	userToEdit, err := userDataStoreHelper.GetUser(username)

	// Check for error when loading user.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Setup the data structure to pass to the page.
	data := struct {
		Action     string
		User       datastore.User
		UserToEdit *datastore.User
		EditUser   bool
		NewUser    bool
	}{
		"userSettings",
		user,
		&userToEdit,
		true,
		false,
	}

	// Parse the page and execute the template.
	tpl, _ := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/users/edit.tpl")
	tpl.Execute(response, data)
}

/**
 * Handles creating a new user.
 */
func (t *User) newuser(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	user := t.Sessions.CurrentUser(response, request)

	// Setup the data structure to pass to the page.
	data := struct {
		Action   string
		User     datastore.User
		NewUser  bool
		EditUser bool
	}{
		"userSettings",
		user,
		true,
		false,
	}

	// Parse the page and execute the template.
	tpl, _ := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/users/edit.tpl")
	tpl.Execute(response, data)
}

/**
 * Handles saving a user.
 */
func (t *User) save(response http.ResponseWriter, request *http.Request) {
	// Try and load a user with the Username.
	existingUsername := request.FormValue("idUsername")
	userDataStoreHelper, err := datastore.GetUserHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	myuser := datastore.User{}

	if existingUsername != "" {
		myuser, err = userDataStoreHelper.GetUser(existingUsername)

		// Check for error when loading user.
		if err != nil {
			http.Error(response, err.Error(), 500)
			return
		}
	}

	// Set the values on the object.
	myuser.Username = request.FormValue("username")
	myuser.DisplayName = request.FormValue("name")
	myuser.SetShortDOB(request.FormValue("dob"))

	if request.FormValue("isAdmin") == "true" {
		myuser.IsAdmin = true
	} else {
		myuser.IsAdmin = false
	}

	err = userDataStoreHelper.SetUser(&myuser)

	// Check for error when saving user.
	if err != nil {
		// There was an error, so report that to the screen.
		http.Redirect(response, request, t.Base()+"?saveerror=true", http.StatusMovedPermanently)
	} else {
		// No error, so report a successful save to the screen and remove the previous version
		// if username has been edited.

		// If the new username is different to the existing username, delete the old
		// version before saving the new version. This prevents duplicates.
		if existingUsername != "" && myuser.Username != existingUsername {
			err = userDataStoreHelper.DeleteUser(existingUsername)

		}
		http.Redirect(response, request, t.Base()+"?savecomplete=true", http.StatusMovedPermanently)
	}
}

/**
 * Handles deleting a user.
 */
func (t *User) remove(response http.ResponseWriter, request *http.Request) {
	// Try and load a user with the Username.
	username := mux.Vars(request)["username"]
	userDataStoreHelper, err := datastore.GetUserHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	err = userDataStoreHelper.DeleteUser(username)
	// Check for error when loading user.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	http.Redirect(response, request, t.Base(), http.StatusMovedPermanently)
}
