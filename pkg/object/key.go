package object

import (
	"fmt"

	"github.com/damienfamed75/rayrem/pkg/common"
	"github.com/damienfamed75/rayrem/pkg/msg"
	"github.com/damienfamed75/rayrem/pkg/physics"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ physics.SpatialAdder = &Key{}
)

// Key is a simple game object that contains a zone and a Lock structure
// to be attached to any other objects that may take Locks as options such as
// a door structure.
type Key struct {
	position r.Vector2
	zone     *physics.Zone
	lock     *Lock
	msgType  string
	pickedUp bool

	*common.BasicEntity
}

// NewKey returns a key with a default key sprite.
// dest: position where the key should spawn.
func NewKey(dest r.Vector2) (*Key, error) {
	k := &Key{
		position: dest,
		msgType:  msg.Key,
		lock: &Lock{
			mailbox: &msg.MessageManager{},
			msgType: msg.Lock,
			locked:  true, // Set the lock's state default to true.
		},
	}

	// Load the spritesheet file.
	ase, err := common.LoadSpritesheet(common.Config.Objects.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("loading spritesheet: %w", err)
	}
	ase.Play("idle") // Play default animation.

	// Create basic entity to draw the key.
	k.BasicEntity, err = common.NewBasicEntity(ase)
	if err != nil {
		return nil, fmt.Errorf("basic entity: %w", err)
	}
	// Create the zone that surrounds the sprite.
	// The zone is used to tell if an entity is colliding with the key's sprite.
	k.zone = physics.NewZone(
		k.position.X, k.position.Y,
		float32(k.Ase.FrameBoundaries().Width), float32(k.Ase.FrameBoundaries().Height),
		k.lock.mailbox, k.msgType,
	)

	return k, nil
}

// Lock returns the key's lock.
// Note: This can be used to give to Lockable objects.
func (k *Key) Lock() *Lock {
	return k.lock
}

// Add is a custom adding function to add the key's zone to the spatial hash
// and then create a handler with its mailbox.
func (k *Key) Add(w *physics.SpatialHashmap) {
	w.Insert(k.zone)

	// Create a listener in the key's mailbox to listen for collisions.
	k.lock.mailbox.ListenOnce(k.msgType, func(m msg.Message) {
		// Cast the message as a zone message.
		if zm, ok := m.(*physics.ZoneMessage); ok {
			// If the entity that's colliding with the zone is a player.
			if zm.Entity.HasTags(common.TagPlayer) {
				// remove zone.
				w.Remove(k.zone)
				k.pickedUp = true

				// Send a message to the lock that it's now unlocked.
				k.lock.mailbox.Dispatch(
					msg.NewGenericMsg(k.lock.msgType, nil),
				)
			}
		}
	})
}

// Draw is used to draw the key's sprite.
func (k *Key) Draw() {
	if k.pickedUp == false {
		srcX, srcY := k.Ase.FrameBoundaries().X, k.Ase.FrameBoundaries().Y
		w, h := k.Ase.FrameBoundaries().Width, k.Ase.FrameBoundaries().Height

		// src represents the cropped rectangle within the spritesheet image.
		src := r.NewRectangle(
			float32(srcX), float32(srcY),
			float32(w), float32(h),
		)

		// Create a destination for the player to be drawn at.
		dest := r.NewRectangle(
			k.position.X, k.position.Y,
			float32(w)*k.Scale, float32(h)*k.Scale,
		)

		// Finally draw the texture.
		r.DrawTexturePro(
			k.Sprite, src, dest, r.NewVector2(0, 0), k.Rotation, k.Color,
		)
	}
}
