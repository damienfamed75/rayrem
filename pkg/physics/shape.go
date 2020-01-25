package physics

import (
	r "github.com/lachee/raylib-goplus/raylib"
)

// Shape is an interface to fulfill any collider in this package to be passed
// into a Space. Even the Space is a Shape so Spaces can go in Spaces.
type Shape interface {
	Tags() []string
	HasTags(tags ...string) bool
	AddTags(tags ...string)
	RemoveTags(tags ...string)
	ClearTags()

	Overlaps(rec r.Rectangle) bool
	Position() r.Vector2
	Center() r.Vector2
	MaxPosition() r.Vector2
	SetPosition(x, y float32)
	Move(x, y float32)
	Width() float32
	Height() float32
}
