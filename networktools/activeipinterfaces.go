package networktools

import (
	"net"
)

/*
 * Common Network Functions...
 */

func ActiveIPInterfaces() ([]net.Addr, error) {
	activeInterfaces := make([]net.Addr, 0)

	localInterfaces, err := net.InterfaceAddrs()
	if err != nil {
		return activeInterfaces, err
	}

	for _, localInterface := range localInterfaces {
		if localInterface, ok := localInterface.(*net.IPAddr); ok {
			if !localInterface.IP.Equal(net.IPv4(0, 0, 0, 0)) && !localInterface.IP.Equal(net.IPv4(127, 0, 0, 1)) && !localInterface.IP.IsLoopback() {
				activeInterfaces = append(activeInterfaces, localInterface)
			}
		}

		if localInterface, ok := localInterface.(*net.IPNet); ok {
			if !localInterface.IP.Equal(net.IPv4(0, 0, 0, 0)) && !localInterface.IP.Equal(net.IPv4(127, 0, 0, 1)) && !localInterface.IP.IsLoopback() {
				activeInterfaces = append(activeInterfaces, localInterface)
			}
		}
	}

	return activeInterfaces, nil
}
