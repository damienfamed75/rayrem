package common

import (
	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	Controls controls
)

type controls struct {
	Left     r.Key
	Right    r.Key
	Jump     r.Key
	Shoot    r.Key
	Interact r.Key
}

func loadControls() {
	Controls = controls{
		Left:     r.Key(PublicConfig.GetInt32("controls.left")),
		Right:    r.Key(PublicConfig.GetInt32("controls.right")),
		Jump:     r.Key(PublicConfig.GetInt32("controls.jump")),
		Shoot:    r.Key(PublicConfig.GetInt32("controls.shoot")),
		Interact: r.Key(PublicConfig.GetInt32("controls.interact")),
	}
}
