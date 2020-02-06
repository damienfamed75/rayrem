package object

import (
	r "github.com/lachee/raylib-goplus/raylib"
)

type OpenOption func(*Openable)

// func WithSprite(s *aseprite.File) OpenOption {
// 	return func(o *Openable) {
// 		o.collider = physics.NewRectangle(c.X, c.Y, c.Width, c.Height)
// 	}
// }

// func WithCollider(c r.Rectangle) OpenOption {
// 	return func(o *Openable) {
// 		o.collider = physics.NewRectangle(c.X, c.Y, c.Width, c.Height)
// 	}
// }

func WithCollider() OpenOption {
	return func(o *Openable) {
		o.hasCollider = true
	}
}

// WithAOA or With Area of Activation creates a zone for the openable to use
// around the sprite to be the activation area.
func WithAOA(c r.Rectangle) OpenOption {
	return func(o *Openable) {
		o.zoneCollider = c
	}
}

func WithLock(l *Lock) OpenOption {
	return func(o *Openable) {
		o.lock = l
	}
}
