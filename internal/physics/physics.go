package physics

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Physics struct {
	Velocity pixel.Vec

	XJustSet    bool
	YJustSet    bool
	RagDollX    bool
	RagDollY    bool
	NearGround  bool
	Grounded    bool
	TopBound    bool
	BottomBound bool
	RightBound  bool
	LeftBound   bool
	WallBound   bool

	// the "Constants"
	Gravity      float64
	Terminal     float64
	Friction     float64
	AirFriction  float64
	Bounciness   float64
	Mass         float64
	WallFriction float64

	// the "settings"

	FrictionOff     bool
	GravityOff      bool
	UseWallFriction bool
}

func New() *Physics {
	return &Physics{
		Gravity:      750.,
		Terminal:     500.,
		Friction:     400.,
		WallFriction: 400.,
		AirFriction:  25.,
		Bounciness:   0.6,
		Mass:         1.0,
	}
}

func (p *Physics) IsMovingX() bool {
	return p.Velocity.X > 0.1 || p.Velocity.X < -0.1
}

func (p *Physics) IsMovingY() bool {
	return p.Velocity.Y > 0.1 || p.Velocity.Y < -0.1
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
	p.Velocity = pixel.ZV
}
