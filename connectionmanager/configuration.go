package connectionmanager

type Configuration struct {
	Disabled bool

	Reader struct {
		Socket    int
		QueueSize int
	}

	Manager struct {
		Agents     int //Number of manager routines
		Timeout    int //Number of seconds before we assume the connection has closed
		MaxPackets int //Number of packets before we...
		QueueSize  int //Number of Closed Connections on the outbound queue.
		GCDisabled bool
		//TODO: Ignore     []net.IP //Any Hosts That Should Be ignored.
	}

	Classification struct {
		Agents    int //Number of classification Routines
		QueueSize int //The maximum outbound queue size.
		Dump      struct {
			All     bool
			Unknown bool
			Path    string
		}
	}
}
