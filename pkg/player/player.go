package player

import (
	"fmt"

	"github.com/damienfamed75/rayrem/pkg/common"
	"github.com/damienfamed75/rayrem/pkg/physics"

	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ common.Entity = &Player{}
)

// Player is an entity that the user can control in the game.
type Player struct {
	*physics.BasicEntity

	doubleJumpPerformed bool
	jumpHeight          float32
	friction            float32

	solids *physics.SpatialHashmap
}

// New creates a player at the position specified.
// the world object is updated dynamically because the ground elements aren't
// stored by value or locally at all.
func New(x, y float32, solids *physics.SpatialHashmap) (*Player, error) {
	p := &Player{
		friction:   common.Config.Player.Friction,
		jumpHeight: common.Config.Player.JumpHeight,
		solids:     solids,
	}

	ase, err := common.LoadSpritesheet(common.Config.Player.Spritesheet)
	if err != nil {
		return nil, fmt.Errorf("aseprite: %w", err)
	}

	ase.Play("idle")

	// Create the collision areas of the player.
	collision := physics.NewSpace()
	collision.Add(
		physics.NewRectangle(x, y, float32(ase.FrameBoundaries().Width), float32(ase.FrameBoundaries().Height)),
	)

	// Prepare the player's basic entity.
	p.BasicEntity, err = physics.NewBasicEntity(
		collision, solids,
		r.NewVector2(common.Config.Player.MaxSpeed.X, common.Config.Player.MaxSpeed.Y),
		ase,
	)
	if err != nil {
		return nil, fmt.Errorf("basic entity: %w", err)
	}

	p.Space.AddTags(common.TagPlayer)

	return p, nil
}
