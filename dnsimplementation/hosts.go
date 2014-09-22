package dnsimplementation

import (
	"github.com/d2g/cloudpathway/datastore"
	"log"
	"net"
)

type Hosts struct {
	Devices map[string]net.IP
}

func (this *Hosts) Add(hostname string, ip net.IP) error {
	this.Devices[hostname] = ip
	return nil
}

func (this *Hosts) Get(hostname string) (bool, net.IP, error) {
	ip := this.Devices[hostname]

	if ip != nil {
		return true, ip, nil
	} else {
		//So we need to check the Device Record to see if the hostname is there also.
		helper, err := datastore.GetDeviceHelper()
		if err != nil {
			log.Println("Error Getting Datastore Helper:" + err.Error())
			return false, ip, err
		}
		device, err := helper.GetDeviceByHostname(hostname)
		if err != nil {
			log.Println("Error Looking Up Hostname:" + hostname + ":Error:" + err.Error())
			return false, ip, err
		}

		return false, device.IPAddress, err
	}
}
