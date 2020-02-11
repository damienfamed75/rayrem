package common

import (
	r "github.com/lachee/raylib-goplus/raylib"
)

// Unloader is used when you have a group of items that need to be unloaded from
// a running scene or from the game itself.
type Unloader interface {
	r.Unloadable
}

// Drawer encompasses any objects that may draw on the screen.
type Drawer interface {
	Draw()
}

// BasicObject is something that may draw and have special treatment such as
// updating its position between frames.
type BasicObject interface {
	Update(dt float32)

	Drawer
}

// Entity is something interacteable such as the player, enemies, or bosses.
type Entity interface {
	TakeDamage()
	Position() r.Vector2

	BasicObject
}

// SceneManager is an object that has the power to change the current scene.
type SceneManager interface {
	SetScene(Mode)
}

// Scene is a current instance in which the game is running. It can cause
// updates, draws, and scene changes with permission from the scene manager.
type Scene interface {
	Update(dt float32)
	Draw()
	Unload()
}
