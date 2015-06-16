package connectionmanager

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/netlink"
	"github.com/d2g/tcp"
	"golang.org/x/net/ipv4"
)

type connectionManager struct {
	active []connection
}
