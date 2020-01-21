package common

// Tags are used to tell what is what when using the world physics space.
// This makes it easy to filter through the objects in the world.
// All tags should be placed in here so it's easy to get to all of them.
const (
	TagCollision   = "collision"
	TagHitbox      = "hitbox"
	TagPhysicsBody = "pbody"
	TagGround      = "ground"
	TagPlayer      = "player"
)
