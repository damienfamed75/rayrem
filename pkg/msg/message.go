package msg

import (
	"sync"
)

// Message is the interface to qualify as a message to be sent to listeners.
type Message interface {
	Type() string
}

// MessageManager routes dispatched messages to a slice of handlers
// based on the same message type (string index.)
type MessageManager struct {
	sync.RWMutex
	listeners        map[string][]HandlerIDPair
	handlersToRemove map[string][]MessageHandlerID
}

// Dispatch sends a message into the manager to the listeners.
func (m *MessageManager) Dispatch(message Message) {
	// Disallow anything to read.
	m.RLock()
	m.clearRemovedHandlers()

	handlers := make([]MessageHandler, len(m.listeners[message.Type()]))
	pairs := m.listeners[message.Type()]
	for i := range pairs {
		handlers[i] = pairs[i].MessageHandler
	}
	m.RUnlock()
	// Send the message to all the handlers.
	for _, handler := range handlers {
		handler(message)
	}
}

// Listen appends the listener into the message manager with a given ID and sets
// it to listen until it is told to stop.
// If you want the handler to stop, then you may use the StopListen function and
// provide the ID given from this function.
func (m *MessageManager) Listen(msgType string, handler MessageHandler) MessageHandlerID {
	m.Lock()
	defer m.Unlock()
	// If the listeners are nil, then create a new map of them.
	if m.listeners == nil {
		m.listeners = make(map[string][]HandlerIDPair)
	}
	// Create the new handler ID pair.
	handlerID := newHandlerID()
	newHandlerIDPair := HandlerIDPair{
		MessageHandlerID: handlerID,
		MessageHandler:   handler,
	}
	// Append the new handler and ID to the listeners.
	m.listeners[msgType] = append(m.listeners[msgType], newHandlerIDPair)

	return handlerID
}

// ListenOnce only accepts one message before being removed from the listeners.
func (m *MessageManager) ListenOnce(msgType string, handler MessageHandler) {
	var handlerID MessageHandlerID
	handlerID = m.Listen(msgType, func(msg Message) {
		handler(msg)
		m.StopListen(msgType, handlerID)
	})
}

// StopListen marks the current handler to be removed from the listeners.
func (m *MessageManager) StopListen(msgType string, handlerID MessageHandlerID) {
	if m.handlersToRemove == nil {
		m.handlersToRemove = make(map[string][]MessageHandlerID)
	}
	m.handlersToRemove[msgType] = append(m.handlersToRemove[msgType], handlerID)
}

// Loop through marked handlers to remove them one by one in the slice
// of listeners.
func (m *MessageManager) clearRemovedHandlers() {
	for msgType, handlerList := range m.handlersToRemove {
		for _, handlerID := range handlerList {
			m.removeHandler(msgType, handlerID)
		}
	}

	m.handlersToRemove = make(map[string][]MessageHandlerID)
}

// Based on the message type and ID, this function will remove a currently
// listening handler from the slice of listeners.
func (m *MessageManager) removeHandler(msgType string, handlerID MessageHandlerID) {
	idxOfHandler := -1
	// Search for the specified listener.
	for i, activeHandler := range m.listeners[msgType] {
		if activeHandler.MessageHandlerID == handlerID {
			idxOfHandler = i
			break
		}
	}

	// If the listener wasn't found, then return.
	if idxOfHandler == -1 {
		return
	}

	// Remove listener from slice of listeners.
	m.listeners[msgType] = append(m.listeners[msgType][:idxOfHandler], m.listeners[msgType][idxOfHandler+1:]...)
}
