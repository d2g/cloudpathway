package web

import (
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/controller"
	template "github.com/d2g/goti/html"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

type Device struct {
	Sessions *ActiveSessions

	base string
}

func (t *Device) Routes() (http.Handler, error) {
	router := mux.NewRouter()

	/*
	 * Setup my Routes.
	 */
	router.HandleFunc(t.Base(), t.index)
	router.HandleFunc(t.Base()+"edit", t.edit)
	router.HandleFunc(t.Base()+"{macAddress}/edit", t.edit)
	router.HandleFunc(t.Base()+"save", t.save).Methods("POST")
	router.HandleFunc(t.Base()+"{macAddress}/delete", t.remove)
	router.HandleFunc(t.Base()+"{macAddress}/removeuser", t.removeuser)

	return router, nil
}

func (t *Device) SetBase(base string) controller.HTTPController {
	t.base = base
	return t
}

func (t *Device) Base() string {
	return t.base
}

/**
 * Displays the device list.
 */
func (t *Device) index(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	myuser := t.Sessions.CurrentUser(request)

	// If there isn't a current session user, redirect to login.
	if myuser.Username == "" {
		http.Redirect(response, request, "/login", http.StatusMovedPermanently)
		return
	}

	// Get the devices.
	deviceHelper, err := datastore.GetDeviceHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	allDevices, err := deviceHelper.GetDevices()

	// Check for error when loading devices.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Retrieve the savecomplete parameter from the URL and determine if it is true or not.
	saveCompleteParam := request.URL.Query().Get("savecomplete")
	saveComplete := saveCompleteParam == "true"

	// Setup the data structure to pass to the page.
	data := struct {
		Action       string
		User         datastore.User
		AllDevices   []datastore.Device
		SaveComplete bool
	}{
		"deviceSettings",
		myuser,
		allDevices,
		saveComplete,
	}

	// Parse the page and execute the template.
	tpl, err := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/devices/index.tpl")

	// Check for error when loading devices.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	tpl.Execute(response, data)
}

func (t *Device) edit(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	myuser := t.Sessions.CurrentUser(request)

	// Get the device.
	macAddress := mux.Vars(request)["macAddress"]
	deviceDataStoreHelper, err := datastore.GetDeviceHelper()
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	MACAddressHardware, err := net.ParseMAC(macAddress)
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	mydevice, err := deviceDataStoreHelper.GetDevice(MACAddressHardware)

	// Check for error when loading device.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Get all users.
	userDataStoreHelper, err := datastore.GetUserHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	allUsers, err := userDataStoreHelper.GetUsers()

	// Check for error when loading users.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// Setup the data structure to pass to the page.
	data := struct {
		Action   string
		Device   datastore.Device
		User     datastore.User
		AllUsers []datastore.User
	}{
		"deviceSettings",
		mydevice,
		myuser,
		allUsers,
	}

	// Parse the page and execute the template.
	tpl, err := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/devices/edit.tpl")
	tpl.Execute(response, data)
}

func (t *Device) save(response http.ResponseWriter, request *http.Request) {
	// Load the device with the MAC Address.
	macAddress := request.FormValue("macAddress")
	deviceHelper, err := datastore.GetDeviceHelper()
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	MACAddressHardware, err := net.ParseMAC(macAddress)
	if err != nil {
		log.Println("Mac Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	device, err := deviceHelper.GetDevice(MACAddressHardware)

	// Check for error when loading device.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	device.Nickname = request.FormValue("nickname")

	if request.FormValue("defaultUser") != "" {
		userHelper, err := datastore.GetUserHelper()
		if err != nil {
			log.Println("Error: Getting User Helper :" + err.Error())
			http.Error(response, err.Error(), 500)
			return
		}

		user, err := userHelper.GetUser(request.FormValue("defaultUser"))
		if err != nil {
			log.Println("Error: Getting User:" + err.Error())
			http.Error(response, err.Error(), 500)
			return
		}

		if user.Username != "" {
			device.DefaultUser = &user
		} else {
			device.DefaultUser = nil
		}
	} else {
		device.DefaultUser = nil
	}

	err = deviceHelper.SetDevice(&device)
	// Check for error when saving device.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	http.Redirect(response, request, t.Base()+"?savecomplete=true", http.StatusMovedPermanently)
}

func (t *Device) remove(response http.ResponseWriter, request *http.Request) {
	// Load the device with the MAC Address
	macAddress := mux.Vars(request)["macAddress"]
	deviceDataStoreHelper, err := datastore.GetDeviceHelper()
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	MACAddressHardware, err := net.ParseMAC(macAddress)
	if err != nil {
		log.Println("Mac Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	err = deviceDataStoreHelper.DeleteDevice(MACAddressHardware)

	// Check for error when deleting device.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	http.Redirect(response, request, t.Base(), http.StatusMovedPermanently)
}

func (t *Device) removeuser(response http.ResponseWriter, request *http.Request) {
	// Load the device with the MAC Address
	macAddress := mux.Vars(request)["macAddress"]
	deviceDataStoreHelper, err := datastore.GetDeviceHelper()
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	MACAddressHardware, err := net.ParseMAC(macAddress)
	if err != nil {
		log.Println("Mac Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	myDevice, err := deviceDataStoreHelper.GetDevice(MACAddressHardware)
	if err != nil {
		log.Println("Error Getting Device:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	myDevice.CurrentUser = nil

	err = deviceDataStoreHelper.SetDevice(&myDevice)
	// Check for error when deleting device.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	http.Redirect(response, request, "/", http.StatusMovedPermanently)
}
