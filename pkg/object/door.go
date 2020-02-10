package object

import (
	"github.com/damienfamed75/rayrem/pkg/common"
	"github.com/damienfamed75/rayrem/pkg/msg"
	"github.com/damienfamed75/rayrem/pkg/physics"

	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	// The door adds itself to the spatial hashmap since there are multiple
	// objects to be added.
	_ physics.SpatialAdder = &Door{}
	// The door is a lockable object and can be configured with WithLock(*Lock)
	_ Lockable = &Door{}
)

// Door is a collection of a solid collider and a zone area where the player
// can open the door from. If the door is opened, then the collider and zone are
// removed from the given spatial hashmap.
type Door struct {
	spriteOpen   r.Texture2D
	spriteClosed r.Texture2D

	open bool
	// openDir is the direction that the door was opened from.
	openDir       common.Direction
	aoaMultiplier float32 // Area of Activation Multiplier

	// This door can be interacted with by the player.
	*interactable
}

// NewDoor returns a door with a fixed zone range.
// The width and height affect the door's collider not the zone.
// func NewDoor(x, y, w, h float32, world *physics.SpatialHashmap, locked bool) *Door {
func NewDoor(pos r.Vector2, oo ...Option) *Door {
	d := &Door{
		openDir:       common.Right,
		aoaMultiplier: 3,
	}

	open, _ := common.LoadPNG("door-open.png")
	closed, _ := common.LoadPNG("door-closed.png")

	d.spriteOpen = r.LoadTextureFromGo(open)
	d.spriteClosed = r.LoadTextureFromGo(closed)

	d.interactable = newInteractable(
		msg.Door,
		// Rectangle for the area of activation zone.
		r.NewRectangle(
			pos.X-(float32(d.spriteClosed.Width)), pos.Y,
			float32(d.spriteClosed.Width)*d.aoaMultiplier, float32(d.spriteClosed.Height),
		),
		// Add a collider for the solid bit of the door.
		withCollider(
			r.NewRectangle(
				pos.X, pos.Y,
				float32(d.spriteClosed.Width), float32(d.spriteClosed.Height),
			),
		),
	)

	// Apply given options.
	for _, o := range oo {
		o(d)
	}

	return d
}

// Add fills the spatial adder interface to be able to custom add itself to
// the world.
func (d *Door) Add(w *physics.SpatialHashmap) {
	// Insert the parts of the door into the world.
	d.interactable.Add(w)

	// Setup a mailbox to listen for door messages.
	d.interactable.mailbox.Listen(d.interactable.msgType, func(m msg.Message) {
		if !d.Lock.locked && r.IsKeyPressed(r.KeyE) {
			zm := m.(*physics.ZoneMessage)
			// If the colliding zone isn't the player, then ignore this message.
			if !zm.Entity.HasTags(common.TagPlayer) {
				return
			}

			d.openDir = common.Right
			// if the overlapping rectangle's max X position is less than
			// the center of the aoa zone then open the door left.
			if zm.Overlap.MaxPosition().X < d.zone.Rectangle.Rectangle.Center().X {
				d.openDir = common.Left
			}

			// Remove the door from the spatial hashmap.
			w.Remove(d.zone)
			w.Remove(d.collider)
			// Set open to true.
			d.open = true
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
		r.DrawTexture(
			d.spriteClosed,
			int(d.collider.Rectangle.X), int(d.collider.Rectangle.Y),
			r.White,
		)
	} else {
		r.DrawTexturePro(
			d.spriteOpen,
			// source rectangle.
			r.NewRectangle(
				0, 0,
				float32(d.spriteOpen.Width*int32(d.openDir)), float32(d.spriteOpen.Height),
			),
			// destination rectangle.
			r.NewRectangle(
				d.collider.Rectangle.X, d.collider.Rectangle.Y,
				float32(d.spriteOpen.Width), float32(d.spriteOpen.Height),
			),
			r.NewVector2(
				// when the door is opened left then the origin X should equal
				// the door's width, and when right then the origin should be zero.
				// -1 >> 0x1 = 1
				// 1 >> 0x1 = 0
				-float32((int(-d.openDir)>>0x1)*int(d.spriteOpen.Width)), 0,
			),
			0, r.White,
		)
	}

	// zone.
	r.DrawRectangleLinesEx(d.zone.Rectangle.Rectangle, 1, r.SkyBlue.Lerp(r.Transparent, 0.5))
}
