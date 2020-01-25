package physics

type point [2]int

func (p point) X() int {
	return p[0]
}

func (p point) Y() int {
	return p[1]
}
