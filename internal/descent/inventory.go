package descent

import (
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/profile"
	"dwarf-sweeper/pkg/timing"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

func UpdateInventory(i *player.Inventory) {
	for _, item := range i.Items {
		if item.Using {
			if item.Timer == nil || item.Timer.UpdateDone() {
				item.Using = false
				item.Timer = nil
				uses, ok := profile.CurrentProfile.ItemLimits.Uses[item.Key]
				if !ok {
					uses = 1
				}
				if uses - item.Uses < 1 {
					item.Count--
					if item.Count < 1 {
						if len(i.Items) > 1 {
							i.Items = append(i.Items[:i.Index], i.Items[i.Index+1:]...)
						} else {
							i.Items = []*player.Item{}
						}
						PrevItem(i)
					}
				}
			}
		}
	}
}

func AddItem(i *player.Inventory, newItem *player.Item) int {
	found := false
	limit, ok := profile.CurrentProfile.ItemLimits.Hold[newItem.Key]
	if !ok {
		limit = 1
	}
	for _, item := range i.Items {
		if item.Key == newItem.Key {
			found = true
			break
		}
	}
	if !found {
		item := newItem.Copy()
		item.Count = 0
		i.Items = append(i.Items, item)
	}
	for _, item := range i.Items {
		if item.Key == newItem.Key {
			t := item.Count + newItem.Count
			r := t - limit
			if r < 0 {
				r = 0
			}
			item.Count = t
			if item.Count > limit {
				item.Count = limit
			}
			return r
		}
	}
	return 0
}

func UseEquipped(p *player.Player, e *ecs.Entity, dPos, tPos pixel.Vec) {
	i := p.Inventory
	if len(i.Items) > 0 && i.Index < len(i.Items) {
		item := i.Items[i.Index]
		sec, ok := profile.CurrentProfile.ItemLimits.Secs[item.Key]
		if !ok {
			sec = 0.
		}
		if !item.Using || item.Timer == nil || sec-item.Timer.Elapsed() < 0.1 {
			if item.OnUseFn(dPos, tPos, e, sec) {
				if sec > 0. {
					item.Timer = timing.New(sec)
				}
				item.Using = true
				item.Uses++
			}
		}
	} else {
		i.Index = 0
	}
}

func DropEquipped(i *player.Inventory, pos pixel.Vec) bool {
	if len(i.Items) > 0 && i.Index < len(i.Items) {
		item := i.Items[i.Index]
		if !item.Using || item.Timer == nil {
			item.Count--
			CreateItemPickUp(pos, item.Key, 1)
			if item.Count < 1 {
				if len(i.Items) > 1 {
					i.Items = append(i.Items[:i.Index], i.Items[i.Index+1:]...)
				} else {
					i.Items = []*player.Item{}
				}
				PrevItem(i)
				return true
			} else {
				return false
			}
		}
	}
	return true
}

func PrevItem(i *player.Inventory) {
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

func NextItem(i *player.Inventory) {
	if len(i.Items) > 0 {
		i.Index = (i.Index + 1) % len(i.Items)
	} else {
		i.Index = 0
	}
}
