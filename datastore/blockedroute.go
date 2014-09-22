package datastore

import (
	"net"
)

type BlockedRoute struct {
	Source      net.IP
	Destination net.IP
}

func (t *BlockedRoute) AsBytes() []byte {
	return append([]byte(t.Source.To4()), []byte(t.Destination.To4())...)
}
