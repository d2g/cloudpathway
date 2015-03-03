package etl

import (
	"log"

	"github.com/d2g/cloudpathway/datastore"
)

type ClassifiedConnectionProcessor interface {
	Process(datastore.ClassifiedConnection) error
}

type ClassifiedConnectionProcessors []ClassifiedConnectionProcessor

func (t *ClassifiedConnectionProcessors) Process(input chan datastore.ClassifiedConnection) error {

	for {
		connection := <-input

		etlprocessors := []ClassifiedConnectionProcessor(*t)

		for i := range etlprocessors {
			err := etlprocessors[i].Process(connection)
			if err != nil {
				log.Printf("Error: ETL Processing %s\n", err)
			}
		}

	}

}
