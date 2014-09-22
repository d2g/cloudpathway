package networktools

import (
	"log"
	"net"
)

func MainActiveIPNetwork() (activeipaddress net.IPNet, err error) {
	interfaces, err := ActiveIPInterfaces()
	if err != nil {
		return
	}

	for i := 0; i < len(interfaces); i++ {
		switch v := interfaces[0].(type) {
		case (*net.IPAddr):
			activeipaddress.IP = v.IP
			break
		case (*net.IPNet):
			activeipaddress = *v
			break
		}
	}

	if len(interfaces) > 1 {
		log.Printf("Warning: Multiple Active IP Interfaces Using %v\n", activeipaddress)
	}

	return
}
