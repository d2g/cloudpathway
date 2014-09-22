package datastore

import (
	"github.com/d2g/cloudpathway/kernelmanager"
	"github.com/d2g/unqlitego"
	"labix.org/v2/mgo/bson"
	"log"
	"net"
)

type blockroutehelper struct {
	collection *unqlitego.Database
}

var blockRouteHelperSingleton *blockroutehelper = nil

func GetBlockedRouteHelper() (*blockroutehelper, error) {
	if blockRouteHelperSingleton == nil {
		var err error

		blockRouteHelperSingleton = new(blockroutehelper)
		blockRouteHelperSingleton.collection, err = unqlitego.NewDatabase("Blocked.unqlite")
		blockRouteHelperSingleton.collection.SetMarshal(bson.Marshal)
		blockRouteHelperSingleton.collection.SetUnmarshal(bson.Unmarshal)
		if err != nil {
			return blockRouteHelperSingleton, err
		}
	}

	return blockRouteHelperSingleton, nil
}

func (t *blockroutehelper) GetBlockedRoutesByIP(clientIP net.IP) (routes []BlockedRoute, err error) {
	cursor, err := t.collection.NewCursor()
	defer cursor.Close()
	if err != nil {
		return
	}

	err = cursor.First()
	if err != nil {
		//You Get -28 When There are no records.
		if err == unqlitego.UnQLiteError(-28) {
			//No Records in the DB.
			err = nil
			return
		} else {
			return
		}
	}

	for {
		if !cursor.IsValid() {
			break
		}

		value, err := cursor.Value()
		if err != nil {
			log.Println("Error: Cursor Get Value Error:" + err.Error())
		} else {
			tmpRecord := BlockedRoute{}
			err = t.collection.Unmarshal()(value, &tmpRecord)
			if err != nil {
				log.Println("Error: Can't Unmarshal Block Record (Deleting)")
				cursor.Delete()
			}

			if tmpRecord.Source.Equal(clientIP) || tmpRecord.Destination.Equal(clientIP) {
				//err = cursor.Delete()
				//We could Optimise this to delete the record from the datastore now rather than using the function.
				routes = append(routes, tmpRecord)
			}
		}

		err = cursor.Next()
		if err != nil {
			break
		}
	}

	err = cursor.Close()
	return
}

func (t *blockroutehelper) IsBlockedRoute(source net.IP, destination net.IP) (bool, error) {

	_, err := t.collection.Fetch(append([]byte(source.To4()), []byte(destination.To4())...))
	if err != nil {
		if err == unqlitego.UnQLiteError(-6) {
			//Not Found
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (t *blockroutehelper) DeleteBlockedRoutesForIP(clientip net.IP) error {
	routes, err := t.GetBlockedRoutesByIP(clientip)
	if err != nil {
		return err
	}

	for i := range routes {
		err := t.DeleteBlockRoute(routes[i].Source, routes[i].Destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *blockroutehelper) DeleteBlockRoute(source net.IP, destination net.IP) error {

	//If it's currently blocked then remove it from the kernel.
	blocked, err := t.IsBlockedRoute(source, destination)
	if err != nil {
		return err
	}

	if blocked {
		//Delete it from store and kernel.
		err := t.collection.Delete(append([]byte(source.To4()), []byte(destination.To4())...))
		if err != nil {
			return err
		}

		kernelManager := kernelmanager.GetKernelManager()

		message := kernelmanager.Message([]byte{0, 0})
		message.SetType(kernelmanager.UNBLOCK)
		message.Append(append([]byte(source.To4()), []byte(destination.To4())...))

		kernelManager.QueueMessage(message)

	}

	return nil
}

func (t *blockroutehelper) AddBlockedRoute(source net.IP, destination net.IP) error {
	blocked, err := t.IsBlockedRoute(source, destination)
	if err != nil {
		return err
	}

	if !blocked {
		route := BlockedRoute{
			Source:      source,
			Destination: destination,
		}

		byteObject, err := t.collection.Marshal()(route)
		if err != nil {
			return err
		}

		err = t.collection.Store(route.AsBytes(), byteObject)
		if err != nil {
			return err
		}

		kernelManager := kernelmanager.GetKernelManager()

		message := kernelmanager.Message([]byte{0, 0})
		message.SetType(kernelmanager.BLOCK)
		message.Append(route.AsBytes())

		kernelManager.QueueMessage(message)
		return nil
	}

	return nil
}
