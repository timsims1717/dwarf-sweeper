package physics

import (
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Physics struct {
	Velocity    pixel.Vec
	interX      *gween.Tween
	interY      *gween.Tween
	XJustSet    bool
	YJustSet    bool
	FrictionOff bool
	GravityOff  bool
	RicochetX   bool
	RicochetY   bool
	Grounded    bool
	Gravity     float64
	Terminal    float64
	Friction    float64
	AirFriction float64
}

func New() *Physics {
	return &Physics{
		Gravity:     750.,
		Terminal:    500.,
		Friction:    400.,
		AirFriction: 25.,
	}
}

func (p *Physics) IsMovingX() bool {
	return p.Velocity.X > 0.01 || p.Velocity.X < -0.01
}

func (p *Physics) IsMovingY() bool {
	return p.Velocity.Y > 0.01 || p.Velocity.Y < -0.01
}

func (p *Physics) SetVelX(vx, spd float64) {
	if spd == 0. {
		p.Velocity.X = vx
	} else {
		p.Velocity.X += spd * timing.DT * (vx - p.Velocity.X)
	}
	p.XJustSet = true
}

func (p *Physics) SetVelY(vy, spd float64) {
	if spd == 0. {
		p.Velocity.Y = vy
	} else {
		p.Velocity.Y += spd * timing.DT * (vy - p.Velocity.Y)
	}
	p.YJustSet = true
}

func (p *Physics) CancelMovement() {
	p.interX = nil
	p.Velocity = pixel.ZV
}