package main

import (
	"github.com/d2g/cloudpathway/connectionmanager"
	"github.com/d2g/cloudpathway/kernelmanager"
	"github.com/d2g/cloudpathway/web"
	"github.com/d2g/dhcp4server"
	"github.com/d2g/dnsforwarder"
)

type Configuration struct {
	HTTP web.Configuration

	ConnectionManager connectionmanager.Configuration

	KernelManager kernelmanager.Configuration

	DNS  dnsforwarder.Configuration
	DHCP struct {
		Disabled bool
		Server   dhcp4server.Configuration
		Leases   struct {
			Start string
			End   string
		}
	}
}
