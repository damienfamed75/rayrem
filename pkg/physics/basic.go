package physics

import (
	"fmt"
	_ "image/png" // allow decoding of PNG files.

	"github.com/damienfamed75/rayrem/pkg/common"

	"github.com/damienfamed75/aseprite"
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	// BasicEntity fills the common.Entity.
	_ common.Entity = &BasicEntity{}
)

// BasicEntity has fundamental parts of a moveable entity such as
// a sprite, facing directions, changeable color, a rigidbody, and a space.
type BasicEntity struct {
	Ase      *aseprite.File
	Sprite   r.Texture2D
	Color    r.Color
	Facing   common.Direction
	Scale    float32
	Rotation float32

	Rigidbody *Body

	Space *Space
}

// NewBasicEntity returns a basic entity that loads in the sprite
// based on the given spritesheet. Also creates and adds the rigidbody.
func NewBasicEntity(collision *Space, solids *SpatialHashmap, maxSpeed r.Vector2, ase *aseprite.File) (*BasicEntity, error) {
	b := &BasicEntity{
		Ase:       ase,
		Color:     r.White,
		Facing:    common.Right,
		Scale:     1.0,
		Rigidbody: NewBody(collision, solids, maxSpeed),
		Space:     NewSpace(),
	}

	// Load the spritesheet image from package.
	sprite, err := common.LoadPNG(ase.Meta.Image)
	if err != nil {
		return nil, fmt.Errorf("loading spritesheet image: %w", err)
	}

	// Load the image into raylib using a Go image.Image.
	img := r.LoadImageFromGo(sprite)

	// Load a texture from a raylib image.
	b.Sprite = r.LoadTextureFromImage(img)

	// Add the rigidbody to the basic entity space.
	b.Space.Add(b.Rigidbody)

	return b, nil
}

// TakeDamage doesn't do anything by default.
func (b *BasicEntity) TakeDamage() {}

// Position returns the position of the rigidbody's collision space.
func (b *BasicEntity) Position() r.Vector2 {
	return b.Rigidbody.Position()
}

// Update is the barebones just update the spritesheet state and rigidbody.
func (b *BasicEntity) Update(dt float32) {
	b.Ase.Update(dt)
	b.Rigidbody.Update(dt)
}

// Draw is used by default if the parent struct doesn't overwrite it.
// it just draws the sprite in the spritesheet, rotates, changes direction,
// scales, and colors as according to the BasicEntity field values.
func (b *BasicEntity) Draw() {
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
