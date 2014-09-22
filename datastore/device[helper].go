package datastore

import (
	"bytes"
	"github.com/d2g/unqlitego"
	"log"
	"net"
)

type devicehelper struct {
	collection   *unqlitego.Database
	ipAddressKey *unqlitego.Database
	hostnameKey  *unqlitego.Database
}

var deviceHelperSingleton *devicehelper = nil

func GetDeviceHelper() (*devicehelper, error) {
	if deviceHelperSingleton == nil {
		var err error

		deviceHelperSingleton = new(devicehelper)
		deviceHelperSingleton.collection, err = unqlitego.NewDatabase("Device.unqlite")
		if err != nil {
			return deviceHelperSingleton, err
		}

		deviceHelperSingleton.ipAddressKey, err = unqlitego.NewDatabase("IPtoDevice.key.unqlite")
		if err != nil {
			return deviceHelperSingleton, err
		}

		deviceHelperSingleton.hostnameKey, err = unqlitego.NewDatabase("HostnametoDevice.key.unqlite")
		if err != nil {
			return deviceHelperSingleton, err
		}
	}

	return deviceHelperSingleton, nil
}

/**
 * Get device by MAC Address.
 */
func (this *devicehelper) GetDevice(macAddress net.HardwareAddr) (device Device, err error) {
	err = this.collection.GetObject(macAddress.String(), &device)
	if !bytes.Equal(device.MACAddress, macAddress) {
		device.MACAddress = macAddress
	}
	return
}

/**
 * Get device by IP Address.
 */
func (this *devicehelper) GetDeviceByIP(IPAddress net.IP) (device Device, err error) {

	var MACAddressString string

	err = this.ipAddressKey.GetObject(IPAddress.String(), &MACAddressString)
	if err != nil {
		return
	}

	if MACAddressString != "" {
		MACAddress, err := net.ParseMAC(MACAddressString)
		if err != nil {
			return device, err
		}

		return this.GetDevice(MACAddress)
	}

	return
}

/**
 * Get device by Hostname.
 */
func (this *devicehelper) GetDeviceByHostname(Hostname string) (device Device, err error) {

	var MACAddressString string

	err = this.hostnameKey.GetObject(Hostname, &MACAddressString)
	if err != nil {
		return
	}

	if MACAddressString != "" {
		MACAddress, err := net.ParseMAC(MACAddressString)
		if err != nil {
			return device, err
		}

		return this.GetDevice(MACAddress)
	}

	return
}

/**
 * Returns an array of devices with the given Username value as CurrentUser and another array
 * with the given Username value as DefaultUser.
 */
func (this *devicehelper) GetDevicesForUser(Username string) ([]Device, []Device, error) {
	currentUserDevices := make([]Device, 0, 0)
	defaultUserDevices := make([]Device, 0, 0)

	// Retrieve all devices in the collection.
	allDevices, err := this.GetDevices()

	// If we have an error, return.
	if err != nil {
		return currentUserDevices, defaultUserDevices, err
	}

	// Check each device...
	for _, device := range allDevices {
		// If the CurrentUser matches, add it to the currentUserDevices array.
		if device.CurrentUser != nil && device.CurrentUser.Username == Username {
			currentUserDevices = append(currentUserDevices, device)
		}

		// If the DefaultUser matches, add it to the defaultUserDevices array.
		if device.DefaultUser != nil && device.DefaultUser.Username == Username {
			defaultUserDevices = append(defaultUserDevices, device)
		}
	}

	return currentUserDevices, defaultUserDevices, nil
}

/**
 * Save the given device to the database.
 */
func (this *devicehelper) SetDevice(device *Device) error {
	//Marshal the item to JSON
	err := this.collection.SetObject(device.MACAddress.String(), device)
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}

	//Update the IPAddress Key if it has one :(
	if !device.IPAddress.Equal(net.IP{}) {
		err = this.ipAddressKey.SetObject(device.IPAddress.String(), device.MACAddress.String())
		if err != nil {
			log.Println("Error Saving IPAddress Key To Datastore")
			return err
		}
	}

	//Update the Hostname Key
	if device.Hostname != "" {
		err = this.hostnameKey.SetObject(device.Hostname, device.MACAddress.String())
		if err != nil {
			log.Println("Error Saving IPAddress Key To Datastore")
			return err
		}
	}
	return nil
}

/**
 * Deletes a device with the given MAC address.
 */
func (this *devicehelper) DeleteDevice(macAddress net.HardwareAddr) error {
	err := this.collection.DeleteObject(macAddress.String())

	if err != nil {
		log.Println("Error Deleting Device from Datastore")
	}

	return err
}

/**
 * Gets an array of all devices within the database.
 */
func (this *devicehelper) GetDevices() ([]Device, error) {

	devices := make([]Device, 0, 0)

	cursor, err := this.collection.NewCursor()
	defer cursor.Close()

	if err != nil {
		return devices, err
	}

	err = cursor.First()
	if err != nil {
		//You Get -28 When There are no records.
		if err == unqlitego.UnQLiteError(-28) {
			return devices, nil
		} else {
			return devices, err
		}
	}

	for {
		if !cursor.IsValid() {
			break
		}

		device := Device{}
		value, err := cursor.Value()

		if err != nil {

			log.Println("Error: Cursor Get Value Error:" + err.Error())

		} else {

			err := this.collection.Unmarshal()(value, &device)
			if err != nil {
				key, err := cursor.Key()
				if err != nil {
					log.Println("Error: Cursor Get Key Error:" + err.Error())
				} else {
					log.Println("Invalid Device in Datastore:" + string(key))
					device.MACAddress, err = net.ParseMAC(string(key))
					if err != nil {
						this.SetDevice(&device)
					}
				}
			}

			devices = append(devices, device)
		}

		err = cursor.Next()
		if err != nil {
			break
		}
	}

	err = cursor.Close()
	if err != nil {
		log.Println("3:" + err.Error())
	}
	return devices, err
}
