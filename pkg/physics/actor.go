package physics

import (
	"fmt"
	_ "image/png" // allow decoding of PNG files.

	"github.com/damienfamed75/rayrem/pkg/common"

	"github.com/damienfamed75/aseprite"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	// Actor fills the common.Entity.
	_ common.Entity = &Actor{}

	_ SpatialAdder = &Actor{}

	_ Moveable = &Actor{}

	_ Entity = &Actor{}
)

// Actor has fundamental parts of a moveable entity such as
// a sprite, facing directions, changeable color, a rigidbody, and a space.
type Actor struct {
	Facing common.Direction

	Rigidbody *Body
	Space     *Space

	*common.BasicEntity
}

// NewActor returns a basic entity that loads in the sprite
// based on the given spritesheet. Also creates and adds the rigidbody.
func NewActor(collision *Space, solids *SpatialHashmap, maxSpeed r.Vector2, ase *aseprite.File) (*Actor, error) {
	b := &Actor{
		Facing:    common.Right,
		Rigidbody: NewBody(collision, solids, maxSpeed),
		Space:     NewSpace(),
	}

	var err error
	b.BasicEntity, err = common.NewBasicEntity(ase)
	if err != nil {
		return nil, fmt.Errorf("basic entity: %w", err)
	}

	// Add the rigidbody to the basic entity space.
	b.Space.Add(b.Rigidbody.Space)

	return b, nil
}

func (b *Actor) ID() uint64 {
	// Since we only really care about the collision space we only return the
	// rigidbody's space ID to check on itself.
	return b.Rigidbody.ID()
}

func (b *Actor) Add(w *SpatialHashmap) {
	w.InsertMoveables(b)
}

func (b *Actor) Velocity() r.Vector2 {
	return b.Rigidbody.velocity
}

func (b *Actor) SetVelocity(x, y float32) {
	b.Rigidbody.velocity.X = x
	b.Rigidbody.velocity.Y = y
}

func (b *Actor) AddVelocity(x, y float32) {
	b.Rigidbody.velocity = b.Rigidbody.velocity.Add(r.NewVector2(x, y))
}

// TakeDamage doesn't do anything by default.
func (b *Actor) TakeDamage() {}

// Position returns the position of the rigidbody's collision space.
func (b *Actor) Position() r.Vector2 {
	return b.Rigidbody.Position()
}

// MaxPosition returns the end position of the rigidbody's collision space.
func (b *Actor) MaxPosition() r.Vector2 {
	return b.Rigidbody.MaxPosition()
}

// Update is the barebones just update the spritesheet state and rigidbody.
func (b *Actor) Update(dt float32) {
	b.Ase.Update(dt)
	b.Rigidbody.Update(dt)
}

// Draw is used by default if the parent struct doesn't overwrite it.
// it just draws the sprite in the spritesheet, rotates, changes direction,
// scales, and colors as according to the Actor field values.
func (b *Actor) Draw() {
	srcX, srcY := b.Ase.FrameBoundaries().X, b.Ase.FrameBoundaries().Y
	w, h := b.Ase.FrameBoundaries().Width, b.Ase.FrameBoundaries().Height

	// src represents the cropped rectangle within the spritesheet image.
	src := r.NewRectangle(
		float32(srcX), float32(srcY),
		float32(int(b.Facing)*w), float32(h),
	)

	// Create a destination for the player to be drawn at.
	dest := r.NewRectangle(
		b.Rigidbody.Position().X, b.Rigidbody.Position().Y,
		float32(w)*b.Scale, float32(h)*b.Scale,
	)

	// Finally draw the texture.
	r.DrawTexturePro(
		b.Sprite, src, dest, r.NewVector2(0, 0), b.Rotation, b.Color,
	)
}
