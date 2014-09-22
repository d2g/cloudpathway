package web

import (
	"crypto/sha1"
	"fmt"
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/controller"
	template "github.com/d2g/goti/html"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

type Root struct {
	Sessions        *ActiveSessions
	NotFoundHandler http.Handler

	base string
}

func (t *Root) Routes() (http.Handler, error) {
	router := mux.NewRouter()

	/*
	 * Setup my Routes.
	 */
	router.HandleFunc(t.Base()+"login", t.login)
	router.HandleFunc(t.Base()+"logoff", t.logoff)
	router.HandleFunc(t.Base(), t.index)

	//The root "/" catches all so if we end up here we need to check it's not a 404.
	router.NotFoundHandler = t.NotFoundHandler

	return router, nil
}

func (t *Root) SetBase(base string) controller.HTTPController {
	t.base = base
	return t
}

func (t *Root) Base() string {
	return t.base
}

func (t *Root) index(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	user := t.Sessions.CurrentUser(response, request)

	// If there isn't a current session user, redirect to login.
	if user.Username == "" {
		http.Redirect(response, request, "/login", http.StatusMovedPermanently)
		return
	}

	// Get the devices for this user.
	deviceDataStoreHelper, err := datastore.GetDeviceHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	currentUserDevices, defaultUserDevices, err := deviceDataStoreHelper.GetDevicesForUser(user.Username)

	// Check for an error when loading devices.
	if err != nil {
		log.Println("Error Getting Devices:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	// Setup the data structure to pass to the page.
	data := struct {
		Action         string
		User           datastore.User
		CurrentDevices []datastore.Device
		DefaultDevices []datastore.Device
	}{
		"index",
		user,
		currentUserDevices,
		defaultUserDevices,
	}

	// Parse the page and execute the template.
	tpl, _ := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/index.tpl")
	tpl.Execute(response, data)
}

func (t *Root) login(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		// Retrieve the loginFailed parameter from the URL and determine if it is true or not.
		loginFailedParam := request.URL.Query().Get("loginFailed")
		loginFailed := loginFailedParam == "true"

		// Add the loginFailed value to a data struct to pass to the template.
		data := struct {
			LoginFailed bool
		}{
			loginFailed,
		}

		tpl, _ := template.ParseFiles("static/main.tpl", "static/login.tpl")
		tpl.Execute(response, data)
	} else {

		request.ParseForm()

		hash := sha1.New()
		hash.Write([]byte(request.FormValue("password")))

		passwordHash := fmt.Sprintf("%x", hash.Sum(nil))

		helper, err := datastore.GetUserHelper()
		if err != nil {
			log.Println("DB Error:" + err.Error())
			http.Error(response, err.Error(), 500)
			return
		}

		myuser, err := helper.GetUser(request.FormValue("username"))
		if err != nil {
			log.Println("Error Looking Up User:" + err.Error())
			http.Error(response, err.Error(), 500)
			return
		}

		if myuser.Username != "" && myuser.Username == request.FormValue("username") && myuser.Password == passwordHash {
			//Setup our session here.
			t.Sessions.SetCurrentUser(response, request, myuser)

			remoteIP, _, err := net.SplitHostPort(request.RemoteAddr)

			//We now need to assosiate the user with the current device :D
			deviceHelper, err := datastore.GetDeviceHelper()
			if err != nil {
				log.Println("DB Error:" + err.Error())
				http.Error(response, err.Error(), 500)
				return
			}

			device, err := deviceHelper.GetDeviceByIP(net.ParseIP(remoteIP))

			if err != nil {
				log.Println("Error Looking Up Device:" + err.Error())
				http.Error(response, err.Error(), 500)
				return
			}

			//Make sure the device is within our network.
			if device.MACAddress.String() != "" {
				//Device is within our network (It might have been external)
				device.CurrentUser = &myuser
				deviceHelper.SetDevice(&device)
			}

			// Login successful, display the user dashboard.
			t.Sessions.SessionInfo.SaveSession(response, request)
			http.Redirect(response, request, "/", http.StatusMovedPermanently)
			return
		} else {
			// Login unsuccessful, return to login screen.
			log.Println("Login Failed:" + myuser.Username)
			http.Redirect(response, request, "/login?loginFailed=true", http.StatusMovedPermanently)
			return
		}

	}
}

/**
 * Logoff.
 */
func (t *Root) logoff(response http.ResponseWriter, request *http.Request) {
	//TODO:
	http.Redirect(response, request, "/", http.StatusMovedPermanently)
}
