package reanimator

import (
	"dwarf-sweeper/pkg/img"
	"github.com/faiface/pixel"
	"time"
)

type Tree struct {
	Root    *Switch
	spr     *pixel.Sprite
	FRate   int
	inter   float64
	timer   time.Time
	animKey string
	update  bool
}

func New(root *Switch, fRate int) *Tree {
	inter := 1. / float64(fRate)
	return &Tree{
		Root:  root,
		FRate: fRate,
		inter: inter,
	}
}

func (t *Tree) Reset() {
	t.timer = time.Now()
}

func (t *Tree) Update() {
	a := t.Root.Choose()
	frameSwitch := time.Since(t.timer).Seconds() > t.inter
	if frameSwitch || t.update {
		t.update = false
		if a.Key != t.animKey {
			a.Step = 0
		} else {
			a.Step++
			if a.Step % len(a.S) == 0 {
				switch a.Finish {
				case Loop:
					a.Step = 0
				case Hold:
					a.Step = len(a.S)-1
				case Tran:
					a.Step = len(a.S)-1
					a.Tran()
					t.update = true
				}
			}
		}
		t.spr = a.S[a.Step]
		t.animKey = a.Key
	}
	if frameSwitch {
		t.Reset()
	}
}

func (t *Tree) CurrentSprite() *pixel.Sprite {
	return t.spr
}

type switchEl struct {
	Switch *Switch
	Anim   *Anim
}

type Switch struct {
	Elements []*switchEl
	Check    func() int
}

func (s *Switch) Choose() *Anim {
	el := s.Elements[s.Check()]
	if el.Switch != nil {
		return el.Switch.Choose()
	} else if el.Anim != nil {
		return el.Anim
	} else {
		return nil
	}
}

type Anim struct {
	Key    string
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

func NewElements(els ...*switchEl) []*switchEl {
	return els
}

func NewAnim(key string, spriteSheet *img.SpriteSheet, rs []int, f Finish, tFn func()) *switchEl {
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
			Key:    key,
			S:      spr,
			Step:   0,
			Finish: f,
			Tran:   tfn,
		},
	}
}

func NewSwitch(s *Switch) *switchEl {
	return &switchEl{
		Switch: s,
		Anim: nil,
	}
}