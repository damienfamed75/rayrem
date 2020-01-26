package msg

import "sync/atomic"

type MessageHandler func(msg Message)

type MessageHandlerID uint64

type HandlerIDPair struct {
	MessageHandlerID
	MessageHandler
}

var currentHandlerID uint64

func newHandlerID() MessageHandlerID {
	atomic.AddUint64(&currentHandlerID, 1)
	return MessageHandlerID(atomic.LoadUint64(&currentHandlerID))
}
