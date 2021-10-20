package data

import (
	"github.com/faiface/pixel"
)

type Collider struct{
	Hitbox     pixel.Rect
	Damage     *Damage
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