package common

// Tag is a representation for spaces and physics objects to be classified as.
// These are useful for getting more information about a collision box or space.
type Tag uint64

// Tags are used to tell what is what when using the world physics space.
// This makes it easy to filter through the objects in the world.
// All tags should be placed in here so it's easy to get to all of them.
const (
	TagCollision = iota + 1
	TagHitbox
	TagPhysicsBody
	TagGround
	TagPlayer
)
