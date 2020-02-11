package object

import (
	"github.com/damienfamed75/rayrem/pkg/physics"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Lockable is inherited by structures that are allowed to be locked.
// Note: To create a Lockable structure you can embed a Lockable structure.
type Lockable interface {
	applyLock(*Lock)
}

// Option is a general use optional function that can be passed into game
// objects to alter them.
type Option func(interface{})

// WithLock will attempt to apply a lock to a Lockable game object.
func WithLock(l *Lock) Option {
	return func(ii interface{}) {
		// If the game object is a Lockable then apply the lock. If not then
		// ignore the configuration entirely.
		if lockable, ok := ii.(Lockable); ok {
			lockable.applyLock(l)
		}
	}
}

// withCollider is an exclusive option for interactable.
func withCollider(c r.Rectangle) Option {
	return func(ii interface{}) {
		if interact, ok := ii.(*interactable); ok {
			interact.collider = physics.NewRectangle(
				c.X, c.Y, c.Width, c.Height,
			)
		}
	}
}
