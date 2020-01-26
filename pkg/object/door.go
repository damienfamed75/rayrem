package object

import (
	"github.com/damienfamed75/rayrem/pkg/msg"
	"github.com/damienfamed75/rayrem/pkg/physics"

	r "github.com/lachee/raylib-goplus/raylib"
)

// Door is a collection of a solid collider and a zone area where the player
// can open the door from. If the door is opened, then the collider and zone are
// removed from the given spatial hashmap.
type Door struct {
	mailbox *msg.MessageManager
	msgType string
	locked  bool
}

// NewDoor returns a door with a fixed zone range.
// The width and height affect the door's collider not the zone.
func NewDoor(x, y, w, h float32, world *physics.SpatialHashmap, locked bool) *Door {
	d := &Door{
		mailbox: &msg.MessageManager{},
		locked:  locked,
		msgType: msg.Door,
	}

	// Create the rectangles and zones for the door.
	collider := physics.NewRectangle(x, y, w, h)
	zone := physics.NewZone(x-w, y, w*3, h, d.mailbox, d.msgType)

	// Insert the parts of the door into the world.
	world.Insert(collider)
	world.Insert(zone)

	// Setup a mailbox to listen for door messages.
	d.mailbox.Listen(d.msgType, func(m msg.Message) {
		if !d.locked {
			if r.IsKeyPressed(r.KeyE) {
				// Remove the door from the spatial hashmap.
				world.Remove(zone)
				world.Remove(collider)
			}
		}
	})

	// Typically you don't need the door object.
	return d
}
