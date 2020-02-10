package object

import (
	"github.com/damienfamed75/rayrem/pkg/msg"
	"github.com/damienfamed75/rayrem/pkg/physics"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ physics.SpatialAdder = &interactable{}
	_ Lockable             = &interactable{}
)

type interactable struct {
	mailbox *msg.MessageManager

	zone     *physics.Zone
	collider *physics.Rectangle

	msgType string
	// TODO toggle
	// toggleable bool

	Lock
}

func newInteractable(msgType string, zone r.Rectangle, oo ...Option) *interactable {
	i := &interactable{
		mailbox: &msg.MessageManager{},
		msgType: msgType,
	}

	i.zone = physics.NewZone(
		zone.X, zone.Y, zone.Width, zone.Height, i.mailbox, i.msgType,
	)

	// Apply given options.
	for _, o := range oo {
		o(i)
	}

	return i
}

func (i *interactable) Add(w *physics.SpatialHashmap) {
	if i.collider != nil {
		w.Insert(i.collider)
	}

	w.Insert(i.zone)
}

func (i *interactable) Draw() {
	if i.collider != nil {
		r.DrawRectangleLinesEx(i.collider.Rectangle, 1, r.White)
	}

	r.DrawRectangleLinesEx(i.zone.Rectangle.Rectangle, 1, r.Gray)
}
