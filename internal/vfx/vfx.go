package vfx

import (
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/reanimator"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type VFX struct {
	Animation *reanimator.Tree
	Matrix    pixel.Matrix
}

var (
	effects     []*VFX
	partBatcher *img.Batcher
)

func Initialize() {
	expSheet, err := img.LoadSpriteSheet("assets/img/explosion.json")
	if err != nil {
		panic(err)
	}
	explosion = reanimator.NewAnimFromSheet("explosion", expSheet, nil, reanimator.Done)
	particleSheet, err := img.LoadSpriteSheet("assets/img/particles.json")
	if err != nil {
		panic(err)
	}
	partBatcher = img.NewBatcher(particleSheet, false, false)
	dazed = reanimator.NewAnimFromSprites("dazed", partBatcher.Animations["dazed"].S, reanimator.Loop)
}

func Update() {
	var drop []int
	for i, effect := range effects {
		effect.Animation.Update()
		if effect.Animation.Done {
			drop = append(drop, i)
		}
	}
	for i := len(drop) - 1; i >= 0; i-- {
		effects = append(effects[:drop[i]], effects[drop[i]+1:]...)
	}
}

func Draw(win *pixelgl.Window) {
	for _, effect := range effects {
		effect.Animation.Draw(win, effect.Matrix)
	}
}

func Clear() {
	effects = []*VFX{}
}
