package etl

import (
	"github.com/d2g/cloudpathway/datastore"
)

type ClassifiedConnectionProcessor interface {
	Process(datastore.ClassifiedConnection) error
}

type ClassifiedConnectionProcessors []ClassifiedConnectionProcessor
