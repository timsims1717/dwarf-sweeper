package descent

import "github.com/faiface/pixel"

var (
	Inventory []*InvItem
	InvIndex  = 0
)

type InvItem struct {
	Name   string
	Sprite *pixel.Sprite
	OnUse  func() bool
	Count  int
	Unique bool
}

func AddToInventory(item *InvItem) bool {
	for _, i := range Inventory {
		if i.Name == item.Name {
			if i.Unique {
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
		if item.OnUse() {
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