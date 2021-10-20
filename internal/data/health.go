package data

import (
	"dwarf-sweeper/internal/vfx"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type BlastHealth struct {
	Dead bool
}

type Health struct {
	Max          int
	Curr         int
	TempHP       int
	TempHPTimer  *timing.FrameTimer
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
	Immune       []DamageType
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
type DamageType int

const (
	Blast = iota
	Shovel
	Enemy
)

type Damage struct {
	Amount    int
	Dazed     float64
	Knockback float64
	Angle     *float64
	Source    pixel.Vec
	Type      DamageType
}

type AreaDamage struct {
	Amount         int
	Dazed          float64
	Knockback      float64
	Type           DamageType
	Source         pixel.Vec
	Center         pixel.Vec
	Radius         float64
	Rect           pixel.Rect
	KnockbackDecay bool
}

type Heal struct {
	Amount    int
	TmpAmount int
}

type TempHP struct {
	Amount int
	Timer  *timing.FrameTimer
}