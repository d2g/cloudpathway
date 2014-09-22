package datastore

import (
	"github.com/d2g/packetclassification"
	"log"
)

type ClassifiedConnection struct {
	Connection

	protocol        string
	classifications []packetclassification.Classification
}

func NewClassifiedConnection(c Connection) (ClassifiedConnection, error) {
	cc := ClassifiedConnection{
		Connection: c,
	}

	cc.classifications = make([]packetclassification.Classification, len(cc.Packets()))
	ct := make(map[string]int)

	for i := range cc.Packets() {
		classified, classification, err := packetclassification.Classify(cc.Packets()[i])
		if err != nil {
			log.Printf("Error Classifying Packet: %v\n", err.Error())
			cc.classifications[i] = packetclassification.Classification{
				Protocol: "Error",
			}
		} else {
			if classified {
				ct[classification.Protocol] += 1
				cc.classifications[i] = classification
			} else {
				cc.classifications[i] = packetclassification.Classification{
					Protocol: "Unknown",
				}
			}
		}
	}

	maxClassifications := 0

	for i := range ct {
		if ct[i] > maxClassifications {
			cc.protocol = i
			maxClassifications = ct[i]
		}
	}

	return cc, nil
}

func (t *ClassifiedConnection) Protocol() string {
	return t.protocol
}

func (t *ClassifiedConnection) Classisications() []packetclassification.Classification {
	return t.classifications
}
