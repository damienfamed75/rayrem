package physics

import r "github.com/lachee/raylib-goplus/raylib"

// Transformer is the bare minimum to be passed into a spatial hashmap.
type Transformer interface {
	Position() r.Vector2
	MaxPosition() r.Vector2
}

type Moveable interface {
	Transformer
	Velocity() r.Vector2
}
