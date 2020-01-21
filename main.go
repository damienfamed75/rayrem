package main

// The main.go file for the game is in the root so then the file packager
// works properly. Until there is a workaround then this file will remain
// here in the root.

import (
	"log"

	"github.com/damienfamed75/rayrem/pkg/common"
	"github.com/damienfamed75/rayrem/pkg/game"

	r "github.com/lachee/raylib-goplus/raylib"
	"github.com/markbates/pkger"
)

func main() {
	// Package the assets and config files into the binary.
	pkger.Include("/assets/")
	pkger.Include("/config/")

	// Load config files.
	err := common.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Initialize the window.
	r.InitWindow(
		common.PublicConfig.GetInt("screen.width"),
		common.PublicConfig.GetInt("screen.height"),
		"rayrem",
	)
	defer r.CloseWindow()

	// Set the target FPS to 60. There is no plan to make this configurable
	// until there is a reason to.
	r.SetTargetFPS(60)

	// If the configuration is set to be fullscreen then toggle fullscreen.
	if common.PublicConfig.GetBool("screen.fullscreen") == true {
		r.ToggleFullscreen()
	}

	// Set the master volume to what the settings say.
	r.SetMasterVolume(float32(common.PublicConfig.GetFloat64("volume.master")))

	// Create a new game structure.
	g := game.NewGame()
	defer g.Unload()

	// Set the default scene to the main menu.
	g.SetScene(common.ModeMainMenu)

	for !r.WindowShouldClose() {
		// Update the game's current scene.
		g.Update(r.GetFrameTime())

		r.BeginDrawing()
		// Draw the game's current scene.
		g.Draw()
		r.EndDrawing()
	}
}
