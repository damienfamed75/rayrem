package physics

import (
	"sort"

	r "github.com/lachee/raylib-goplus/raylib"
)

var _ Shape = &Slope{}

// Slope is just a line, but should be primarily used for sloped platforms.
// See slopeplatform.go for more information.
type Slope struct {
	*BasicShape
	p1 r.Vector2
	p2 r.Vector2
}

// NewSlope returns a barebones line from A to B.
func NewSlope(p1, p2 r.Vector2) *Slope {
	return &Slope{
		BasicShape: NewBasicShape(),
		p1:         p1,
		p2:         p2,
	}
}

// Overlaps checks for intersection points against a raylib rectangle and returns
// true if any collisions are detected.
func (l *Slope) Overlaps(re r.Rectangle) bool {
	rec := NewRectangle(re.X, re.Y, re.Width, re.Height)
	intersection := l.GetIntersectionPoints(rec)

	// If there were more than 0 intersection points then there was a collision.
	colliding := len(intersection) > 0

	if !colliding {
		// Catch any corner cases where the line may be colliding, but wasn't caught.
		return (l.p1.X >= re.X && l.p1.Y >= re.Y && l.p1.X < re.X+re.Width && l.p1.Y < re.Y+re.Height) ||
			(l.p2.X >= re.X && l.p2.Y >= re.Y && l.p2.X < re.X+re.Width && l.p2.Y < re.Y+re.Height)
	}

	return colliding
}

// Points returns raylib vectors of the two points of the slope.
func (l *Slope) Points() (r.Vector2, r.Vector2) {
	return l.p1, l.p2
}

// Delta returns the delta (or difference) between the start and end point of a Slope.
func (l *Slope) Delta() (float32, float32) {
	dx := l.p2.X - l.p1.X
	dy := l.p2.Y - l.p1.Y
	return dx, dy
}

// GetLength returns the length of the Slope.
func (l *Slope) GetLength() float32 {
	return Distance(l.p1.X, l.p1.Y, l.p2.X, l.p2.Y)
}

// Center returns the center X and Y values of the Slope.
func (l *Slope) Center() r.Vector2 {
	x := l.p1.X + ((l.p2.X - l.p1.X) / 2)
	y := l.p1.Y + ((l.p2.Y - l.p1.Y) / 2)

	return r.NewVector2(x, y)
}

// Move moves the Slope by the values specified.
func (l *Slope) Move(x, y float32) {
	l.p1.X += x
	l.p1.Y += y

	l.p2.X += x
	l.p2.Y += y
}

// Position is just filling the Shape interface and returns the first point.
func (l *Slope) Position() r.Vector2 {
	return l.p1
}

// MaxPosition is just filling the Shape interface and returns the second point.
func (l *Slope) MaxPosition() r.Vector2 {
	return l.p2
}

// SetPosition updates based on a difference.
func (l *Slope) SetPosition(x, y float32) {
	diff := r.NewVector2(l.Width(), l.Height())

	l.p1.X = x
	l.p1.Y = y

	l.p2.X = x + diff.X
	l.p2.Y = y + diff.Y
}

// Height returns the Y difference between point 1 and 2.
func (l *Slope) Height() float32 {
	if l.p1.Y > l.p2.Y {
		return l.p1.Y - l.p2.Y
	}

	return l.p2.Y - l.p1.Y
}

// Width returns the X difference between point 1 and 2.
func (l *Slope) Width() float32 {
	if l.p1.X > l.p2.X {
		return l.p1.X - l.p2.X
	}

	return l.p2.X - l.p1.X
}

// IntersectionPoint represents a point of intersection from a Slope with another Shape.
type IntersectionPoint struct {
	X, Y  float32
	Shape Shape
}

// GetIntersectionPoints returns points based on another shape.
func (l *Slope) GetIntersectionPoints(other Shape) []IntersectionPoint {
	intersections := []IntersectionPoint{}

	switch b := other.(type) {

	case *Slope:
		// When calculating collisions against another slope/line we find the
		// determinant to represent the minor of a 2x2 matrix and with that

		// determinant is the volume scaling factor.
		det := (l.p2.X-l.p1.X)*(b.p2.Y-b.p1.Y) - (b.p2.X-b.p1.X)*(l.p2.Y-l.p1.Y)

		if det != 0 {

			// MAGIC MATH; the extra + 1 here makes it so that corner cases work.
			// These two lines calculate the distance to intersection points.
			// In other words these are the cross products of the matrix.
			lambda := (float32(((l.p1.Y-b.p1.Y)*(b.p2.X-b.p1.X))-((l.p1.X-b.p1.X)*(b.p2.Y-b.p1.Y))) + 1) / float32(det)
			gamma := (float32(((l.p1.Y-b.p1.Y)*(l.p2.X-l.p1.X))-((l.p1.X-b.p1.X)*(l.p2.Y-l.p1.Y))) + 1) / float32(det)

			// Detect coincident lines that have a collision.
			if (0 < lambda && lambda < 1) && (0 < gamma && gamma < 1) {
				dx, dy := l.Delta()
				intersections = append(intersections, IntersectionPoint{l.p1.X + float32(lambda*float32(dx)), l.p1.Y + float32(lambda*float32(dy)), other})
			}

		}
	case *Rectangle:
		// For rectangles we section it off into its individual sides and then
		// recursively solve it out through as slopes.

		side := NewSlope(r.NewVector2(b.Rectangle.X, b.Rectangle.Y), r.NewVector2(b.Rectangle.X, b.Rectangle.Y+b.Rectangle.Height))
		intersections = append(intersections, l.GetIntersectionPoints(side)...)

		side.p1.Y = b.Rectangle.Y + b.Rectangle.Height
		side.p2.X = b.Rectangle.X + b.Rectangle.Width
		side.p2.Y = b.Rectangle.Y + b.Rectangle.Height
		intersections = append(intersections, l.GetIntersectionPoints(side)...)

		side.p1.X = b.Rectangle.X + b.Rectangle.Width
		side.p2.Y = b.Rectangle.Y
		intersections = append(intersections, l.GetIntersectionPoints(side)...)

		side.p1.Y = b.Rectangle.Y
		side.p2.X = b.Rectangle.X
		side.p2.Y = b.Rectangle.Y
		intersections = append(intersections, l.GetIntersectionPoints(side)...)
	case *Space:
		for _, shape := range *b {
			intersections = append(intersections, l.GetIntersectionPoints(shape)...)
		}
	}

	// Sort the slice by distance.
	sort.Slice(intersections, func(i, j int) bool {
		return Distance(l.p1.X, l.p1.Y, intersections[i].X, intersections[i].Y) < Distance(l.p1.X, l.p1.Y, intersections[j].X, intersections[j].Y)
	})

	return intersections

}

// Draw is used for debugging and draws a triangle.
func (l *Slope) Draw() {
	// Draw the slope.
	r.DrawLineEx(l.p1, l.p2, 1, r.Green)
	// Draw each side of the triangle.
	r.DrawLineEx(l.p1, r.NewVector2(l.p2.X, l.p1.Y), 1, r.Green)
	r.DrawLineEx(l.p2, r.NewVector2(l.p2.X, l.p1.Y), 1, r.Green)
}
