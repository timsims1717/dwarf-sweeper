package vfx

import (
	"dwarf-sweeper/pkg/animation"
	"dwarf-sweeper/pkg/img"
	"github.com/faiface/pixel/pixelgl"
)

var effects []*animation.AnimationInstance

func Initialize() {
	expSheet, err := img.LoadSpriteSheet("assets/img/explosion.json")
	if err != nil {
		panic(err)
	}
	explosion = animation.NewAnimation(expSheet, false, false, 0.5)
}

func Update() {
	var drop []int
	for i, effect := range effects {
		effect.Update()
		if effect.Done {
			drop = append(drop, i)
		}
	}
	for i := len(drop)-1; i >= 0; i-- {
		effects = append(effects[:drop[i]], effects[drop[i]+1:]...)
	}
}

func Draw(win *pixelgl.Window) {
	for _, effect := range effects {
		effect.Draw(win)
	}
}