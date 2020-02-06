package object

import (
	"github.com/damienfamed75/rayrem/pkg/common"
	"github.com/damienfamed75/rayrem/pkg/msg"
	"github.com/damienfamed75/rayrem/pkg/physics"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ physics.SpatialAdder = &Openable{}
)

// Openable is any toggleable item that can be interacted with.
type Openable struct {
	*common.BasicEntity

	position     r.Vector2
	zoneCollider r.Rectangle
	mailbox      *msg.MessageManager
	collider     *physics.Rectangle
	zone         *physics.Zone
	lock         *Lock // if there isn't a lock, then o.lock = nil

	msgType     string
	hasCollider bool
	open        bool // if the openable is currently open.
	toggleable  bool // if player is able to re-close the openable.
	locked      bool // if the openable is locked from opening.
}

type Lock struct {
	mailbox *msg.MessageManager
	msgType string
}

func NewOpenable(det r.Rectangle, toggle bool, spritePath string, w *physics.SpatialHashmap, oo ...OpenOption) (*Openable, error) {
	o := &Openable{
		mailbox:      &msg.MessageManager{},
		toggleable:   toggle,
		msgType:      msg.Open,
		zoneCollider: det,
	}
	// Loop through the possible options given.
	for _, opt := range oo {
		opt(o) // apply option to openable.
	}

	// Load spritesheet into the openable.
	// ase, err := common.LoadSpritesheet(spritePath)
	// if err != nil {
	// 	return nil, fmt.Errorf("loading spritesheet: %w", err)
	// }
	// Play the closed animation by default.
	// ase.Play("closed")

	// Create the basic entity to draw.
	// o.BasicEntity, err = common.NewBasicEntity(ase)
	// if err != nil {
	// 	return nil, fmt.Errorf("basic entity: %w", err)
	// }

	o.zone = physics.NewZone(
		o.zoneCollider.X, o.zoneCollider.Y,
		o.zoneCollider.Width, o.zoneCollider.Height,
		o.mailbox, o.msgType,
	)

	// lock option was applied.
	if o.lock != nil {
		// Listen for when the openable should be unlocked.
		o.lock.mailbox.ListenOnce(o.lock.msgType, func(m msg.Message) {
			o.locked = false
		})
	}

	return o, nil
}

// Add is a custom function used to let this object add itself to the spatial
// hashmap.
func (o *Openable) Add(w *physics.SpatialHashmap) {
	if o.hasCollider {
		// o.collider = physics.NewRectangle()

		w.Insert(o.collider)
	}
	// if o.collider != nil {
	// 	w.Insert(o.collider)
	// }

	// haha... o zone...
	w.Insert(o.zone)

	// Setup a mailbox to listen for door messages.
	o.mailbox.Listen(o.msgType, func(m msg.Message) {
		if !o.locked {
			if r.IsKeyPressed(r.KeyE) {
				// Remove the door from the spatial hashmap.
				w.Remove(o.zone)
				w.Remove(o.collider)
				// Set open to true.
				o.open = true
			}
		}
	})
}

func (o *Openable) Draw() {
	srcX, srcY := o.Ase.FrameBoundaries().X, o.Ase.FrameBoundaries().Y
	w, h := o.Ase.FrameBoundaries().Width, o.Ase.FrameBoundaries().Height

	// src represents the cropped rectangle within the spritesheet image.
	src := r.NewRectangle(
		float32(srcX), float32(srcY),
		float32(w), float32(h),
	)

	// Create a destination for the player to be drawn at.
	dest := r.NewRectangle(
		o.position.X, o.position.Y,
		float32(w)*o.Scale, float32(h)*o.Scale,
	)

	// Finally draw the texture.
	r.DrawTexturePro(
		o.Sprite, src, dest, r.NewVector2(0, 0), o.Rotation, o.Color,
	)
}
