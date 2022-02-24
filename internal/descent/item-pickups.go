package descent

import (
	"dwarf-sweeper/internal/constants"
	"dwarf-sweeper/internal/data"
	"dwarf-sweeper/internal/data/player"
	"dwarf-sweeper/internal/menus"
	"dwarf-sweeper/internal/myecs"
	"dwarf-sweeper/internal/random"
	"dwarf-sweeper/internal/util"
	"dwarf-sweeper/pkg/img"
	"dwarf-sweeper/pkg/timing"
	"github.com/faiface/pixel"
	"math"
)


func CreateBeerItem(pos pixel.Vec) {
	spr := img.Batchers[constants.EntityKey].Sprites["beer"]
	fn := func(pos pixel.Vec, d *Dwarf) bool {
		return d.Player.Inventory.AddItem(&player.InvItem{
			Name:   "beer",
			Sprite: spr,
			OnUse: func(_ pixel.Vec) {
				d.Entity.AddComponent(myecs.Healing, &data.Heal{
					TmpAmount: 2,
				})
			},
			Count: 1,
			Limit: 2,
		})
	}
	CreateItemPickUp(pos, fn, spr)
}


func CreateBubbleItem(pos pixel.Vec) {
	spr := img.Batchers[constants.EntityKey].Sprites["bubble_item"]
	fn := func(pos pixel.Vec, d *Dwarf) bool {
		return d.Player.Inventory.AddItem(&player.InvItem{
			Name:   "bubble",
			Sprite: spr,
			OnUse: func(_ pixel.Vec) {
				if d.Bubble != nil {
					d.Bubble.Pop()
				}
				bubble := &Bubble{
					Dwarf: d,
				}
				bubble.Create(pixel.Vec{})
			},
			Count: 1,
			Limit: 1,
			Sec:   BubbleSec,
		})
	}
	CreateItemPickUp(pos, fn, spr)
}


func CreateXRayItem(pos pixel.Vec) {
	spr := img.Batchers[constants.EntityKey].Sprites["x-ray-helmet"]
	fn := func(_ pixel.Vec, d *Dwarf) bool {
		return d.Player.Inventory.AddItem(&player.InvItem{
			Name:   "xray",
			Sprite: spr,
			OnUse: func(_ pixel.Vec) {
				StartXRayVision(d)
			},
			Count: 1,
			Limit: 1,
			Sec:   XRaySec,
		})
	}
	CreateItemPickUp(pos, fn, spr)
}


func CreateItemPickUp(pos pixel.Vec, fn func(pixel.Vec, *Dwarf) bool, spr *pixel.Sprite) {
	e := myecs.Manager.NewEntity()
	i := NewInteract(fn, spr.Frame().W(), true)
	popUp := menus.NewPopUp("{symbol:interact}:pick up")
	popUp.Dist = spr.Frame().W()
	phys, trans := util.RandomPosAndVel(pos, 0., 0., math.Pi*0.5, math.Pi*0.25, 125., 10., random.Effects)
	coll := data.NewCollider(pixel.R(0., 0., spr.Frame().W(), spr.Frame().H()), true, false)
	hp := &data.SimpleHealth{Immune: data.ItemImmunity}
	e.AddComponent(myecs.Transform, trans).
		AddComponent(myecs.Physics, phys).
		AddComponent(myecs.Collision, coll).
		AddComponent(myecs.Health, hp).
		AddComponent(myecs.Temp, timing.New(10.)).
		AddComponent(myecs.Interact, i).
		AddComponent(myecs.PopUp, popUp).
		AddComponent(myecs.Drawable, spr).
		AddComponent(myecs.Batch, constants.EntityKey).
		AddComponent(myecs.Func, data.NewTimerFunc(func() bool {
			myecs.AddEffect(e, data.NewBlink(2.))
			return true
		}, 8.)).
		AddComponent(myecs.Update, data.NewFrameFunc(func() bool {
			if hp.Dead {
				e.AddComponent(myecs.Temp, myecs.ClearFlag(true))
			}
			return false
		}))
}
