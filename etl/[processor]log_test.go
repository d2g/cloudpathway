// DO NOT USE: Due to Issue: 8702
package etl

//import (
//	"database/sql"
//	_ "github.com/mattn/go-sqlite3"

//"github.com/d2g/cloudpathway/datastore"
//)

//type Log struct {
//	DB       *sql.DB
//	statment *sql.Stmt
//}

//func (t *Log) Process(connection datastore.ClassifiedConnection) (err error) {

//	if t.statment == nil {
//		t.statment, err = t.DB.Prepare(
//			`INSERT INTO tbl_log(
//							log_source_ip,
//							log_source_port,
//							log_destination_ip,
//							log_destination_port,
//							log_packets,
//							log_updated,
//							log_device_id,
//							log_username,
//							log_protocol
//							) VALUES(?,?,?,?,?,?,?,?,?)`,
//		)
//		if err != nil {
//			return
//		}
//	}

//	_, err = t.statment.Exec(connection.SourceIP().To4().String(),
//		int(connection.SourcePort()),
//		connection.DestinationIP().To4().String(),
//		int(connection.DestinationPort()),
//		connection.Packets(), //array of array of bytes?
//		connection.Updated(),
//		connection.DeviceID().String(),
//		connection.Username(),
//		connection.Protocol(),
//	)

//	return
//}
