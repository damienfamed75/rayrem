package physics

import (
	"github.com/damienfamed75/rayrem/pkg/common"

	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ Shape = &Body{}

	_ Moveable = &Body{}

	_ Entity = &Body{}
)

// Body returns a physics rigidbody that reads ground as elements
// that are tagged as common.TagGround and gravity is defaulted to what is in
// the settings.
type Body struct {
	velocity r.Vector2

	gravity  float32
	onGround bool
	maxSpeed r.Vector2
	solids   *SpatialHashmap

	*Space
}

// NewBody creates a default rigidbody and tags the Rigidbody space as a
// physics body.
func NewBody(collision *Space, solids *SpatialHashmap, maxSpeed r.Vector2) *Body {
	b := &Body{
		Space:    collision,
		maxSpeed: maxSpeed,
		gravity:  common.Config.Game.Gravity,
		solids:   solids,
	}

	b.AddTags(common.TagPhysicsBody)

	return b
}

func (b *Body) Velocity() r.Vector2 {
	return b.velocity
}

func (b *Body) SetVelocity(x, y float32) {
	b.velocity.X = x
	b.velocity.Y = y
}

func (b *Body) AddVelocity(x, y float32) {
	b.velocity = b.velocity.Add(r.NewVector2(x, y))
}

// SetGravity overrides the default gravity.
func (b *Body) SetGravity(g float32) {
	b.gravity = g
}

// OnGround returns if the collision space is touching ground elements on the
// Y axis.
func (b *Body) OnGround() bool {
	return b.onGround
}

// Update checks for collisions in the world against the colliders in the
// rigidbody space and adds gravity.
func (b *Body) Update(dt float32) {
	b.velocity.Y += b.gravity * dt

	b.maxVelocityCheck()

	if b.velocity.Y < -(b.gravity * dt) {
		b.onGround = false
	}

	b.ResolveForces(dt)

	b.Move(b.velocity.X, b.velocity.Y)
}

// ResolveForces loops through the collision shapes and checks if they are
// currently colliding with anything with the given velocity. If so then
// the velocity is manipulated so that the potential collision will never happen.
func (b *Body) ResolveForces(dt float32) {
	// Limit the player to touching one object on each axis at a time.
	// This means that numbers won't get messed up when touching two
	// ground objects at the same time.
	// var colx, coly bool
	var col colCheck

	for i := range *b.Space {
		collider := (*b.Space)[i].(*Rectangle).Rectangle
		tmpXRec := collider.Move(b.velocity.X, 0)
		tmpYRec := collider.Move(0, b.velocity.Y)
		original := b.velocity

		// Get all possible collision boxes.
		possible := b.solids.Retrieve(collider.Move(b.velocity.X, b.velocity.Y))

		b.resolveShapes(col, original, collider, tmpXRec, tmpYRec, possible...)
	}
}

func (b *Body) resolveShapes(col colCheck, original r.Vector2, collider, tmpXRec, tmpYRec r.Rectangle, possible ...interface{}) colCheck {
	for _, p := range possible {
		// Skip from colliding against itself.
		if ss, ok := p.(Entity); ok {
			if ss.ID() == b.ID() {
				continue
			}
		}

		switch t := p.(type) {
		case *Actor:
			b.resolveSpace(
				t.Space, func(tt ...interface{}) {
					col = b.resolveShapes(col, original, collider, tmpXRec, tmpYRec, tt...)
				},
			)
		case *Space:
			b.resolveSpace(
				t, func(tt ...interface{}) {
					col = b.resolveShapes(col, original, collider, tmpXRec, tmpYRec, tt...)
				},
			)
		case *Rectangle:
			// Resolve the X collisions.
			col.SetX(b.resolveRectangle(
				t, tmpXRec, col.X(), b.velocity.X, r.NewVector2(1, 0),
			))
			// Resolve the Y collisions.
			col.SetY(b.resolveRectangle(
				t, tmpYRec, col.Y(), b.velocity.Y, r.NewVector2(0, 1),
			))
		// SlopePlatform is just three slopes.
		case *SlopePlatform:
			if t.Overlaps(tmpYRec) {
				// Ignore coly check, because slopes take higher priority.
				// Resolve all the slopes within the platform.
				col.SetY(b.resolveSlope(t.landingZone1, tmpYRec, original))
				col.SetY(b.resolveSlope(t.landingZone2, tmpYRec, original))
				col.SetY(b.resolveSlope(t.slope, tmpYRec, original))
			}
		case *Slope:
			if t.Overlaps(tmpYRec) {
				col.SetY(b.resolveSlope(t, tmpYRec, original))
			}
		case *Platform:
			if !col.Y() && b.velocity.Y > 0 && t.Overlaps(tmpYRec) {
				overlap := t.Rectangle.GetOverlapRec(tmpYRec)

				// If the overlapped rectangle is more than halfway down
				// the entity, then count this as a collision.
				if overlap.Y > tmpYRec.Center().Y {
					col.SetY(true)
					b.onGround = true
					b.velocity.Y -= overlap.Height
				}
			}
		case *Zone:
			if t.Overlaps(collider) {
				overlap := t.Rectangle.GetOverlapRec(collider)
				t.dispatchMessage(b.Space, overlap)
			}
		}
	}

	return col
}

func (b *Body) maxVelocityCheck() {
	// Cap player movement speed.
	if b.velocity.X > b.maxSpeed.X {
		b.velocity.X = b.maxSpeed.X
	}
	if b.velocity.X < -b.maxSpeed.X {
		b.velocity.X = -b.maxSpeed.X
	}

	if b.velocity.Y > b.maxSpeed.Y {
		b.velocity.Y = b.maxSpeed.Y
	}
	if b.velocity.Y < -b.maxSpeed.Y {
		b.velocity.Y = -b.maxSpeed.Y
	}
}
