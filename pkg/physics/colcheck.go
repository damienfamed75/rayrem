package physics

// colCheck is used for collision checking every frame in the physics Body.
// indexes 0-X, 1-Y are by default false.
type colCheck [2]bool

const (
	colXIdx = iota
	colYIdx
)

func (c *colCheck) X() bool {
	return c[0]
}
func (c *colCheck) Y() bool {
	return c[1]
}
func (c *colCheck) SetX(val bool) {
	c[0] = val
}
func (c *colCheck) SetY(val bool) {
	c[1] = val
}
