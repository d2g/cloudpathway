package datastore

import (
	"net"
	"sync"
)

type connectionhelper struct {
	connections []Connection
	lock        sync.Mutex
}

var connectionHelperSingleton *connectionhelper = nil

func GetConnectionHelper() (*connectionhelper, error) {
	if connectionHelperSingleton == nil {
		connectionHelperSingleton = new(connectionhelper)
		connectionHelperSingleton.connections = make([]Connection, 0)
	}

	return connectionHelperSingleton, nil
}

func (t *connectionhelper) Lock() {
	t.lock.Lock()
}

func (t *connectionhelper) Unlock() {
	t.lock.Unlock()
}

func (t *connectionhelper) GetConnection(sIP net.IP, sPort uint16, dIP net.IP, dPort uint16) (Connection, error) {

	for i := range t.connections {
		if t.connections[i].SourceIP().Equal(sIP) &&
			t.connections[i].SourcePort() == sPort &&
			t.connections[i].DestinationIP().Equal(dIP) &&
			t.connections[i].DestinationPort() == dPort {

			return t.connections[i], nil
		}
	}
	return Connection{}, nil
}

func (t *connectionhelper) SetConnection(c Connection) error {
	for i := range t.connections {
		if t.connections[i].SourceIP().Equal(c.SourceIP()) &&
			t.connections[i].SourcePort() == c.SourcePort() &&
			t.connections[i].DestinationIP().Equal(c.DestinationIP()) &&
			t.connections[i].DestinationPort() == c.DestinationPort() {

			t.connections[i] = c
			return nil
		}
	}

	t.connections = append(t.connections, c)
	return nil
}

func (t *connectionhelper) RemoveConnection(c Connection) error {
	for i := range t.connections {
		if t.connections[i].SourceIP().Equal(c.SourceIP()) &&
			t.connections[i].SourcePort() == c.SourcePort() &&
			t.connections[i].DestinationIP().Equal(c.DestinationIP()) &&
			t.connections[i].DestinationPort() == c.DestinationPort() {

			//Replace this item with the last.
			t.connections[i] = t.connections[len(t.connections)-1]
			//Reduce the size and lose the last (Now Duplicate) item
			t.connections = t.connections[:len(t.connections)-1]
			return nil
		}
	}

	return nil
}

func (t *connectionhelper) GetConnections() ([]Connection, error) {
	return t.connections, nil
}
