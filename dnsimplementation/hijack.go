package dnsimplementation

import (
	"bytes"
	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/cloudpathway/networktools"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

func Hijack(response dns.ResponseWriter, message *dns.Msg) (bool, error) {

	if (message.Question[0].Qtype == dns.TypeA || message.Question[0].Qtype == dns.TypeAAAA) && message.Question[0].Qclass == dns.ClassINET {

		ipString, _, err := net.SplitHostPort(response.RemoteAddr().String())
		if err != nil {
			log.Println("Warning: Received No IP4 Address:" + response.RemoteAddr().String())
			return false, err
		}

		ip := net.ParseIP(ipString)
		deviceHelper, err := datastore.GetDeviceHelper()
		if err != nil {
			log.Println("Error: Datastore Error (Open Access):" + err.Error())
			return false, err
		}

		device, err := deviceHelper.GetDeviceByIP(ip)
		if err != nil {
			log.Println("Error: Datastore Error (Open Access):" + err.Error())
			return false, err
		}

		if bytes.Equal(device.MACAddress, net.HardwareAddr{}) {
			//Device Doesn't exists?
			log.Println("Warning: IP:" + ipString + " MAC:" + device.MACAddress.String() + " Doesn't Exists As A Local Device But is Using Our DNS??")
			return false, nil
		} else {
			if activeUser := device.GetActiveUser(); activeUser != nil {
				//User is active check they have access.
				userFiltercollectionsHelper, err := datastore.GetUserFilterCollectionsHelper()
				if err != nil {
					return false, err
				}

				usersFiltercollection, err := userFiltercollectionsHelper.GetUserFilterCollections(activeUser.Username)
				if err != nil {
					return false, err
				}

				if usersFiltercollection.Username != "" {
					//Get the collections helper.
					filterCollectionsHelper, err := datastore.GetFilterCollectionHelper()
					if err != nil {
						return false, err
					}

					//Start Moving down the domain (i.e. removing subdomains etc)
					collections := make([]datastore.FilterCollection, 0)

					domain := strings.TrimSuffix(message.Question[0].Name, ".")

					for {
						additionalCollections, err := filterCollectionsHelper.GetFilterCollectionsWithDomain(domain)
						if err != nil {
							return false, err
						}

						if len(additionalCollections) > 0 {
							collections = append(collections, additionalCollections...)
						}

						if len(strings.SplitAfterN(domain, ".", 2)) == 2 {
							domain = strings.SplitAfterN(domain, ".", 2)[1]
						} else {
							break
						}
					}

					//Does the domain appear in any collections??
					for _, collection := range collections {
						if usersFiltercollection.ContainsCollection(collection.Name) {
							//The Url Should be blocked...
							//For Now lets just not respond with a DNS record...

							//localResponse := new(dns.Msg)
							//localResponse.SetReply(message)

							//rr_header := dns.RR_Header{Name: message.Question[0].Name, Class: dns.ClassINET, Ttl: 0}

							//Main Active IP
							//mainIPNetwork, err := networktools.MainActiveIPNetwork()
							//if err != nil {
							//	return false, err
							//}

							//switch message.Question[0].Qtype {

							//case dns.TypeA:
							//	rr_header.Rrtype = dns.TypeA
							//	a := &dns.A{rr_header, mainIPNetwork.IP}
							//	localResponse.Answer = append(localResponse.Answer, a)

							//case dns.TypeAAAA:
							//	rr_header.Rrtype = dns.TypeAAAA
							//}

							//err = response.WriteMsg(localResponse)
							//if err != nil {
							//	return false, err
							//}
							return true, nil
						}
					}

					return false, nil
				}
			} else {
				//User is not active redirect to login page.
				localResponse := new(dns.Msg)
				localResponse.SetReply(message)
				rr_header := dns.RR_Header{Name: message.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}
				//Main Active IP
				mainIPNetwork, err := networktools.MainActiveIPNetwork()
				if err != nil {
					return false, err
				}
				a := &dns.A{rr_header, mainIPNetwork.IP}
				localResponse.Answer = append(localResponse.Answer, a)
				err = response.WriteMsg(localResponse)
				if err != nil {
					return false, err
				}
				return true, nil
			}
		}
	}

	return false, nil
}
