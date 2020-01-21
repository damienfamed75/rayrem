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
	p.BasicEntity.Update(dt)

	if p.doubleJumpPerformed {
		p.doubleJumpPerformed = !p.Rigidbody.OnGround()
	}

	if p.Rigidbody.Velocity.X > p.friction {
		p.Rigidbody.Velocity.X -= p.friction
		p.Facing = common.Right
		p.Ase.Play("run")
	} else if p.Rigidbody.Velocity.X < -p.friction {
		p.Rigidbody.Velocity.X += p.friction
		p.Facing = common.Left
		p.Ase.Play("run")
	} else {
		p.Rigidbody.Velocity.X = 0
		p.Ase.Play("idle")
	}

	if r.IsKeyDown(r.KeyRight) {
		p.Rigidbody.Velocity.X++
	}

	if r.IsKeyDown(r.KeyLeft) {
		p.Rigidbody.Velocity.X--
	}

	if r.IsKeyPressed(r.KeyUp) {
		if p.Rigidbody.OnGround() {
			p.Rigidbody.Velocity.Y = -p.jumpHeight
		} else if !p.Rigidbody.OnGround() && !p.doubleJumpPerformed {
			p.Rigidbody.Velocity.Y = -p.jumpHeight
			p.doubleJumpPerformed = true
		}
	}
}
