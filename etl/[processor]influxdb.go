package etl

import (
	"log"
	"strconv"

	"github.com/d2g/cloudpathway/datastore"
	"github.com/influxdb/influxdb/client"
)

type InfluxDB struct {
	Client *client.Client
}

func (t *InfluxDB) Process(connection datastore.ClassifiedConnection) (err error) {

	w := client.Write{
		Database:        "d2g",
		RetentionPolicy: "inf",
		Points: []client.Point{
			client.Point{
				Name:      "Log",
				Timestamp: client.Timestamp(connection.Updated()),
				Values: map[string]interface{}{"Source": connection.SourceIP().To4().String() + ":" + strconv.FormatUint(uint64(connection.SourcePort()), 10),
					"Destination": connection.DestinationIP().To4().String() + ":" + strconv.FormatUint(uint64(connection.DestinationPort()), 10),
					"Username":    connection.Username(),
					"DeviceID":    connection.DeviceID().String(),
					"Protocol":    connection.Protocol(),
				},
				Precision: "u",
			},
		},
	}

	results, err := t.Client.Write(w)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("Results: %v\n", results)
	return
}
