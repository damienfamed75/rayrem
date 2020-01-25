package game

import (
	"log"

	"github.com/damienfamed75/rayrem/pkg/common"
	"github.com/damienfamed75/rayrem/pkg/physics"
	"github.com/damienfamed75/rayrem/pkg/player"
	"github.com/damienfamed75/rayrem/pkg/scene"

	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ common.SceneManager = &Game{}
)

// Game is the scene manager and holder of the player and world.
type Game struct {
	mode common.Mode

	player *player.Player
	solids *physics.SpatialHashmap

	scenes map[common.Mode]common.Scene
}

// NewGame sets up the game and creates the player and world.
// Also is used to create all the scenes and assign them in the scenes map.
func NewGame() *Game {
	g := &Game{
		solids: physics.NewSpatialHashmap(6),
	}

	// Create the player.
	player, err := player.New(0, 0, g.solids)
	if err != nil {
		log.Fatal(err)
	}

	g.player = player

	// Setup all the scenes in the game.
	g.scenes = map[common.Mode]common.Scene{
		common.ModeTesting:  scene.NewTestingScene(g, g.player, g.solids),
		common.ModeMainMenu: scene.NewMenu(g),
	}

	return g
}

// SetScene changes the scene mode.
func (g *Game) SetScene(mode common.Mode) {
	// g.world.Clear() // Clear the world of any remaining shapes.
	g.mode = mode
}

// Update updates whatever scene is currently set.
func (g *Game) Update(dt float32) {
	g.scenes[g.mode].Update(dt)
}

// Draw draws the current scene that's set.
func (g *Game) Draw() {
	g.scenes[g.mode].Draw()
}

// Unload unloads all assets loaded in by raylib.
func (g *Game) Unload() {
	// Thinking about removing until we have other things that aren't raylib
	// to unload in the scenes.
	// for _, s := range g.scenes {
	// 	s.Unload()
	// }

	// Unload all raylib assets.
	r.UnloadAll()
}
