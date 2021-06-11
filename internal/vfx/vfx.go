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

var effects []*VFX

func Initialize() {
	expSheet, err := img.LoadSpriteSheet("assets/img/explosion.json")
	if err != nil {
		panic(err)
	}
	explosion = reanimator.NewAnimFromSheet("explosion", expSheet, nil, reanimator.Done, nil).Anim
}

func Update() {
	var drop []int
	for i, effect := range effects {
		effect.Animation.Update()
		if effect.Animation.Done {
			drop = append(drop, i)
		}
	}
	for i := len(drop)-1; i >= 0; i-- {
		effects = append(effects[:drop[i]], effects[drop[i]+1:]...)
	}
}

func Draw(win *pixelgl.Window) {
	for _, effect := range effects {
		effect.Animation.CurrentSprite().Draw(win, effect.Matrix)
	}
}

func Clear() {
	effects = []*VFX{}
}