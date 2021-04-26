package physics

import (
	"dwarf-sweeper/pkg/animation"
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
	"math"
)

type Physics struct {
	*animation.Transform
	Velocity pixel.Vec
	interX   *gween.Tween
	Off      bool
}

func (p *Physics) Update() {
	if p.interX != nil {
		vx, fin := p.interX.Update(timing.DT)
		p.Velocity.X = vx
		if fin {
			p.interX = nil
		}
	}
	p.Pos.X += timing.DT * p.Velocity.X
	p.Pos.Y += timing.DT * p.Velocity.Y
	if !p.Off {
		if p.Velocity.Y > -500. {
			p.Velocity.Y -= 5.
		}
		if p.Velocity.X > 75. {
			p.Velocity.X -= 10. * timing.DT
		} else if p.Velocity.X < -75. {
			p.Velocity.X += 10. * timing.DT
		}
	}
	p.Transform.Update(pixel.Rect{})
}

func (p *Physics) SetVelX(vx, dur float64) {
	diff := math.Abs(p.Velocity.X - vx)
	if diff != 0. {
		p.interX = gween.New(p.Velocity.X, vx, dur, ease.Linear)
	}
}

func (p *Physics) CancelMovement() {
	p.interX = nil
	p.Velocity = pixel.ZV
}