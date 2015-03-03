package dnsimplementation

import (
	"log"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/d2g/unqlitego"
	"github.com/miekg/dns"
	"gopkg.in/mgo.v2/bson"
)

type cache struct {
	collection   *unqlitego.Database
	ipToHostname *unqlitego.Database
}

type CacheRecord struct {
	Expiry time.Time
	Record []byte
}

var cacheSingleton *cache = nil

func GetCache() (*cache, error) {
	if cacheSingleton == nil {

		var err error = nil
		cacheSingleton = new(cache)

		cacheSingleton.collection, err = unqlitego.NewDatabase("userdata/DNSCache.unqlite")
		if err != nil {
			return nil, err
		}

		//Commit After 100 Changes
		cacheSingleton.collection.CommitAfter = 100

		cacheSingleton.collection.SetMarshal(bson.Marshal)
		cacheSingleton.collection.SetUnmarshal(bson.Unmarshal)

		cacheSingleton.ipToHostname, err = unqlitego.NewDatabase("userdata/IPtoHostname.unqlite")
		//Commit After 100 Changes
		cacheSingleton.ipToHostname.CommitAfter = 100

		return cacheSingleton, err
	}

	return cacheSingleton, nil
}

func (this *cache) Add(message *dns.Msg) error {
	if message.Question[0].Qtype == dns.TypeA && message.Question[0].Qclass == dns.ClassINET {
		byteMessage, err := message.Pack()
		if err != nil {
			return err
		}

		if len(message.Answer) > 0 {
			err = this.collection.SetObject(message.Question[0].Name, CacheRecord{Expiry: time.Now().Add(time.Duration(message.Answer[0].Header().Ttl) * time.Second), Record: byteMessage})
			if err != nil {
				return err
			}

			for _, part := range message.Answer {
				switch part.(type) {
				case *dns.A:
					cacheErr := this.ipToHostname.SetObject(part.(*dns.A).A.String(), strings.TrimSuffix(message.Question[0].Name, "."))
					if cacheErr != nil {
						log.Println("Warning: Error Adding/Updating Cache IP:" + part.(*dns.A).A.String() + " Hostname:" + strings.TrimSuffix(message.Question[0].Name, ".") + " Error:" + cacheErr.Error())
					}
				case *dns.CNAME:
					//CNAME Don't contain the IP for reverse lookups.
					//log.Printf("CNAME: Type:%s Value:%v", reflect.TypeOf(part).Name(), part)
					//TODO: We should probably Add the CNAME as a hostname??
				default:
					log.Printf("Debug: DNS Message Answer Type:%s Value:%v", reflect.TypeOf(part).Name(), part)
				}
			}
		}

		return err
	}
	return nil
}

func (this *cache) Get(message *dns.Msg) (bool, *dns.Msg, error) {
	cachedMessage := CacheRecord{}
	err := this.collection.GetObject(message.Question[0].Name, &cachedMessage)
	if err != nil {
		return false, nil, err
	}

	cachedMessageMsg := dns.Msg{}
	err = cachedMessageMsg.Unpack(cachedMessage.Record)
	if err != nil {
		return false, nil, err
	}

	if cachedMessageMsg.Id != 0 && cachedMessage.Expiry.After(time.Now()) {
		return true, &cachedMessageMsg, nil
	} else {
		return false, &dns.Msg{}, nil
	}
}

func (t *cache) GetHostname(ip net.IP) (string, error) {
	var hostname string
	err := t.ipToHostname.GetObject(ip.To4().String(), &hostname)
	return hostname, err
}

func (t *cache) GC() error {

	cursor, err := t.collection.NewCursor()
	defer cursor.Close()

	if err != nil {
		return err
	}

	err = cursor.First()
	if err != nil {
		if err == unqlitego.UnQLiteError(-28) {
			return nil
		} else {
			return err
		}
	}

	for {
		if !cursor.IsValid() {
			break
		}

		key, err := cursor.Key()
		if err != nil {
			log.Println("Error: Cursor Get Key Error:" + err.Error())
			continue
		}

		value, err := cursor.Value()
		if err != nil {
			log.Println("Error: Cursor Get Value Error:" + err.Error())
			continue
		}

		cacherecord := CacheRecord{}
		err = t.collection.Unmarshal()(value, &cacherecord)
		if err != nil {
			log.Println("Error: Unmarshalling \"" + string(key) + "\":" + err.Error())
		}

		if cacherecord.Expiry.Before(time.Now()) {
			//Record Has Expired
			err = cursor.Delete()
			if err != nil {
				log.Println("Error: Unable to Delete Cache Record \"" + string(key) + "\":" + err.Error())
			}
		}

		err = cursor.Next()
		if err != nil {
			break
		}
	}

	err = cursor.Close()

	return err
}
