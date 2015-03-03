package etl

import (
	"encoding/csv"
	"log"
	"strconv"

	"github.com/d2g/cloudpathway/datastore"
	"github.com/d2g/rotatingfile"
)

type Log struct {
	File *rotatingfile.File
}

func (t *Log) Process(connection datastore.ClassifiedConnection) (err error) {

	w := csv.NewWriter(t.File)
	log.Printf("Trace: Logging Connection:%v\n", connection)
	w.Write(
		[]string{
			connection.Updated().Format("2006-01-02 15:04:05"),
			connection.SourceIP().To4().String() + ":" + strconv.FormatUint(uint64(connection.SourcePort()), 10),
			connection.DestinationIP().To4().String() + ":" + strconv.FormatUint(uint64(connection.DestinationPort()), 10),
			connection.Username(),
			connection.DeviceID().String(),
			connection.Protocol(),
		},
	)
	w.Flush()

	err = w.Error()
	return
}
