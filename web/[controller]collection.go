package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/controller"
	template "github.com/d2g/goti/html"
	"github.com/gorilla/mux"
)

type Collection struct {
	Sessions *ActiveSessions

	base string
}

func (t *Collection) Routes() (http.Handler, error) {
	router := mux.NewRouter()

	/*
	 * Setup my Routes.
	 */
	router.HandleFunc(t.Base(), t.index)
	router.HandleFunc(t.Base()+"create", t.create)

	router.HandleFunc(t.Base()+"{collection}/edit", t.edit)
	//router.HandleFunc(t.Base()+"{collection}/url/add/{url}", t.urladd)
	//router.HandleFunc(t.Base()+"{collection}/url/remove/{url}", t.urlremove)

	router.HandleFunc(t.Base()+"save", t.save).Methods("POST")

	router.HandleFunc(t.Base()+"{collection}/delete", t.remove)

	return router, nil
}

func (t *Collection) SetBase(base string) controller.HTTPController {
	t.base = base
	return t
}

func (t *Collection) Base() string {
	return t.base
}

/**
 * Displays the list of collections.
 */
func (t *Collection) index(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	myuser := t.Sessions.CurrentUser(request)

	// Get all collections.
	domainCollectionDataStoreHelper, err := datastore.GetDomainCollectionHelper()
	allCollections, err := domainCollectionDataStoreHelper.GetDomainCollections()

	// Check for error when loading collections.
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
		Action         string
		User           datastore.User
		AllCollections []datastore.DomainCollection
		SaveComplete   bool
		SaveError      bool
	}{
		"collectionSettings",
		myuser,
		allCollections,
		saveComplete,
		saveError,
	}

	// Parse the page and execute the template.
	tpl, _ := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/collections/index.tpl")
	tpl.Execute(response, data)
}

/**
 * Handles adding and removing urls for the collection.
 */
func (t *Collection) edit(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	myuser := t.Sessions.CurrentUser(request)

	// Get the collection.
	collectionName := mux.Vars(request)["collection"]
	domainCollectionDataStoreHelper, err := datastore.GetDomainCollectionHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	collection, err := domainCollectionDataStoreHelper.GetDomainCollection(collectionName)

	// Check for error when loading collection.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Setup the data structure to pass to the page.
	data := struct {
		Action     string
		User       datastore.User
		Collection datastore.DomainCollection
	}{
		"collectionSettings",
		myuser,
		collection,
	}

	// Parse the page and execute the template.
	tpl, _ := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/collections/edit.tpl")
	tpl.Execute(response, data)
}

/**
 * Handles creating a new collection.
 */
func (t *Collection) create(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	myuser := t.Sessions.CurrentUser(request)

	// Setup the data structure to pass to the page.
	data := struct {
		Action string
		User   datastore.User
	}{
		"collectionSettings",
		myuser,
	}

	// Parse the page and execute the template.
	tpl, _ := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/collections/create.tpl")
	tpl.Execute(response, data)
}

/**
 * Handles saving a collection.
 */
func (t *Collection) save(response http.ResponseWriter, request *http.Request) {
	domainCollectionDataStoreHelper, err := datastore.GetDomainCollectionHelper()

	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	//collection := datastore.DomainCollection(request.FormValue("name"))
	collection := datastore.DomainCollection{}

	// Save the collection.
	err = domainCollectionDataStoreHelper.SetDomainCollection(collection)

	// Check for error when saving collection.
	if err != nil {
		// There was an error, so report that to the screen.
		http.Redirect(response, request, "/settings_collections.html?saveerror=true", http.StatusMovedPermanently)
	} else {
		// No error, so report a successful save to the screen
		http.Redirect(response, request, t.Base()+"?savecomplete=true", http.StatusMovedPermanently)
	}
}

/**
 * Handles deleting a collection.
 */
func (t *Collection) remove(response http.ResponseWriter, request *http.Request) {
	// Try and load a filter collection using the collection name.
	collectionName := mux.Vars(request)["collection"]
	domainCollectionDataStoreHelper, err := datastore.GetDomainCollectionHelper()

	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	err = domainCollectionDataStoreHelper.DeleteDomainCollection(collectionName)

	// Check for error when deleting collection.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	http.Redirect(response, request, t.Base(), http.StatusMovedPermanently)
}

/**
 * Gets the list of blocked sites that match the given filter and returns as JSON response.
 */
func (t *Collection) getFilterSiteList(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("collection")
	filter := request.FormValue("filter")
	log.Println(name)
	log.Println(filter)
	//domainCollectionDataStoreHelper, err := datastore.GetDomainCollectionHelper()
	//collection, err := domainCollectionDataStoreHelper.GetDomainCollection(name)

	// Check for error when loading collection.
	//if err != nil {
	//	http.Error(response, err.Error(), 500)
	//	return
	//}

	// Setup the data structure to return.
	data := struct {
		Sites []string
	}{
		[]string{}, //TODO: Fix
	}

	enc := json.NewEncoder(response)
	enc.Encode(data)
}
