package kernelmanager

type MessageType uint8

const (
	CLEAR   MessageType = 0
	CLEARIP MessageType = 1

	BLOCK   MessageType = 10
	BLOCKIP MessageType = 11

	UNBLOCK   MessageType = 20
	UNBLOCKIP MessageType = 21
)
