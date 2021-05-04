package img

import (
	gween "dwarf-sweeper/pkg/gween64"
	"dwarf-sweeper/pkg/gween64/ease"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Instance struct {
	*Animation
	inter  *gween.Tween
	Matrix pixel.Matrix
	Done   bool
	draw   bool
	step   int
}

type Animation struct {
	Loop  bool
	Hold  bool
	S     []*pixel.Sprite
	dur   float64
}

func NewAnimation(spriteSheet *SpriteSheet, a []pixel.Rect, loop, hold bool, dur float64) *Animation {
	var spr []*pixel.Sprite
	for _, r := range a {
		spr = append(spr, pixel.NewSprite(spriteSheet.Img, r))
	}
	return &Animation{
		Loop:  loop,
		Hold:  hold,
		S:     spr,
		dur:   dur,
	}
}

func (a *Animation) NewInstance() *Instance {
	return &Instance{
		Animation: a,
		inter:     gween.New(0., float64(len(a.S)), a.dur, ease.Linear),
		draw:      true,
	}
}

func (a *Instance) Update() {
	if !a.Done {
		var step float64
		step, a.Done = a.inter.Update(timing.DT)
		a.step = int(step)
		a.Done = a.Done || a.step >= len(a.S)
		if a.Done {
			if a.Loop {
				a.Done = false
				a.inter = gween.New(0., float64(len(a.S)), a.dur, ease.Linear)
				a.step = 0
			} else if a.Hold {
				a.step = len(a.S) - 1
				a.inter = nil
			} else {
				a.step = 0
				a.inter = nil
				a.draw = false
			}
		}
	}
}

func (a *Instance) Draw(target pixel.Target) {
	if a.draw {
		a.S[a.step].Draw(target, a.Matrix)
	}
}

func (a *Instance) SetMatrix(mat pixel.Matrix) {
	a.Matrix = mat
}

func (a *Instance) Reset() {
	a.Done = false
	a.inter = gween.New(0., float64(len(a.S)), a.dur, ease.Linear)
	a.step = 0
}