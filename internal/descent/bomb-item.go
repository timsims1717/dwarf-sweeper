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
	"dwarf-sweeper/pkg/world"
	"github.com/faiface/pixel"
	"math"
)

func CreateBombItem(pos pixel.Vec) {
	e := myecs.Manager.NewEntity()
	spr := img.Batchers[constants.EntityKey].Sprites["bomb_item"]
	i := NewInteract(func(pos pixel.Vec, d *Dwarf) bool {
		return d.Player.Inventory.AddItem(&player.InvItem{
			Name:   "bomb",
			Sprite: spr,
			OnUse: func(pos pixel.Vec) {
				tile := Descent.Cave.GetTile(pos)
				CreateBomb(tile.Transform.Pos)
			},
			Count: 3,
			Limit: 3,
		})
	}, spr.Frame().W() * 0.5, true)
	popUp := menus.NewPopUp("{symbol:player-interact}:pick up")
	popUp.Dist = (spr.Frame().W() + world.TileSize) * 0.5
	phys, trans := util.RandomPosAndVel(pos, 0., 0., math.Pi*0.5, math.Pi*0.25, 5., 2., random.Effects)
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
				CreateBomb(trans.Pos)
			}
			return false
		}))
}