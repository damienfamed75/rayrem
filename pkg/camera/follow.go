// +build !windows

// Until we find out that linux breaks, this version of the camera will build on
// anything that's NOT Windows.

package camera

import (
	"github.com/damienfamed75/rayrem/pkg/common"
	"github.com/damienfamed75/rayrem/pkg/physics"
	r "github.com/lachee/raylib-goplus/raylib"
)

// FollowCamera follows the player's position.
type FollowCamera struct {
	LerpAmount float32
	r.Camera2D
}

// NewFollow creates a default offset of the player's position.
func NewFollow(playerColl *physics.Space) *FollowCamera {
	return &FollowCamera{
		Camera2D: r.Camera2D{
			// Center the camera on the player's position.
			Offset: r.NewVector2(
				float32(r.GetScreenWidth()/2)-float32(playerColl.Width()/2),
				float32(r.GetScreenHeight()/2)-float32(playerColl.Height()/2),
			),
			Rotation: 0,
			Zoom:     common.Config.Camera.Zoom,
		},
		LerpAmount: common.Config.Camera.Lerp,
	}
}

// Update changes the offset position of the camera and the target.
func (e *FollowCamera) Update(curr r.Vector2) {
	// Update camera offset coordinates for it to move.
	xOff, yOff := -float32(curr.X+4)*e.Zoom, -float32(curr.Y+8)*e.Zoom
	e.Offset = r.NewVector2(
		xOff+float32(r.GetScreenWidth()),
		yOff+float32(r.GetScreenHeight()),
	)

	// Reset the camera's target to the player's current position.
	// Using a lerp to make the camera movement smoother.
	e.Target = e.Target.Lerp(curr, e.LerpAmount)
}
