package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/profile"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"dwarf-sweeper/pkg/world"
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
	"math"
)

func InitItems() {
	bombItem.Sprite = img.Batchers[constants.EntityKey].Sprites["bomb_item"]
	beerItem.Sprite = img.Batchers[constants.EntityKey].Sprites["beer"]
	xrayItem.Sprite = img.Batchers[constants.EntityKey].Sprites["x-ray-helmet"]
	throwShovelItem.Sprite = img.Batchers[constants.EntityKey].Sprites["throwing_shovel_item"]
	throwShovelSpr = img.Batchers[constants.EntityKey].Sprites["throwing_shovel"]
	pickaxeItem.Sprite = img.Batchers[constants.EntityKey].Sprites["pickaxe"]
	detectorItem.Sprite = img.Batchers[constants.EntityKey].Sprites["detector"]
}

var (
	throwShovelSpr *pixel.Sprite
	items = map[string]*player.Item{
		"bomb_item":    bombItem,
		"beer":         beerItem,
		"xray":         xrayItem,
		"throw_shovel": throwShovelItem,
		"pickaxe":      pickaxeItem,
		"detector":     detectorItem,
	}
	bombItem = &player.Item{
		Key:     "bomb_item",
		Name:    "Bomb",
		Temp:    true,
		DeadFn:  CreateBomb,
		OnUseFn: func(pos, _ pixel.Vec, _ *ecs.Entity, _ float64) bool {
			tile := Descent.Cave.GetTile(pos)
			CreateBomb(tile.Transform.Pos)
			return true
		},
	}
	beerItem = &player.Item{
		Key:     "beer",
		Name:    "Beer",
		Temp:    true,
		OnUseFn: func(_, _ pixel.Vec, e *ecs.Entity, _ float64) bool {
			e.AddComponent(myecs.Healing, &data.Heal{
				TmpAmount: 2,
			})
			return true
		},
	}
	pickaxeItem = &player.Item{
		Key:  "pickaxe",
		Name: "Pickaxe",
		Temp: true,
		OnUseFn: func(dPos, tPos pixel.Vec, e *ecs.Entity, _ float64) bool {
			tTile := Descent.GetTile(tPos)
			if tTile != nil {
				digLegal := math.Abs(dPos.X-tTile.Transform.Pos.X) < world.TileSize*DigRange &&
					math.Abs(dPos.Y-tTile.Transform.Pos.Y) < world.TileSize*DigRange
				if digLegal && tTile.Solid() {
					var p *player.Player
					if pl, ok := e.GetComponentData(myecs.Player); ok {
						p, _ = pl.(*player.Player)
					}
					if p != nil {
						profile.CurrentProfile.Stats.BlocksDug++
						p.Stats.BlocksDug++
					}
					tTile.Destroy(p, true)
					return true
				}
			}
			return false
		},
	}
	detectorItem = &player.Item{
		Key:     "detector",
		Name:    "Metal Detector",
		Temp:    true,
		OnUseFn: func(dPos, tPos pixel.Vec, e *ecs.Entity, _ float64) bool {
			tTile := Descent.GetTile(tPos)
			if tTile != nil {
				digLegal := math.Abs(dPos.X-tTile.Transform.Pos.X) < world.TileSize*DigRange &&
					math.Abs(dPos.Y-tTile.Transform.Pos.Y) < world.TileSize*DigRange
				if digLegal && tTile.Solid() {
					var p *player.Player
					if pl, ok := e.GetComponentData(myecs.Player); ok {
						p, _ = pl.(*player.Player)
					}
					if p != nil && tTile.Bomb {
						FlagTile(p, tTile)
					}
					return true
				}
			}
			return false
		},
	}
)

func CreateInvItem(inv *player.Inventory, key string, count int) {
	_, ok := items[key]
	if !ok {
		fmt.Printf("error: no item key named '%s'\n", key)
		return
	}
	item := items[key].Copy()
	item.Count = count
	AddItem(inv, item)
}

func CreateItemPickUp(pos pixel.Vec, key string, count int) {
	_, ok := items[key]
	if !ok {
		fmt.Printf("error: no item key named '%s'\n", key)
		return
	}
	item := items[key].Copy()
	item.Count = count
	e := myecs.Manager.NewEntity()
	i := NewInteract(func(_ pixel.Vec, d *Dwarf) bool {
		return AddItem(d.Player.Inventory, item) == 0
	}, item.Sprite.Frame().W(), true)
	popUp := menus.NewPopUp("{symbol:player-interact}:pick up")
	popUp.Dist = item.Sprite.Frame().W()
	phys, trans := data.RandomPosAndVel(pos, 0., 0., math.Pi*0.5, math.Pi*0.25, 125., 10., random.Effects)
	coll := data.NewCollider(pixel.R(0., 0., item.Sprite.Frame().W(), item.Sprite.Frame().H()), data.Item)
	hp := &data.SimpleHealth{Immune: data.ItemImmunity1}
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Health, hp).
		AddComponent(myecs.Interact, i).
		AddComponent(myecs.PopUp, popUp).
		AddComponent(myecs.Drawable, item.Sprite).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
			if i.Interacted {
				myecs.Manager.DisposeEntity(e)
			} else if hp.Dead {
				if item.DeadFn != nil {
					item.DeadFn(trans.Pos)
				}
				myecs.Manager.DisposeEntity(e)
			}
			return false
		}))
	if item.Temp {
		e.AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
			myecs.AddEffect(e, data.NewBlink(2.))
			return true
		}, 8.))
		e.AddComponent(myecs.Temp, timing.New(10.))
	} else {
		e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	}
}