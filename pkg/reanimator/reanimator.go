package reanimator

import (
	"dwarf-sweeper/pkg/img"
	"github.com/faiface/pixel"
)

type Tree struct {
	Root *Switch
}

type switchEl struct {
	Switch *Switch
	Anim   *Anim
}

type Switch struct {
	Elements []*switchEl
	index    int
	Check    func() int
}

type Anim struct {
	S      []*pixel.Sprite
	Step   int
	Finish Finish
	Tran   func()
}

type Finish int

const (
	Loop = iota
	Hold
	Tran
)

func SetSwitch(s *Switch, i int) {
	s.index = i
}

func NewElements(els ...*switchEl) []*switchEl {
	return els
}

func NewAnim(spriteSheet *img.SpriteSheet, rs []int, f Finish, tFn func()) *switchEl {
	var spr []*pixel.Sprite
	for _, r := range rs {
		spr = append(spr, pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[r]))
	}
	var tfn func()
	if f == Tran {
		tfn = tFn
	}
	return &switchEl{
		Switch: nil,
		Anim: &Anim{
			S:      spr,
			Step:   0,
			Finish: f,
			Tran:   tfn,
		},
	}
}

func NewSwitch() *switchEl {
	return &switchEl{
		Switch: nil,
		Anim: nil,
	}
}