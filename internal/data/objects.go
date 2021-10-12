package data

import (
	"github.com/faiface/pixel"
)

type Collider struct{
	Hitbox     pixel.Rect
	GroundOnly bool
	CanPass    bool
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