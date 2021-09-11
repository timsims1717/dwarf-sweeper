package data

import "github.com/faiface/pixel"

type Collider struct{
	GroundOnly bool
	CanPass    bool
}

type Collectible struct{
	OnCollect func(pos pixel.Vec) bool
	Sprite    *pixel.Sprite
	Collected bool
}