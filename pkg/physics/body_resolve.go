package physics

import (
	r "github.com/lachee/raylib-goplus/raylib"
)

func (b *Body) resolveSpace(space *Space, fn func(...interface{})) {
	tt := make([]interface{}, len(*space))
	for i := range *space {
		tt[i] = (*space)[i]
	}

	fn(tt...)
}

func (b *Body) resolveRectangle(
	rec *Rectangle, tmpRec r.Rectangle, prevCol bool, velocity float32, velocityMult r.Vector2,
) bool {
	// If the player hasn't collided with anything on the axis yet and is
	// overlapping with the physics rectangle.
	if !prevCol && rec.Overlaps(tmpRec) {
		overlap := rec.Rectangle.GetOverlapRec(tmpRec)
		prevCol = true

		if velocity*velocityMult.Y > 0 {
			b.onGround = true
		}

		if velocity > 0 {
			b.velocity.X -= (overlap.Width * velocityMult.X)
			b.velocity.Y -= (overlap.Height * velocityMult.Y)
		} else {
			b.velocity.X += (overlap.Width * velocityMult.X)
			b.velocity.Y += (overlap.Height * velocityMult.Y)
		}
	}

	return prevCol
}

func (b *Body) resolveSlope(t *Slope, tmpYRec r.Rectangle, original r.Vector2) bool {
	// Get the intersection points given a temporary Y rectangle.
	intersections := t.GetIntersectionPoints(
		NewRectangle(tmpYRec.X, tmpYRec.Y, tmpYRec.Width, tmpYRec.Height),
	)
	// If there are no intersections then return.
	if len(intersections) == 0 {
		return false
	}

	// Create a line of best fit across the intersection points.
	tmpLine := NewSlope(
		r.NewVector2(intersections[0].X, intersections[0].Y),
		r.NewVector2(intersections[len(intersections)-1].X, intersections[len(intersections)-1].Y),
	)

	dy := tmpLine.p2.Y - tmpLine.p1.Y
	colBox := r.NewRectangle(tmpYRec.X, intersections[0].Y+(dy/2), tmpYRec.Width, tmpYRec.Height/2)
	overlap := colBox.GetOverlapRec(tmpYRec)

	// Since there was a collision, set the onGround to true right away.
	b.onGround = true
	b.velocity.Y = original.Y - overlap.Height

	return true
}
