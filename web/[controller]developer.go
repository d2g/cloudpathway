package web

import (
	"bufio"
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/controller"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

type Developer struct {
	Sessions *ActiveSessions

	base string
}

func (t *Developer) Routes() (http.Handler, error) {
	router := mux.NewRouter()

	/*
	 * Setup my Routes.
	 */
	//Create Admin User
	router.HandleFunc(t.Base()+"create/user/{username}", t.createadmin)
	router.HandleFunc(t.Base()+"create/device/{macaddress}/{hostname}", t.createdevice)
	router.HandleFunc(t.Base()+"assign/device/{macaddress}/{ip}", t.assignip)
	router.HandleFunc(t.Base()+"import/collections", t.importcollections)
	router.HandleFunc(t.Base()+"dump/connections", t.dumpconnections)
	router.HandleFunc(t.Base()+"access/{username}/{domain}", t.access)
	router.HandleFunc(t.Base()+"quickaccess/{username}/{domain}", t.quickaccess)

	router.HandleFunc(t.Base()+"poke/block/{source}/{destination}", t.block)
	//router.HandleFunc(t.Base()+"poke/unblock/{source}/{destination}", t.unblock)

	return router, nil
}

func (t *Developer) SetBase(base string) controller.HTTPController {
	t.base = base
	return t
}

func (t *Developer) Base() string {
	return t.base
}

/*
 * Create a User in the Database with the username specified in the URL and the password "administrator".
 */
func (t *Developer) createadmin(response http.ResponseWriter, request *http.Request) {
	myuser := datastore.User{}
	myuser.Username = mux.Vars(request)["username"]
	myuser.Password = "b3aca92c793ee0e9b1a9b0a5f5fc044e05140df3"
	myuser.IsAdmin = true
	helper, err := datastore.GetUserHelper()
	if err != nil {
		log.Println("DB Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	err = helper.SetUser(&myuser)
	if err != nil {
		log.Println("Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	return
}

/*
 * Create Dummy Devices
 */
func (t *Developer) createdevice(response http.ResponseWriter, request *http.Request) {
	//We Need to Create or update the Device for this Lease...
	deviceHelper, err := datastore.GetDeviceHelper()
	if err != nil {
		log.Println("Error: Getting Helper Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	macAddress, err := net.ParseMAC(mux.Vars(request)["macaddress"])
	if err != nil {
		log.Println("Error: Invalid MACAddress :" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	device, err := deviceHelper.GetDevice(macAddress)
	if err != nil {
		log.Println("Error: Getting Device Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	device.Hostname = mux.Vars(request)["hostname"]

	err = deviceHelper.SetDevice(&device)
	if err != nil {
		log.Println("Error: Setting Device To Datastore:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}
}

/*
 * Import the collections directory and Import the files as collection lists.
 */
func (t *Developer) importcollections(response http.ResponseWriter, request *http.Request) {
	filepath.Walk("import/collections/", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file, err := os.Open(path)
			defer file.Close()

			scanner := bufio.NewScanner(file)

			collectionHelper, err := datastore.GetFilterCollectionHelper()
			if err != nil {
				log.Println("Error: Getting Helper Error:" + err.Error())
				http.Error(response, err.Error(), 500)
				return err
			}

			collection, err := collectionHelper.GetFilterCollection(info.Name())
			if err != nil {
				log.Println("Error: Getting Collection:" + err.Error())
				http.Error(response, err.Error(), 500)
				return err
			}

			collection.Name = info.Name()

			if collection.Domains == nil {
				collection.Domains = make([]string, 0, 0)
				for scanner.Scan() {
					collection.Domains = append(collection.Domains, scanner.Text())
				}
			} else {
				for scanner.Scan() {
					for _, existingDomain := range collection.Domains {
						if existingDomain == scanner.Text() {
							//Yeah lets goto! Aka Break 2
							goto Break2
						}
					}
					collection.Domains = append(collection.Domains, scanner.Text())

				Break2:
				}
			}

			if err := scanner.Err(); err != nil {
				log.Println("Warning: Importing Collection \"" + path + "\":" + err.Error())
			}

			err = collectionHelper.SetFilterCollection(collection)
			if err != nil {
				log.Println("Warning: Saving Collection \"" + collection.Name + "\":" + err.Error())
			}
		}
		return nil
	})
}

/*
 * Check if a user has access to a domain.
 */
func (t *Developer) access(response http.ResponseWriter, request *http.Request) {
	log.Println("Started")

	userFilterCollectionsHelper, err := datastore.GetUserFilterCollectionsHelper()
	if err != nil {
		log.Println("Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	userFilterCollections, err := userFilterCollectionsHelper.GetUserFilterCollections(mux.Vars(request)["username"])
	if err != nil {
		log.Println("Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	filterCollectionHelper, err := datastore.GetFilterCollectionHelper()
	if err != nil {
		log.Println("Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	for _, collectionName := range userFilterCollections.Collections {
		collection, err := filterCollectionHelper.GetFilterCollection(collectionName)
		if err != nil {
			log.Println("Error:" + err.Error())
			http.Error(response, err.Error(), 500)
			return
		}

		for _, domain := range collection.Domains {
			if domain == mux.Vars(request)["domain"] {
				log.Println("Blocked")
				return
			}
		}
	}

	log.Println("Not Blocked")
}

/*
 * Check if a user has access to a domain using the domain look up (Which should be quicker)
 */
func (t *Developer) quickaccess(response http.ResponseWriter, request *http.Request) {
	log.Println("Started")

	userFilterCollectionsHelper, err := datastore.GetUserFilterCollectionsHelper()
	if err != nil {
		log.Println("Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	userFilterCollections, err := userFilterCollectionsHelper.GetUserFilterCollections(mux.Vars(request)["username"])
	if err != nil {
		log.Println("Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	filterCollectionHelper, err := datastore.GetFilterCollectionHelper()
	if err != nil {
		log.Println("Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	inCollections, err := filterCollectionHelper.GetFilterCollectionsWithDomain(mux.Vars(request)["domain"])
	if err != nil {
		log.Println("Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	for _, collection := range inCollections {
		for _, userCollectionName := range userFilterCollections.Collections {
			if userCollectionName == collection.Name {
				log.Println("Blocked")
				return
			}
		}
	}

	log.Println("Not Blocked")
}

/*
 * Manually Assign an IP to a Mac Address
 */
func (t *Developer) assignip(response http.ResponseWriter, request *http.Request) {
	deviceHelper, err := datastore.GetDeviceHelper()
	if err != nil {
		log.Println("Error: Getting Helper Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	macAddress, err := net.ParseMAC(mux.Vars(request)["macaddress"])
	if err != nil {
		log.Println("Error: Invalid MACAddress :" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	device, err := deviceHelper.GetDevice(macAddress)
	if err != nil {
		log.Println("Error: Getting Device Error:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}

	device.IPAddress = net.ParseIP(mux.Vars(request)["ip"])

	err = deviceHelper.SetDevice(&device)
	if err != nil {
		log.Println("Error: Setting Device To Datastore:" + err.Error())
		http.Error(response, err.Error(), 500)
		return
	}
}

/*
 * Dump the current Connections Object to See the current open connections.
 */
func (t *Developer) dumpconnections(response http.ResponseWriter, request *http.Request) {
	connectionhelper, err := datastore.GetConnectionHelper()
	if err != nil {
		log.Printf("Error: Failed to get Connections Helper \"%v\"", err.Error())
		return
	}

	connections, err := connectionhelper.GetConnections()
	if err != nil {
		log.Printf("Error: Failed to get Connections \"%v\"", err.Error())
		return
	}

	log.Printf("Connections: \n %v \n", connections)
}

/*
 * Poke the kernel telling it to block connections from source to destination.
 */
func (t *Developer) block(response http.ResponseWriter, request *http.Request) {

	sourceip := []byte(net.ParseIP(mux.Vars(request)["source"]).To4())
	destinationip := []byte(net.ParseIP(mux.Vars(request)["destination"]).To4())

	blockedroutehelper, err := datastore.GetBlockedRouteHelper()
	if err != nil {
		log.Println("Error:" + err.Error())
		return
	}

	err = blockedroutehelper.AddBlockedRoute(sourceip, destinationip)

	if err != nil {
		log.Println("Error:" + err.Error())
	}

}
