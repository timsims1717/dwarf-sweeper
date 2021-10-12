package descent

import (
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
)

var (
	Inventory []*InvItem
	InvIndex  = 0
)

type InvItem struct {
	Name   string
	Sprite *pixel.Sprite
	OnUse  func()
	Count  int
	Limit  int
	Using  bool
	Sec    float64
	Timer  *timing.FrameTimer
}

func UpdateInventory() {
	for _, item := range Inventory {
		if item.Using {
			if item.Timer == nil || item.Timer.UpdateDone() {
				item.Using = false
				item.Timer = nil
				item.Count--
				if item.Count < 1 {
					if len(Inventory) > 1 {
						Inventory = append(Inventory[:InvIndex], Inventory[InvIndex+1:]...)
					} else {
						Inventory = []*InvItem{}
					}
					PrevItem()
				}
			}
		}
	}
}

func AddToInventory(item *InvItem) bool {
	for _, i := range Inventory {
		if i.Name == item.Name {
			if i.Count >= i.Limit {
				return false
			} else {
				i.Count += item.Count
				return true
			}
		}
	}
	Inventory = append(Inventory, item)
	return true
}

func UseEquipped() {
	if len(Inventory) > 0 && InvIndex < len(Inventory) {
		item := Inventory[InvIndex]
		if !item.Using || item.Timer == nil || item.Sec - item.Timer.Elapsed() < 0.1 {
			item.OnUse()
			if item.Sec > 0. {
				item.Timer = timing.New(item.Sec)
			}
			item.Using = true
		}
	} else {
		InvIndex = 0
	}
}

func PrevItem() {
	if len(Inventory) > 0 {
		newInv := InvIndex - 1
		if newInv < 0 {
			newInv += len(Inventory)
		}
		InvIndex = newInv% len(Inventory)
	} else {
		InvIndex = 0
	}
}

func NextItem() {
	if len(Inventory) > 0 {
		InvIndex = (InvIndex + 1) % len(Inventory)
	} else {
		InvIndex = 0
	}
}