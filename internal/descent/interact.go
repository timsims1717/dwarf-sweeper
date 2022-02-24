package descent

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Interact struct {
	OnInteract func(pixel.Vec, *Dwarf) bool
	Distance   float64
	Interacted bool
	Remove     bool
}

func NewInteract(fn func(pixel.Vec, *Dwarf) bool, dist float64, remove bool) *Interact {
	return &Interact{
		OnInteract: fn,
		Distance:   dist,
		Remove:     remove,
	}
}

type Collectible struct {
	OnCollect   func(pixel.Vec, *Dwarf) bool
	Sprite      *pixel.Sprite
	Collected   bool
	AutoCollect bool
	Timer       *timing.FrameTimer
}