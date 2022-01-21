package vfx

import (
	"dwarf-sweeper/pkg/reanimator"
	"github.com/faiface/pixel"
)

var explosion *reanimator.Anim

func CreateExplosion(vec pixel.Vec) {
	exp := &VFX{
		Animation: reanimator.NewSimple(explosion),
		Matrix:    pixel.IM.Moved(vec),
	}
	effects = append(effects, exp)
}
