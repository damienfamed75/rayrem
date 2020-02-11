// +build debug

package player

import (
	"fmt"

	"github.com/damienfamed75/rayrem/pkg/physics"

	r "github.com/lachee/raylib-goplus/raylib"
)

// Draw is used to debug the game.
// This function is only used when the game is ran in debug mode.
func (p *Player) Draw() {
	p.Actor.Draw()

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

	for i := range *p.Rigidbody.Space {
		collider := (*p.Rigidbody.Space)[i].(*physics.Rectangle).Rectangle
		possible := p.solids.Retrieve(collider.Move(
			p.Velocity().X, p.Velocity().Y,
		))

		var tpossible []physics.Transformer
		for _, pp := range possible {
			if tt, ok := pp.(physics.Transformer); ok {
				tpossible = append(tpossible, tt)
			}
		}

		p.debugPotential(tpossible)
	}
}

func (p *Player) debugPotential(potential []physics.Transformer) {
	for _, obj := range potential {

		if s, ok := obj.(*physics.Space); ok {
			var tt []physics.Transformer
			for i := range *s {
				tt = append(tt, (*s)[i])
			}

			p.debugPotential(tt)

			r.DrawRectangleLinesEx(r.NewRectangle(
				obj.Position().X, obj.Position().Y,
				obj.MaxPosition().X-obj.Position().X, obj.MaxPosition().Y-obj.Position().Y,
			), 1, r.Blue)
		} else {
			r.DrawRectangleLinesEx(r.NewRectangle(
				obj.Position().X, obj.Position().Y,
				obj.MaxPosition().X-obj.Position().X, obj.MaxPosition().Y-obj.Position().Y,
			), 1, r.Green)
		}

	}
}
