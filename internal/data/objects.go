package data

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type CollisionClass int

const (
	Critter = iota
	PlayerC
	Stacker
	ItemC
	GroundOnly
)

type Collider struct {
	Hitbox       pixel.Rect
	Class        CollisionClass
	NoClip       bool // ignores all tiles
	ThroughWalls bool // ignores RightBound and LeftBound tile colliders
	Fallthrough  bool // can fall through bridges
	Collided     bool // encountered another collider this frame
	Damage       *Damage

	BottomBound bool
	TopBound    bool
	RightBound  bool
	LeftBound   bool

	UL    bool
	UR    bool
	RU    bool
	RD    bool
	DL    bool
	DR    bool
	LU    bool
	LD    bool
	CUL   bool
	CUR   bool
	CDL   bool
	CDR   bool
	Debug bool
}

func NewCollider(hitbox pixel.Rect, class CollisionClass) *Collider {
	return &Collider{
		Hitbox: hitbox,
		Class:  class,
	}
}

type TimerFunc struct {
	Timer *timing.Timer
	Func  func() bool
}

func NewTimerFunc(fn func() bool, dur float64) *TimerFunc {
	return &TimerFunc{
		Timer: timing.New(dur),
		Func:  fn,
	}
}

type FrameFunc struct {
	Func func() bool
}

func NewFrameFunc(fn func() bool) *FrameFunc {
	return &FrameFunc{Func: fn}
}

type TriggerFunc struct {
	Func func(*Player) bool
}

func NewTriggerFunc(fn func(*Player) bool) *TriggerFunc {
	return &TriggerFunc{Func: fn}
}
