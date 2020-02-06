package object

import (
	"github.com/damienfamed75/rayrem/pkg/msg"
	"github.com/damienfamed75/rayrem/pkg/physics"

	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ physics.SpatialAdder = &Door{}
)

// Door is a collection of a solid collider and a zone area where the player
// can open the door from. If the door is opened, then the collider and zone are
// removed from the given spatial hashmap.
type Door struct {
	mailbox  *msg.MessageManager
	collider *physics.Rectangle
	zone     *physics.Zone

	msgType string
	open    bool
	locked  bool
}

// NewDoor returns a door with a fixed zone range.
// The width and height affect the door's collider not the zone.
// func NewDoor(x, y, w, h float32, world *physics.SpatialHashmap, locked bool) *Door {
func NewDoor(col, det r.Rectangle, locked bool) *Door {
	d := &Door{
		mailbox: &msg.MessageManager{},
		locked:  locked,
		msgType: msg.Door,
	}

	// Create the rectangles and zones for the door.
	d.collider = physics.NewRectangle(col.X, col.Y, col.Width, col.Height)
	d.zone = physics.NewZone(det.X, det.Y, det.Width, det.Height, d.mailbox, d.msgType)

	return d
}

// Add fills the spatial adder interface to be able to custom add itself to
// the world.
func (d *Door) Add(w *physics.SpatialHashmap) {
	// Insert the parts of the door into the world.
	w.Insert(d.collider, d.zone)

	// Setup a mailbox to listen for door messages.
	d.mailbox.Listen(d.msgType, func(m msg.Message) {
		if !d.locked {
			if r.IsKeyPressed(r.KeyE) {
				// Remove the door from the spatial hashmap.
				w.Remove(d.zone)
				w.Remove(d.collider)
				// Set open to true.
				d.open = true
			}
		}
	})
}

// Open returns a boolean of if the door is open or closed.
func (d *Door) Open() bool {
	return d.open
}

// Draw is meant for debugging if nothing else is passed.
func (d *Door) Draw() {
	if !d.open {
		r.DrawRectangleLinesEx(d.collider.Rectangle, 1, r.Orange)
	}

	r.DrawRectangleLinesEx(d.zone.Rectangle.Rectangle, 1, r.SkyBlue.Lerp(r.Transparent, 0.5))
}
