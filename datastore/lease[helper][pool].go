package datastore

import (
	"github.com/d2g/dhcp4server/leasepool"
	"log"
	"net"
)

/*
 * Lease Pool Implementation
 * ......
 * ......
 * ......
 */
//Add A Lease To The Pool
func (t *leasehelper) AddLease(lease leasepool.Lease) error {

	existingLease, err := t.GetLeaseByIP(lease.IP)
	if err != nil {
		return err
	}

	if existingLease.IP.Equal(net.IP{}) {
		//Lease Doesn't Exists
		tmpLease := leasepool.Lease{IP: lease.IP}
		err = t.SetLease(&tmpLease)
		if err != nil {
			return err
		}
	}

	return nil
}

//Remove
func (t *leasehelper) RemoveLease(ipaddress net.IP) error {
	return t.DeleteLease(ipaddress)
}

//Remove All Leases from the Pool (Required for Persistant LeaseManagers)
func (t *leasehelper) PurgeLeases() error {
	leases, err := t.GetLeases()
	if err != nil {
		return err
	}

	for _, lease := range leases {
		err = t.DeleteLease(lease.IP)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
 * Get the Lease
 * -Found
 * -Copy Of the Lease
 * -Any Error
 */
func (t *leasehelper) GetLease(ipaddress net.IP) (bool, leasepool.Lease, error) {
	databaseLease, err := t.GetLeaseByIP(ipaddress)
	if err != nil {
		return false, leasepool.Lease{}, err
	}

	if databaseLease.IP.Equal(ipaddress) {
		return true, databaseLease, nil
	} else {
		return false, databaseLease, nil
	}
}

//Get the lease already in use by that hardware address.
func (t *leasehelper) GetLeaseForHardwareAddress(hardwareAddress net.HardwareAddr) (bool, leasepool.Lease, error) {
	currentLease, err := t.GetLeaseByMacaddress(hardwareAddress)

	if err != nil {
		return false, currentLease, err
	}

	if currentLease.IP.Equal(net.IP{}) {
		return false, leasepool.Lease{}, nil
	}

	return true, currentLease, nil
}

/*
 * -Lease Available
 * -Lease
 * -Error
 */
func (t *leasehelper) GetNextFreeLease() (bool, leasepool.Lease, error) {
	leases, err := t.GetLeases()
	if err != nil {
		return false, leasepool.Lease{}, err
	}

	for i := 0; i < len(leases); i++ {
		if t.currentLease >= len(leases) {
			t.currentLease = 0
		}

		if leases[t.currentLease].Status == leasepool.Free {
			t.currentLease++
			return true, leases[t.currentLease-1], nil
		} else {
			t.currentLease++
		}
	}

	return false, leasepool.Lease{}, nil
}

/*
 * Update Lease
 * - Has Updated
 * - Error
 */
func (t *leasehelper) UpdateLease(lease leasepool.Lease) (bool, error) {
	currentLease, err := t.GetLeaseByIP(lease.IP)
	if err != nil {
		return false, err
	}

	if currentLease.IP.Equal(net.IP{}) {
		return false, nil
	} else {
		if lease.Status == leasepool.Active {

			//Attach the IP to the Device
			devicehelper, err := GetDeviceHelper()
			if err != nil {
				return false, err
			}

			//If OK
			device, err := devicehelper.GetDevice(lease.MACAddress)
			if err != nil {
				log.Println("Error: Unable to get Device for:" + lease.MACAddress.String())
				return false, err
			}

			device.IPAddress = lease.IP
			device.Hostname = lease.Hostname

			blockedroutehelper, err := GetBlockedRouteHelper()
			if err != nil {
				log.Println("Error: Unable to get Blocked Route Helper:" + err.Error())
				return false, err
			}

			err = blockedroutehelper.DeleteBlockedRoutesForIP(device.IPAddress)
			if err != nil {
				log.Println("Error: Unable to Delete Blocked Routes for IP:" + err.Error())
				return false, err
			}

			//TODO: Send Blocked Sites for this IP.

			err = devicehelper.SetDevice(&device)
			if err != nil {
				log.Println("Error: Setting Device with New IP:" + err.Error())
				return false, err
			}

		} else {
			//All other cases make sure the IP is not assigned to a device.

		}

		err := t.SetLease(&lease)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}
