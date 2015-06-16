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
	configuration *Configuration

	reader struct {
		connection netlink.Connection
		output     chan []byte
	}

	manager struct {
		output chan datastore.Connection
	}

	classification struct {
		output chan datastore.ClassifiedConnection
	}
}

func NewConnectionManager(c *Configuration) (connectionManager, error) {
	cm := connectionManager{}
	cm.configuration = c

	if cm.configuration.Reader.QueueSize <= 0 {
		cm.configuration.Reader.QueueSize = 100
	}
	cm.reader.output = make(chan []byte, cm.configuration.Reader.QueueSize)

	if cm.configuration.Manager.QueueSize <= 0 {
		cm.configuration.Manager.QueueSize = 100
	}
	cm.manager.output = make(chan datastore.Connection, cm.configuration.Manager.QueueSize)

	if cm.configuration.Classification.QueueSize <= 0 {
		cm.configuration.Classification.QueueSize = 100
	}
	cm.classification.output = make(chan datastore.ClassifiedConnection, cm.configuration.Classification.QueueSize)

	cm.reader.connection = netlink.GetNetlinkSocket(cm.configuration.Reader.Socket, netlink.Broadcast)
	cm.reader.connection.SetHandleFunc(func(message []byte) error {
		select {
		case cm.reader.output <- message:
		default:
			//If were not keeping up lets not make things worse..
			log.Println("Error: Reader Queue is full! Discarding Packet..")
		}

		return nil
	})

	return cm, nil
}

func (t *connectionManager) Listen() error {
	return t.reader.connection.ListenAndServe()
}

func (t *connectionManager) Process() error {

	var wg sync.WaitGroup

	for i := 0; i < t.configuration.Manager.Agents; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				message := <-t.reader.output

				ipheader, err := ipv4.ParseHeader(message)
				if err != nil {
					log.Println("Error: Parsing IPHeader:" + err.Error())
					log.Printf("Message:%v\n", message)
					continue
				}

				if int(ipheader.Flags) < 2 {
					log.Printf("Warning: More Fragments Flag Is Set\n")
				}

				tcpheader, err := tcp.ParseHeader(message[ipheader.Len:])
				if err != nil {
					log.Println("Error: Parsing TCPHeader:" + err.Error())
					log.Printf("Message:%v\n", message)
					continue
				}

				connectionhelper, err := datastore.GetConnectionHelper()
				if err != nil {
					log.Printf("Error: Getting Connection Helper: %v\n", err.Error())
					return
				}

				connectionhelper.Lock()
				connection, err := connectionhelper.GetConnection(ipheader.Src, tcpheader.SourcePort, ipheader.Dst, tcpheader.DestinationPort)
				if err != nil {
					log.Printf("Error: Getting Connection: %v\n", err.Error())
					connectionhelper.Unlock()
					return
				}

				if bytes.Equal(connection.SourceIP(), net.IP{}) {
					//No connection found
					connection, err = datastore.NewConnectionFromPacketAndIPHeaderAndTCPHeader(message, ipheader, tcpheader)
					if err != nil {
						log.Printf("Error: Creating New Connection: %v\n", err.Error())
						connectionhelper.Unlock()
						continue
					}

					err = connectionhelper.SetConnection(connection)
					if err != nil {
						log.Printf("Error: Saving New Connection: %v\n", err.Error())
						connectionhelper.Unlock()
						continue
					}

				} else {
					//Connection found
					//We need to find out if the connection is still for this device and user
					devicehelper, err := datastore.GetDeviceHelper()
					if err != nil {
						log.Printf("Error: Getting Device Helper: %v\n", err.Error())
						connectionhelper.Unlock()
						return
					}

					device, err := devicehelper.GetDeviceByIP(ipheader.Src)
					if err != nil {
						log.Printf("Error: Getting Source Device By IP(%v) \"%v\"\n", ipheader.Src.String(), err.Error())
					}

					if bytes.Equal(device.MACAddress, net.HardwareAddr{}) {
						device, err = devicehelper.GetDeviceByIP(ipheader.Dst)
						if err != nil {
							log.Printf("Error: Getting Destination Device By IP(%v) \"%v\"\n", ipheader.Dst.String(), err.Error())
						}
					}

					//If we have a device, is now or it's empty
					currentuser := ""
					if device.GetActiveUser() != nil {
						currentuser = device.GetActiveUser().Username
					}

					//Device, User or Number of In Memory Packets Breached...
					if !bytes.Equal(device.MACAddress, connection.DeviceID()) || (connection.Username() != currentuser) || connection.NumberOfPackets() >= t.configuration.Manager.MaxPackets {
						//Device or User has changed.
						err = connectionhelper.RemoveConnection(connection)
						if err != nil {
							log.Printf("Error: Removing Connection: %v\n", err.Error())
							connectionhelper.Unlock()
							continue
						}

						//If a channel is registered and we had no errors closing it.
						if t.manager.output != nil {
							select {
							case t.manager.output <- connection:
							default:
								log.Println("Error: Manager Queue is full! Discarding Closed Connection..")
							}

						}

						//Create A New Connection
						connection, err := datastore.NewConnectionFromPacketAndIPHeaderAndTCPHeader(message, ipheader, tcpheader)
						if err != nil {
							log.Printf("Error: Creating New Connection After Remove: %v\n", err.Error())
							connectionhelper.Unlock()
							continue
						}

						err = connectionhelper.SetConnection(connection)
						if err != nil {
							log.Printf("Error: Saving New Connection After Remove: %v\n", err.Error())
							connectionhelper.Unlock()
							continue
						}
					} else {
						err = connection.AddPacket(message)
						if err != nil {
							log.Printf("Error: Adding Packet: %v\n", err.Error())
							connectionhelper.Unlock()
							continue
						}

						err = connectionhelper.SetConnection(connection)
						if err != nil {
							log.Printf("Error: Saving Packet after adding packet: %v\n", err.Error())
							connectionhelper.Unlock()
							continue
						}
					}

				}
				connectionhelper.Unlock()
			}

		}()
	}

	wg.Wait()
	log.Println("Error: We seem to be existing packet processing??")
	return errors.New("Error: Packet Processing Exited?")
}

