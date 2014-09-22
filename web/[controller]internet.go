package web

import (
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/cloudpathway/dnsimplementation"
	"github.com/d2g/controller"
	template "github.com/d2g/goti/html"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"strconv"
)

type Internet struct {
	Sessions *ActiveSessions

	base string
}

func (t *Internet) Routes() (http.Handler, error) {
	router := mux.NewRouter()

	/*
	 * Setup my Routes.
	 */
	router.HandleFunc(t.Base(), t.index)
	router.HandleFunc(t.Base()+"{source}/{sourcePort}/{destination}/{destinationPort}/", t.connection)

	return router, nil
}

func (t *Internet) SetBase(base string) controller.HTTPController {
	t.base = base
	return t
}

func (t *Internet) Base() string {
	return t.base
}

func (t *Internet) index(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	myuser := t.Sessions.CurrentUser(response, request)

	connectionhelper, err := datastore.GetConnectionHelper()
	if err != nil {
		log.Printf("Error: Failed to get Connections Helper \"%v\"", err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	// Get the devices.
	deviceHelper, err := datastore.GetDeviceHelper()
	if err != nil {
		log.Println("Error: DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	connections, err := connectionhelper.GetConnections()
	if err != nil {
		log.Printf("Error: Failed to get Connections \"%v\"", err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	//TODO: Remove.
	connections = datastore.GetExampleConnections()

	devices, err := deviceHelper.GetDevices()
	// Check for error when loading devices.
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	mappedDevices := make(map[string]datastore.Device)

	for _, device := range devices {
		mappedDevices[device.MACAddress.String()] = device
	}

	hostnames := make(map[string]string)
	cache, err := dnsimplementation.GetCache()
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	for _, connection := range connections {
		//Lookup Source IP
		hostname, err := cache.GetHostname(connection.SourceIP())

		if err != nil {
			log.Printf("Error: Looking Up Hostname \"%v\"", err.Error())
		} else if hostname == "" {
			hostnames[connection.SourceIP().To4().String()] = connection.SourceIP().To4().String()
		} else {
			hostnames[connection.SourceIP().To4().String()] = hostname
		}

		//Lookup Destination IP
		hostname, err = cache.GetHostname(connection.DestinationIP())

		if err != nil {
			log.Printf("Error: Looking Up Hostname \"%v\"", err.Error())
		} else if hostname == "" {
			hostnames[connection.DestinationIP().To4().String()] = connection.DestinationIP().To4().String()
		} else {
			hostnames[connection.DestinationIP().To4().String()] = hostname
		}
	}

	connections = datastore.Connections(connections)

	// Setup the data structure to pass to the page.
	data := struct {
		Action      string
		User        datastore.User
		Connections datastore.Connections
		Devices     map[string]datastore.Device
		Hostnames   map[string]string
	}{
		"internet",
		myuser,
		connections,
		mappedDevices,
		hostnames,
	}

	// Parse the page and execute the template.
	tpl, err := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/internet/index.tpl")
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	err = tpl.Execute(response, data)
	if err != nil {
		log.Println(err)
		return
	}

}

func (t *Internet) connection(response http.ResponseWriter, request *http.Request) {
	// Get the current session user.
	user := t.Sessions.CurrentUser(response, request)

	connectionhelper, err := datastore.GetConnectionHelper()
	if err != nil {
		log.Printf("Error: Failed to get Connections Helper \"%v\"", err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	sourcePort, err := strconv.ParseUint(mux.Vars(request)["sourcePort"], 10, 16)
	if err != nil {
		log.Printf("Error: Failed to convert sourceport to uint16 \"%v\"", err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	destinationPort, err := strconv.ParseUint(mux.Vars(request)["destinationPort"], 10, 16)
	if err != nil {
		log.Printf("Error: Failed to convert sourceport to uint16 \"%v\"", err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	connection, err := connectionhelper.GetConnection(net.ParseIP(mux.Vars(request)["source"]), uint16(sourcePort), net.ParseIP(mux.Vars(request)["destination"]), uint16(destinationPort))
	if err != nil {
		log.Printf("Error: Getting Connection \"%v\"", err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	// Parse the page and execute the template.
	tpl, err := template.ParseFiles("static/main.tpl", "static/main_authenticated.tpl", "static/internet/connection.tpl")
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	data := struct {
		Action     string
		User       datastore.User
		Connection *datastore.Connection
	}{
		"internet",
		user,
		&connection,
	}

	err = tpl.Execute(response, data)
	if err != nil {
		log.Println(err)
		return
	}
}
