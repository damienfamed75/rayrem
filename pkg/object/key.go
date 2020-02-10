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
	pickedUp bool

	*common.BasicEntity
}

// NewKey returns a key with a default key sprite.
func NewKey(dest r.Vector2) (*Key, error) {
	k := &Key{
		position: dest,
		lock: &Lock{
			mailbox: &msg.MessageManager{},
			msgType: msg.Lock,
			locked:  true,
		},
	}

	// Load the spritesheet file.
	ase, err := common.LoadSpritesheet(common.Config.Objects.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("loading spritesheet: %w", err)
	}
	ase.Play("idle")

	// Create basic entity to draw the key.
	k.BasicEntity, err = common.NewBasicEntity(ase)
	if err != nil {
		return nil, fmt.Errorf("basic entity: %w", err)
	}

	k.zone = physics.NewZone(
		k.position.X, k.position.Y,
		float32(k.Ase.FrameBoundaries().Width), float32(k.Ase.FrameBoundaries().Height),
		k.lock.mailbox, k.lock.msgType,
	)

	return k, nil
}

func (k *Key) Add(w *physics.SpatialHashmap) {
	w.Insert(k.zone)

	k.lock.mailbox.ListenOnce(k.lock.msgType, func(m msg.Message) {
		fmt.Println("picked up")
		// remove zone.
		w.Remove(k.zone)
		k.pickedUp = true
	})
}

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

func (k *Key) Lock() *Lock {
	return k.lock
}
