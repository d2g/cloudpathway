package web

import (
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/controller"
	template "github.com/d2g/goti/html"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type UsersCollections struct {
	Sessions *ActiveSessions

	base string
}

func (t *UsersCollections) Routes() (http.Handler, error) {
	router := mux.NewRouter()

	/*
	 * Setup my Routes.
	 */
	router.HandleFunc(t.Base(), t.index)
	router.HandleFunc(t.Base()+"{username}/edit", t.edit)
	router.HandleFunc(t.Base()+"save", t.save).Methods("POST")

	return router, nil
}

func (t *UsersCollections) SetBase(base string) controller.HTTPController {
	t.base = base
	return t
}

func (t *UsersCollections) Base() string {
	return t.base
}

/**
 * Displays the user access list.
 */
func (t *UsersCollections) index(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	myuser := t.Sessions.CurrentUser(response, request)

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
		"userAccessSettings",
		myuser,
		allUsers,
		saveComplete,
		saveError,
	}

	// Parse the page and execute the template.
	tpl, _ := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/access/index.tpl")
	tpl.Execute(response, data)
}

/**
 * Handles editing a user filter collections object.
 */
func (t *UsersCollections) edit(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	myuser := t.Sessions.CurrentUser(response, request)

	// Get the user filter collections.
	username := mux.Vars(request)["username"]
	userFilterCollectionsHelper, err := datastore.GetUserFilterCollectionsHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	userFilterCollections, err := userFilterCollectionsHelper.GetUserFilterCollections(username)

	// Check for error when loading user filter collection.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Get the user relating to the filter collection.
	userDataStoreHelper, err := datastore.GetUserHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	filterCollectionsUser, err := userDataStoreHelper.GetUser(username)

	// Check for error when loading user.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Get all collections to display on the page.
	filterCollectionHelper, err := datastore.GetFilterCollectionHelper()
	allCollections, err := filterCollectionHelper.GetFilterCollections()

	// Check for error when loading collections.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Setup the data structure to pass to the page.
	data := struct {
		Action                string
		User                  datastore.User
		UserFilterCollections datastore.UserFilterCollections
		FilterCollectionsUser *datastore.User
		AllCollections        []datastore.FilterCollection
	}{
		"userAccessSettings",
		myuser,
		userFilterCollections,
		&filterCollectionsUser,
		allCollections,
	}

	// Parse the page and execute the template.
	tpl, _ := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/access/edit.tpl")
	tpl.Execute(response, data)
}

/**
 * Handles saving user access settings.
 */
func (t *UsersCollections) save(response http.ResponseWriter, request *http.Request) {
	// Try and load a user filter collection using the user name.
	username := request.FormValue("idName")
	userFilterCollectionsHelper, err := datastore.GetUserFilterCollectionsHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	userFilterCollections, err := userFilterCollectionsHelper.GetUserFilterCollections(username)

	// Check for error when loading collection.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Set the name of the collection.
	userFilterCollections.Username = username

	// Retrieve the collections from the HTML form.
	if err := request.ParseForm(); err != nil {
		http.Error(response, err.Error(), 500)
		return
	}
	collections := request.Form["collections[]"]
	userFilterCollections.Collections = collections

	// Save the user filter collection.
	err = userFilterCollectionsHelper.SetUserFilterCollections(userFilterCollections)

	// Check for error when saving collection.
	if err != nil {
		// There was an error, so report that to the screen.
		http.Redirect(response, request, t.Base()+"?saveerror=true", http.StatusMovedPermanently)
	} else {
		// No error, so report a successful save to the screen.
		http.Redirect(response, request, t.Base()+"?savecomplete=true", http.StatusMovedPermanently)
	}
}
