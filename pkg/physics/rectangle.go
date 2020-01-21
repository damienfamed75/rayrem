package physics

import r "github.com/lachee/raylib-goplus/raylib"

var (
	_ Shape = &Rectangle{}
)

// Rectangle is a four sided polygon that uses a raylib.Rectangle to represent
// itself.
type Rectangle struct {
	*BasicShape
	r.Rectangle
}

// NewRectangle creates a rectangle and basic object.
func NewRectangle(x, y, w, h float32) *Rectangle {
	return &Rectangle{
		BasicShape: &BasicShape{},
		Rectangle:  r.NewRectangle(x, y, w, h),
	}
}

// Overlaps checks if the rectangle is overlapping another raylib rectangle.
func (rec *Rectangle) Overlaps(re r.Rectangle) bool {
	return rec.Rectangle.Overlaps(re)
}

// Position gets the small point of the rectangle.
func (rec *Rectangle) Position() r.Vector2 {
	return rec.Rectangle.Position()
}

// Center gets the center point of the rectangle.
func (rec *Rectangle) Center() r.Vector2 {
	return rec.Rectangle.Position()
}

// SetPosition sets the placement of the rectangle to the provided coordinates.
func (rec *Rectangle) SetPosition(x, y float32) {
	rec.Rectangle = rec.Rectangle.SetPosition(r.NewVector2(x, y))
}

// Move moves the rectangle based on the delta provided.
func (rec *Rectangle) Move(x, y float32) {
	rec.Rectangle = rec.Rectangle.Move(x, y)
}

// Width returns the width of the rectangle.
func (rec *Rectangle) Width() float32 {
	return rec.Rectangle.Width
}

// Height returns the height of the rectangle.
func (rec *Rectangle) Height() float32 {
	return rec.Rectangle.Height
}
