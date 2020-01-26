package physics

import (
	"github.com/damienfamed75/rayrem/pkg/msg"
	r "github.com/lachee/raylib-goplus/raylib"
)

type Zone struct {
	mailbox *msg.MessageManager
	msgType string
	*Rectangle
}

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

type ZoneMessage struct {
	Entity  *Space
	Overlap r.Rectangle

	msgType string
}

func (z *ZoneMessage) Type() string {
	return z.msgType
}
