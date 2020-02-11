package physics

import (
	"github.com/damienfamed75/rayrem/pkg/msg"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Zone is a rectangle that sends a signal to the mailbox given with the
// specified message type.
type Zone struct {
	mailbox *msg.MessageManager
	msgType string
	*Rectangle
}

// NewZone returns a new zone rectangle that sends messages through a mailbox.
func NewZone(x, y, w, h float32, mailbox *msg.MessageManager, msgType string) *Zone {
	z := &Zone{
		mailbox:   mailbox,
		msgType:   msgType,
		Rectangle: NewRectangle(x, y, w, h),
	}

	return z
}

func (z *Zone) dispatchMessage(entity *Space, overlap r.Rectangle) {
	msg := &ZoneMessage{
		Entity:  entity,
		Overlap: overlap,
		msgType: z.msgType,
	}

	z.mailbox.Dispatch(msg)
}

// ZoneMessage stores the information through the message in the zone's mailbox.
type ZoneMessage struct {
	Entity  *Space
	Overlap r.Rectangle

	msgType string
}

// Type returns the message type.
func (z *ZoneMessage) Type() string {
	return z.msgType
}
