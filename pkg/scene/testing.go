package scene

import (
	"github.com/damienfamed75/rayrem/pkg/camera"
	"github.com/damienfamed75/rayrem/pkg/common"
	"github.com/damienfamed75/rayrem/pkg/object"
	"github.com/damienfamed75/rayrem/pkg/physics"
	"github.com/damienfamed75/rayrem/pkg/player"

	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ common.Scene = &Testing{}
)

// Testing is a debug scene to test new features or functionality.
type Testing struct {
	sceneManager common.SceneManager
	player       *player.Player
	solids       *physics.SpatialHashmap
	camera       *camera.FollowCamera

	ground []physics.Transformer
}

// NewTestingScene returns a new scene. Not much to say here because this will
// often change.
func NewTestingScene(sceneManager common.SceneManager, player *player.Player, solids *physics.SpatialHashmap) common.Scene {
	t := &Testing{
		sceneManager: sceneManager,
		player:       player,
		solids:       solids,
	}

	// Set the player's position to its spawn position.
	t.player.SetPosition(100, 100)

	// Somehow add it to a list of things to draw, but not in the ground.
	object.NewDoor(225, 170, 7, 30, solids, false)

	// Add ground elements.
	t.ground = []physics.Transformer{
		// Ground elements.
		physics.NewRectangle(0, 200, 50, 50),
		physics.NewRectangle(50, 200, 50, 50),
		physics.NewRectangle(100, 200, 50, 50),
		physics.NewRectangle(150, 200, 50, 50),
		physics.NewRectangle(200, 200, 50, 50),
		physics.NewRectangle(250, 200, 50, 50),
		physics.NewRectangle(300, 200, 50, 50),
		physics.NewRectangle(350, 200, 50, 50),
		physics.NewRectangle(400, 200, 200, 200),

		// Floating platforms.
		physics.NewRectangle(200, 130, 50, 40),
		physics.NewPlatform(168, 160, 32, 10),

		physics.NewRectangle(375, 180, 100, 20),

		// Slope platform.
		physics.NewSlopePlatform(r.NewVector2(300, 200), r.NewVector2(350, 180), 25),
	}

	// Insert the solids into the world.
	t.solids.Insert(t.ground...)

	// Create the scene camera.
	t.camera = camera.NewFollow(t.player.Rigidbody.Space)

	return t
}

// Update takes delta time and updates objects in the scene.
func (t *Testing) Update(dt float32) {
	t.player.Update(dt)
	t.camera.Update(t.player.Rigidbody.Position())
}

// Draw draws to the screen.
func (t *Testing) Draw() {
	r.BeginMode2D(t.camera.Camera2D)
	r.ClearBackground(r.Black)

	for _, g := range t.ground {
		switch t := g.(type) {
		case *physics.Space:
			for _, obj := range *t {
				r.DrawRectangleLinesEx(
					r.NewRectangle(
						obj.Position().X, obj.Position().Y,
						obj.Width(), obj.Height(),
					),
					1, r.SkyBlue,
				)
			}
		case interface{ Draw(float32, r.Color) }:
			t.Draw(1.0, r.Green)
		case *physics.Rectangle:
			r.DrawRectangleLinesEx(t.Rectangle, 1, r.Orange)
		case physics.Transformer:
			r.DrawRectangleLinesEx(r.NewRectangle(
				t.Position().X, t.Position().Y,
				t.MaxPosition().X-t.Position().X, t.MaxPosition().Y-t.Position().Y,
			), 1, r.Gold)
		}
	}

	t.player.Draw()

	r.EndMode2D()
}

// Unload doesn't do much right now.
func (t *Testing) Unload() {
}
