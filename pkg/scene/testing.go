package scene

import (
	"github.com/damienfamed75/rayrem/pkg/camera"
	"github.com/damienfamed75/rayrem/pkg/common"
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
	ground       *physics.Space
	space        *physics.SpatialHashmap
	camera       *camera.FollowCamera

	world *physics.Space
}

// NewTestingScene returns a new scene. Not much to say here because this will
// often change.
func NewTestingScene(sceneManager common.SceneManager, player *player.Player, world *physics.Space, space *physics.SpatialHashmap) common.Scene {
	t := &Testing{
		sceneManager: sceneManager,
		player:       player,
		ground:       physics.NewSpace(),
		world:        world,
		space:        space,
	}

	// Set the player's position to its spawn position.
	t.player.SetPosition(100, 100)

	// Add ground elements.
	// t.ground.Add(
	t.space.InsertMulti(
		physics.NewRectangle(0, 200, 50, 50),
		physics.NewRectangle(50, 200, 50, 50),
		physics.NewRectangle(100, 200, 50, 50),
		physics.NewRectangle(150, 200, 50, 50),
		physics.NewRectangle(200, 200, 50, 50),
		physics.NewRectangle(250, 200, 50, 50),
		physics.NewRectangle(300, 200, 50, 50),
		physics.NewRectangle(350, 200, 50, 50),
		physics.NewRectangle(400, 200, 200, 200),
		physics.NewRectangle(200, 130, 50, 40),
		physics.NewRectangle(375, 180, 100, 20),
	)

	t.ground.Add(
		physics.NewRectangle(0, 200, 100, 200),
		physics.NewRectangle(100, 200, 100, 200),
		physics.NewRectangle(200, 200, 100, 200),
		physics.NewRectangle(300, 200, 100, 200),
		physics.NewRectangle(400, 200, 200, 200),
		physics.NewRectangle(200, 130, 50, 40),
		physics.NewRectangle(375, 180, 100, 20),
	)

	// Tag the ground elements as so.
	t.ground.AddTags(common.TagGround)

	// Create the scene camera.
	t.camera = camera.NewFollow(t.player.Rigidbody.Space)

	// Add the ground and player to the world space.
	world.Add(
		t.ground, t.player.Space,
	)

	// t.space.Insert(t.player.Rigidbody)

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

	// for _, g := range *t.ground {
	// 	switch t := g.(type) {
	// 	case *physics.Rectangle:
	// 		r.DrawRectangleLinesEx(t.Rectangle, 1, r.Orange)
	// 	case interface{ Draw(float32, r.Color) }:
	// 		t.Draw(1.0, r.Green)
	// 	}
	// }
	dest := r.NewRectangle(
		t.player.Rigidbody.Position().X, t.player.Rigidbody.Position().Y,
		float32(t.player.Ase.FrameBoundaries().Width), float32(t.player.Ase.FrameBoundaries().Height),
	)

	potential := t.space.Retrieve(dest)
	for _, p := range potential {
		r.DrawRectangleLinesEx(r.NewRectangle(
			p.Position().X, p.Position().Y,
			p.MaxPosition().X-p.Position().X, p.MaxPosition().Y-p.Position().Y,
		), 1, r.Red)
	}

	r.DrawRectangleLinesEx(dest, 1, r.Green)

	t.player.Draw()

	r.EndMode2D()
}

// Unload doesn't do much right now.
func (t *Testing) Unload() {
}
