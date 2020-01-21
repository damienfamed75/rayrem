// +build debug

package player

import (
	"fmt"
	"sort"

	"github.com/damienfamed75/rayrem/pkg/common"
	"github.com/damienfamed75/rayrem/pkg/physics"

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

	ground := p.world.FilterByTags(common.TagGround)
	var colx, coly bool

	for i := range *p.Rigidbody.Space {
		collider := (*p.Rigidbody.Space)[i].(*physics.Rectangle).Rectangle
		tmpXRec := collider.Move(p.Rigidbody.Velocity.X+1, 0)
		tmpYRec := collider.Move(0, p.Rigidbody.Velocity.Y+1)

		for _, g := range *((*ground)[0].(*physics.Space)) {
			switch t := g.(type) {
			case *physics.Rectangle:
				// If the player hasn't collided with anything on the x-axis yet.
				if !colx {
					if g.Overlaps(tmpXRec) {
						colx = true
						r.DrawRectangleLinesEx(t.Rectangle, 1, r.Green)
					}
				}

				// If the player hasn't collided with anything on the y-axis yet.
				if !coly {
					if g.Overlaps(tmpYRec) {
						coly = true
						r.DrawRectangleLinesEx(t.Rectangle, 1, r.Green)
					}
				}
			case *physics.Slope:
				if g.Overlaps(tmpYRec) {
					r.DrawRectangleLinesEx(
						r.NewRectangle(collider.X, collider.MaxPosition().Y, collider.Width, collider.Height/2),
						1, r.Green,
					)
					intersections := t.GetIntersectionPoints(physics.NewRectangle(tmpYRec.X, tmpYRec.Y, tmpYRec.Width, tmpYRec.Height))
					if len(intersections) == 0 {
						continue
					}

					sort.Slice(intersections, func(i, j int) bool {
						return intersections[i].Y > intersections[j].Y
					})

					tmpL := physics.NewSlope(
						r.NewVector2(intersections[0].X, intersections[0].Y),
						r.NewVector2(intersections[len(intersections)-1].X, intersections[len(intersections)-1].Y),
					)

					_, dy := tmpL.Delta()

					colBox := r.NewRectangle(collider.X, intersections[0].Y+(dy/2), collider.Width, collider.Height/2)

					r.DrawRectangleLinesEx(
						colBox,
						1, r.Purple,
					)

					r.DrawRectangleLinesEx(
						colBox.GetOverlapRec(tmpYRec),
						1, r.White,
					)
				}
			}
		}
	}
}
