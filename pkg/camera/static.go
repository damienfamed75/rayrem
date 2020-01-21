package camera

import (
	r "github.com/lachee/raylib-goplus/raylib"
)

// StaticCamera doesn't move at all.
type StaticCamera struct {
	r.Camera2D
}

// NewStatic creates a camera that doesn't move
func NewStatic(offset r.Vector2, zoom float32) *StaticCamera {
	return &StaticCamera{
		Camera2D: r.Camera2D{
			Offset:   offset,
			Rotation: 0,
			Zoom:     zoom,
		},
	}
}
