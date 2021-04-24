package vfx

import (
	"dwarf-sweeper/pkg/animation"
	"github.com/faiface/pixel"
)

var explosion *animation.Animation

func CreateExplosion(vec pixel.Vec) {
	exp := explosion.NewInstance()
	exp.Matrix = pixel.IM.Moved(vec)
	effects = append(effects, exp)
}