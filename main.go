package main

import (
	"bytes"
	//"database/sql"
	"encoding/json"
	"github.com/d2g/cloudpathway/connectionmanager"
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/cloudpathway/dnsimplementation"
	"github.com/d2g/cloudpathway/kernelmanager"
	"github.com/d2g/cloudpathway/networktools"
	"github.com/d2g/cloudpathway/web"
	"github.com/d2g/dhcp4"
	"github.com/d2g/dhcp4server"
	"github.com/d2g/dhcp4server/leasepool"
	"github.com/d2g/dnsforwarder"
	"io/ioutil"
	"log"
	"net"
	_ "net/http/pprof"
	"time"
)

func main() {
	configuration_file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("Error Loading config.json")
		panic(err)
	}

	configuration := Configuration{}

	err = json.Unmarshal(configuration_file, &configuration)
	if err != nil {
		log.Println("Error In JSON Format of config.json")
		panic(err)
	}

	//Start the DNS Service
	go StartDNSService(&configuration.DNS)

	//Start The DHCP Service
	go StartDHCPService(&configuration)

	go StartConnectionManager(&configuration)

	go StartKernelManager(&configuration.KernelManager)

	//Start The HTTP Service
	StartHTTPService(&configuration.HTTP)

}

func StartHTTPService(configuration *web.Configuration) {
	//Setup The Local HTTP Server
	httpServer, err := web.NewServer(configuration)
	if err != nil {
		log.Println("Unable to Start HTTP Service.. Error:" + err.Error())
	}

	//Launch the HTTP service in a new thread.
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Println("Unable to Start HTTP Service.. Continuing...")
	}
}

func StartDNSService(configuration *dnsforwarder.Configuration) {
	cache, err := dnsimplementation.GetCache()
	if err != nil {
		log.Printf("Error: DNS Cache:%v\n", err)
	}

	hosts := dnsimplementation.Hosts{}
	hosts.Devices = make(map[string]net.IP)

	activeIPs, err := networktools.ActiveIPInterfaces()
	if err == nil {
		for _, activeIP := range activeIPs {
			switch activeIP.(type) {
			case *net.IPAddr:
				hosts.Add("cloudpathway", activeIP.(*net.IPAddr).IP)
				break
			case *net.IPNet:
				hosts.Add("cloudpathway", activeIP.(*net.IPNet).IP)
				break
			}
		}
	} else {
		log.Printf("Warning: Error Getting Active Interfaces:" + err.Error())
	}

	server := dnsforwarder.Server{}
	server.Configuration = configuration
	server.Cache = cache
	server.Hosts = &hosts
	server.Hijacker = dnsimplementation.Hijack

	//Start GC on the DNS Cache.
	go func() {
		for {
			select {
			case <-time.After(24 * time.Hour):
			}

			err := cache.GC()
			if err != nil {
				log.Printf("Error: DNS GC Failed Error: \"%v\"\n", err)
				return
			}
		}
	}()

	go func() {
		err := server.ListenAndServeUDP(net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 53})
		if err != nil {
			log.Printf("UDP Error:%v\n", err)
		}
	}()

	err = server.ListenAndServeTCP(net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 53})
	if err != nil {
		log.Printf("TCP Error:%v\n", err)
	}
}

