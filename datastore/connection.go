package datastore

import (
	"bytes"
	"github.com/d2g/tcp"
	"golang.org/x/net/ipv4"
	"net"
	"time"
)

type Connection struct {
	/*
		A connection is from one point to another
	*/
	sourceIP   net.IP
	sourcePort uint16

	destinationIP   net.IP
	destinationPort uint16

	packets [][]byte
	updated time.Time

	/*
	 * The Device and User of That Device at the time are important to us but a connection is still a connection (i.e. you can't have 2 of them)
	 */
	deviceid net.HardwareAddr
	username string
}

func NewConnectionFromPacket(packet []byte) (Connection, error) {

	ipheader, err := ipv4.ParseHeader(packet)
	if err != nil {
		return Connection{}, err
	}

	tcpheader, err := tcp.ParseHeader(packet[ipheader.Len:])
	if err != nil {
		return Connection{}, err
	}

	return NewConnectionFromPacketAndIPHeaderAndTCPHeader(packet, ipheader, tcpheader)
}

func NewConnectionFromPacketAndIPHeaderAndTCPHeader(packet []byte, ipheader *ipv4.Header, tcpheader *tcp.Header) (Connection, error) {
	c := Connection{}
	c.sourceIP = ipheader.Src
	c.sourcePort = tcpheader.SourcePort
	c.destinationIP = ipheader.Dst
	c.destinationPort = tcpheader.DestinationPort
	err := c.AddPacket(packet)
	if err != nil {
		return c, err
	}

	deviceHelper, err := GetDeviceHelper()
	if err != nil {
		return c, err
	}

	device, err := deviceHelper.GetDeviceByIP(c.sourceIP)
	if err != nil {
		return c, err
	}

	if bytes.Equal(device.MACAddress, net.HardwareAddr{}) {
		device, err = deviceHelper.GetDeviceByIP(c.sourceIP)
		if err != nil {
			return c, err
		}
	}

	c.deviceid = device.MACAddress
	if device.GetActiveUser() != nil {
		c.username = device.GetActiveUser().Username
	}

	return c, nil
}

func (t *Connection) AddPacket(packet []byte) error {
	t.packets = append(t.packets, packet)
	t.updated = time.Now()
	return nil
}

func (t *Connection) NumberOfPackets() int {
	return len(t.packets)
}

func (t *Connection) Packets() [][]byte {
	return t.packets
}

func (t *Connection) SourceIP() net.IP {
	return t.sourceIP
}

func (t *Connection) SourcePort() uint16 {
	return t.sourcePort
}

func (t *Connection) DestinationIP() net.IP {
	return t.destinationIP
}

func (t *Connection) DestinationPort() uint16 {
	return t.destinationPort
}

func (t *Connection) DeviceID() net.HardwareAddr {
	return t.deviceid
}

func (t *Connection) Username() string {
	return t.username
}

func (t *Connection) Updated() time.Time {
	return t.updated
}

type Connections []Connection

func (t Connections) Usernames() []string {
	tmpUsernames := make(map[string]bool)
	usernames := make([]string, 0)

	for _, v := range []Connection(t) {
		tmpUsernames[v.Username()] = true
	}

	for k := range tmpUsernames {
		usernames = append(usernames, k)
	}

	return usernames
}

func (t Connections) DeviceIDs() []net.HardwareAddr {
	devices := make([]net.HardwareAddr, 0)

	for _, v := range []Connection(t) {
		needed := true

		for i := range devices {
			if bytes.Equal(devices[i], v.DeviceID()) {
				needed = false
			}
		}

		if needed {
			devices = append(devices, v.DeviceID())
		}
	}

	return devices
}

func (t Connections) DeviceIDsForUsername(username string) []net.HardwareAddr {
	devices := make([]net.HardwareAddr, 0)

	for _, v := range []Connection(t) {
		if v.Username() != username {
			continue
		}
		needed := true

		for i := range devices {
			if bytes.Equal(devices[i], v.DeviceID()) {
				needed = false
			}
		}

		if needed {
			devices = append(devices, v.DeviceID())
		}
	}

	return devices
}

func (t Connections) ConnectionsForUsernameAndDeviceID(username, deviceid string) []Connection {
	connections := []Connection{}

	for _, connection := range t {
		if connection.Username() == username &&
			connection.DeviceID().String() == deviceid {
			connections = append(connections, connection)
		}
	}

	return connections
}

func GetExampleConnections() []Connection {
	collections := []Connection{}

	testDevice1Mac, _ := net.ParseMAC("5c:26:0a:3a:08:16")
	testDevice2Mac, _ := net.ParseMAC("9f:7b:03:a6:b2:9c")

	//First Device On google.co.uk
	test := Connection{
		sourceIP:        net.IPv4(192, 168, 1, 100),
		sourcePort:      80,
		destinationIP:   net.IPv4(173, 194, 34, 151), //google.co.uk
		destinationPort: 80,
		updated:         time.Now(),
		username:        "goldsmithd",
		deviceid:        testDevice1Mac,
	}
	collections = append(collections, test)

	//Second Device Also On google.co.uk
	test = Connection{
		sourceIP:        net.IPv4(192, 168, 1, 101),
		sourcePort:      80,
		destinationIP:   net.IPv4(173, 194, 34, 151), //google.co.uk
		destinationPort: 80,
		updated:         time.Now(),
		username:        "goldsmithd",
		deviceid:        testDevice2Mac,
	}
	collections = append(collections, test)

	//First Device Also On slashdot.org
	test = Connection{
		sourceIP:        net.IPv4(192, 168, 1, 100),
		sourcePort:      80,
		destinationIP:   net.IPv4(216, 34, 181, 45), //slashdot.org
		destinationPort: 80,
		updated:         time.Now(),
		username:        "smith",
		deviceid:        testDevice1Mac,
	}
	collections = append(collections, test)

	return collections
}
