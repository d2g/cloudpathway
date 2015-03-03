package web

import (
	"github.com/d2g/cloudpathway/networktools"
	template "github.com/d2g/goti/html"
	"log"
	"net"
	"net/http"
	"strings"
)

type NotFound struct {
	NotFoundHandler http.Handler
}

func (t NotFound) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	//Parse the Host into an IP
	host, _, err := net.SplitHostPort(request.Host)
	if err != nil {
		host = request.Host
	}

	if isHostLocal(host) {
		t.NotFoundHandler.ServeHTTP(response, request)
	} else {
		log.Printf("Trace: Handling 511 %v\n", host)
		response.WriteHeader(511)
		tpl, err := template.ParseFiles("static/authentication.tpl")
		if err != nil {
			log.Printf("Error: %s\n", err)
			return
		}
		tpl.Execute(response, nil)
	}

	return
}

func isHostLocal(host string) bool {
	if host == "127.0.0.1" || host == "::1" || host == "localhost" {
		return true
	}

	requestIP := net.ParseIP(host)

	if requestIP != nil {

		//Check We're not using a local IP etc.
		localAddresses, err := networktools.ActiveIPInterfaces()

		if err != nil {
			log.Printf("Debug: Failed To Get Active Interfaces When looking Up HTTP Request \"%s\"\n", err)
		} else {
			for i := range localAddresses {
				if ip, ok := localAddresses[i].(*net.IPAddr); ok {
					if requestIP.Equal(ip.IP) {
						return true
					}
				}
				if ip, ok := localAddresses[i].(*net.IPNet); ok {
					if requestIP.Equal(ip.IP) {
						return true
					}
				}
			}
		}
	} else {
		log.Printf("Debug: \"%s\" is Not A Valid IP\n", host)
	}

	//If the hostname is cloudpathway
	if strings.ToLower(host) == "cloudpathway.d2g.org.uk" {
		return true
	}
	return false
}
