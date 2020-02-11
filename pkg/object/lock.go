package object

import (
	"github.com/damienfamed75/rayrem/pkg/msg"
)

// Lock contains enough information to give to any Lockable items
// to have a dynamic state lock.
type Lock struct {
	mailbox *msg.MessageManager
	msgType string
	locked  bool
}

func (l *Lock) applyLock(ll *Lock) {
	*l = *ll

	// Create a mailbox handler to unlock the door.
	l.mailbox.ListenOnce(l.msgType, func(m msg.Message) {
		l.locked = false
	})
}
