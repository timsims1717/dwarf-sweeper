package reanimator

import (
	"dwarf-sweeper/pkg/img"
	"github.com/faiface/pixel"
	"image/color"
	"time"
)

var (
	Timer       time.Time
	FRate       int
	inter       float64
	FrameSwitch bool
)

type Tree struct {
	Root    *Switch
	spr     *pixel.Sprite
	animKey string
	frame   int
	update  bool
	Done    bool
	Default string
}

func SetFrameRate(fRate int) {
	FRate = fRate
	inter = 1. / float64(fRate)
}

func Reset() {
	Timer = time.Now()
}

func Update() {
	FrameSwitch = time.Since(Timer).Seconds() > inter
	if FrameSwitch {
		Reset()
	}
}

func NewSimple(anim *Anim) *Tree {
	t := &Tree{
		Root: NewSwitch().
			AddAnimation(anim).
			SetChooseFn(func() int {
				return 0
			}),
	}
	t.Update()
	return t
}

func New(root *Switch, def string) *Tree {
	t := &Tree{
		Root:    root,
		update:  true,
		Default: def,
	}
	t.Update()
	return t
}

func (t *Tree) ForceUpdate() {
	t.update = true
}

func (t *Tree) Update() {
	if !t.Done {
		if FrameSwitch || t.update {
			t.update = false
			a := t.Root.choose()
			if a == nil {
				t.spr = nil
				t.animKey = ""
				t.frame = 0
			} else {
				pKey := t.animKey
				pFrame := t.frame
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
						fn(a, pKey, pFrame)
					}
				}
				t.spr = a.S[a.Step]
				if t.update {
					t.animKey = t.Default
					t.frame = a.Step
				} else {
					t.animKey = a.Key
					t.frame = a.Step
				}
			}
		}
	}
}

func (t *Tree) CurrentSprite() *pixel.Sprite {
	return t.spr
}

func (t *Tree) Draw(target pixel.Target, mat pixel.Matrix) {
	if t.spr != nil {
		t.spr.Draw(target, mat)
	}
}

func (t *Tree) DrawColorMask(target pixel.Target, mat pixel.Matrix, col color.RGBA) {
	if t.spr != nil {
		t.spr.DrawColorMask(target, mat, col)
	}
}

type switchEl struct {
	Switch *Switch
	Anim   *Anim
}

type Switch struct {
	Elements []*switchEl
	Choose   func() int
}

func NewSwitch() *Switch {
	return &Switch{}
}

func (s *Switch) AddNull() *Switch {
	s.Elements = append(s.Elements, &switchEl{})
	return s
}

func (s *Switch) AddAnimation(anim *Anim) *Switch {
	s.Elements = append(s.Elements, &switchEl{
		Anim: anim,
	})
	return s
}

func (s *Switch) AddSubSwitch(ss *Switch) *Switch {
	s.Elements = append(s.Elements, &switchEl{
		Switch: ss,
	})
	return s
}

func (s *Switch) SetChooseFn(fn func() int) *Switch {
	s.Choose = fn
	return s
}

func (s *Switch) choose() *Anim {
	el := s.Elements[s.Choose()]
	if el.Switch != nil {
		return el.Switch.choose()
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
	Triggers map[int]func(*Anim, string, int)
}

type Finish int

const (
	Loop = iota
	Hold
	Tran
	Done
)

func NewAnimFromSprite(key string, spr *pixel.Sprite, f Finish) *Anim {
	return &Anim{
		Key:    key,
		S:      []*pixel.Sprite{spr},
		Step:   0,
		Finish: f,
	}
}

func NewAnimFromSprites(key string, spr []*pixel.Sprite, f Finish) *Anim {
	return &Anim{
		Key:    key,
		S:      spr,
		Step:   0,
		Finish: f,
	}
}

func NewAnimFromSheet(key string, spriteSheet *img.SpriteSheet, rs []int, f Finish) *Anim {
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
	return &Anim{
		Key:    key,
		S:      spr,
		Step:   0,
		Finish: f,
	}
}

func (anim *Anim) SetTrigger(i int, fn func(*Anim, string, int)) *Anim {
	if anim.Triggers == nil {
		anim.Triggers = map[int]func(*Anim, string, int){}
	}
	anim.Triggers[i] = fn
	return anim
}

func (anim *Anim) Copy() *Anim {
	return &Anim{
		Key:      anim.Key,
		S:        anim.S,
		Step:     anim.Step,
		Finish:   anim.Finish,
		Triggers: anim.Triggers,
	}
}