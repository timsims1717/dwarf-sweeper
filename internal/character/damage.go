package character

import (
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Health struct {
	Max          int
	Curr         int
	Dead         bool
	Dazed        bool
	DazedO       bool
	DazeOverride bool
	DazedTimer   *timing.FrameTimer
	DazedVFX     *vfx.VFX
	Inv          bool
	TempInv      bool
	TempInvTimer *timing.FrameTimer
	TempInvSec   float64
	Override     bool
}

func (h *Health) Delete() {
	if h == nil {
		return
	}
	if h.DazedVFX != nil {
		h.DazedVFX.Animation.Done = true
		h.DazedVFX = nil
	}
}

type Damage struct {
	Amount    int
	Dazed     float64
	Knockback float64
	Angle     *float64
	Source    pixel.Vec
	Override  bool
}

type AreaDamage struct {
	Area           []pixel.Vec
	Amount         int
	Dazed          float64
	Knockback      float64
	KnockbackDecay bool
	Source         pixel.Vec
	Override       bool
}