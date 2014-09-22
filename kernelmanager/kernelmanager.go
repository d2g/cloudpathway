package kernelmanager

import (
	"github.com/d2g/netlink"
	"log"
)

type kernelManager struct {
	configuration *Configuration

	connection netlink.Connection
	input      chan Message
}

var kernelManagerSingleton *kernelManager = nil

func CreateKernelManager(c *Configuration) (*kernelManager, error) {
	if kernelManagerSingleton == nil {

		kernelManagerSingleton = new(kernelManager)
		kernelManagerSingleton.configuration = c

		if kernelManagerSingleton.configuration.QueueSize <= 0 {
			kernelManagerSingleton.configuration.QueueSize = 100
		}

		kernelManagerSingleton.input = make(chan Message, kernelManagerSingleton.configuration.QueueSize)

		kernelManagerSingleton.connection = netlink.GetNetlinkSocket(kernelManagerSingleton.configuration.Socket, netlink.Unicast)
	}

	return kernelManagerSingleton, nil
}

func GetKernelManager() *kernelManager {
	return kernelManagerSingleton
}

func (t *kernelManager) Listen() error {
	err := t.connection.Connect()
	if err != nil {
		return err
	}
	defer t.connection.Close()

	for {
		message := <-t.input

		log.Printf("Debug: Message:%v\n", message)

		err = t.connection.Write(message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *kernelManager) QueueMessage(m Message) {

	select {
	case t.input <- m:
	default:
		//If were not keeping up lets not make things worse..
		log.Println("Error: Queue For Kernel Is Full.. Discarding Message")
	}
}
