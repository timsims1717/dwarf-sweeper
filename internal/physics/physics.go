package physics

import (
	"dwarf-sweeper/pkg/animation"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Physics struct {
	*animation.Transform
	Velocity pixel.Vec
}

func (p *Physics) Update() {
	p.Pos.X += p.Velocity.X
	p.Pos.Y += p.Velocity.Y
	p.Velocity.Y -= timing.DT * 5.0
	p.Transform.Update(pixel.Rect{})
}