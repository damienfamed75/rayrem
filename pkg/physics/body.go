package physics

import (
	"sort"

	"github.com/damienfamed75/rayrem/pkg/common"

	r "github.com/lachee/raylib-goplus/raylib"
)

var (
	_ Shape = &Body{}
)

// Body returns a physics rigidbody that reads ground as elements
// that are tagged as common.TagGround and gravity is defaulted to what is in
// the settings.
type Body struct {
	Velocity r.Vector2

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
	b.Velocity.Y += b.gravity * dt

	b.maxVelocityCheck()

	if b.Velocity.Y < -(b.gravity * dt) {
		b.onGround = false
	}

	b.ResolveForces(dt)

	b.Move(b.Velocity.X, b.Velocity.Y)
}

// ResolveForces loops through the collision shapes and checks if they are
// currently colliding with anything with the given velocity. If so then
// the velocity is manipulated so that the potential collision will never happen.
func (b *Body) ResolveForces(dt float32) {
	// Limit the player to touching one object on each axis at a time.
	// This means that numbers won't get messed up when touching two
	// ground objects at the same time.
	var colx, coly bool

	for i := range *b.Space {
		collider := (*b.Space)[i].(*Rectangle).Rectangle
		tmpXRec := collider.Move(b.Velocity.X, 0)
		tmpYRec := collider.Move(0, b.Velocity.Y)
		original := b.Velocity

		// Get all possible collision boxes.
		// TODO check for uniqueness.
		var possible []Transformer
		possible = b.solids.Retrieve(tmpXRec)
		possible = append(possible, b.solids.Retrieve(tmpYRec)...)

		for _, p := range possible {
			switch t := p.(type) {
			case *Rectangle:
				// If the player hasn't collided with anything on the x-axis yet.
				if !colx {
					if t.Overlaps(tmpXRec) {
						overlap := t.Rectangle.GetOverlapRec(tmpXRec)
						colx = true

						if b.Velocity.X > 0 {
							b.Velocity.X -= overlap.Width
						} else {
							b.Velocity.X += overlap.Width
						}
					}
				}

				// If the player hasn't collided with anything on the y-axis yet.
				if !coly {
					if t.Overlaps(tmpYRec) {
						overlap := t.Rectangle.GetOverlapRec(tmpYRec)
						coly = true

						// If the player is moving downward and colliding.
						if b.Velocity.Y > 0 {
							b.onGround = true
						}

						if b.Velocity.Y > 0 {
							b.Velocity.Y -= overlap.Height
						} else {
							b.Velocity.Y += overlap.Height
						}
					}
				}
			// SlopePlatform is just three slopes.
			case *SlopePlatform:
				if t.Overlaps(tmpYRec) {
					// Ignore coly check, because slopes take higher priority.
					// Resolve all the slopes within the platform.
					coly = b.resolveSlope(t.landingZone1, tmpYRec, original)
					coly = b.resolveSlope(t.landingZone2, tmpYRec, original)
					coly = b.resolveSlope(t.slope, tmpYRec, original)
				}
			case *Slope:
				if t.Overlaps(tmpYRec) {
					coly = b.resolveSlope(t, tmpYRec, original)
				}
			}
		}
	}
}

func (b *Body) resolveSlope(t *Slope, tmpYRec r.Rectangle, original r.Vector2) (coly bool) {
	var intersections []IntersectionPoint

	side := NewSlope(
		r.NewVector2(tmpYRec.X, tmpYRec.Y),
		r.NewVector2(tmpYRec.X, tmpYRec.MaxPosition().Y),
	)
	intersections = append(intersections, t.GetIntersectionPoints(side)...)

	side.p1 = r.NewVector2(tmpYRec.MaxPosition().X, tmpYRec.MaxPosition().Y)
	intersections = append(intersections, t.GetIntersectionPoints(side)...)

	side.p2 = r.NewVector2(tmpYRec.MaxPosition().X, tmpYRec.Y)
	intersections = append(intersections, t.GetIntersectionPoints(side)...)

	// intersections := t.GetIntersectionPoints(NewRectangle(tmpYRec.X, tmpYRec.Y, tmpYRec.Width, tmpYRec.Height))
	if len(intersections) == 0 {
		return
	}

	// Sort the intersections based on Y.
	sort.Slice(intersections, func(i, j int) bool {
		return intersections[i].Y > intersections[j].Y
	})

	tmpL := NewSlope(
		r.NewVector2(intersections[0].X, intersections[0].Y),
		r.NewVector2(intersections[len(intersections)-1].X, intersections[len(intersections)-1].Y),
	)

	dy := tmpL.p2.Y - tmpL.p1.Y
	colBox := r.NewRectangle(tmpYRec.X, intersections[0].Y+(dy/2), tmpYRec.Width, tmpYRec.Height/2)
	overlap := colBox.GetOverlapRec(tmpYRec)

	b.onGround = true
	coly = true

	b.Velocity.Y = original.Y - overlap.Height

	return
}

func (b *Body) maxVelocityCheck() {
	// Cap player movement speed.
	if b.Velocity.X > b.maxSpeed.X {
		b.Velocity.X = b.maxSpeed.X
	}
	if b.Velocity.X < -b.maxSpeed.X {
		b.Velocity.X = -b.maxSpeed.X
	}

	if b.Velocity.Y > b.maxSpeed.Y {
		b.Velocity.Y = b.maxSpeed.Y
	}
	if b.Velocity.Y < -b.maxSpeed.Y {
		b.Velocity.Y = -b.maxSpeed.Y
	}
}
