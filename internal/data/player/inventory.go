package player

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

type Inventory struct {
	Items []*InvItem
	Index int
}

type InvItem struct {
	Name   string
	Sprite *pixel.Sprite
	OnUse  func(pos pixel.Vec)
	Count  int
	Limit  int
	Using  bool
	Sec    float64
	Timer  *timing.Timer
}

func (i *Inventory) Update() {
	for _, item := range i.Items {
		if item.Using {
			if item.Timer == nil || item.Timer.UpdateDone() {
				item.Using = false
				item.Timer = nil
				item.Count--
				if item.Count < 1 {
					if len(i.Items) > 1 {
						i.Items = append(i.Items[:i.Index], i.Items[i.Index+1:]...)
					} else {
						i.Items = []*InvItem{}
					}
					i.PrevItem()
				}
			}
		}
	}
}

func (i *Inventory) AddItem(newItem *InvItem) bool {
	for _, item := range i.Items {
		if item.Name == newItem.Name {
			if item.Count >= item.Limit {
				return false
			} else {
				item.Count += newItem.Count
				return true
			}
		}
	}
	i.Items = append(i.Items, newItem)
	return true
}

func (i *Inventory) UseEquipped(pos pixel.Vec) {
	if len(i.Items) > 0 && i.Index < len(i.Items) {
		item := i.Items[i.Index]
		if !item.Using || item.Timer == nil || item.Sec-item.Timer.Elapsed() < 0.1 {
			item.OnUse(pos)
			if item.Sec > 0. {
				item.Timer = timing.New(item.Sec)
			}
			item.Using = true
		}
	} else {
		i.Index = 0
	}
}

func (i *Inventory) PrevItem() {
	if len(i.Items) > 0 {
		newInv := i.Index - 1
		if newInv < 0 {
			newInv += len(i.Items)
		}
		i.Index = newInv % len(i.Items)
	} else {
		i.Index = 0
	}
}

func (i *Inventory) NextItem() {
	if len(i.Items) > 0 {
		i.Index = (i.Index + 1) % len(i.Items)
	} else {
		i.Index = 0
	}
}
