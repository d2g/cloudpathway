package datastore

import (
	"bytes"
	"github.com/d2g/dhcp4server/leasepool"
	"github.com/d2g/unqlitego"
	"log"
	"net"
)

type leasehelper struct {
	collection    *unqlitego.Database
	macaddressKey *unqlitego.Database

	currentLease int
}

var leaseHelperSingleton *leasehelper = nil

func GetLeaseHelper() (*leasehelper, error) {
	if leaseHelperSingleton == nil {
		var err error

		leaseHelperSingleton = new(leasehelper)
		leaseHelperSingleton.collection, err = unqlitego.NewDatabase("Lease.unqlite")
		if err != nil {
			return leaseHelperSingleton, err
		}

		leaseHelperSingleton.macaddressKey, err = unqlitego.NewDatabase("MacaddressToLease.key.unqlite")
		if err != nil {
			return leaseHelperSingleton, err
		}
	}
	return leaseHelperSingleton, nil
}

func (t *leasehelper) GetLeaseByIP(ip net.IP) (databaseLease leasepool.Lease, err error) {
	err = t.collection.GetObject(ip.String(), &databaseLease)
	return
}

func (t *leasehelper) GetLeaseByMacaddress(hardwareAddress net.HardwareAddr) (databaseLease leasepool.Lease, err error) {
	var leaseip string
	err = t.macaddressKey.GetObject(hardwareAddress.String(), &leaseip)
	if err != nil {
		return
	}

	if leaseip != "" {
		return t.GetLeaseByIP(net.ParseIP(leaseip))
	}

	return
}

func (t *leasehelper) GetLeases() ([]leasepool.Lease, error) {
	leases := make([]leasepool.Lease, 0, 0)

	cursor, err := t.collection.NewCursor()
	defer cursor.Close()

	if err != nil {
		return leases, err
	}

	err = cursor.First()
	if err != nil {
		if err == unqlitego.UnQLiteError(-28) {
			return leases, nil
		} else {
			return leases, err
		}
	}

	for {
		if !cursor.IsValid() {
			break
		}

		lease := leasepool.Lease{}
		value, err := cursor.Value()

		if err != nil {

			log.Println("Error: Cursor Get Value Error:" + err.Error())

		} else {

			err := t.collection.Unmarshal()(value, &lease)
			if err != nil {
				key, err := cursor.Key()
				if err != nil {
					log.Println("Error: Cursor Get Key Error:" + err.Error())
				} else {
					log.Println("Invalid Lease in Datastore:" + string(key))
					lease.IP = net.ParseIP(string(key))
					t.SetLease(&lease)
				}
			}

			leases = append(leases, lease)
		}

		err = cursor.Next()
		if err != nil {
			break
		}
	}

	err = cursor.Close()

	return leases, err
}

func (t *leasehelper) SetLease(lease *leasepool.Lease) error {
	err := t.collection.SetObject(lease.IP.String(), lease)
	if err != nil {
		log.Println("Error Saving Lease To Datastore")
		return err
	}

	if !bytes.Equal(lease.MACAddress, net.HardwareAddr{}) {
		err := t.macaddressKey.SetObject(lease.MACAddress.String(), lease.IP.String())
		if err != nil {
			log.Println("Error Saving Lease To Datastore")
			return err
		}
	}

	return nil
}

/**
 * Deletes a Lease with the given IP.
 */
func (t *leasehelper) DeleteLease(ip net.IP) error {
	return t.collection.Delete([]byte(ip.String()))
}
