package player

import (
	"github.com/damienfamed75/rayrem/pkg/common"

	r "github.com/lachee/raylib-goplus/raylib"
)

// SetPosition is here so then throughout different scenes, the player can just
// be moved around instead of remade each time.
func (p *Player) SetPosition(x, y float32) {
	p.Rigidbody.SetPosition(x, y)
}

// TakeDamage doesn't do anything by default.
// TODO
func (p *Player) TakeDamage() {}

// Update updates the default basic entity and checks for movement and sends it
// to the Rigidbody.
func (p *Player) Update(dt float32) {
	p.Actor.Update(dt)

	if p.doubleJumpPerformed {
		p.doubleJumpPerformed = !p.Rigidbody.OnGround()
	}

	if p.Velocity().X > p.friction {
		p.AddVelocity(-p.friction, 0)
		p.Facing = common.Right
		p.Ase.Play("run")
	} else if p.Velocity().X < -p.friction {
		p.AddVelocity(+p.friction, 0)
		p.Facing = common.Left
		p.Ase.Play("run")
	} else {
		p.SetVelocity(0, p.Velocity().Y)
		p.Ase.Play("idle")
	}

	if r.IsKeyDown(r.KeyRight) {
		p.AddVelocity(1, 0)
	}

	if r.IsKeyDown(r.KeyLeft) {
		p.AddVelocity(-1, 0)
	}

	if r.IsKeyPressed(r.KeyUp) {
		if p.Rigidbody.OnGround() {
			p.SetVelocity(p.Velocity().X, -p.jumpHeight)
		} else if !p.Rigidbody.OnGround() && !p.doubleJumpPerformed {
			p.SetVelocity(p.Velocity().X, -p.jumpHeight)
			p.doubleJumpPerformed = true
		}
	}
}
