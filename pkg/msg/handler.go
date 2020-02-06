package msg

import "sync/atomic"

// MessageHandler is a function to take a message and use it however.
type MessageHandler func(msg Message)

// MessageHandlerID is the unique ID associated with a handler.
type MessageHandlerID uint64

// HandlerIDPair is a simple structure to pair an ID to a handler.
type HandlerIDPair struct {
	MessageHandlerID
	MessageHandler
}

// currentHandlerID is the tracker for what ID is being assigned to handlers.
// This iterates over the course of the game depending on how many handlers
// are being added, and should not be manipulated outside of this package.
var currentHandlerID uint64

func newHandlerID() MessageHandlerID {
	// Iterates the global handler ID by 1.
	atomic.AddUint64(&currentHandlerID, 1)
	return MessageHandlerID(atomic.LoadUint64(&currentHandlerID))
}