/*
Garbage Collection service for detecting connections that have closed.
*/
func (t *connectionManager) GC() error {

	nextrun := time.Now().Add(time.Second * time.Duration(t.configuration.Manager.Timeout))

	for {
		select {
		case <-time.After(time.Second * time.Duration(t.configuration.Manager.Timeout)):
		case <-time.After(nextrun.Sub(time.Now())):
		}

		//Run GC
		connectionhelper, err := datastore.GetConnectionHelper()
		if err != nil {
			log.Printf("Error: Getting GC Connection Helper: %v\n", err.Error())
			return err
		}

		//We don't want to lock the connections if we can help it.
		connections, err := connectionhelper.GetConnections()
		if err != nil {
			log.Printf("Error: Getting GC Connections: %v\n", err.Error())
			return err
		}

		for i := range connections {
			if connections[i].Updated().Add(time.Second * time.Duration(t.configuration.Manager.Timeout)).Before(time.Now()) {
				//We should GC this record.
				//Lock The records
				connectionhelper.Lock()
				//We Have to get the connection again as it might have changed.
				//Use the get connection as although it's just a loop again we only optimise one place.
				connection, err := connectionhelper.GetConnection(connections[i].SourceIP(), connections[i].SourcePort(), connections[i].DestinationIP(), connections[i].DestinationPort())
				if err != nil {
					log.Printf("Error: Getting GC Connection: %v\n", err.Error())
					connectionhelper.Unlock()
					continue
				}

				if !connection.SourceIP().Equal(connections[i].SourceIP()) {
					//The connection was removed by another thread.
					connectionhelper.Unlock()
					continue
				}

				//If the connection hasn't been updated by another thread.
				if connection.Updated().Add(time.Second * time.Duration(t.configuration.Manager.Timeout)).Before(time.Now()) {
					err := connectionhelper.RemoveConnection(connection)
					if err != nil {
						log.Printf("Error: GCing Connection: %v\n", err.Error())
					}

					//Pass closed connections on if they are wanted.
					if t.manager.output != nil {
						select {
						case t.manager.output <- connection:
						default:
							log.Println("Error: Manager Queue is full! Discarding Closed Connection..")
						}

					}

				}

				connectionhelper.Unlock()
			}

			nextrun = time.Now().Add(time.Second * time.Duration(t.configuration.Manager.Timeout))
			//We need to workout if nextrun should be less that the normal 15 seconds.
			if connections[i].Updated().Add(time.Second * time.Duration(t.configuration.Manager.Timeout)).Before(nextrun) {
				nextrun = connections[i].Updated().Add(time.Second * time.Duration(t.configuration.Manager.Timeout))
			}
		}
	}
	return nil
}

/*
Routine Function for Processing Messages off the Channel
*/
func (t *connectionManager) Classify() error {

	var wg sync.WaitGroup

	for i := 0; i < t.configuration.Classification.Agents; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				connection := <-t.manager.output

				isConnectionDumped := false

				if t.configuration.Classification.Dump.All {
					err := t.DumpConnectionToDisk(&connection)
					if err != nil {
						log.Println("Error: Dumping file Error:" + err.Error())
					} else {
						isConnectionDumped = true
					}
				}

				//Classify
				classifiedconnection, err := datastore.NewClassifiedConnection(connection)
				if err != nil {
					log.Println("Error: Unable Create Classified Connection" + err.Error())
					return
				}

				if t.configuration.Classification.Dump.Unknown && !isConnectionDumped && classifiedconnection.Protocol() == "Unknown" {
					err := t.DumpConnectionToDisk(&connection)
					if err != nil {
						log.Println("Error: Dumping file Error:" + err.Error())
					} else {
						isConnectionDumped = true
					}
				}

				log.Printf("Debug: Connection Classified As \"%s\"", classifiedconnection.Protocol())
				if t.classification.output != nil {
					select {
					case t.classification.output <- classifiedconnection:
					default:
						log.Println("Error: Classified Connection Queue is full! Discarding Classified Connection..")
					}
				}
			}
		}()
	}

	wg.Wait()
	log.Println("Error: We seem to be existing packet processing??")
	return errors.New("Error: Packet Processing Exited?")
}

func (t *connectionManager) Output() chan datastore.ClassifiedConnection {
	return t.classification.output
}

func (t *connectionManager) DumpConnectionToDisk(connection *datastore.Connection) error {
	currentTime := time.Now()
	randomid := strconv.Itoa(rand.Int())
	for i := range connection.Packets() {
		err := ioutil.WriteFile(t.configuration.Classification.Dump.Path+currentTime.Format("20060102150405")+"_"+connection.SourceIP().String()+"_"+strconv.Itoa(int(connection.SourcePort()))+"_"+connection.DestinationIP().String()+"_"+strconv.Itoa(int(connection.DestinationPort()))+"_"+strconv.Itoa(i)+"_"+randomid+".dump", connection.Packets()[i], 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rand.Seed(time.Now().Unix())
}
