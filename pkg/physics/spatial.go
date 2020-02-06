package physics

import (
	"fmt"
)

// SpatialHashmap is a data structure to tell what objects are close, if not
// already colliding with each other.
// More info: http://hhoppe.com/perfecthash.pdf
type SpatialHashmap struct {
	list             []Transformer
	cellSize         int
	lastTotalCleared int
	getKeys          func(Transformer) []point
	hash             map[point][]Transformer
}

// SpatialAdder is a custom adder to the spatial hashmap.
type SpatialAdder interface {
	Add(w *SpatialHashmap) // custom add function
}

// NewSpatialHashmap returns a hashmap with the sensitivity given.
// The larger the cellSize, the more sensitive the hashmap is.
func NewSpatialHashmap(cellSize int) *SpatialHashmap {
	s := &SpatialHashmap{
		cellSize: cellSize,
		getKeys:  makeKeys(cellSize),
		hash:     make(map[point][]Transformer),
	}

	return s
}

// makeKeys returns a function to store into memory, this is faster when
// calculating the shift sensitivity to newly inserted objects.
func makeKeys(shift int) func(Transformer) []point {
	return func(t Transformer) []point {
		sx := int(t.Position().X) >> shift
		sy := int(t.Position().Y) >> shift
		ex := int(t.MaxPosition().X) >> shift
		ey := int(t.MaxPosition().Y) >> shift

		var x, y int
		var keys []point

		for y = sy; y <= ey; y++ {
			for x = sx; x <= ex; x++ {
				keys = append(keys, point{x, y})
			}
		}

		return keys
	}
}

// Clear empties the spatial hashmap.
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

// InsertI allows from interfaces to be placed in.
func (s *SpatialHashmap) InsertI(objects ...interface{}) error {
	for i := range objects {
		switch t := objects[i].(type) {
		case Transformer:
			s.insertSingle(t)
		case SpatialAdder:
			t.Add(s)
		default:
			return fmt.Errorf("invalid insert: %v", t)
		}
	}

	return nil
}

// Insert loops through the given transformers and inserts them into the map.
func (s *SpatialHashmap) Insert(t ...Transformer) {
	for i := range t {
		s.insertSingle(t[i])
	}
}

func (s *SpatialHashmap) insertSingle(t Transformer) {
	// Get the hash key of the object.
	keys := s.getKeys(t)
	s.list = append(s.list, t)

	for i := 0; i < len(keys); i++ {
		key := keys[i]

		// If the key already exists, then append.
		// If not then create the key location.
		if s.hash[key] != nil {
			s.hash[key] = append(s.hash[key], t)
		} else {
			s.hash[key] = []Transformer{t}
		}
	}
}

func (s *SpatialHashmap) Remove(t Transformer) {
	keys := s.getKeys(t)
	idxOfT := -1

	for i := range s.list {
		if s.list[i] == t {
			idxOfT = i
			break
		}
	}

	for i := 0; i < len(keys); i++ {
		key := keys[i]

		if s.hash[key] != nil {
			for j := range s.hash[key] {
				if s.hash[key][j] == t {
					s.hash[key] = append(s.hash[key][:j], s.hash[key][j+1:]...)
					break
				}
			}
		}
	}

	if idxOfT == -1 {
		return
	}

	s.list = append(s.list[:idxOfT], s.list[idxOfT+1:]...)
}

// NumBuckets returns the number of key locations (buckets.)
func (s *SpatialHashmap) NumBuckets() int {
	var count int
	for key := range s.hash {
		if len(s.hash[key]) > 0 {
			count++
		}
	}

	return count
}

// Retrieve queries the spatial hashmap for nearby transforms to the given.
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
