package vfx

import (
	"dwarf-sweeper/pkg/img"
	"github.com/faiface/pixel"
)

var explosion *img.Animation

func CreateExplosion(vec pixel.Vec) {
	exp := explosion.NewInstance()
	exp.Matrix = pixel.IM.Moved(vec)
	effects = append(effects, exp)
}