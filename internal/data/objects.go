package data

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Collider struct{
	Hitbox     pixel.Rect
	GroundOnly bool
	CanPass    bool
	UL bool
	UR bool
	RU bool
	RD bool
	DL bool
	DR bool
	LU bool
	LD bool
	CUL bool
	CUR bool
	CDL bool
	CDR bool
}

func NewCollider(hitbox pixel.Rect, groundOnly, canPass bool) *Collider {
	return &Collider{
		Hitbox:     hitbox,
		GroundOnly: groundOnly,
		CanPass:    canPass,
	}
}

type Collectible struct{
	OnCollect   func(pos pixel.Vec) bool
	Sprite      *pixel.Sprite
	Collected   bool
	AutoCollect bool
}

type Interact struct {
	OnInteract func(pos pixel.Vec) bool
	Distance   float64
	//Timer      *timing.FrameTimer
	//Sec        float64
	Interacted bool
	Remove     bool
}

type TimerFunc struct {
	Timer *timing.FrameTimer
	Func  func()
}

type FrameFunc struct {
	Func func() bool
}