package data

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type Inventory struct {
	Items []*Item
	Index int
}

type Item struct {
	Key    string
	Name   string
	Sprite *pixel.Sprite

	Temp   bool
	DeadFn func(pixel.Vec)

	OnUseFn func(pixel.Vec, pixel.Vec, *ecs.Entity, float64) bool
	Using   bool
	Count   int
	Uses    int
	Timer   *timing.Timer
}

func (i *Item) Copy() *Item {
	return &Item{
		Key:     i.Key,
		Name:    i.Name,
		Sprite:  i.Sprite,
		Temp:    i.Temp,
		DeadFn:  i.DeadFn,
		OnUseFn: i.OnUseFn,
	}
}