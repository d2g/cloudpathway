package datastore

import (
	"github.com/d2g/dhcp4server/leasepool"
	"log"
	"net"
)

/*
Lease Pool Implementation
*/

//Add A Lease To The Pool
func (t *leasehelper) AddLease(lease leasepool.Lease) error {
	log.Printf("Trace: Adding lease \"%s\"\n", lease.IP.String())

	existingLease, err := t.GetLeaseByIP(lease.IP)
	if err != nil {
		log.Printf("Debug: Error adding lease \"%s\" because of \"%s\"\n", lease.IP.String(), err.Error())
		return err
	}

	if existingLease.IP.Equal(net.IP{}) {
		log.Printf("Trace: Lease \"%s\" doesn't exist in pool adding it.\n", lease.IP.String())

		//Lease Doesn't Exists
		tmpLease := leasepool.Lease{IP: lease.IP}
		err = t.SetLease(&tmpLease)
		if err != nil {
			log.Printf("Debug: Error Adding Lease \"%s\" to pool because of \"%s\"\n", lease.IP.String(), err.Error())
			return err
		}
	}

	log.Printf("Trace: Lease \"%s\" is in pool\n", lease.IP.String())
	return nil
}

//Remove
func (t *leasehelper) RemoveLease(ipaddress net.IP) error {
	log.Printf("Trace: Removing Lease \"%s\" from pool\n", ipaddress.String())
	return t.DeleteLease(ipaddress)
}

//Remove All Leases from the Pool (Required for Persistant LeaseManagers)
func (t *leasehelper) PurgeLeases() error {
	log.Printf("Trace: Purging leases\n")
	leases, err := t.GetLeases()
	if err != nil {
		log.Printf("Debug: Error getting leases because of \"%s\"\n", err.Error())
		return err
	}

	for _, lease := range leases {
		err = t.DeleteLease(lease.IP)
		if err != nil {
			log.Printf("Debug: Error deleting lease \"%s\" from pool because of \"%s\"\n", lease.IP.String(), err.Error())
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
	log.Printf("Trace: GetLease \"%s\"\n", ipaddress.String())
	databaseLease, err := t.GetLeaseByIP(ipaddress)
	if err != nil {
		log.Printf("Debug: Error getting Lease \"%s\" because of \"%s\"\n", ipaddress.String(), err.Error())
		return false, leasepool.Lease{}, err
	}

	if databaseLease.IP.Equal(ipaddress) {
		log.Printf("Trace: Lease Found \"%s\"\n", ipaddress.String())
		return true, databaseLease, nil
	} else {
		log.Printf("Trace: Lease Not Found \"%s\"\n", ipaddress.String())
		return false, databaseLease, nil
	}
}

//Get the lease already in use by that hardware address.
func (t *leasehelper) GetLeaseForHardwareAddress(hardwareAddress net.HardwareAddr) (bool, leasepool.Lease, error) {
	log.Printf("Trace: Getting Lease For Hardware Address \"%s\"\n", hardwareAddress.String())
	currentLease, err := t.GetLeaseByMacaddress(hardwareAddress)

	if err != nil {
		return false, currentLease, err
	}

	if currentLease.IP.Equal(net.IP{}) {
		log.Printf("Trace: Didn't find lease for \"%s\"\n", hardwareAddress.String())
		return false, leasepool.Lease{}, nil
	}

	log.Printf("Trace: Find lease for \"%s\" it's \"%s\"\n", hardwareAddress.String(), currentLease.IP.String())
	return true, currentLease, nil
}

/*
 * -Lease Available
 * -Lease
 * -Error
 */
func (t *leasehelper) GetNextFreeLease() (bool, leasepool.Lease, error) {
	log.Printf("Trace: Get net free lease\n")
	leases, err := t.GetLeases()
	if err != nil {
		return false, leasepool.Lease{}, err
	}

	log.Printf("Trace: We have %d Leases\n", len(leases))
	for i := 0; i < len(leases); i++ {
		if t.currentLease >= len(leases) {
			t.currentLease = 0
		}

		if leases[t.currentLease].Status == leasepool.Free {
			log.Printf("Trace: Free Lease \"%s\"\n", leases[t.currentLease].IP.String())
			t.currentLease++
			return true, leases[t.currentLease-1], nil
		} else {
			t.currentLease++
		}
	}

	log.Printf("Debug: No Free Leases \n")
	return false, leasepool.Lease{}, nil
}

/*
 * Update Lease
 * - Has Updated
 * - Error
 */
func (t *leasehelper) UpdateLease(lease leasepool.Lease) (bool, error) {
	log.Printf("Trace: Updating lease \"%s\"\n", lease.IP.String())
	currentLease, err := t.GetLeaseByIP(lease.IP)
	if err != nil {
		return false, err
	}

	if currentLease.IP.Equal(net.IP{}) {
		log.Printf("Debug: Updating lease \"%s\" doesn't exist in our pool\n", lease.IP.String())
		return false, nil
	} else {
		if lease.Status == leasepool.Active {
			log.Printf("Trace: Updating lease \"%s\" setting active.\n", lease.IP.String())

			//Attach the IP to the Device
			devicehelper, err := GetDeviceHelper()
			if err != nil {
				log.Printf("Debug: Updating lease \"%s\" failed getting device because of \"%s\".\n", lease.IP.String(), err.Error())
				return false, err
			}

			//If OK
			device, err := devicehelper.GetDevice(lease.MACAddress)
			if err != nil {
				log.Printf("Debug: Updating lease \"%s\" failed getting device \"%s\" because of \"%s\".\n", lease.IP.String(), lease.MACAddress.String(), err.Error())
				return false, err
			}

			device.IPAddress = lease.IP
			device.Hostname = lease.Hostname

			blockedroutehelper, err := GetBlockedRouteHelper()
			if err != nil {
				log.Printf("Debug: Updating lease \"%s\" failed getting blocked route because of \"%s\".\n", lease.IP.String(), err.Error())
				return false, err
			}

			err = blockedroutehelper.DeleteBlockedRoutesForIP(device.IPAddress)
			if err != nil {
				log.Printf("Debug: Updating lease \"%s\" failed removing blocked routes because of \"%s\".\n", lease.IP.String(), err.Error())
				return false, err
			}

			//TODO: Send Blocked Sites for this IP.

			err = devicehelper.SetDevice(&device)
			if err != nil {
				log.Printf("Debug: Updating lease \"%s\" failed device \"%s\" update because of \"%s\".\n", lease.IP.String(), device.MACAddress.String(), err.Error())
				return false, err
			}

		} else {
			log.Printf("Trace: Updating lease \"%s\" setting something other than active.\n", lease.IP.String())
			//All other cases make sure the IP is not assigned to a device.
			//TODO:
		}

		err := t.SetLease(&lease)
		if err != nil {
			log.Printf("Debug: Updating lease \"%s\" failed because of \"%s\".\n", lease.IP.String(), err.Error())
			return false, err
		} else {
			log.Printf("Trace: Updated lease \"%s\"\n", lease.IP.String())
			return true, nil
		}
	}
}
