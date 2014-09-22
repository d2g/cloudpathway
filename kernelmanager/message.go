package kernelmanager

import (
	"errors"
)

type Message []byte

func (t *Message) Type() MessageType {
	return MessageType((*t)[0])
}

func (t *Message) SetType(messageType MessageType) {
	(*t)[0] = byte(messageType)
}

func (t *Message) Records() uint8 {
	return (*t)[1]
}

func (t *Message) Append(data []byte) error {
	if (*t)[1] < 255 {
		//Append the record to the END.
		(*t)[1] += 1
		(*t) = append((*t), data...)
	} else {
		//Error
		return errors.New("Unable to Add Record Message Is Full")
	}
	return nil
}
