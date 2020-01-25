// +build debug

package player

import (
	"fmt"

	r "github.com/lachee/raylib-goplus/raylib"
)

// Draw is used to debug the game.
// This function is only used when the game is ran in debug mode.
func (p *Player) Draw() {
	p.BasicEntity.Draw()

	r.DrawRectangleLines(
		int(p.Position().X), int(p.Position().Y),
		p.Ase.FrameBoundaries().Width, p.Ase.FrameBoundaries().Height,
		r.Red,
	)

	r.DrawText(
		fmt.Sprintf("pos[%.2f, %.2f]", p.Rigidbody.Position().X, p.Rigidbody.Position().Y),
		int(p.Position().X), int(p.Position().Y)+p.Ase.FrameBoundaries().Height,
		5, r.White,
	)

	potential := p.solids.Retrieve(p.Rigidbody)
	for _, obj := range potential {
		r.DrawRectangleLinesEx(r.NewRectangle(
			obj.Position().X, obj.Position().Y,
			obj.MaxPosition().X-obj.Position().X, obj.MaxPosition().Y-obj.Position().Y,
		), 1, r.Green)
	}
}
