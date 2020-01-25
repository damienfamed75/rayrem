package physics

import (
	r "github.com/lachee/raylib-goplus/raylib"
)

type SpatialHashmap struct {
	list             []Transformer
	cellSize         int
	lastTotalCleared int
	getKeys          func(Transformer) [][2]int
	hash             map[[2]int][]Transformer
}

type Transformer interface {
	Position() r.Vector2
	MaxPosition() r.Vector2
}

func NewSpatialHashmap(cellSize int) *SpatialHashmap {
	s := &SpatialHashmap{
		cellSize: cellSize,
		getKeys:  makeKeys(cellSize),
		hash:     make(map[[2]int][]Transformer),
	}

	return s
}

func makeKeys(shift int) func(Transformer) [][2]int {
	return func(t Transformer) [][2]int {
		sx := int(t.Position().X) >> shift
		sy := int(t.Position().Y) >> shift
		ex := int(t.MaxPosition().X) >> shift
		ey := int(t.MaxPosition().Y) >> shift

		var x, y int
		var keys [][2]int

		for y = sy; y <= ey; y++ {
			for x = sx; x <= ex; x++ {
				keys = append(keys, [2]int{x, y})
			}
		}

		return keys
	}
}

func (s *SpatialHashmap) Clear() {
	for key := range s.hash {
		if len(s.hash[key]) == 0 {
			delete(s.hash, key)
		} else {
			s.hash[key] = []Transformer{}
		}
	}

	s.list = []Transformer{}
}

func (s *SpatialHashmap) InsertMulti(tt ...Transformer) {
	for i := range tt {
		s.Insert(tt[i])
	}
}

func (s *SpatialHashmap) Insert(t Transformer) {
	keys := s.getKeys(t)
	s.list = append(s.list, t)

	for i := 0; i < len(keys); i++ {
		key := keys[i]
		if s.hash[key] != nil {
			s.hash[key] = append(s.hash[key], t)
		} else {
			s.hash[key] = []Transformer{t}
		}
	}
}

func (s *SpatialHashmap) NumBuckets() int {
	var count int
	for key := range s.hash {
		if len(s.hash[key]) > 0 {
			count++
		}
	}

	return count
}

func (s *SpatialHashmap) Retrieve(t Transformer) []Transformer {
	var res []Transformer

	if t == nil {
		return s.list
	}

	keys := s.getKeys(t)
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		if s.hash[key] != nil {
			res = append(res, s.hash[key]...)
		}
	}

	return res
}
