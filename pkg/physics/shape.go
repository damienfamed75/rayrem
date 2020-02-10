package physics

import (
	"sync/atomic"

	"github.com/damienfamed75/rayrem/pkg/common"
	r "github.com/lachee/raylib-goplus/raylib"
)

// Shape is an interface to fulfill any collider in this package to be passed
// into a Space. Even the Space is a Shape so Spaces can go in Spaces.
type Shape interface {
	Tags() []common.Tag
	HasTags(tags ...common.Tag) bool
	AddTags(tags ...common.Tag)
	RemoveTags(tags ...common.Tag)
	ClearTags()

	Overlaps(rec r.Rectangle) bool
	Position() r.Vector2
	Center() r.Vector2
	MaxPosition() r.Vector2
	SetPosition(x, y float32)
	Move(x, y float32)
	Width() float32
	Height() float32

	Entity
}

type Entity interface {
	ID() uint64
}

var idTracker uint64

func newID() uint64 {
	atomic.AddUint64(&idTracker, 1)
	return atomic.LoadUint64(&idTracker)
}
