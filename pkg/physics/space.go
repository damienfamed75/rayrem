package physics

import r "github.com/lachee/raylib-goplus/raylib"

var (
	_ Shape = &Space{}
)

// Space is just a collection of shapes.
type Space []Shape

// NewSpace returns a reference to an empty slice of shapes.
func NewSpace() *Space {
	return &Space{}
}

// Clear empties the slice.
func (s *Space) Clear() {
	*s = []Shape{}
}

// Add appends a shape without uniqueness.
func (s *Space) Add(shapes ...Shape) {
	*s = append(*s, shapes...)
}

// Center gets the center of all the shapes within the space.
func (s *Space) Center() r.Vector2 {
	min := s.Position()
	max := s.MaxPosition()

	return r.NewVector2(
		min.X-max.X,
		min.Y-max.Y,
	)
}

// Width gets the overall width from all the shapes combined.
func (s *Space) Width() float32 {
	min := s.Position()
	max := s.MaxPosition()

	return min.X - max.X
}

// Height gets the overall height from all the shapes combined.
func (s *Space) Height() float32 {
	min := s.Position()
	max := s.MaxPosition()

	return min.Y - max.Y
}

// Move loops through all the shapes in the space and moves them.
func (s *Space) Move(x, y float32) {
	for i := range *s {
		(*s)[i].Move(x, y)
	}
}

// MaxPosition gets the highest coordinate in the space.
func (s *Space) MaxPosition() r.Vector2 {
	var (
		res r.Vector2
		set r.Vector2
	)

	for i := range *s {
		tmp := (*s)[i].Position()
		if tmp.X > res.X || set.X == 0 {
			set.X = 1
			res.X = tmp.X
		}
		if tmp.Y > res.Y || set.Y == 0 {
			set.Y = 1
			res.Y = tmp.Y
		}
	}

	return res
}

// Position gets the smallest coordinate in the space.
func (s *Space) Position() r.Vector2 {
	var (
		res r.Vector2
		set r.Vector2
	)

	for i := range *s {
		tmp := (*s)[i].Position()
		if tmp.X < res.X || set.X == 0 {
			set.X = 1
			res.X = tmp.X
		}
		if tmp.Y < res.Y || set.Y == 0 {
			set.Y = 1
			res.Y = tmp.Y
		}
	}

	return res
}

// SetPosition sets the position of the space, which moves every shape
// according to their delta.
func (s *Space) SetPosition(x, y float32) {
	if len(*s) > 0 {

		delta := (*s)[0].Position()
		delta.X = x - delta.X
		delta.Y = y - delta.Y

		for _, shape := range *s {
			shape.Move(delta.X, delta.Y)
		}

	}
}

// Overlaps sees if the provided rectangle overlaps any shapes within the space.
func (s *Space) Overlaps(rec r.Rectangle) bool {
	for i := range *s {
		if (*s)[i].Overlaps(rec) {
			return true
		}
	}

	return false
}

// Remove removes specified shapes from the space if they exist.
// This function checks by value not reference.
func (s *Space) Remove(rec ...Shape) {
	for i := range rec {
		for j := len(*s) - 1; i >= 0; i-- {
			if rec[i] == (*s)[j] {
				*s = append((*s)[:j], (*s)[j+1:]...)
			}
		}
	}
}

// AddTags adds tags to all shapes.
func (s *Space) AddTags(tags ...string) {
	for i := range *s {
		for _, t := range tags {
			if !(*s)[i].HasTags(t) {
				(*s)[i].AddTags(t)
			}
		}
	}
}

// ClearTags removes all tags from the space's shapes.
func (s *Space) ClearTags() {
	for i := range *s {
		(*s)[i].ClearTags()
	}
}

// Tags gets all tags from its shapes and returns a big list of them.
func (s *Space) Tags() []string {
	var tmp = &Rectangle{}
	for i := range *s {
		tt := (*s)[i].Tags()
		for _, t := range tt {
			if !tmp.HasTags(t) {
				tmp.AddTags(t)
			}
		}
	}

	return tmp.Tags()
}

// HasTags sees if any shapes have all those tags.
func (s *Space) HasTags(tags ...string) bool {
	for i := range *s {
		for _, t := range tags {
			if (*s)[i].HasTags(t) {
				return true
			}
		}
	}
	return false
}

// RemoveTags removes tags provided from the shapes if they exist.
func (s *Space) RemoveTags(tags ...string) {
	for i := range *s {
		(*s)[i].RemoveTags(tags...)
	}
}

// Filter is a custom filterer to remove shapes from a list based on a
// certain property specified by the user.
func (s *Space) Filter(filter func(Shape) bool) *Space {
	subSpace := &Space{}
	for i := range *s {
		if filter((*s)[i]) {
			subSpace.Add((*s)[i])
		}
	}

	return subSpace
}

// FilterByTags filters out any shapes that don't have all the tags provided.
func (s *Space) FilterByTags(tags ...string) *Space {
	return s.Filter(func(r Shape) bool {
		if r.HasTags(tags...) {
			return true
		}
		return false
	})
}

// FilterOutByTags filters out any shapes that have all the tags provided.
func (s *Space) FilterOutByTags(tags ...string) *Space {
	return s.Filter(func(r Shape) bool {
		if r.HasTags(tags...) {
			return false
		}
		return true
	})
}
