package common

import (
	"fmt"

	"github.com/damienfamed75/aseprite"
	r "github.com/lachee/raylib-goplus/raylib"
)

// BasicEntity is a non-rigidbody effected sprite that can be moved around
// and manipulated by just adding a custom Draw function.
type BasicEntity struct {
	Ase      *aseprite.File // Spritesheet
	Sprite   r.Texture2D    // Raylib sprite
	Color    r.Color        // Default: White
	Rotation float32        // Default: 0
	Scale    float32        // Set in config
}

// NewBasicEntity creates a very basic drawable sprite sheet. The animation
// state has to be provided that will be default when the sheet will be loaded.
func NewBasicEntity(ase *aseprite.File) (*BasicEntity, error) {
	b := &BasicEntity{
		Ase:   ase,
		Color: r.White,
		Scale: Config.Game.EntityScale,
	}

	// Load the spritesheet image from package.
	sprite, err := LoadPNG(ase.Meta.Image)
	if err != nil {
		return nil, fmt.Errorf("loading spritesheet image: %w", err)
	}

	// Load the image into raylib using a Go image.Image.
	img := r.LoadImageFromGo(sprite)
	// Load a texture from a raylib image.
	b.Sprite = r.LoadTextureFromImage(img)

	return b, nil
}

// Update lets the spritesheet update.
func (b *BasicEntity) Update(dt float32) {
	b.Ase.Update(dt)
}
