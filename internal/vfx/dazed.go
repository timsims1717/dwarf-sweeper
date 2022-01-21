package vfx

import (
	"dwarf-sweeper/pkg/reanimator"
	"github.com/faiface/pixel"
)

var dazed *reanimator.Anim

func CreateDazed(vec pixel.Vec) *VFX {
	exp := &VFX{
		Animation: reanimator.NewSimple(dazed),
		Matrix:    pixel.IM.Moved(vec),
	}
	effects = append(effects, exp)
	return exp
}