func StartDHCPService(configuration *Configuration) {
	//This allows us to deactivate DHCP
	if configuration.DHCP.Disabled {
		log.Println("Warning: DHCP is disabled by configuration (config.json)")
		return
	}

	leasePool, err := datastore.GetLeaseHelper()
	if err != nil {
		log.Fatal("Unable to Start Lease Pool")
	}

	//Lets work out the number of leases
	numLeases := dhcp4.IPRange(net.ParseIP(configuration.DHCP.Leases.Start), net.ParseIP(configuration.DHCP.Leases.End))

	//Lets add a list of IPs to the pool these will be served to the clients so make sure they work for you.
	//So Create Array of IPs 192.168.1.1 to 192.168.1.30
	for i := 0; i < numLeases; i++ {
		err = leasePool.AddLease(leasepool.Lease{IP: dhcp4.IPAdd(net.ParseIP(configuration.DHCP.Leases.Start), i)})
		if err != nil {
			log.Fatalln("Error Adding IPs to pool:" + err.Error())
		}
	}

	//If we don't have an IP configured then use the first Active IP interface (warn if there is more than 1)
	if configuration.DHCP.Server.IP == nil || configuration.DHCP.Server.IP.Equal(net.IPv4(0, 0, 0, 0)) {
		//Get IP And Subnet (Windows doesn't seem to give subnet)
		mainActiveIPNetwork, err := networktools.MainActiveIPNetwork()
		if err != nil {
			log.Fatalln("Error Getting Local Interfaces:" + err.Error())
		}

		if mainActiveIPNetwork.IP.Equal(net.IP{}) {
			log.Fatal("No Interface Found With a Valid IP")
		}

		configuration.DHCP.Server.IP = mainActiveIPNetwork.IP

		if !bytes.Equal(mainActiveIPNetwork.Mask, net.IPMask{}) && (configuration.DHCP.Server.SubnetMask == nil || configuration.DHCP.Server.SubnetMask.Equal(net.IPv4(0, 0, 0, 0))) {
			configuration.DHCP.Server.SubnetMask = net.ParseIP(mainActiveIPNetwork.Mask.String())
		}
	}

	//Make Sure we Now have a Subnet Mask
	if configuration.DHCP.Server.SubnetMask == nil || configuration.DHCP.Server.SubnetMask.Equal(net.IPv4(0, 0, 0, 0)) {
		log.Fatalln("DHCP SubnetMask Missing in Configuration and Unable to Calculate...")
	}

	//If For some reason We've setup alternative DNS servers (I Can't think why we should) use them otherwise add ourself.
	if configuration.DHCP.Server.DNSServers == nil || len(configuration.DHCP.Server.DNSServers) == 0 {
		configuration.DHCP.Server.DNSServers = append(configuration.DHCP.Server.DNSServers, configuration.DHCP.Server.IP)
	} else {
		log.Printf("Warning: Using Preconfigured DNS Servers %v\n", configuration.DHCP.Server.DNSServers)
	}

	//If For some reason We've setup alternative Gateway (I Can't think why we should) use them otherwise add ourself.
	if configuration.DHCP.Server.DefaultGateway == nil {
		configuration.DHCP.Server.DefaultGateway = configuration.DHCP.Server.IP
	} else {
		log.Printf("Warning: Using Preconfigured Default Gateway %v\n", configuration.DHCP.Server.DefaultGateway)
	}

	//Create The Server
	myServer := dhcp4server.Server{}

	myServer.Configuration = &configuration.DHCP.Server
	myServer.LeasePool = leasePool

	//Start GC on the Leases.
	go func() {
		for {
			select {
			case <-time.After(configuration.DHCP.Server.LeaseDuration):
			}

			err := myServer.GC()
			if err != nil {
				log.Printf("Error: DHCP GC Failed Error: \"%v\"\n", err)
				return
			}
		}
	}()

	//Start the Server...
	err = myServer.ListenAndServe()
	if err != nil {
		log.Fatalln("Error Starting Server:" + err.Error())
	}
}

func StartConnectionManager(configuration *Configuration) {
	if configuration.ConnectionManager.Disabled {
		log.Println("Warning: Connection Manager Service is disabled by configuration (config.json)")
		return
	}

	connectionManager, err := connectionmanager.NewConnectionManager(&configuration.ConnectionManager)
	if err != nil {
		log.Fatal("Connection Manager Init Error:" + err.Error())
	}

	if !configuration.ConnectionManager.Manager.GCDisabled {
		go func() {
			err := connectionManager.GC()
			if err != nil {
				log.Fatal("Connection Manager GC Error:" + err.Error())
			}
		}()
	} else {
		log.Println("Warning: Connection Manager Garbage Collection is disabled by configuration (config.json)")
	}

	if configuration.ConnectionManager.Classification.Agents > 0 {
		go func() {
			err := connectionManager.Classify()
			if err != nil {
				log.Fatal("Connection Classification Error:" + err.Error())
			}
		}()
	} else {
		log.Println("Warning: Connection Manager Classification is disabled by configuration (config.json)")
	}

	if configuration.ConnectionManager.Manager.Agents > 0 {
		go func() {
			err = connectionManager.Process()
			if err != nil {
				log.Fatal("Connection Manager Process Error:" + err.Error())
			}
		}()
	} else {
		log.Println("Warning: Connection Manager Processing is disabled by configuration (config.json)")
	}

	err = connectionManager.Listen()
	if err != nil {
		log.Fatal("Connection Manager Reader Error:" + err.Error())
	}
}

func StartKernelManager(configuration *kernelmanager.Configuration) {
	if configuration.Disabled {
		log.Println("Warning: Kernel Manager Service is disabled by configuration (config.json)")
		return
	}

	kernelManager, err := kernelmanager.CreateKernelManager(configuration)
	if err != nil {
		log.Fatal("Kernel Manager Init Error:" + err.Error())
	}

	err = kernelManager.Listen()
	if err != nil {
		log.Fatal("Kernel Manager Error:" + err.Error())
	}
}

func StartETLProcessing() {
	// Issue: 8702 Is currently blocking me moving this forward.
	// https://code.google.com/p/go/issues/detail?id=8702

	//dbconnection, err := sql.Open("sqlite3", "./reports.sqlite")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer dbconnection.Close()

	//processing := etl.ClassifiedConnectionProcessors{
	//	&etl.Log{
	//		DB: dbconnection,
	//	},
	//}

	//log.Printf("%v\n", processing)
}
