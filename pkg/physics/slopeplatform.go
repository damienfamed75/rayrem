package physics

import r "github.com/lachee/raylib-goplus/raylib"

// SlopePlatform is a fully functioning sloped platform for use in game.
//
// It includes two landing platforms and the slope itself. The reasoning behind
// having landing zones is because the player would often get stuck on the apex
// of the slope and bang against the next wall on the X-axis. To solve this
// issue we use a landing zone on the upper portion to raise the player high
// enough to land properly on the next platform.
//
// The second landing zone being on the bottom of the slope to catch the player
// when they're walking on the ground. Since they don't have an immediate
// downward Y velocity, the slope won't register the player right away unless
// there was this bottom landing zone.
//
//           ___ landingzone2
//         /
//       /  slope
//     /
// --- landingzone1
//
type SlopePlatform struct {
	landingZone1 *Slope // hangs left
	landingZone2 *Slope // hangs right
	slope        *Slope

	*Space
}

// NewSlopePlatform expects p1 to be the left side of the slope
// and p2 to be the right side.
// landingWidth is the X-width of each landing platform.
func NewSlopePlatform(p1, p2 r.Vector2, landingWidth float32) *SlopePlatform {
	sp := &SlopePlatform{
		slope: NewSlope(p1, p2),
		Space: NewSpace(),
	}

	sp.landingZone1 = NewSlope(r.NewVector2(p1.X-landingWidth, p1.Y), p1)
	sp.landingZone2 = NewSlope(p2, r.NewVector2(p2.X+landingWidth, p2.Y))

	sp.Add(sp.landingZone1, sp.landingZone2, sp.slope)

	return sp
}

// LandingZones returns references to each landing zone.
func (sp *SlopePlatform) LandingZones() (*Slope, *Slope) {
	return sp.landingZone1, sp.landingZone2
}

// Slope returns a reference to the slope alone.
func (sp *SlopePlatform) Slope() *Slope {
	return sp.slope
}

// Draw is used for debugging when needing to draw the whole platform.
func (sp *SlopePlatform) Draw() {
	// Draw each landing zone.
	r.DrawLineEx(sp.landingZone1.p1, sp.landingZone1.p2, 1, r.Green)
	r.DrawLineEx(sp.landingZone2.p1, sp.landingZone2.p2, 1, r.Green)
	// Draw the slope using its own draw function. (which will draw a triangle.)
	sp.slope.Draw()
}
