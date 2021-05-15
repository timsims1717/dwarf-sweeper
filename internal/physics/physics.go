package physics

import (
	"dwarf-sweeper/pkg/transform"
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math"
	"math/rand"
)

type Physics struct {
	*transform.Transform
	Velocity    pixel.Vec
	interX      *gween.Tween
	interY      *gween.Tween
	MovingX     bool
	MovingY     bool
	XJustSet    bool
	YJustSet    bool
	FrictionOff bool
	GravityOff  bool
	RicochetX   bool
	RicochetY   bool
	Grounded    bool
}

func (p *Physics) Update() {
	if p.interX != nil {
		vx, fin := p.interX.Update(timing.DT)
		p.Velocity.X = vx
		if fin {
			p.interX = nil
		}
	}
	if p.interY != nil {
		vy, fin := p.interY.Update(timing.DT)
		p.Velocity.Y = vy
		if fin {
			p.interY = nil
		}
	}
	p.Pos.X += timing.DT * p.Velocity.X
	p.Pos.Y += timing.DT * p.Velocity.Y
	if !p.GravityOff && !p.YJustSet {
		if p.Velocity.Y > -500. {
			p.Velocity.Y -= 750. * timing.DT
		}
		p.YJustSet = false
	}
	if !p.FrictionOff && !p.XJustSet {
		friction := 10.
		if p.Grounded {
			friction = 25.
		}
		if p.Velocity.X > 0. {
			p.Velocity.X -= friction * timing.DT
			if p.Velocity.X < 0. {
				p.Velocity.X = 0
			}
		} else if p.Velocity.X < 0. {
			p.Velocity.X += friction * timing.DT
			if p.Velocity.X > 0. {
				p.Velocity.X = 0
			}
		}
		p.XJustSet = false
	}
	p.MovingX = p.Velocity.X != 0.
	p.MovingY = p.Velocity.Y != 0.
	p.Transform.Update()
}

func (p *Physics) SetVelX(vx, dur float64) {
	diff := math.Abs(p.Velocity.X - vx)
	if diff != 0. {
		p.interX = gween.New(p.Velocity.X, vx, dur, ease.Linear)
	}
}

func (p *Physics) SetVelY(vy, dur float64) {
	diff := math.Abs(p.Velocity.Y - vy)
	if diff != 0. {
		p.interY = gween.New(p.Velocity.Y, vy, dur, ease.Linear)
	}
}

func (p *Physics) CancelMovement() {
	p.interX = nil
	p.Velocity = pixel.ZV
}

func RandomVelocity(orig pixel.Vec, variance float64) *Physics {
	tran := transform.NewTransform()
	physicsT := &Physics{Transform: tran}
	physicsT.Pos = orig
	actVar := variance * world.TileSize
	//if square {
	xVar := (rand.Float64() - 0.5) * actVar
	yVar := (rand.Float64() - 0.5) * actVar
	physicsT.Pos.X += xVar
	physicsT.Pos.Y += yVar
	physicsT.Velocity.X = xVar * 2.
	physicsT.Velocity.Y = 10.
	//}
	return physicsT
}