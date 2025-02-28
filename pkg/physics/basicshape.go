package physics

import (
	"github.com/damienfamed75/rayrem/pkg/common"
)

// BasicShape has all essentials of a Shape.
type BasicShape struct {
	uid  uint64
	tags []common.Tag
}

func NewBasicShape() *BasicShape {
	return &BasicShape{
		uid: newID(),
	}
}

func (s *BasicShape) ID() uint64 {
	return s.uid
}

// AddTags appends new tags without checking for uniqueness.
func (s *BasicShape) AddTags(tags ...common.Tag) {
	s.tags = append(s.tags, tags...)
}

// ClearTags empties all the tags in this shape.
func (s *BasicShape) ClearTags() {
	s.tags = []common.Tag{}
}

// Tags returns a list of the tags.
func (s *BasicShape) Tags() []common.Tag {
	return s.tags
}

// RemoveTags removes all of the tags provided.
func (s *BasicShape) RemoveTags(tags ...common.Tag) {
	for _, t := range tags {
		for i := len(s.tags) - 1; i >= 0; i-- {
			if t == s.tags[i] {
				s.tags = append(s.tags[:i], s.tags[i+1:]...)
			}
		}
	}
}

// HasTags returns whether or not this shape has all the tags provided.
func (s *BasicShape) HasTags(tags ...common.Tag) bool {
	hasTags := true

	for _, wanted := range tags {
		found := false
		for _, shapeTag := range s.tags {
			if wanted == shapeTag {
				found = true
				continue
			}
		}
		if !found {
			hasTags = false
			break
		}
	}

	return hasTags
}
