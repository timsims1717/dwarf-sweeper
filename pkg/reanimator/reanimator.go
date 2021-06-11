package reanimator

import (
	"dwarf-sweeper/pkg/img"
	"github.com/faiface/pixel"
	"time"
)

var (
	Timer       time.Time
	FRate       int
	inter       float64
	frameSwitch bool
)

type Tree struct {
	Root    *Switch
	spr     *pixel.Sprite
	animKey string
	update  bool
	Done    bool
}

func SetFrameRate(fRate int) {
	FRate = fRate
	inter = 1. / float64(fRate)
}

func Reset() {
	Timer = time.Now()
}

func Update() {
	frameSwitch = time.Since(Timer).Seconds() > inter
	if frameSwitch {
		Reset()
	}
}

func NewSimple(anim *Anim) *Tree {
	t := &Tree{
		Root: &Switch{
			Elements: NewElements(
				&switchEl{
					Switch: nil,
					Anim: anim,
				},
			),
			Check:    func() int {
				return 0
			},
		},
	}
	t.Update()
	return t
}

func New(root *Switch) *Tree {
	t := &Tree{
		Root:   root,
		update: true,
	}
	t.Update()
	return t
}

func (t *Tree) ForceUpdate() {
	t.update = true
}

func (t *Tree) Update() {
	if !t.Done {
		a := t.Root.Choose()
		if frameSwitch || t.update {
			t.update = false
			var trigger int
			if a.Key != t.animKey {
				a.Step = 0
				trigger = 0
			} else {
				a.Step++
				trigger = a.Step
				if a.Step%len(a.S) == 0 {
					switch a.Finish {
					case Loop:
						a.Step = 0
					case Hold:
						a.Step = len(a.S) - 1
					case Tran:
						a.Step = len(a.S) - 1
						t.update = true
					case Done:
						a.Step = len(a.S) - 1
						t.Done = true
					}
				}
			}
			if a.Triggers != nil {
				if fn, ok := a.Triggers[trigger]; ok {
					fn()
				}
			}
			t.spr = a.S[a.Step]
			t.animKey = a.Key
		}
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
	Key      string
	S        []*pixel.Sprite
	Step     int
	Finish   Finish
	Triggers map[int]func()
}

type Finish int

const (
	Loop = iota
	Hold
	Tran
	Done
)

func NewElements(els ...*switchEl) []*switchEl {
	return els
}

func NewAnimFromSprites(key string, spr []*pixel.Sprite, f Finish, triggers map[int]func()) *switchEl {
	return &switchEl{
		Switch: nil,
		Anim: &Anim{
			Key:      key,
			S:        spr,
			Step:     0,
			Finish:   f,
			Triggers: triggers,
		},
	}
}

func NewAnimFromSheet(key string, spriteSheet *img.SpriteSheet, rs []int, f Finish, triggers map[int]func()) *switchEl {
	var spr []*pixel.Sprite
	if len(rs) > 0 {
		for _, r := range rs {
			spr = append(spr, pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[r]))
		}
	} else {
		for _, s := range spriteSheet.Sprites {
			spr = append(spr, pixel.NewSprite(spriteSheet.Img, s))
		}
	}
	return &switchEl{
		Switch: nil,
		Anim: &Anim{
			Key:      key,
			S:        spr,
			Step:     0,
			Finish:   f,
			Triggers: triggers,
		},
	}
}

func NewSwitch(s *Switch) *switchEl {
	return &switchEl{
		Switch: s,
		Anim: nil,
	}
}