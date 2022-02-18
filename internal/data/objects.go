package data

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Collider struct {
	Hitbox       pixel.Rect
	GroundOnly   bool
	ThroughWalls bool
	CanPass      bool
	Collided     bool
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

func NewCollider(hitbox pixel.Rect, groundOnly, canPass bool) *Collider {
	return &Collider{
		Hitbox:     hitbox,
		GroundOnly: groundOnly,
		CanPass:    canPass,
	}
}

type Collectible struct {
	OnCollect   func(pos pixel.Vec) bool
	Sprite      *pixel.Sprite
	Collected   bool
	AutoCollect bool
	Timer       *timing.FrameTimer
}

type Interact struct {
	OnInteract func(pos pixel.Vec) bool
	Distance   float64
	Interacted bool
	Remove     bool
}

func NewInteract(fn func(pos pixel.Vec) bool, dist float64, remove bool) *Interact {
	return &Interact{
		OnInteract: fn,
		Distance:   dist,
		Remove:     remove,
	}
}

type TimerFunc struct {
	Timer *timing.FrameTimer
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
