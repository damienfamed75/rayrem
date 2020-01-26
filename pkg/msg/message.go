package msg

import (
	"sync"
)

type Message interface {
	Type() string
}

type MessageManager struct {
	sync.RWMutex
	listeners        map[string][]HandlerIDPair
	handlersToRemove map[string][]MessageHandlerID
}

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

func (m *MessageManager) ListenOnce(msgType string, handler MessageHandler) {
	var handlerID MessageHandlerID
	handlerID = m.Listen(msgType, func(msg Message) {
		handler(msg)
		m.StopListen(msgType, handlerID)
	})
}

func (m *MessageManager) StopListen(msgType string, handlerID MessageHandlerID) {
	if m.handlersToRemove == nil {
		m.handlersToRemove = make(map[string][]MessageHandlerID)
	}
	m.handlersToRemove[msgType] = append(m.handlersToRemove[msgType], handlerID)
}

func (m *MessageManager) clearRemovedHandlers() {
	for msgType, handlerList := range m.handlersToRemove {
		for _, handlerID := range handlerList {
			m.removeHandler(msgType, handlerID)
		}
	}

	m.handlersToRemove = make(map[string][]MessageHandlerID)
}

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
